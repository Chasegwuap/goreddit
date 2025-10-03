package main

import (
	"log"
	"net/http"

	"github.com/Chasegwuap/goreddit/postgres"
	"github.com/Chasegwuap/goreddit/web"
)

func main() {
	// Connect to your Postgres store
	store, err := postgres.NewStore("postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	h := web.NewHandler(store)

	// Print a message to confirm the server started
	log.Println("Server started on http://localhost:3000")

	// Start the server and log fatal errors if ListenAndServe fails
	err = http.ListenAndServe(":3000", h)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
