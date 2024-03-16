package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func privateKey() string {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("PRIVATE_KEY")

}

func sepolia_rpc() string {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("SEPOLIA_RPC")

}
