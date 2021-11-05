package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	apiServer *http.Server
)

func apiRouter() http.Handler {
	r := gin.New()

	// Middlewares.
	r.Use(gin.Recovery())
	r.Use(apiMiddlewareCors())

	r.POST("/mails", apiControllersMailsCreate)
	r.POST("/mails/refresh", apiControllersMailsRefresh)
	r.GET("/mails/:mailHost", apiControllersMailsGet)
	r.POST("/mails/:mailHost/inboxes", apiControllersMailInboxesCreate)
	r.GET("/mails/:mailHost/inboxes/:mailInboxAddr", apiControllersMailInboxes)
	r.GET("/mails/:mailHost/messages/:mailMessageID", apiControllersMailMessages)

	// Routes.
	smtp := r.Group("/")
	smtp.Use(apiMiddlewareAuthSmtp())
	{
		smtp.POST("/smtp/inbound", apiControllersSmtpInbound)
		smtp.POST("/smtp/outbound", apiControllersSmtpOutbound)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusCreated, map[string]interface{}{
			"success": false,
			"error":   "Not found",
		})
	})

	logger.Debugln("api router created")
	return r
}

// listenAPI starts the listener and is a blocking function.
// Has to be called async.
func listenAPI(sslCrt, sslKey string, srv *http.Server) {
	if config.API.TLS.Status {
		// Only if the API is using TLS, add the TLS config.
		srv.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				// TLS 1.2 cipher suites.
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,

				// TLS 1.3 cipher suites.
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		}

		logger.Printf("api with tls running on: %s", srv.Addr)
		if err := srv.ListenAndServeTLS(sslCrt, sslKey); err != nil {
			if err != http.ErrServerClosed {
				err = fmt.Errorf("failed to start api listener with tls: %w", err)
				logger.Fatalln("Error:", err)
			}
		}
	} else {
		logger.Printf("api w/o tls running on: %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				err = fmt.Errorf("failed to start api listener w/o tls: %w", err)
				logger.Fatalln("Error:", err)
			}
		}
	}
}

// api creates the api listener. This is the main function
// called by the main function.
func api() {
	if apiServer != nil {
		logger.Println("shutting down api server")

		if err := apiServer.Shutdown(context.Background()); err != nil {
			err = fmt.Errorf("failed to shutdown api server: %w", err)
			logger.Fatalln("Error:", err)
		}
	}

	apiServer = &http.Server{
		Addr:    config.API.Addr,
		Handler: apiRouter(),
	}

	go listenAPI(config.API.TLS.Crt, config.API.TLS.Key, apiServer)
}
