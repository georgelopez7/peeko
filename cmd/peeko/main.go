package main

import (
	"log"
	"strings"
	"time"

	"github.com/georgelopez7/peeko/api/http"
	"github.com/georgelopez7/peeko/internal/service"
	"github.com/georgelopez7/peeko/internal/store"
	"github.com/georgelopez7/peeko/internal/webhook"
)

func main() {
	// CONFIG
	config := NewConfig()

	// STORE
	st := store.NewStore(100)

	// WEBHOOK
	wh := webhook.NewVerifier(webhook.Config{
		Secret:          []byte(config.Env.WebhookSecret),
		SignatureHeader: config.Env.WebhookSignatureHeader,
		TimestampHeader: config.Env.WebhookTimestampHeader,
		Encoding:        webhook.Encoding(strings.ToLower(config.Env.WebhookEncoding)),
		Prefix:          config.Env.WebhookSignaturePrefix,
		Tolerance:       time.Duration(config.Env.WebhookTolerance) * time.Millisecond,
	})

	// SERVICE
	svc := service.New(st, wh)

	// SERVER
	server := http.NewServer(":"+config.Port, svc)
	defer server.Stop()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
