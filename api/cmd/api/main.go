package main

import (
	"log"
	"time"

	"github.com/meetaayush/secure-auth/db"
	"github.com/meetaayush/secure-auth/store"
)

func main() {
	dbCfg := dbConfig{
		addr:         "postgres://admin:password@localhost/secure_auth?sslmode=disable",
		maxOpenConns: 20,
		maxIdleConns: 20,
		maxIdleTime:  "10m",
	}
	redisCfg := redisConfig{
		addr:     "localhost:6379",
		password: "",
		db:       0,
	}
	tokenCfg := tokenConfig{
		secret: "a-very-secret-token",
		exp:    time.Hour * 24,
		iss:    "secure-downtask-auth",
	}
	db := db.New(dbCfg.addr, dbCfg.maxOpenConns, dbCfg.maxIdleConns, dbCfg.maxIdleTime)
	defer db.Close()
	log.Println("DB pool connected")

	redisClient := store.NewRedisClient(redisCfg.addr, redisCfg.password, redisCfg.db)
	defer redisClient.Close()
	log.Println("Connected to redis database")

	sessionManager := store.NewSessionManager(redisClient, time.Hour*24)

	authenticator := store.NewJWTAuthenticator(tokenCfg.secret, tokenCfg.iss, tokenCfg.iss)

	app := application{
		config: config{
			Env:             "development",
			Addr:            ":3001",
			DbConfig:        dbCfg,
			RedisConfig:     redisCfg,
			AuthTokenConfig: tokenCfg,
		},
		store: *store.NewStore(db, sessionManager, authenticator),
	}

	routes := app.NewRoutes()

	log.Fatal(app.run(routes))
}
