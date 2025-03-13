package main

import (
	"socialone/internal/db"
	"socialone/internal/env"
	"socialone/internal/mailer"
	"socialone/internal/store"
	"time"
	"log"
	"github.com/joho/godotenv" // Import the godotenv package
	"go.uber.org/zap"
)

const version = "2.0.0"

// @title SocialOne API

// @description This is a server for social One.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1

//@securityDefinitions.apikey ApiKeyAuth
//@in			header
//@name			Authorization
//@description

// Load environment variables from the .env file
func loadEnv() {
	// Load the environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Load environment variables from .env file
	loadEnv()

	cfg := config{
		addr: env.Getstring("ADDR", ":8080"),
		apiUrl: env.Getstring("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.Getstring("FRONTEND_URL", "localhost:4000"),
		db: dbConfig{
			addr: env.Getstring("DB_ADDR", "host=localhost port=5433 user=admin password=adminpassword dbname=postgres sslmode=disable"),
			maxOpenConns: env.Getint("DB_MAX_OPEN_CONNS", 25),
			maxidealConns: env.Getint("DB_MAX_IDEAL_CONNS", 25),
			maxidealTime: env.Getstring("DB_MAX_IDEAL_TIME", "15m"),
		},
		env: env.Getstring("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour*24*3,
			fromEmail: env.Getstring("MAIL_FROM_EMAIL", ""),
			smtp2go: smtp2goConfig{
				apiKey: env.Getstring("SMTP2GO_API_KEY", ""),
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	logger.Infow("MAIL_FROM_EMAIL: ", "value", cfg.mail.fromEmail)

	// Initialize the database connection
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxidealConns, cfg.db.maxidealTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("Database connection successful")

	// Initialize the store and mailer
	store := store.NewPostgress(db)
	mailer := mailer.NewSMTP2GoMailer(cfg.mail.smtp2go.apiKey, cfg.mail.fromEmail)

	// Initialize the application with config, store, logger, and mailer
	app := &application{
		config: cfg,
		store:  *store,
		logger: logger,
		mailer: mailer,
	}

	// Set up the HTTP mux and run the application
	Mux := app.mount()
	logger.Fatal(app.run(Mux))
}
