package main

import (
	"flag"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

var (
	address string
	static  string
	build   string
)

func init() {
	flag.StringVar(&address, "address", ":80", "web server's address")
	flag.StringVar(&static, "static", "./build/static", "path to static files")
	flag.StringVar(&build, "build", "./build", "path to build files")
}

func main() {
	router := gin.Default()
	router.StaticFS("/static", http.Dir(static))
	router.NoRoute(func(c *gin.Context) {
		c.File(path.Join(build, "index.html"))
	})

	log.Fatalln(router.Run(address))
}
