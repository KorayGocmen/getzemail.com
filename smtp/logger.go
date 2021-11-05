package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	loggerModeFile    = "file"
	loggerModeConsole = "console"

	loggerLevelDebug   = 3
	loggerLevelNormal  = 2
	loggerLevelMinimal = 1
)

var (
	logger      *loggerType
	loggerMutex = &sync.Mutex{}
)

type loggerType struct {
	out       *os.File
	timestamp string
}

func initLogger() {
	// Initialize the logger right away.
	loggerCreate()

	// Check if it's time to rotate every CheckEvery seconds only
	// when log mode is file.
	if config.Logger.Mode == loggerModeFile {

		// If the logger mode is file, write stdout and stderr to file.
		var err error
		os.Stdout, err = os.Create(filepath.Join(config.Logger.Path, "out.log"))
		if err != nil {
			log.Fatalln("Failed to create std out file")
		}
		os.Stderr, err = os.Create(filepath.Join(config.Logger.Path, "err.log"))
		if err != nil {
			log.Fatalln("Failed to create std err file")
		}

		go func() {
			every := timeDuration(config.Logger.CheckEvery)
			ticker := time.NewTicker(every)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// Timestamp is used by logger only when the logger mode
					// is "file". Log file name changes and rotates daily.
					timestamp := time.Now().Format("2006-01-02")

					// Initialize file if the logger is nil or if the timestamp
					// which is YYYY-MM-DD has changed. That means the day has changed,
					// and it's time to rotate logs.
					if logger == nil || logger.timestamp != timestamp {
						loggerCreate()
					}
				}
			}
		}()
	}
}

func loggerCreate() {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// Cleanup logger if the logger is not nil.
	// Close the out socket.
	if logger != nil && logger.out != nil {
		logger.out.Close()
	}

	// Timestamp is used by logger only when the logger mode
	// is "file". Log file name changes and rotates daily.
	timestamp := time.Now().Format("2006-01-02")

	// Initialize out as a file writer if the log mode
	// is set as "file".
	if config.Logger.Mode == loggerModeFile {
		pathEnsure(config.Logger.Path)

		loggerFileName := fmt.Sprintf("%s.log", timestamp)
		loggerFilePath := config.Logger.Path + "/" + loggerFileName

		fileFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
		out, err := os.OpenFile(loggerFilePath, fileFlags, os.ModePerm)
		if err != nil {
			log.Fatalln("Failed to open logger file", err)
		}

		logger = &loggerType{
			out:       out,
			timestamp: timestamp,
		}
	}

	// Initialize the os.StdOut as the output destination
	// if the logger mode is set as "console".
	if config.Logger.Mode == loggerModeConsole {
		logger = &loggerType{
			out:       os.Stdout,
			timestamp: timestamp,
		}
	}
}

func (l loggerType) prefix(level string) string {
	return fmt.Sprintf("%s [%s] ", timestamp(), level)
}

// Write prints into the logs when level is >= 3. Used by
// the http logger.
func (l loggerType) Write(p []byte) (n int, err error) {
	if config.Logger.Level >= loggerLevelDebug {
		line := l.prefix("WRITE") + string(p)
		return l.out.Write([]byte(line))
	}
	return 0, nil
}

// Debugln prints into the logs when level is >= 3
func (l loggerType) Debugln(args ...interface{}) error {
	if config.Logger.Level >= loggerLevelDebug {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		line := l.prefix("DEBUG") + fmt.Sprintln(args...)
		_, err := l.out.Write([]byte(line))
		return err
	}

	return nil
}

// Debugf prints into the logs when level is >= 3
func (l loggerType) Debugf(format string, args ...interface{}) error {
	if config.Logger.Level >= loggerLevelDebug {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		line := l.prefix("DEBUG") + fmt.Sprintf(format, args...)
		if !strings.HasSuffix(line, "\n") {
			line += "\n"
		}
		_, err := l.out.Write([]byte(line))
		return err
	}

	return nil
}

// Println prints into the logs when level is >= 2
func (l loggerType) Println(args ...interface{}) error {
	if config.Logger.Level >= loggerLevelNormal {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		line := l.prefix("INFOR") + fmt.Sprintln(args...)
		_, err := l.out.Write([]byte(line))
		return err
	}

	return nil
}

// Printf prints into the logs when level is >= 2
func (l loggerType) Printf(format string, args ...interface{}) error {
	if config.Logger.Level >= loggerLevelNormal {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		line := l.prefix("INFOR") + fmt.Sprintf(format, args...)
		if !strings.HasSuffix(line, "\n") {
			line += "\n"
		}
		_, err := l.out.Write([]byte(line))
		return err
	}

	return nil
}

// Errorln prints into the logs when level is >= 1
func (l loggerType) Errorln(args ...interface{}) error {
	if config.Logger.Level >= loggerLevelMinimal {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		line := l.prefix("ERROR") + fmt.Sprintln(args...)
		_, err := l.out.Write([]byte(line))
		return err
	}

	return nil
}

// Errorf prints into the logs when level is >= 1
func (l loggerType) Errorf(format string, args ...interface{}) error {
	if config.Logger.Level >= loggerLevelMinimal {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		line := l.prefix("ERROR") + fmt.Sprintf(format, args...)
		if !strings.HasSuffix(line, "\n") {
			line += "\n"
		}
		_, err := l.out.Write([]byte(line))
		return err
	}

	return nil
}

// Fatalln prints into the logs regardless of level and exits.
func (l loggerType) Fatalln(args ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	line := l.prefix("FATAL") + fmt.Sprintln(args...)
	l.out.Write([]byte(line))
	log.Fatalln(args...)
}

// Fatalf prints into the logs regardless of level and exits.
func (l loggerType) Fatalf(format string, args ...interface{}) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	line := l.prefix("FATAL") + fmt.Sprintf(format, args...)
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	l.out.Write([]byte(line))
	log.Fatalln(args...)
}
