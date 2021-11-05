package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/emersion/go-smtp"
	"github.com/violetnorth/smtplib"
)

// Mail is mail from session. Mail command starts the
// mail exchange between the two SMTP servers.
func (s *smtpSession) Mail(from string, opts smtp.MailOptions) error {
	logger.Debugf("Session %s, from: %s", s.UUID, from)

	if strings.LastIndex(from, "@") <= 0 {
		logger.Errorf("Failed to parse from for %s, recipient with bad format %s", s.UUID, from)
		return smtpError(
			smtplib.StatusActionNotTakenMailboxNameNotAllowed,
			fmt.Sprintf(`Email Receiver: mailbox name error "%s"`, from),
		)
	}

	s.Opts = opts
	s.From = from
	return nil
}

// Rcpt is recipient from session.
func (s *smtpSession) Rcpt(recipient string) error {
	logger.Debugf("Session %s, rcpt: %s", s.UUID, recipient)

	if strings.LastIndex(recipient, "@") <= 0 {
		return smtpError(
			smtplib.StatusActionNotTakenMailboxNameNotAllowed,
			fmt.Sprintf(`Email Receiver: mailbox name not allowed "%s"`, recipient),
		)
	}

	mailHost := strings.Split(recipient, "@")[1]
	if mail, ok := mailsFind(mailHost); ok {
		// If the mail is in relay only mode, add the recipient to
		// the recipients list otherwise if the mail is hosted by
		// API, check if the inbox exists and save the inbox id.
		if mail.Relay {
			s.Recipients = append(s.Recipients, smtpRecipient{
				Address: recipient,
			})
		} else {
			for _, inbox := range mail.Inboxes {
				if inbox.Address == recipient {
					s.Recipients = append(s.Recipients, smtpRecipient{
						InboxID: inbox.ID,
						Address: recipient,
					})
					return nil
				}
			}

			return smtpError(
				smtplib.StatusActionNotTakenMailboxInaccessible,
				fmt.Sprintf(`Email Receiver: mailbox name unknown "%s"`, recipient),
			)
		}
	}

	return nil
}

// Data is data from session.
func (s *smtpSession) Data(r io.Reader) error {
	messageRaw, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Errorf("Failed to read data for %s, read error %v", s.UUID, err)
		return smtpError(
			smtplib.StatusActionAbortedLocalError,
			fmt.Sprintf("Email Receiver: email reading failed"),
		)
	}

	logger.Debugf("Session %s, data:\n%s", s.UUID, string(messageRaw))

	// Iterate all recipients and apply plugins and rules
	// for those mails. Send the email to upstream.
	for _, recipient := range s.Recipients {
		// Get the recipient host to find the mail.
		recipientSplit := strings.Split(recipient.Address, "@")
		if len(recipientSplit) != 2 {
			err := fmt.Errorf("address format error: %s", recipient.Address)
			logger.Errorf("Failed get recipient host %s, split host error %v", s.UUID, err)
			return smtpError(
				smtplib.StatusActionAbortedLocalError,
				fmt.Sprintf(`Email Receiver: format error for "%s"`, recipient.Address),
			)
		}

		mailHost := recipientSplit[1]
		mail, ok := mailsFind(mailHost)
		if !ok {
			return smtpError(
				smtplib.StatusActionNotTakenMailboxInaccessible,
				fmt.Sprintf(`Email Receiver: mail unknown "%s"`, mailHost),
			)
		}

		// Recreate the message from the raw message data in
		// order to start from scratch for each mail rules.
		message, err := smtpMessageParse(s, messageRaw)
		if err != nil {
			logger.Errorf("Failed to read envelople for %s, read error %v", s.UUID, err)
			return smtpError(
				smtplib.StatusActionAbortedLocalError,
				fmt.Sprintf("Email Receiver: email parsing failed"),
			)
		}

		// Relay the email message to upstreams, only if the mail is in
		// the firewall only configuration.
		if mail.Relay {
			messageEncoded := smtpMessageBuild(message)
			if err := smtpSend(messageEncoded, mail.Upstreams); err != nil {
				logger.Errorf("Failed to relay message for %s, all upstreams failed", s.UUID)
				return smtpError(
					smtplib.StatusActionNotTakenMailboxInaccessible,
					fmt.Sprintf(`Email Receiver: email relaying failed %s`, err.Error()),
				)
			}
		} else {
			msg, err := messageParse(recipient.InboxID, message)
			if err != nil {
				logger.Errorf("Failed to save message for %s, %v", s.UUID, err)
				return smtpError(
					smtplib.StatusActionAbortedLocalError,
					fmt.Sprintf(`Email: email receive failed due to internal error`),
				)
			}

			if err := messagesSave(msg); err != nil {
				logger.Errorf("Failed to send message for %s, %v", s.UUID, err)
				return smtpError(
					smtplib.StatusActionAbortedLocalError,
					fmt.Sprintf(`Email: email receive failed due to internal error`),
				)
			}
		}

	}

	return nil
}

// Reset is reset from session.
func (s *smtpSession) Reset() {
	logger.Debugf("Session %s, reset", s.UUID)
}

// Logout is log out from session.
func (s *smtpSession) Logout() error {
	logger.Debugf("Session %s, logout", s.UUID)
	return nil
}
