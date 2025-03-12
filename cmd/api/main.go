package main

import (
	"socialone/internal/db"
	"socialone/internal/env"
	"socialone/internal/mailer"
	"socialone/internal/store"
	"time"

	"go.uber.org/zap"
)
const version = "2.0.0"

// @title SocialOne API

// @description This is a  server for social One.
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



func main(){
		cfg:= config{
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
				sendGrid: sendGridConfig{
					apiKey: env.Getstring("SENDGRID_API_KEY", ""),

				},
			},
		}

//logger
		logger:= zap.Must(zap.NewProduction()).Sugar()
		defer logger.Sync()


		db, err:= db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxidealConns, cfg.db.maxidealTime)
		if err!=nil{
			logger.Fatal(err)
		}
		defer db.Close()
		logger.Info("Database connection successful")
		store:= store.NewPostgress(db)
		mailer:= mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
		app:= &application{
			config: cfg,
			store: *store,
			logger: logger,
			mailer: mailer,
		}

		Mux:= app.mount()
		logger.Fatal(app.run(Mux))


}