package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConns, maxIdealConns int, maxIdealTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdealConns)

	duration,err:= time.ParseDuration(maxIdealTime)
	if err!=nil{
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)
	ctx,cancel:= context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err:= db.PingContext(ctx); err!=nil{
		return nil, err
	}

	return db, nil
}
