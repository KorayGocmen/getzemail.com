package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

const (
	version = "1.0.0"
)

func init() {
	versionAsked := pflag.BoolP("version", "v", false, "Print the version")
	pflag.StringVarP(&flagConfigPath, "config", "c", "config.toml", flagUsageConfigPath)
	pflag.BoolVarP(&flagDBMigrate, "db-migrate", "m", false, flagUsageDBMigrate)
	pflag.Parse()

	// If the version argument passed,
	// print the version and exit.
	if *versionAsked {
		fmt.Printf("v%s\n", version)
		os.Exit(0)
	}

	initConfig()

	initLogger()

	initCLI()

	initDB()

	initAWS()
}

func main() {
	api()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			logger.Fatalln("Signal received, stopping")
		}
	}
}
