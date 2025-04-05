package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/y3g0r/modern-full-stack-blog-go/api"
	"github.com/y3g0r/modern-full-stack-blog-go/configs"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/logger"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/repo"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

func main() {
	config := configs.Load()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// TODO: get connection parameters from config
	db, err := sqlx.Connect("postgres", "user=admin password=CHANGEME dbname=blog sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	logger := logger.NewZapLogger()

	postsRepo := repo.NewPostgresPostsRepo(db, logger)
	// postsRepo := repo.NewInMemoryPostsRepo()
	postsService := service.NewPostsService(postsRepo)

	// Create an instance of our handler which satisfies the generated interface
	blogApi := api.NewBlog(postsService)

	blogStictHandler := api.NewStrictHandler(blogApi, nil)

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	api.HandlerFromMux(blogStictHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", config.PORT),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
