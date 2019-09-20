package main

import (
	"context"
	"log"

	language "cloud.google.com/go/language/apiv1"
	"github.com/windy-tech/recommendation/pkg/apiserver"
	"github.com/windy-tech/recommendation/pkg/database"
)

func main() {
	log.Println("Prepare Database")
	ds, err := database.PrepareDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Prepare AI")
	ctx := context.Background()
	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create ai client: %v", err)
	}
	log.Println("Run apiserver")
	server := &apiserver.APIServer{DB: ds, LangClient: client}
	server.Run()
}
