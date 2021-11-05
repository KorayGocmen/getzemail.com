package main

import "time"

var (
	dispositionAttachment = "attachment"
	dispositionInline     = "inline"
)

// typeMailUpstream is the upstream associated with
// the mail and mail's MX records. It is the
// type used by the API.
type typeMailUpstream struct {
	Target   string `json:"target"`
	Priority int    `json:"priority"`
}

// typeMailInbox is the inboxes associated with
// the mail. Inboxes only exists when mail
// is not firewall only, meaning that it's
// hosted by API.
type typeMailInbox struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	Address     string `json:"address"`
}

// typeMail is the main mail struct.
type typeMail struct {
	ID        uint               `json:"id"`
	Host      string             `json:"host"`
	Relay     bool               `json:"relay"`
	Version   int                `json:"version,omitempty"`
	Inboxes   []typeMailInbox    `json:"mail_inboxes,omitempty"`
	Upstreams []typeMailUpstream `json:"mail_upstreams,omitempty"`
}

// Message related structs.

// typeMailMessageError is the mail message error
// struct that is used to track message errors.
type typeMailMessageError struct {
	ID            uint   `json:"id,omitempty"`
	MailMessageID uint   `json:"mail_message_id,omitempty"`
	Error         string `json:"error"`
}

// MailMessageRelation is the main mail message relation struct
// used by API mail message relation.
type typeMailMessageRelation struct {
	DisplayName string `json:"display_name,omitempty"`
	Address     string `json:"address"`
}

// MailMessageFile is the main mail message file struct
// used by API mail message file.
type typeMailMessageFile struct {
	ID            uint   `json:"id,omitempty"`
	MailMessageID uint   `json:"mail_message_id,omitempty"`
	Disposition   string `json:"disposition"`
	FileName      string `json:"file_name,omitempty"`
	ContentID     string `json:"content_id,omitempty"`
	ContentType   string `json:"content_type,omitempty"`
	Key           string `json:"key"`
	URL           string `json:"url"`
}

// MailMessage is the main mail message struct
// used by API mail message.
type typeMailMessage struct {
	ID      uint `json:"id,omitempty"`
	InboxID uint `json:"inbox_id,omitempty"`

	MessageID   string `json:"message_id"`
	InReplyToID string `json:"in_reply_to_id"`

	From typeMailMessageRelation   `json:"from"`
	To   []typeMailMessageRelation `json:"to"`
	Cc   []typeMailMessageRelation `json:"cc"`
	Bcc  []typeMailMessageRelation `json:"bcc"`

	Date    time.Time `json:"date"`
	Subject string    `json:"subject"`
	Text    string    `json:"text"`
	HTML    string    `json:"html"`

	IsDraft     bool `json:"is_draft"`
	IsDelivered bool `json:"is_delivered"`

	Files []typeMailMessageFile `json:"mail_message_files"`
}
