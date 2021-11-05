package main

import (
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
)

// The smtpBackend implements SMTP server methods.
type smtpBackend struct{}

// Login handles a login command with username and password.
func (bkd *smtpBackend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	session := smtpSession{
		UUID: uuid.New().String(),

		Conn: smtpConn{
			LocalAddr:  state.LocalAddr,
			RemoteAddr: state.RemoteAddr,
			TLS:        state.TLS,
		},

		Auth: smtpAuth{
			Anonymous: false,
			Username:  username,
			Password:  password,
		},
	}

	return &session, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *smtpBackend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	session := smtpSession{
		UUID: uuid.New().String(),

		Conn: smtpConn{
			LocalAddr:  state.LocalAddr,
			RemoteAddr: state.RemoteAddr,
			TLS:        state.TLS,
		},

		Auth: smtpAuth{
			Anonymous: true,
		},
	}

	return &session, nil
}
