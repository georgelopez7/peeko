package main

import (
	"log"

	"github.com/georgelopez7/peeko/api/http"
	"github.com/georgelopez7/peeko/internal/store"
)

func main() {
	config := NewConfig()

	st := store.NewStore(100)
	server := http.NewServer(":"+config.Port, st)

	log.Fatal(server.Start())
}
