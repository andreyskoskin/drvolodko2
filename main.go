package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/andreyskoskin/drvolodko2/localenv"
)

func main() {
	var ldb, err = localenv.NewLocalDB(localenv.DefaultConfig().DB)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = ldb.Close()
	}()

	if err := ldb.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := ldb.Test(); err != nil {
		log.Fatalln(err)
	}
}
