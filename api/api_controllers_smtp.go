package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// apiControllersSmtpInbound receives an inbound mail message from the SMTP server.
func apiControllersSmtpInbound(c *gin.Context) {
	var req typeApiReqMailMessagesInbound
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf("failed to save mail message inbound: bind json error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	var mailInbox MailInbox
	err := db.First(&mailInbox, "id = ?", req.MailMessage.InboxID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Mail inbox not found",
			})
			return
		}

		logger.Errorf("failed to save mail message inbound: find inbox error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	db.Transaction(func(tx *gorm.DB) error {
		mailMessage := MailMessage{
			MailInboxID: mailInbox.ID,

			MessageID:   req.MailMessage.MessageID,
			InReplyToID: req.MailMessage.InReplyToID,

			Subject: req.MailMessage.Subject,
			Text:    req.MailMessage.Text,
			HTML:    req.MailMessage.HTML,
		}

		if err := tx.Create(&mailMessage).Error; err != nil {
			logger.Errorf("failed to create mail message: db create error: %v", err)
			return err
		}

		var mailMessageFiles []MailMessageFile
		for _, file := range req.MailMessage.Files {
			mailMessageFile := MailMessageFile{
				MailMessageID: mailMessage.ID,
				URL:           file.URL,
				Disposition:   file.Disposition,
				Key:           file.Key,
				FileName:      file.FileName,
				ContentID:     file.ContentID,
				ContentType:   file.ContentType,
			}

			mailMessageFiles = append(mailMessageFiles, mailMessageFile)
		}

		if err := tx.CreateInBatches(mailMessageFiles, len(mailMessageFiles)).Error; err != nil {
			logger.Errorf("failed to create mail message files: db create files error: %v", err)
			return err
		}

		var mailMessageRelations []MailMessageRelation
		for _, to := range req.MailMessage.To {
			to.MailMessageID = mailMessage.ID
			to.Type = mailMessageRelationTypeTo
			mailMessageRelations = append(mailMessageRelations, to)
		}

		for _, cc := range req.MailMessage.Cc {
			cc.MailMessageID = mailMessage.ID
			cc.Type = mailMessageRelationTypeCc
			mailMessageRelations = append(mailMessageRelations, cc)
		}

		for _, bcc := range req.MailMessage.Bcc {
			bcc.MailMessageID = mailMessage.ID
			bcc.Type = mailMessageRelationTypeBcc
			mailMessageRelations = append(mailMessageRelations, bcc)
		}

		if err := tx.CreateInBatches(mailMessageRelations, len(mailMessageRelations)).Error; err != nil {
			logger.Errorf("failed to create mail message files: db create relations error: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
	})
}

// apiControllersSmtpOutbound returns outbound mail messages
func apiControllersSmtpOutbound(c *gin.Context) {
	// Not implemented.

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":       true,
		"mail_messages": nil,
	})
}
