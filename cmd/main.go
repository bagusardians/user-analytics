package main

import (
	"log"
	"net/http"
	"time"
	"user-analytics/internal/httpapi"
	"user-analytics/internal/repo"
	"user-analytics/internal/service"
)

func main() {
	r := repo.NewUserPgRepo()
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
