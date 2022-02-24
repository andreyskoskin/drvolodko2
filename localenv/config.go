package localenv

import (
	"github.com/andreyskoskin/drvolodko2/webapi"
)

var defaultConfig = webapi.Config{
	HTTP: webapi.HTTPConfig{
		Address: ":8080",
	},
	DB: webapi.DBConfig{
		Name:     "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
	},
}

func DefaultConfig() webapi.Config {
	return defaultConfig
}
