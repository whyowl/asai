package main

import (
	"asai/cmd/cli"
	"asai/cmd/http"
	"asai/cmd/telegram"
	"asai/internal/core"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	mode := flag.String("mode", "telegram", "Interface mode: cli | http | telegram")
	flag.Parse()

	ctx := context.Background()

	agent := core.NewAgent()
	switch *mode {
	case "cli":
		cli.Run(ctx, agent)
	case "http":
		http.Run(ctx, agent)
	case "telegram":
		telegram.Run(ctx, agent, os.Getenv("TELEGRAM_TOKEN"))
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
