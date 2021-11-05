package main

import (
	"crypto/tls"
	"net"

	"github.com/emersion/go-smtp"
)

const (
	uploadFileNameMIME = "mime"
	uploadFileNameText = "text"
	uploadFileNameHTML = "html"

	contentTypeMIME = "message/rfc822"
	contentTypeText = "text/plain"
	contentTypeHTML = "text/html"
)

// smtpConn is the conn type in smtp session.
type smtpConn struct {
	LocalAddr  net.Addr
	RemoteAddr net.Addr
	TLS        tls.ConnectionState
}

// smtpAuth is the auth type in smtp session.
type smtpAuth struct {
	Anonymous bool
	Username  string
	Password  string
}

// smtpSession is returned after successful login.
type smtpSession struct {
	UUID string
	Opts smtp.MailOptions

	From       string
	Recipients []smtpRecipient

	Conn smtpConn
	Auth smtpAuth
}

type smtpRecipient struct {
	Address string
	InboxID uint
}
