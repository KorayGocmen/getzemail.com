package main

import "time"

type typeApiReqMailsCreate struct {
	Host  string `json:"host"`
	Relay bool   `json:"relay"`
}

type typeApiReqMailMessagesInbound struct {
	MailMessage struct {
		InboxID uint `json:"inbox_id,omitempty"`

		MessageID   string `json:"message_id"`
		InReplyToID string `json:"in_reply_to_id"`

		From MailMessageRelation   `json:"from"`
		To   []MailMessageRelation `json:"to"`
		Cc   []MailMessageRelation `json:"cc"`
		Bcc  []MailMessageRelation `json:"bcc"`

		Date    time.Time `json:"date"`
		Subject string    `json:"subject"`
		Text    string    `json:"text"`
		HTML    string    `json:"html"`

		Files []MailMessageFile `json:"mail_message_files"`
	} `json:"mail_message"`
}

type typeApiReqMailsRefresh struct {
	MailVersions map[int]int `json:"mail_versions"`
}

type typeApiReqMailInboxesCreate struct {
	Address     string `json:"address"`
	DisplayName string `json:"display_name"`
}
