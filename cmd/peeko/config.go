package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type Env struct {
	Port                   string
	WebhookSecret          string
	WebhookSignatureHeader string
	WebhookTimestampHeader string
	WebhookEncoding        string
	WebhookSignaturePrefix string
	WebhookTolerance       int // milliseconds
}

type Config struct {
	Port string
	Env  Env
}

func NewConfig() Config {
	env := Env{}

	flag.StringVar(&env.Port, "port", "", "port to listen on")
	flag.StringVar(&env.WebhookSecret, "webhook-secret", "", "webhook HMAC secret (empty disables verification)")
	flag.StringVar(&env.WebhookSignatureHeader, "webhook-signature-header", "", "header carrying the signature")
	flag.StringVar(&env.WebhookTimestampHeader, "webhook-timestamp-header", "", "header carrying the timestamp (empty = body-only)")
	flag.StringVar(&env.WebhookEncoding, "webhook-signature-encoding", "", "signature encoding: hex or base64")
	flag.StringVar(&env.WebhookSignaturePrefix, "webhook-signature-prefix", "", "signature prefix, e.g. \"sha256=\" or \"v1,\"")
	flag.IntVar(&env.WebhookTolerance, "webhook-tolerance", 0, "timestamp tolerance in milliseconds (e.g. 300000 = 5m)")
	flag.Parse()

	if env.Port == "" {
		env.Port = os.Getenv("PORT")
	}

	if env.WebhookSecret == "" {
		env.WebhookSecret = os.Getenv("WEBHOOK_SECRET")
	}

	if env.WebhookSignatureHeader == "" {
		env.WebhookSignatureHeader = envWithDefault("WEBHOOK_SIGNATURE_HEADER", "X-Webhook-Signature")
	}

	if env.WebhookTimestampHeader == "" {
		env.WebhookTimestampHeader = envWithDefault("WEBHOOK_TIMESTAMP_HEADER", "X-Webhook-Timestamp")
	}

	if env.WebhookEncoding == "" {
		env.WebhookEncoding = os.Getenv("WEBHOOK_SIGNATURE_ENCODING")
	}

	if env.WebhookSignaturePrefix == "" {
		env.WebhookSignaturePrefix = os.Getenv("WEBHOOK_SIGNATURE_PREFIX")
	}

	if env.WebhookTolerance == 0 {
		if v, ok := os.LookupEnv("WEBHOOK_TOLERANCE"); ok {
			parsed, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("invalid WEBHOOK_TOLERANCE %q - must be an integer (milliseconds)", v)
			}

			env.WebhookTolerance = parsed
		}
	}

	return Config{Port: env.Port, Env: env}
}

// envWithDefault - reads env, falling back to a default when unset or empty.
func envWithDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
