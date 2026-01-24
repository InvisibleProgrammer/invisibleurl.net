package main

import (
	"context"
	"encoding/gob"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	slogfiber "github.com/samber/slog-fiber"
	"gopkg.in/natefinch/lumberjack.v2"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
	repository "invisibleprogrammer.com/invisibleurl/db"
	"invisibleprogrammer.com/invisibleurl/environment"
	"invisibleprogrammer.com/invisibleurl/routing"
	"invisibleprogrammer.com/invisibleurl/security"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/users"
)

func main() {
	logDir := "log"

	logRotator := &lumberjack.Logger{
		Filename:   logDir + "/invisibleurl.log",
		MaxSize:    100,  // Max size in MB
		MaxBackups: 5,    // Number of backups
		MaxAge:     30,   // Days
		Compress:   true, // Enable compression
	}

	multiWriter := slog.NewJSONHandler(
		io.MultiWriter(os.Stdout, logRotator),
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		},
	)

	// logger := slog.New(tint.NewHandler(logFile, nil))
	log := slog.New(multiWriter)

	environment.Init()

	log.Info("Starting invisibleurl.net",
		slog.String("version", "0.0.0"),
		slog.String("address", environment.HOST),
		slog.Int("pid", os.Getpid()))

	// Initialize repositories
	repository, err := repository.NewRepository()
	if err != nil {
		log.Error("Failed to initialize db: %v", slog.Any("err", err))
	}

	userRepository := users.NewUserRepository(repository)
	urlShortenerRepository := urlshortener.NewUrlShortenerRepository(repository)
	auditLogRepository := auditlog.NewAuditLogRepository(repository)

	// initialize audit logger
	auditLogger := auditlog.NewAuditLogService(auditLogRepository)

	// Initialize server session
	store := session.New()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	// HTML templates
	engine := html.New("./views", ".html")

	// Initialize engine
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	app.Use(requestid.New())

	app.Use(slogfiber.New(log))

	// security
	security.RegisterRateLimitingMiddleware(app)

	// Show authenticated user name on header partial
	users.RegisterUsernameMiddleware(app, store)

	// Set up routing
	routing.RegisterRoutes(app, store, userRepository, urlShortenerRepository, auditLogger)

	go func() {
		if err := app.Listen("127.0.0.1:8080"); err != nil {
			log.Error("Error on start: %v", slog.Any("err", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Info("Shutdown error: %v", slog.Any("err", err))
	}
}
