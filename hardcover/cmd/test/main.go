package main

import (
	"fmt"

	"github.com/KelvinMcClean/palimpsest/hardcover/api"
	"github.com/KelvinMcClean/palimpsest/hardcover/tomlConfig"
	"github.com/KelvinMcClean/palimpsest/hardcover/structs"
)

func main() {
	cfg := tomlConfig.GetConfig("../../config.toml")
	fmt.Println(cfg.Hardcover.Key)
	apiClient, err := api.NewAPIClient("https://api.hardcover.app/v1/graphql", cfg)
	if err != nil {
		fmt.Println("Error creating API client:", err)
		return
	}
	authors, ok := apiClient.GetFollowedAuthors(cfg)
	if !ok {
		fmt.Println("Error getting followed authors")
		return
	}
	fmt.Println("Followed authors:", authors)

	books := structs.AuthorsToBooks(authors)
	fmt.Println("Books:\n", books)
}
