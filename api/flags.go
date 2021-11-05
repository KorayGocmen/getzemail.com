package main

var (
	flagConfigPath string
	flagDBMigrate  bool
)

var (
	flagUsageConfigPath = "path to config file. Default: config.toml"
	flagUsageDBMigrate  = "migrate database when supplied"
)
