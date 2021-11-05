BINARY_NAME = smtp
CONFIG_FILE = config.dev.toml

.PHONY: all clean build build_linux run 

all: clean build run

clean: 
	rm -f $(BINARY_NAME)
	rm -f *.log

build:
	go build -o $(BINARY_NAME)

build_linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)

run:
	./$(BINARY_NAME) --config=$(CONFIG_FILE)
