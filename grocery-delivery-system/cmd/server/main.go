package main

import (
	"grocery-delivery/internal/application/commands"
	"grocery-delivery/internal/infrastructure/persistence/eventstore"
	"log"
	"net/http"
)

func main() {
	eventStore := eventstore.NewPostgresEventStore()
	cartRepo := cart.NewEventSourcedRepository(eventStore)
	addToCartHandler := commands.NewAddToCartHandler(cartRepo)

	router := rest.NewRouter(addToCartHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
