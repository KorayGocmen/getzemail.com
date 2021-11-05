package main

import (
	"net/url"
	"strings"

	"github.com/emersion/go-smtp"
)

// smtpError creates an smtp error with the provided code and message.
func smtpError(code int, message string) error {
	return &smtp.SMTPError{
		Code:         code,
		Message:      message,
		EnhancedCode: smtp.EnhancedCodeNotSet,
	}
}

func smtpIDHeaderDecode(v string) string {
	if v == "" {
		return v
	}

	v = strings.TrimLeft(v, "<")
	v = strings.TrimRight(v, ">")
	if r, err := url.QueryUnescape(v); err == nil {
		v = r
	}

	return v
}

func smtpIDHeaderEncode(v string) string {
	v = url.QueryEscape(v)
	return "<" + strings.Replace(v, "%40", "@", -1) + ">"
}
