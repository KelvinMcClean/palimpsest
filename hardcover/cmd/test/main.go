package main

import (
	"fmt"

	"github.com/KelvinMcClean/palimpsest/hardcover/api"
	"github.com/KelvinMcClean/palimpsest/hardcover/tomlConfig"
)

func main() {
	cfg := tomlConfig.GetConfig("../../config.toml")
	fmt.Println(cfg.Hardcover.Key)
	apiClient, err := api.NewAPIClient("https://api.hardcover.app/v1/graphql", cfg)
	if err != nil {
		fmt.Println("Error creating API client:", err)
		return
	}
	apiClient.GetFollowedAuthors(cfg)
}
