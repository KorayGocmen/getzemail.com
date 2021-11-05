BINARY_NAME = web

.PHONY: all clean build build_linux 

all: clean build run

clean: 
	rm -f $(BINARY_NAME)

build:
	go build -o $(BINARY_NAME)

build_linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)
