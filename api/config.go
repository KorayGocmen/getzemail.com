package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	API struct {
		Addr   string `toml:"addr"`
		Secret string `toml:"secret"`
		TLS    struct {
			Status bool   `toml:"status"`
			Crt    string `toml:"crt"`
			Key    string `toml:"key"`
		} `toml:"tls"`
	} `toml:"api"`

	Database struct {
		Driver string `toml:"driver"`
		DSN    string `toml:"dsn"`
	} `toml:"database"`

	S3 struct {
		Region          string `toml:"region"`
		AccessKeyID     string `toml:"access_key_id"`
		SecretAccessKey string `toml:"secret_access_key"`
	} `toml:"s3"`

	S3Emails struct {
		Bucket string `toml:"bucket"`
		ACL    string `toml:"acl"`
	} `toml:"s3_emails"`

	Logger struct {
		Level      int    `toml:"level"`
		Mode       string `toml:"mode"`
		Path       string `toml:"path"`
		CheckEvery int    `toml:"check_every"`
	} `toml:"logger"`
}

var (
	config tomlConfig

	pathExecutable    string
	pathExecutableDir string
)

func initConfig() {
	var err error
	if pathExecutable, err = os.Executable(); err != nil {
		log.Fatalln("Failed to get path of executable", err)
	}
	pathExecutableDir = path.Dir(pathExecutable)

	flagConfigPath = filepath.Join(pathExecutableDir, flagConfigPath)
	if _, err := toml.DecodeFile(flagConfigPath, &config); err != nil {
		log.Fatalln("Reading config failed", err)
	}
}
