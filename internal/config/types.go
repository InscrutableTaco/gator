package config

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	config *Config
}

type Command struct {
	name string
	args []string
}

type Commands struct {
	cmdMap map[string]func(*State, Command) error
}
