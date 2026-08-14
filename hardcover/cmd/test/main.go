package main

import (
	"context"
	"log"

	"github.com/KelvinMcClean/palimpsest/database"
	"github.com/KelvinMcClean/palimpsest/hardcover/api"
	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
	"github.com/KelvinMcClean/palimpsest/hardcover/tomlConfig"
)

func main() {
	ctx := context.Background()
	cfg := tomlConfig.GetConfig("../../config.toml")

	db, err := database.Connect(ctx)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	defer db.Close(ctx)
	log.Println("Connected to database")
	databaseName, userName, err := db.Info(ctx)
	if err != nil {
		log.Println("Error getting database info:", err)
		return
	}
	log.Println("Connected to database:", databaseName, "as user:", userName)

	apiClient, err := api.NewAPIClient("https://api.hardcover.app/v1/graphql", cfg)
	if err != nil {
		log.Fatalf("Error creating API client: %v", err)
		return
	}
	authors, ok := apiClient.GetFollowedAuthors(cfg)
	if !ok {
		log.Println("Error getting followed authors")
		return
	}
	log.Println("Followed authors:", authors)

	books := structs.AuthorsToBooks(authors)
	log.Println("Books:\n", books)

	if err := db.SaveAuthors(ctx, authors); err != nil {
		log.Fatal(err)
	}

	log.Println("Saved Hardcover data to postgres")
}
