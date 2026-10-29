package main

import (
	"fmt"
	"nxctl/cmd"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load .env variables: %s\n", err)
		os.Exit(1)
	}

	cmd.Execute()
}
