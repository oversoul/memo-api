package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"memo/api"
	"memo/api/auth"
	"memo/api/notes/repository"
	"memo/pkg/database"
	"memo/pkg/logger"

	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	return os.Getenv(key)
}

func run(ctx context.Context, _ io.Reader, out, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := logger.New(out)

	db, err := database.New(getEnv)
	if err != nil {
		return err
	}

	defer database.Close(db)

	di := api.DI{
		Logger:    logger,
		NoteRepo:  repository.NewNotes(db),
		TodoRepo:  repository.NewTodoNotes(db),
		MovieRepo: repository.NewMovieNotes(db),
		AuthStore: auth.NewStore(auth.NewRepo(db)),
	}

	srv := api.New(di)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(getEnv("HOST"), getEnv("PORT")),
		Handler: srv,
	}

	go func() {
		log.Printf("Listening on %s\n", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		// make new context for the shutdown
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "Error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("err loading: %v", err)
	}

	ctx := context.Background()
	if err := run(ctx, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
