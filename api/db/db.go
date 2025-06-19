package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var count int64

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) *sql.DB {
	db := connectToDB(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	if db == nil {
		log.Panicln("error connecting db")
	}
	return db
}

func connectToDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) *sql.DB {
	for {
		connection, err := openDB(addr, maxOpenConns, maxIdleConns, maxIdleTime)
		if err != nil {
			log.Println("Failed to connect DB")
			count++
		} else {
			log.Println("Connected to DB")
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Sleeping for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		log.Println("error opening database")
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
