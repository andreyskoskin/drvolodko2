package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/andreyskoskin/drvolodko2/webapi"
)

func main() {
	var err = webapi.Start(webapi.Config{
		HTTP: webapi.HTTPConfig{
			Address: ":8080",
		},
		DB: webapi.DBConfig{
			Name:     "postgres",
			Host:     "dbserver",
			Port:     5432,
			User:     "user",
			Password: "password",
		},
	})

	if err != nil {
		log.Fatalln(err)
	}
}
