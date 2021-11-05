package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// apiControllersMailInboxes returns a mail inbox and all mail
// messages in that mail inbox with the provided host and address.
func apiControllersMailMessages(c *gin.Context) {
	mailHost := c.Param("mailHost")
	mailMessageID := c.Param("mailMessageID")

	var mail Mail
	if err := db.First(&mail, "host = ?", mailHost).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Mail not found",
		})
		return
	}

	var mailMessage MailMessage
	err := db.
		Preload("MailMessageFiles").
		Preload("MailMessageRelations").
		First(&mailMessage, "id = ?", mailMessageID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "Mail inbox not found",
			})
			return
		}

		logger.Errorf("failed to get mail message: %s: %w", mailMessageID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	textreq, _ := awsS3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.S3Emails.Bucket),
		Key:    aws.String(mailMessage.MessageID + "/text"),
	})

	mailMessage.TextURL, err = textreq.Presign(15 * time.Minute)
	if err != nil {
		logger.Errorf("failed to sign text url for mail message: %s: %w", mailMessageID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	htmlreq, _ := awsS3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.S3Emails.Bucket),
		Key:    aws.String(mailMessage.MessageID + "/html"),
	})

	mailMessage.HtmlURL, err = htmlreq.Presign(15 * time.Minute)
	if err != nil {
		logger.Errorf("failed to sign html url for mail message: %s: %w", mailMessageID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Something went wrong",
		})
		return
	}

	for _, mailMessageFile := range mailMessage.MailMessageFiles {
		filereq, _ := awsS3.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(config.S3Emails.Bucket),
			Key:    aws.String(mailMessage.MessageID + "/" + mailMessageFile.Key),
		})

		mailMessageFile.URL, err = filereq.Presign(15 * time.Minute)
		if err != nil {
			logger.Errorf("failed to sign file url for mail message: %s: %w", mailMessageID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"error":   "Something went wrong",
			})
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success":      true,
		"mail_message": mailMessage,
	})
}
