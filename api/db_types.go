package main

import (
	"time"

	"gorm.io/gorm"
)

const (
	dbDriverMysql    = "mysql"
	dbDriverPostgres = "postgres"
	dbDriverSqlite   = "sqlite"

	mailMessageRelationTypeTo  = "to"
	mailMessageRelationTypeCc  = "cc"
	mailMessageRelationTypeBcc = "bcc"
)

type Mail struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	Host    string `gorm:"column:host" json:"host"`
	Relay   bool   `gorm:"column:relay" json:"relay"`
	Version int    `gorm:"column:version" json:"version"`

	MailUpstreams []MailUpstream `gorm:"foreignkey:mail_id" json:"mail_upstreams,omitempty"`
	MailInboxes   []MailInbox    `gorm:"foreignkey:mail_id" json:"mail_inboxes,omitempty"`
}

type MailUpstream struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	MailID   uint   `gorm:"column:mail_id" json:"mail"`
	Target   string `gorm:"column:target" json:"target"`
	Priority int    `gorm:"column:priority" json:"priority"`
}

type MailInbox struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	MailID      uint   `gorm:"column:mail_id" json:"mail"`
	Address     string `gorm:"column:address" json:"address"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`

	MailMessages []MailMessage `gorm:"foreignkey:mail_inbox_id" json:"mail_messages,omitempty"`
}

type MailMessage struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	MailInboxID uint `gorm:"column:mail_inbox_id" json:"mail_inbox"`

	MessageID   string `gorm:"column:message_id" json:"message_id"`
	InReplyToID string `gorm:"column:in_reply_to_id" json:"in_reply_to_id"`

	Subject string `gorm:"column:subject" json:"subject"`
	Text    string `gorm:"column:text" json:"text"`
	HTML    string `gorm:"column:html" json:"html"`

	MailMessageRelations []MailMessageRelation `gorm:"foreignkey:mail_message_id" json:"mail_message_relations,omitempty"`
	MailMessageFiles     []MailMessageFile     `gorm:"foreignkey:mail_message_id" json:"mail_message_files,omitempty"`
	MailMessageErrors    []MailMessageError    `gorm:"foreignkey:mail_message_id" json:"mail_message_errors,omitempty"`

	TextURL string `json:"text_url,omitempty"`
	HtmlURL string `json:"html_url,omitempty"`
}

type MailMessageRelation struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	MailMessageID uint   `gorm:"column:mail_message_id" json:"mail_message"`
	Type          string `gorm:"column:type" json:"type"`
	Address       string `gorm:"column:address" json:"address"`
	DisplayName   string `gorm:"column:display_name" json:"display_name"`
}

type MailMessageFile struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	MailMessageID uint   `gorm:"column:mail_message_id" json:"mail_message"`
	URL           string `gorm:"column:url" json:"url"`
	Disposition   string `gorm:"column:disposition" json:"disposition"`
	Key           string `gorm:"column:key" json:"key"`
	FileName      string `gorm:"column:file_name" json:"file_name"`
	ContentID     string `gorm:"column:content_id" json:"content_id"`
	ContentType   string `gorm:"column:content_type" json:"content_type"`
}

type MailMessageError struct {
	ID        uint           `gorm:"primaryKey,column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`

	MailMessageID uint   `gorm:"column:mail_message_id" json:"mail_message"`
	Error         string `gorm:"column:error" json:"error"`
}
