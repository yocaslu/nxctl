package main

import (
	"fmt"
	"nxctl/cmd"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load .env variables: %s\n", err)
	}

	cmd.Execute()
}
