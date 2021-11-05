package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// apiControllersMailsCreate creates a mail from the provided host.
func apiControllersMailsCreate(c *gin.Context) {
	var req typeApiReqMailsCreate
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf("failed to create mail: bind json error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	if req.Host == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Host is required",
		})
		return
	}

	var mailFound Mail
	if err := db.First(&mailFound, "host = ?", req.Host).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Mail with the same host already exists",
		})
		return
	}

	mail := Mail{
		Host:    req.Host,
		Relay:   req.Relay,
		Version: 1,
	}

	if err := db.Create(&mail).Error; err != nil {
		logger.Errorf("failed to create mail: db create error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"mail":    mail,
	})
}

// apiControllersMailsRefresh creates a mail from the provided host.
func apiControllersMailsRefresh(c *gin.Context) {
	var req typeApiReqMailsRefresh
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf("failed to create mail: bind json error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	var mails []Mail
	for mailID, mailVersion := range req.MailVersions {
		var mail Mail
		err := db.
			Preload("MailInboxes").
			Preload("MailUpstreams").
			First(&mail, "id = ? and version > ?", mailID, mailVersion).Error

		if err != nil {
			err = fmt.Errorf("get mail refresh with id error: %d: %w", mailID, err)
			logger.Errorf("failed to refresh mail with version: %v", err)
			continue
		}

		mails = append(mails, mail)
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"mails":   mails,
	})
}

func apiControllersMailsGet(c *gin.Context) {
	mailHost := c.Param("mailHost")

	var mail Mail
	err := db.
		Preload("MailInboxes").
		Preload("MailUpstreams").
		First(&mail, "host = ?", mailHost).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"found":   false,
			})
			return
		}

		err = fmt.Errorf("get mail with host error: %s: %w", mailHost, err)
		logger.Errorf("failed to get mail: %v", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"found":   true,
		"mail":    mail,
	})
}
