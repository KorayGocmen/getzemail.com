package main

import (
	"errors"
	"fmt"
	"net/http"

	gomail "net/mail"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// apiControllersMailInboxes returns a mail inbox and all mail
// messages in that mail inbox with the provided host and address.
func apiControllersMailInboxes(c *gin.Context) {
	mailHost := c.Param("mailHost")
	mailInboxAddr := c.Param("mailInboxAddr")
	mailInboxFullAddr := fmt.Sprintf("%s@%s", mailInboxAddr, mailHost)

	var mail Mail
	if err := db.First(&mail, "host = ?", mailHost).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Mail not found",
		})
		return
	}

	var mailInbox MailInbox
	err := db.
		Preload("MailMessages", func(db *gorm.DB) *gorm.DB {
			return db.Order("mail_messages.id DESC")
		}).
		Preload("MailMessages.MailMessageFiles").
		Preload("MailMessages.MailMessageRelations").
		First(&mailInbox, "mail_id = ? and address = ?", mail.ID, mailInboxFullAddr).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Mail inbox not found",
			})
			return
		}

		logger.Errorf("failed to get mail messages: %s: %w", mailHost, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"mail_inbox": mailInbox,
	})
}

// apiControllersMailInboxesCreate creates a mail from the provided host.
func apiControllersMailInboxesCreate(c *gin.Context) {
	mailHost := c.Param("mailHost")

	var mail Mail
	if err := db.First(&mail, "host = ?", mailHost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Mail not found",
			})
			return
		}

		logger.Errorf("failed to find mail via host: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	var req typeApiReqMailInboxesCreate
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf("failed to create mail inbox: bind json error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	inboxAddress, err := gomail.ParseAddress(fmt.Sprintf("%s@%s", req.Address, mail.Host))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Invalid email address: %v", err),
		})
		return
	}

	var mailInboxFound MailInbox
	if err := db.First(&mailInboxFound, "address = ?", inboxAddress.Address).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Mail inbox with the same address already exists",
		})
		return
	}

	mailInbox := MailInbox{
		MailID:      mail.ID,
		Address:     inboxAddress.Address,
		DisplayName: req.DisplayName,
	}

	if err := db.Create(&mailInbox).Error; err != nil {
		logger.Errorf("failed to create mail inbox: db create error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success":    true,
		"mail_inbox": mailInbox,
	})
}
