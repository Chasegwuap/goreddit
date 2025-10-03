package main

import (
	"log"
	"net/http"

	"github.com/Chasegwuap/goreddit/web"

	"github.com/Chasegwuap/goreddit/postgres"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	http.ListenAndServe(":3000", h)
}
