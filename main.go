package main

import (
	"log"
	"nxctl/cmd"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env variables: %s\n", err)
	}

	cmd.Execute()
}
