package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	awsSess *session.Session
	awsS3   *s3.S3
)

func initAWS() {
	logger.Println("Creating AWS session")

	awsConfig := aws.Config{
		Region: aws.String(config.S3.Region),
		Credentials: credentials.NewStaticCredentials(
			config.S3.AccessKeyID,
			config.S3.SecretAccessKey,
			"", // a token will be created when the session it's used.
		),
	}

	var err error
	awsSess, err = session.NewSession(&awsConfig)
	if err != nil {
		logger.Fatalln("Failed to create AWS session", err)
	}

	awsS3 = s3.New(awsSess)
}
