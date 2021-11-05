package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jhillyerd/enmime"
)

// API and SMTP message communication:
// 1. SMTP --- inbound --> API:
// 	- SMTP receives SMTP message on the SMTP server and
//	sends it to the API right away.
// 2. SMTP <-- outbound --- API :
// 	- SMTP receives outbound mail messages from the API
//	every "outbound_every" ticker and tries to deliver them.

func initMessages() {
	go func() {
		logger.Println("Creating messages send timer")

		outboundEvery := timeDuration(config.Messages.OutboundEvery)
		outboundTicker := time.NewTicker(outboundEvery)
		defer outboundTicker.Stop()

		// Receive the messages first time when the application
		// starts, the timer will refresh after this one.
		messagesOutbound()

		for {
			select {
			case <-outboundTicker.C:
				messagesOutbound()
			}
		}
	}()
}

// messagesSave sends a message to the API to be saved.
func messagesSave(message typeMailMessage) error {
	logger.Println("Sending message", message.MessageID)

	if err := apiRequestMessagesInbound(message); err != nil {
		logger.Errorln("Failed to send inbound message, api request error", err)
		return err
	}

	return nil
}

// messagesOutbound is called every messages receive every
// timer and receive the outbound messages from the API.
func messagesOutbound() {
	logger.Println("Receiving outbound messages")

	messages, err := apiRequestMessagesOutbounds()
	if err != nil {
		logger.Errorln("Failed to get outbound messages, api request error", err)
		return
	}

	logger.Debugf("Attempting to send %d messages", len(messages))

	var (
		messageIDs    []uint
		messageErrors []typeMailMessageError
	)

	for _, message := range messages {
		relations := message.To
		relations = append(relations, message.Cc...)
		relations = append(relations, message.Bcc...)

		messageEncoded, err := messageBuild(message)
		if err != nil {
			logger.Errorln("Failed to convert message to smtp", message.MessageID, err)
			messageErrors = append(messageErrors, typeMailMessageError{
				MailMessageID: message.ID,
				Error:         fmt.Sprintf("email format error: %s", err.Error()),
			})
			continue
		}

		mailHostsSent := make(map[string]bool)
		for _, relation := range relations {
			address := relation.Address

			if strings.LastIndex(address, "@") <= 0 {
				err := fmt.Errorf(`email address "%s" is not valid`, address)
				logger.Errorln("Failed to get email address host", message.MessageID, err)
				messageErrors = append(messageErrors, typeMailMessageError{
					MailMessageID: message.ID,
					Error:         fmt.Sprintf(`email address format error: %s`, err.Error()),
				})
				continue
			}

			mailHost := strings.Split(address, "@")[1]
			if _, ok := mailHostsSent[mailHost]; ok {
				// Mail already sent to the mail host. If there are multiple
				// emails with the same mail host, only send it once.
				continue
			}
			mailHostsSent[mailHost] = true

			mxRecords, err := net.LookupMX(mailHost)
			if err != nil {
				err = fmt.Errorf(`email address "%s" host's MX lookup failed due to %s`, address, err.Error())
				logger.Errorln("Failed to lookup MX records for", message.MessageID, err)
				messageErrors = append(messageErrors, typeMailMessageError{
					MailMessageID: message.ID,
					Error:         fmt.Sprintf("email address host error: %s", err.Error()),
				})
				continue
			}

			var upstreams []typeMailUpstream
			for _, mx := range mxRecords {
				upstreams = append(upstreams, typeMailUpstream{
					Target:   mx.Host,
					Priority: int(mx.Pref),
				})
			}

			if err := smtpSend(messageEncoded, upstreams); err != nil {
				err = fmt.Errorf(`email address "%s" delivery failed due to %s`, address, err.Error())
				logger.Errorln("Failed to lookup MX records", message.MessageID, err)
				messageErrors = append(messageErrors, typeMailMessageError{
					MailMessageID: message.ID,
					Error:         fmt.Sprintf("email address delivery error: %s", err.Error()),
				})
				continue
			}
		}

		messageIDs = append(messageIDs, message.ID)
	}
}

// messageParse parses an SMTP message to API Message
// which can be used by the api handler.
func messageParse(inboxID uint, message smtpMessage) (typeMailMessage, error) {
	logger.Debugln("Parsing api message", message.MessageID)

	msg := typeMailMessage{
		InboxID: inboxID,

		MessageID:   message.MessageID,
		InReplyToID: message.InReplyTo,

		From: typeMailMessageRelation{
			DisplayName: message.From.Name,
			Address:     message.From.Address,
		},

		Date:    message.Date.UTC(),
		Subject: message.Subject,

		Text: stringsFirstNChars(message.Text, 255),
		HTML: stringsFirstNChars(message.HTML, 255),
	}

	// Upload the MIME format.
	s3UploadOptsMIME := s3UploadOpts{
		Bucket: config.S3Emails.Bucket,
		Key:    fmt.Sprintf("%s/%s", msg.MessageID, uploadFileNameMIME),

		ACL:         config.S3Emails.ACL,
		ContentType: contentTypeMIME,

		MetaData: map[string]string{
			"Message-Id": msg.MessageID,
		},
	}

	if _, err := s3Upload(s3UploadOptsMIME, bytes.NewReader(message.Raw)); err != nil {
		logger.Errorln("Failed to upload mime file to S3", err)
		return msg, err
	}

	// Upload the Text.
	s3UploadOptsText := s3UploadOpts{
		Bucket: config.S3Emails.Bucket,
		Key:    fmt.Sprintf("%s/%s", msg.MessageID, uploadFileNameText),

		ACL:         config.S3Emails.ACL,
		ContentType: contentTypeText,

		MetaData: map[string]string{
			"Message-Id": msg.MessageID,
		},
	}

	if _, err := s3Upload(s3UploadOptsText, strings.NewReader(message.Text)); err != nil {
		logger.Errorln("Failed to upload text file to S3", err)
		return msg, err
	}

	// Upload the HTML.
	s3UploadOptsHTML := s3UploadOpts{
		Bucket: config.S3Emails.Bucket,
		Key:    fmt.Sprintf("%s/%s", msg.MessageID, uploadFileNameHTML),

		ACL:         config.S3Emails.ACL,
		ContentType: contentTypeHTML,

		MetaData: map[string]string{
			"Message-Id": msg.MessageID,
		},
	}

	if _, err := s3Upload(s3UploadOptsHTML, strings.NewReader(message.HTML)); err != nil {
		logger.Errorln("Failed to upload HTML file to S3", err)
		return msg, err
	}

	// Message parts are the inlines and attachments.
	messageParts := message.Inlines
	messageParts = append(messageParts, message.Attachments...)

	for _, part := range messageParts {
		key := fmt.Sprintf("%s/%s/%s", msg.MessageID, part.Disposition, part.ContentID)
		s3UploadOpts := s3UploadOpts{
			Bucket: config.S3Emails.Bucket,
			Key:    key,

			ACL:         config.S3Emails.ACL,
			ContentType: part.ContentType,

			MetaData: map[string]string{
				"Message-Id":          msg.MessageID,
				"Content-Id":          part.ContentID,
				"Content-Disposition": part.Disposition,
			},
		}

		s3UploadOut, err := s3Upload(s3UploadOpts, bytes.NewReader(part.Content))
		if err != nil {
			logger.Errorln("Failed to upload message part file to S3", err)
			return msg, err
		}

		messageFile := typeMailMessageFile{
			Disposition: part.Disposition,
			FileName:    part.FileName,
			ContentID:   part.ContentID,
			ContentType: part.ContentType,
			Key:         key,
			URL:         s3UploadOut.Location,
		}
		msg.Files = append(msg.Files, messageFile)
	}

	for _, to := range message.To {
		msg.To = append(msg.To, typeMailMessageRelation{
			DisplayName: to.Name,
			Address:     to.Address,
		})
	}

	for _, cc := range message.Cc {
		msg.Cc = append(msg.Cc, typeMailMessageRelation{
			DisplayName: cc.Name,
			Address:     cc.Address,
		})
	}

	for _, bcc := range message.Cc {
		msg.Bcc = append(msg.Bcc, typeMailMessageRelation{
			DisplayName: bcc.Name,
			Address:     bcc.Address,
		})
	}

	return msg, nil
}

// messageBuild converts Violetnorth Message to SMTP message
// which can be used by the SMTP handler.
func messageBuild(message typeMailMessage) (enmime.MailBuilder, error) {
	builder := enmime.MailBuilder{}

	builder = builder.Header("Message-Id", smtpIDHeaderEncode(message.MessageID))
	builder = builder.Header("In-Reply-To", smtpIDHeaderEncode(message.InReplyToID))

	builder = builder.
		From(message.From.DisplayName, message.From.Address).
		Date(message.Date).
		Subject(message.Subject)

	for _, to := range message.To {
		builder = builder.To(to.DisplayName, to.Address)
	}

	for _, cc := range message.Cc {
		builder = builder.CC(cc.DisplayName, cc.Address)
	}

	for _, bcc := range message.Bcc {
		builder = builder.BCC(bcc.DisplayName, bcc.Address)
	}

	// Download the text, assign it to the message.
	text, err := s3Download(s3DownloadOpts{
		Bucket: config.S3Emails.Bucket,
		Key:    fmt.Sprintf("%s/%s", message.MessageID, uploadFileNameText),
	})
	if err != nil {
		logger.Errorln("Failed to download text from S3", err)
		return builder, err
	}

	if strings.TrimSpace(string(text)) != "" {
		builder = builder.Text(text)
	}

	// Download the HTML, assign it to the message.
	html, err := s3Download(s3DownloadOpts{
		Bucket: config.S3Emails.Bucket,
		Key:    fmt.Sprintf("%s/%s", message.MessageID, uploadFileNameHTML),
	})
	if err != nil {
		logger.Errorln("Failed to download html from S3", err)
		return builder, err
	}

	if strings.TrimSpace(string(html)) != "" {
		builder = builder.HTML(html)
	}

	// Download the message parts.
	for _, file := range message.Files {
		data, err := s3Download(s3DownloadOpts{
			Bucket: config.S3Emails.Bucket,
			Key:    fmt.Sprintf("%s/%s/%s", message.MessageID, file.Disposition, file.ContentID),
		})
		if err != nil {
			logger.Errorln("Failed to download message part from S3", err)
			return builder, err
		}

		if file.Disposition == dispositionAttachment {
			builder = builder.AddAttachment(data, file.ContentType, file.FileName)
		} else if file.Disposition == dispositionInline {
			builder = builder.AddInline(data, file.ContentType, file.FileName, file.ContentID)
		}
	}

	return builder, nil
}
