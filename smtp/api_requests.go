package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	apiClient = &http.Client{
		Timeout: timeDuration(config.API.Timeout),
	}
)

// apiRequestMailsRefresh gets all mail configs with the provided version map.
func apiRequestMailsRefresh(mailVersions map[uint]int) ([]typeMail, error) {
	logger.Printf("Api request mails refresh, refreshing %d mails", len(mailVersions))

	reqURL := fmt.Sprintf("%s/mails/refresh",
		config.API.BaseURL,
	)

	type reqBodyType struct {
		MailVersions map[uint]int `json:"mail_versions"`
	}

	reqBody := reqBodyType{
		MailVersions: mailVersions,
	}

	reqBodyMarshalled, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorln("Failed to request refresh mails, marshal request body error", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqBodyMarshalled))
	if err != nil {
		logger.Errorln("Failed to request refresh mails, create request error", err)
		return nil, err
	}
	req.Header.Set(headerContentType, applicationJSON)
	req.Header.Set(headerAuth, config.API.Secret)

	res, err := apiClient.Do(req)
	if err != nil {
		logger.Errorln("Failed to request refresh mails, do request error", err)
		return nil, err
	}
	defer res.Body.Close()

	type resBodyType struct {
		Success bool       `json:"success"`
		Error   string     `json:"error"`
		Mails   []typeMail `json:"mails"`
	}
	var b resBodyType

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorln("Failed to request refresh mails, read response body error", err)
		return nil, err
	}

	if err := json.Unmarshal(body, &b); err != nil {
		logger.Errorln("Failed to request refresh mails, unmarshal response body error", err)
		return nil, err
	}

	if !b.Success {
		err := errors.New(b.Error)
		logger.Errorln("Failed to request refresh mails, api returned error", b.Error)
		return nil, err
	}

	return b.Mails, nil
}

// apiRequestMailsGet gets mail the provided host.
func apiRequestMailsGet(host string) (typeMail, bool, error) {
	logger.Printf("Api request mails get, %s", host)

	reqURL := fmt.Sprintf("%s/mails/%s",
		config.API.BaseURL,
		host,
	)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		logger.Errorln("Failed to request get mail via host, create request error", err)
		return typeMail{}, false, err
	}
	req.Header.Set(headerAuth, config.API.Secret)

	res, err := apiClient.Do(req)
	if err != nil {
		logger.Errorln("Failed to request get mail via host, do request error", err)
		return typeMail{}, false, err
	}
	defer res.Body.Close()

	type resBodyType struct {
		Success bool     `json:"success"`
		Found   bool     `json:"found"`
		Error   string   `json:"error"`
		Mail    typeMail `json:"mail"`
	}
	var b resBodyType

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorln("Failed to request get mail via host, read response body error", err)
		return typeMail{}, false, err
	}

	if err := json.Unmarshal(body, &b); err != nil {
		logger.Errorln("Failed to request get mail via host, unmarshal response body error", err)
		return typeMail{}, false, err
	}

	if !b.Success {
		err := errors.New(b.Error)
		logger.Errorln("Failed to request get mail via host, api returned error", b.Error)
		return typeMail{}, false, err
	}

	// Mail is unknown to the API.
	if !b.Found {
		return typeMail{}, false, nil
	}

	return b.Mail, true, nil
}

// apiRequestMessagesInbound sends all inbound messages
// that were received from SMTP to the API.
func apiRequestMessagesInbound(mailMessage typeMailMessage) error {
	logger.Printf("Api send inbound mail message")

	reqURL := fmt.Sprintf("%s/smtp/inbound",
		config.API.BaseURL,
	)

	type reqBodyType struct {
		MailMessage typeMailMessage `json:"mail_message"`
	}

	reqBody := reqBodyType{
		MailMessage: mailMessage,
	}

	reqBodyMarshalled, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorln("Failed to send inbound mail message, marshal request body error", err)
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqBodyMarshalled))
	if err != nil {
		logger.Errorln("Failed to send inbound mail message, create request error", err)
		return err
	}
	req.Header.Set(headerContentType, applicationJSON)
	req.Header.Set(headerAuth, config.API.Secret)

	res, err := apiClient.Do(req)
	if err != nil {
		logger.Errorln("Failed to send inbound mail message, do request error", err)
		return err
	}
	defer res.Body.Close()

	type resBodyType struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	var b resBodyType

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorln("Failed to send inbound mail message, read response body error", err)
		return err
	}

	if err := json.Unmarshal(body, &b); err != nil {
		logger.Errorln("Failed to send inbound mail message, unmarshal response body error", err)
		return err
	}

	if !b.Success {
		err := errors.New(b.Error)
		logger.Errorln("Failed to send inbound mail message, api returned error", b.Error)
		return err
	}

	return nil
}

// apiRequestMessagesOutbounds gets all outbound mail messages
// from the API and send them via the SMTP server.
func apiRequestMessagesOutbounds() ([]typeMailMessage, error) {
	logger.Printf("Api get outbound mail messages")

	reqURL := fmt.Sprintf("%s/smtp/outbound",
		config.API.BaseURL,
	)

	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		logger.Errorln("Failed to get outbound mail messages, create request error", err)
		return nil, err
	}
	req.Header.Set(headerContentType, applicationJSON)
	req.Header.Set(headerAuth, config.API.Secret)

	res, err := apiClient.Do(req)
	if err != nil {
		logger.Errorln("Failed to get outbound mail messages, do request error", err)
		return nil, err
	}
	defer res.Body.Close()

	type resBodyType struct {
		Success      bool              `json:"success"`
		Error        string            `json:"error"`
		MailMessages []typeMailMessage `json:"mail_messages"`
	}
	var b resBodyType

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorln("Failed to get outbound mail messages, read response body error", err)
		return nil, err
	}

	if err := json.Unmarshal(body, &b); err != nil {
		logger.Errorln("Failed to get outbound mail messages, unmarshal response body error", err)
		return nil, err
	}

	if !b.Success {
		err := errors.New(b.Error)
		logger.Errorln("Failed to get outbound mail messages, api returned error", b.Error)
		return nil, err
	}

	return b.MailMessages, nil
}
