package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func initDB() {
	logger.Debugln("connecting to database")

	var (
		dbDriver = config.Database.Driver
		dbDSN    = config.Database.DSN
	)

	var dialector gorm.Dialector
	if dbDriver == dbDriverPostgres {
		dialector = postgres.Open(dbDSN)
	} else if dbDriver == dbDriverMysql {
		dialector = mysql.Open(dbDSN)
	} else if dbDriver == dbDriverSqlite {
		dialector = sqlite.Open(dbDSN)
	} else {
		err := fmt.Errorf("unsupported database driver: %s", dbDriver)
		logger.Fatalln("Error:", err)
	}

	var err error
	db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("failed to open database: %w", err)
		logger.Fatalln("Error:", err)
	}
}
