package config

import "github.com/wfcornelissen/blogag/internal/database"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	err := Write(Config{
		DbUrl:           c.DbUrl,
		CurrentUserName: user,
	})
	if err != nil {
		return err
	}

	return nil
}

type State struct {
	Db    *database.Queries
	State *Config
}

const configFilePath = ".gatorconfig.json"
