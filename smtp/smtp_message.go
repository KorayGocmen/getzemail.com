package main

import (
	"bytes"
	"net/mail"
	"time"

	"github.com/jhillyerd/enmime"
)

type smtpMessageSession struct {
	UUID string
	From string
}

type smtpMessage struct {
	Raw     []byte
	Session smtpMessageSession

	MessageID string
	InReplyTo string

	From        mail.Address
	ReplyTo     mail.Address
	To, Cc, Bcc []mail.Address

	Date    time.Time
	Subject string
	Text    string
	HTML    string

	Inlines     []*enmime.Part
	Attachments []*enmime.Part
}

// smtpMessageParse decodes a raw message buffer into an smtp message.
func smtpMessageParse(sess *smtpSession, messageRaw []byte) (smtpMessage, error) {
	message := smtpMessage{
		Session: smtpMessageSession{
			UUID: sess.UUID,
			From: sess.From,
		},
		Raw: messageRaw,
	}

	messageEnvelope, err := enmime.ReadEnvelope(bytes.NewReader(messageRaw))
	if err != nil {
		return message, err
	}

	message.MessageID = smtpIDHeaderDecode(messageEnvelope.GetHeader("Message-Id"))
	message.InReplyTo = smtpIDHeaderDecode(messageEnvelope.GetHeader("In-Reply-To"))

	from, err := mail.ParseAddress(messageEnvelope.GetHeader("From"))
	if from != nil && err == nil {
		message.From = *from
	}

	replyTo, err := mail.ParseAddress(messageEnvelope.GetHeader("Reply-To"))
	if replyTo != nil && err == nil {
		message.ReplyTo = *replyTo
	}

	toAddrs, _ := messageEnvelope.AddressList("To")
	for _, toAddr := range toAddrs {
		if toAddr != nil {
			message.To = append(message.To, *toAddr)
		}
	}

	ccAddrs, _ := messageEnvelope.AddressList("Cc")
	for _, ccAddr := range ccAddrs {
		if ccAddr != nil {
			message.Cc = append(message.Cc, *ccAddr)
		}
	}

	bccAddrs, _ := messageEnvelope.AddressList("Bcc")
	for _, bccAddr := range bccAddrs {
		if bccAddr != nil {
			message.Bcc = append(message.Bcc, *bccAddr)
		}
	}

	date, err := mail.ParseDate(messageEnvelope.GetHeader("Date"))
	if err == nil {
		message.Date = date
	}

	message.Subject = messageEnvelope.GetHeader("Subject")
	message.Text = messageEnvelope.Text
	message.HTML = messageEnvelope.HTML

	message.Inlines = messageEnvelope.Inlines
	message.Attachments = messageEnvelope.Attachments

	return message, nil
}

// smtpMessageBuild encodes an smtp message into an enmime mailbuilder.
// The encoded mail builder can be used to mail to upstream.
func smtpMessageBuild(message smtpMessage) enmime.MailBuilder {
	builder := enmime.MailBuilder{}

	builder = builder.From(message.From.Name, message.From.Address).
		ReplyTo(message.ReplyTo.Name, message.ReplyTo.Address).
		ToAddrs(message.To).
		CCAddrs(message.Cc).
		BCCAddrs(message.Bcc).
		Date(message.Date).
		Subject(message.Subject).
		Text([]byte(message.Text)).
		HTML([]byte(message.HTML))

	for _, i := range message.Inlines {
		builder = builder.AddInline(i.Content, i.ContentType, i.FileName, i.ContentID)
	}

	for _, a := range message.Attachments {
		builder = builder.AddAttachment(a.Content, a.ContentType, a.FileName)
	}

	return builder
}
