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
var RECAPTCHA_SITE string

// Email settings
var SMTP_HOST string
var SMTP_PORT string
var SMTP_PASSWWORD string
var EMAIL_FROM string

// Database settings
var DB_PASSWORD string
var DB_USER string
var DB_PORT string
var DB_HOST string
var DB_NAME string

func Init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load the env vars: %v", slog.Any("err", err))
	}
	godotenv.Load(".env")

	// Captcha
	RECAPTCHA_SECRET = os.Getenv("RECAPTCHA_SECRET")
	RECAPTCHA_VERIFY_URL = os.Getenv("RECAPTCHA_VERIFY_URL")
	RECAPTCHA_SITE = os.Getenv("RECAPTCHA_SITE")

	// Email
	SMTP_HOST = os.Getenv("SMTP_HOST")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	SMTP_PASSWWORD = os.Getenv("SMTP_PASSWWORD")
	EMAIL_FROM = os.Getenv("EMAIL_FROM")

	// Database
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_USER = os.Getenv("DB_USER")
	DB_PORT = os.Getenv("DB_PORT")
	DB_HOST = os.Getenv("DB_HOST")
	DB_NAME = os.Getenv("DB_NAME")
}
