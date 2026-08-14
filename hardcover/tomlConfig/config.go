package tomlConfig

import (
	"log"

	"github.com/pelletier/go-toml"
)

func GetConfig(filePath string) Config {
	config, err := toml.LoadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	err = config.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

type Config struct {
	Hardcover Hardcover
}

type Hardcover struct {
	Key                   string
	FollowedAuthorsListID int `toml:"list_id"`
	MinUsersCount 		  int `toml:"min_users_count"`
}
