package main

import (
	"errors"
	"net/http"
	"time"

	"flag"
	"log"

	"github.com/EwRvp7LV7/48170360shop/internal/storage/postgres"

	"github.com/EwRvp7LV7/48170360shop/api"
	"github.com/EwRvp7LV7/48170360shop/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func run() (err error) {
	log.Println("Loading config at", config.FileName)
	err = config.Load(config.FileName)

	if err != nil {
		err = errors.New("error load config - Rename and setup config.toml.example in configs/")
		return
	}

	postgres.OpenConnectDB()
	defer postgres.CloseConnectionDB()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	api.AddRouteAuthentication(r)
	api.AddRouteInputUserBasket(r)
	api.AddRouteInputManager(r)


	addr := config.GetServerAddress()
	log.Println("Server listening at", addr)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(addr, r)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
