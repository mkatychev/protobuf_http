package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	port := os.Getenv("NTBK_PORT")
	if port == "" {
		// default port entry
		port = "8080"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	log.Println(fmt.Sprintf("Starting notebook server on %s", addr))
	r := mux.NewRouter()
	repo := NewNotebookRepo()

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	r.HandleFunc("/note", repo.CreateNote).Methods("POST")
	r.HandleFunc("/note", repo.DeleteNote).Methods("DELETE")
	r.HandleFunc("/note", repo.GetNote).Methods("GET")
	r.HandleFunc("/note", repo.UpdateNote).Methods("UPDATE")
	r.HandleFunc("/notebook", repo.CreateNotebook).Methods("POST")
	r.HandleFunc("/notebook", repo.GetNotebook).Methods("GET")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
