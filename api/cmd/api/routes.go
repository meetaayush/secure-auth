package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/meetaayush/secure-auth/store"
)

type application struct {
	config config
	store  store.Store
}

type config struct {
	Env             string
	Addr            string
	DbConfig        dbConfig
	RedisConfig     redisConfig
	AuthTokenConfig tokenConfig
	WebAddr         string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type redisConfig struct {
	addr     string
	password string
	db       int
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

const (
	userCtx    = "user_ctx"
	sessionCtx = "session_ctx"
)

func (app *application) NewRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.config.WebAddr},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/health", app.healthHandler)

			r.Route("/users", func(r chi.Router) {
				r.Post("/register", app.userRegistrationHandler)
				r.Post("/auth", app.userLoginHandler)
			})

			// private route
			r.Route("/auth", func(r chi.Router) {
				r.Use(app.authMiddleware)

				r.Get("/me", app.getLoggedinUser)
				r.Post("/logout", app.logoutHandler)
			})
		})
	})

	return mux
}

func (app *application) run(routes *chi.Mux) error {
	srv := http.Server{
		Addr:    app.config.Addr,
		Handler: routes,
	}

	log.Println("Server is running on port", app.config.Addr)

	return srv.ListenAndServe()
}
