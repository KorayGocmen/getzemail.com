package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	address string
	static  string
)

func init() {
	flag.StringVar(&address, "address", ":80", "web server's address")
	flag.StringVar(&static, "static", "./build", "path to static files built")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(static)))
	log.Fatalln(http.ListenAndServe(address, nil))
}
