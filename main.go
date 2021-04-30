package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go-proxy/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	opts := cased.CurrentPublisher().Options()
	if opts.PublishKey == "" {
		panic("CASED_PUBLISH_KEY is missing, please configure and start the proxy again: https://github.com/cased/cased-go#configuration")
	}

	if !strings.HasPrefix(opts.PublishKey, "publish_") {
		panic(`Please provide a Cased publish key that starts with "publish_" and start the proxy again: https://github.com/cased/cased-go#configuration`)
	}

	http.HandleFunc("/publish", handlers.AuditEvents)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
