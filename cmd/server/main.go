package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"user-analytics/internal/httpapi"
	"user-analytics/internal/repo"
	"user-analytics/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		log.Fatal("DATABASE_URL required")
	}

	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		log.Fatal(err)
	}
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := repo.NewUserPgRepo(db)
	s := service.NewUserService(r)
	h := httpapi.NewHandlers(s)
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           httpapi.NewRouter(h),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
