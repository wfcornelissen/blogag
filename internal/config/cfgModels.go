package config

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	state *Config
}

type Command struct {
	name string
	args []string
}

const configFilePath = ".gatorConfig.json"
