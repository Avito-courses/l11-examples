package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"

	"github.com/Avito-courses/l11-examples/internal/handler/common"
	handler "github.com/Avito-courses/l11-examples/internal/handler/user"
	repository "github.com/Avito-courses/l11-examples/internal/repository/user"
	"github.com/Avito-courses/l11-examples/pkg/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Note: Could not load .env file: %v", err)
		log.Println("Continuing with environment variables...")
	}

	dbPool := db.MustInitDB()
	startServer(
		dbPool,
		handler.NewUserController(
			repository.NewUserRepository(dbPool),
		),
		resolvePort(),
	)
}

func startServer(
	dbPool *pgxpool.Pool,
	handler *handler.Controller,
	port string,
) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: initRouter(handler),
	}

	log.Printf("Server started on %s\n", srv.Addr)
	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	waitGracefulShutdown(srv, dbPool, serverErr)

	log.Println("Shutting down service")
}

func resolvePort() string {
	port := os.Getenv("PORT")

	var portFlag = flag.String("port", "", "укажите порт")
	flag.Parse()

	if portFlag != nil && *portFlag != "" {
		port = *portFlag
	}

	if port == "" {
		port = "8080"
	}

	return port
}

func waitGracefulShutdown(srv *http.Server, dbPool *pgxpool.Pool, serverErr <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var reason string
	select {
	case <-ctx.Done():
		reason = "signal"
	case err := <-serverErr:
		reason = "server error: " + err.Error()
	}
	log.Printf("Shutdown initiated (%s)", reason)

	log.Println("Shutting down HTTP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server stopped")
	}

	log.Println("Closing DB pool...")
	dbPool.Close()
	log.Println("DB pool closed")
}

func initRouter(controller *handler.Controller) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", common.Ping)
	r.Head("/healthcheck", common.HealthCheck)

	r.Route("/user", func(r chi.Router) {
		r.Get("/{id}", controller.Get)
	})

	return r
}
