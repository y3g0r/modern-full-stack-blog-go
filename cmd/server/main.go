package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lmittmann/tint"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/y3g0r/modern-full-stack-blog-go/api"
	"github.com/y3g0r/modern-full-stack-blog-go/configs"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/repo"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

func main() {
	// logHandler := slog.NewTextHandler(os.Stdout, nil)
	logHandler := tint.NewHandler(os.Stdout, nil)
	logger := slog.New(logHandler)
	logLogger := slog.NewLogLogger(logHandler, slog.LevelError)
	slog.SetDefault(logger)

	config := configs.Load()

	clerk.SetKey(config.CLERK_SK)

	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Error("Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// TODO: get connection parameters from config
	// db, err := sqlx.Connect("postgres", "user=admin password=CHANGEME dbname=blog sslmode=disable")
	db, err := sqlx.Connect("postgres", config.DATABASE_URL)
	if err != nil {
		logger.Error(err.Error())
	}

	postsRepo := repo.NewPostgresPostsRepo(db, logger)
	// postsRepo := repo.NewInMemoryPostsRepo()
	postsService := service.NewPostsService(postsRepo)

	// jamsRepo := repo.NewInMemoryJams()
	jamsRepo := repo.NewPostgresJamsRepo(db, logger)
	jamsService := service.NewJams(jamsRepo)

	// Create an instance of our handler which satisfies the generated interface
	restApi := api.NewApi(postsService, jamsService)

	stictHandler := api.NewStrictHandler(restApi, nil)

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
	r.Use(clerkhttp.WithHeaderAuthorization())

	// We now register our petStore above as the handler for the interface
	api.HandlerFromMux(stictHandler, r)

	addr := net.JoinHostPort("0.0.0.0", config.PORT)

	s := &http.Server{
		Handler:  r,
		Addr:     addr,
		ErrorLog: logLogger,
	}

	logger.Info("Listening on " + addr)
	// And we serve HTTP until the world ends.
	err = s.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
