package main

import (
	"github.com/andreyskoskin/drvolodko2/api"
	"github.com/andreyskoskin/drvolodko2/datasource"
)

type Config struct {
	Echo     api.EchoConfig            `toml:"echo"`
	Postgres datasource.PostgresConfig `toml:"postgres"`
}
