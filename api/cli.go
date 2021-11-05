package main

import (
	"fmt"
	"os"
)

func initCLI() {
	if flagDBMigrate {
		initDB()

		err := db.AutoMigrate(
			&Mail{},
			&MailInbox{},
			&MailUpstream{},
			&MailMessage{},
			&MailMessageRelation{},
			&MailMessageFile{},
			&MailMessageError{},
		)
		if err != nil {
			err = fmt.Errorf("failed to migrate database: %w", err)
			logger.Fatalln("Error:", err)
		}

		logger.Println("Success: database migrated")
		os.Exit(0)
	}
}
