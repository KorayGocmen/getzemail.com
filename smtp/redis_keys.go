package main

import (
	"fmt"
	"strings"
)

// Redis Schema
// `known:<host>` <string>
// 		- "true" / "false" (is mail known)
// `mail:<host>` <string>
// 		- mail details (marshalled json string)

// redisKeyMailKnown is used to check if a mail is known.
func redisKeyMailKnown(host string) string {
	host = strings.ToLower(host)
	host = strings.TrimSpace(host)
	return fmt.Sprintf("known:%s", host)
}

// redisKeyMails is used to get all mails.
func redisKeyMails() string {
	return fmt.Sprintf("mail:*")
}

// redisKeyMail is used to find a mail with the provided host.
func redisKeyMail(host string) string {
	host = strings.ToLower(host)
	host = strings.TrimSpace(host)
	return fmt.Sprintf("mail:%s", host)
}
