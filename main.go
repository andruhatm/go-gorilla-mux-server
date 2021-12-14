package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	//init server
	r := mux.NewRouter()
	r.Handle("/livenessProbe", healthChecker())
	//var err error
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		log.Printf("starting server on %v", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("error starting server on %v", server.Addr)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}

func healthChecker() http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("OK"))
		},
	)
}
