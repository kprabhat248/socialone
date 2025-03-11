package main

import (
	"log"
	"socialone/internal/db"
	"socialone/internal/env"
	"socialone/internal/store"
)
const version = "1.0.0"

// @title SocialOne API

// @description This is a sample server Petstore server.
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
			db: dbConfig{
				addr: env.Getstring("DB_ADDR", "host=localhost port=5433 user=admin password=adminpassword dbname=postgres sslmode=disable"),
				maxOpenConns: env.Getint("DB_MAX_OPEN_CONNS", 25),
				maxidealConns: env.Getint("DB_MAX_IDEAL_CONNS", 25),
				maxidealTime: env.Getstring("DB_MAX_IDEAL_TIME", "15m"),
			},
			env: env.Getstring("ENV", "development"),
		}




		db, err:= db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxidealConns, cfg.db.maxidealTime)
		if err!=nil{
			log.Panic(err)
		}
		defer db.Close()
		log.Println("Database connection successful")
		store:= store.NewPostgress(db)
		app:= &application{
			config: cfg,
			store: *store,
		}

		Mux:= app.mount()
		log.Fatal(app.run(Mux))


}