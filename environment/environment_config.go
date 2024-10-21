package environment

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// Recaptcha
var RECAPTCHA_SECRET string
var RECAPTCHA_VERIFY_URL string

func Init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load the env vars: %v", slog.Any("err", err))
	}
	godotenv.Load(".env")

	RECAPTCHA_SECRET = os.Getenv("RECAPTCHA_SECRET")
	RECAPTCHA_VERIFY_URL = os.Getenv("RECAPTCHA_VERIFY_URL")
}
