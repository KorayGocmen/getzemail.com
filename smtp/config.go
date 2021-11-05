package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Server struct {
		Host              string `toml:"host"`
		Port              int    `toml:"port"`
		Domain            string `toml:"domain"`
		TimeoutRead       int    `toml:"timeout_read"`
		TimeoutWrite      int    `toml:"timeout_write"`
		MaxMessageBytes   int    `toml:"max_message_bytes"`
		MaxRecipients     int    `toml:"max_recipients"`
		AllowInsecureAuth bool   `toml:"allow_insecure_auth"`
		FirewallOnly      bool   `toml:"firewall_only"`
	} `toml:"server"`

	Mails struct {
		RefreshEvery int `toml:"refresh_every"`
		TTL          int `toml:"ttl"`
	} `toml:"mails"`

	Messages struct {
		OutboundEvery int `toml:"outbound_every"`
	} `toml:"messages"`

	API struct {
		BaseURL string `toml:"base_url"`
		Secret  string `toml:"secret"`
		Timeout int    `toml:"timeout"`
	} `toml:"api"`

	S3 struct {
		Region          string `toml:"region"`
		AccessKeyID     string `toml:"access_key_id"`
		SecretAccessKey string `toml:"secret_access_key"`
	} `toml:"s3"`

	S3Emails struct {
		Bucket string `toml:"bucket"`
		ACL    string `toml:"acl"`
	} `toml:"s3_emails"`

	Redis struct {
		Addr         string `toml:"addr"`
		Pass         string `toml:"pass"`
		DB           int    `toml:"db"`
		PoolSize     int    `toml:"pool_size"`
		MinIdleConns int    `toml:"min_idle_conns"`
		IdleTimeout  int    `toml:"idle_timeout"`
	} `toml:"redis"`

	Logger struct {
		Level      int    `toml:"level"`
		Mode       string `toml:"mode"`
		Path       string `toml:"path"`
		CheckEvery int    `toml:"check_every"`
	} `toml:"logger"`
}

var (
	config     tomlConfig
	configPath string

	pathExecutable    string
	pathExecutableDir string
)

func initConfig() {
	var err error
	if pathExecutable, err = os.Executable(); err != nil {
		log.Fatalln("Failed to get path of executable", err)
	}
	pathExecutableDir = path.Dir(pathExecutable)

	configPath = filepath.Join(pathExecutableDir, configPath)
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatalln("Reading config failed", err)
	}
}
