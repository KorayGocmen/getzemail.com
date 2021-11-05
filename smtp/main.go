package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

const (
	version = "1.3.1"
)

func init() {
	versionAsked := pflag.BoolP("version", "v", false, "Print the version")
	pflag.StringVarP(&configPath, "config", "c", "config.toml", "Path to config file")
	pflag.Parse()

	// If the version argument passed,
	// print the version and exit.
	if *versionAsked {
		fmt.Printf("v%s\n", version)
		os.Exit(0)
	}

	initConfig()

	initLogger()

	initRedis()

	initAWS()

	initMails()

	initMessages()
}

func main() {
	smtpRelay()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			logger.Fatalln("Signal received, stopping")
		}
	}
}
