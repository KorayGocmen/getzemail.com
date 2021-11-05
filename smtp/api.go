package main

import (
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
)

// initMails initializes the refresh timer which will
// refresh all known mails if their version is changed.
func initMails() {
	go func() {
		logger.Println("Creating mails refresh timer")

		every := timeDuration(config.Mails.RefreshEvery)
		ticker := time.NewTicker(every)
		defer ticker.Stop()

		// Refresh the mails first time when the application
		// starts, the timer will refresh after this one.
		mailsRefresh()

		for {
			select {
			case <-ticker.C:
				mailsRefresh()
			}
		}
	}()
}

// mailsFind finds a mail with the provided host.
// If the mail is not found, check API for the provided host.
func mailsFind(host string) (typeMail, bool) {
	if host == "" {
		logger.Errorln("Failed to get mail, host is nil")
		return typeMail{}, false
	}

	// Check if mail is known.
	mailKnown, err := redisdb.Get(redisKeyMailKnown(host)).Result()
	if err != nil && err != redis.Nil {
		logger.Errorln("Failed to get mail known from redis", host, err)
		return typeMail{}, false
	}

	// If mail is unknown, mail has been requested recently.
	// In order to reduce the number of requests made to the API
	// unknown mails won't be requested until the key expires.
	// When the key is expired mailKnown will be nil and function
	// will continue down the line and make the API request.
	if mailKnown == "false" {
		return typeMail{}, false
	}

	mailRaw, err := redisdb.Get(redisKeyMail(host)).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Errorln("Failed to get mail from redis", host, err)
			return typeMail{}, false
		}

		mail, known, err := apiRequestMailsGet(host)
		if err != nil {
			logger.Errorln("Failed to request mail from api", host, err)
			return typeMail{}, false
		}

		mailsAdd(host, mail, known)
		return mail, known
	}

	var mail typeMail
	if err := json.Unmarshal([]byte(mailRaw), &mail); err != nil {
		logger.Errorln("Failed to unmarshal mail from redis", err)
		return typeMail{}, false
	}

	return mail, true
}

// mailsRefresh is called every mails refresh timer and
// requests the known mails with versions from API. Only
// the mails which has a different version is returned by the
// API, therefore these mails are refreshed.
func mailsRefresh() {
	logger.Println("Refreshing mails")

	mailKeys, err := redisdb.Keys(redisKeyMails()).Result()
	if err != nil {
		logger.Errorln("Failed to get mail keys from redis", err)
		return
	}

	if len(mailKeys) == 0 {
		logger.Debugln("No mails to refresh")
		return
	}

	mailsRawInterf, err := redisdb.MGet(mailKeys...).Result()
	if err != nil {
		logger.Errorln("Failed to get multiple mails from redis", err)
		return
	}

	var mails []typeMail
	for _, mailRawInterf := range mailsRawInterf {
		if mailRaw, ok := mailRawInterf.(string); ok {
			var mail typeMail
			json.Unmarshal([]byte(mailRaw), &mail)
			mails = append(mails, mail)
		}
	}

	mailVersions := make(map[uint]int)
	for _, mail := range mails {
		mailVersions[mail.ID] = mail.Version
	}

	if len(mailVersions) > 0 {
		mails, err := apiRequestMailsRefresh(mailVersions)
		if err != nil {
			logger.Errorln("Failed to refresh mails, api request error", err)
			return
		}

		for _, mail := range mails {
			logger.Debugln("Mails refresh, mail has a new version", mail.Host)
			mailsAdd(mail.Host, mail, true)
		}
	}
}

// mailsAdd initializes the mail data and registers the
// mail under Mails.
func mailsAdd(host string, mail typeMail, known bool) error {
	logger.Println("Adding mail", host)

	var (
		mailKnownString = fmt.Sprintf("%t", known)
		expiration      = timeDuration(config.Mails.TTL)
	)

	if err := redisdb.Set(redisKeyMailKnown(host), mailKnownString, expiration).Err(); err != nil {
		logger.Errorln("Failed to set mail on redis", host, err)
		return err
	}

	if known {
		mailMarshalled, err := json.Marshal(mail)
		if err != nil {
			logger.Errorln("Failed to marshal mail to string", host, err)
			return err
		}

		if err := redisdb.Set(redisKeyMail(host), string(mailMarshalled), expiration).Err(); err != nil {
			logger.Errorln("Failed to set mail on redis", host, err)
			return err
		}
	}

	return nil
}

// mailsDel deletes the provided mail from the known Mails.
func mailsDel(mail typeMail) {
	logger.Println("Deleting mail", mail.Host)

	if err := redisdb.Del(redisKeyMail(mail.Host)).Err(); err != nil {
		logger.Errorln("Failed to delete mail on redis", mail.Host, err)
	}
}
