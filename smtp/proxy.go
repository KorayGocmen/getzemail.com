package main

import (
	"fmt"

	"github.com/emersion/go-smtp"
)

var (
	serverSMTP *smtp.Server
)

func smtpRelay() {
	be := &smtpBackend{}
	serverSMTP = smtp.NewServer(be)

	serverSMTP.Addr = fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	serverSMTP.Domain = config.Server.Domain

	serverSMTP.ReadTimeout = timeDuration(config.Server.TimeoutRead)
	serverSMTP.WriteTimeout = timeDuration(config.Server.TimeoutWrite)

	serverSMTP.MaxMessageBytes = config.Server.MaxMessageBytes
	serverSMTP.MaxRecipients = config.Server.MaxRecipients
	serverSMTP.AllowInsecureAuth = config.Server.AllowInsecureAuth

	go smtpListen(serverSMTP)
}

func smtpListen(s *smtp.Server) {
	logger.Println("(SMTP) Listening on", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatalln("Failed to serve smtp", err)
	}
}
