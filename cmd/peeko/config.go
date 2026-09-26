package main

import (
	"flag"
	"os"
)

type Config struct {
	Port string
}

func NewConfig() Config {
	var port string

	flag.StringVar(&port, "port", "", "port to listen on")
	flag.Parse()

	if port == "" {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "8080"
	}

	return Config{port}
}
