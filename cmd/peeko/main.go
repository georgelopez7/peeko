package main

import (
	"flag"
	"log"
	"os"

	"github.com/george-lopez/peeko/api/http"
	"github.com/george-lopez/peeko/internal/store"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "", "port to listen on")
	flag.Parse()

	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	st := store.NewStore(100)
	server := http.NewServer(":"+port, st)

	log.Fatal(server.Start())
}
