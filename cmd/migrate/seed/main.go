package main

import (
	"log"
	"socialone/internal/db"
	"socialone/internal/env"
	"socialone/internal/store"
)

func main(){
	addr := env.Getstring("DB_ADDR", "host=localhost port=5433 user=admin password=adminpassword dbname=postgres sslmode=disable")
	conn,err := db.New(addr, 30,30, "15m")
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()
	store:= store.NewPostgress(conn)


	db.Seed(*store)
}