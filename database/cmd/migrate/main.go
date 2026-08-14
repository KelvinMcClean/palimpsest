package main

import (
	"context"
	"log"

	"github.com/KelvinMcClean/palimpsest/database"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	if err := db.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Database migration complete")
}
