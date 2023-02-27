package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GrpcPort string
	MongoURI string
}

func GetConfig() Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		GrpcPort: ":50051",
		MongoURI: os.Getenv("MONGO_URI"),
	}
}
