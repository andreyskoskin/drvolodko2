package main

import (
	"fmt"
	"os"

	"github.com/andreyskoskin/drvolodko2/database"
)

func main() {
	if err := database.Demo(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("Done")
}
