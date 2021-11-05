package main

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	awsSess      *session.Session
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
)

type s3UploadOpts struct {
	Bucket      string
	Key         string
	ACL         string
	ContentType string
	MetaData    map[string]string
}

type s3DownloadOpts struct {
	Bucket string
	Key    string
}

func initAWS() {
	if config.Server.FirewallOnly {
		logger.Println("Server is firewall only, AWS session is not created")
		return
	}

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

	s3Uploader = s3manager.NewUploader(awsSess)
	s3Downloader = s3manager.NewDownloader(awsSess)
}

// s3Upload uploads the provided file to S3.
func s3Upload(opts s3UploadOpts, file io.Reader) (*s3manager.UploadOutput, error) {
	if s3Uploader == nil {
		err := fmt.Errorf("S3 upload is called when S3 uploader is nil, is server on firewall only mode?")
		logger.Errorln("Faied to upload file to S3", err)
		return nil, err
	}

	up, err := s3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(opts.Bucket),
		Key:         aws.String(opts.Key),
		ACL:         aws.String(opts.ACL),
		ContentType: aws.String(opts.ContentType),
		Metadata:    aws.StringMap(opts.MetaData),
		Body:        file,
	})

	return up, err
}

func s3Download(opts s3DownloadOpts) ([]byte, error) {
	if s3Uploader == nil {
		err := fmt.Errorf("S3 upload is called when S3 uploader is nil, is server on firewall only mode?")
		logger.Errorln("Faied to upload file to S3", err)
		return nil, err
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := s3Downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(opts.Bucket),
		Key:    aws.String(opts.Key),
	})

	return buf.Bytes(), err
}
