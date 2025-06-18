package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
}

var (
	jsonConfigPath = "config/development.json"
	socketAddress  = ""
)

func main() {
	ctx := context.Background()

	flag.StringVar(&jsonConfigPath, "config", "config/development.json", "config json path")
	flag.StringVar(&socketAddress, "socket-address", "", "where to run the service")
	flag.Parse()

	// init config
	initConfig(ctx, jsonConfigPath)
	CONFIG.SetString("service.socket.address", socketAddress)

	address := CONFIG.GetString("service.socket.address")

	router := httprouter.New()

	enableBasicAPI(router)

	server := http.Server{
		Addr:    address,
		Handler: router,
		ReadTimeout: time.Duration(
			CONFIG.GetInt("service.socket.read-timeout-seconds")) * time.Second,
		WriteTimeout: time.Duration(
			CONFIG.GetInt("service.socket.write-timeout-seconds"),
		) * time.Second,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		// Start interrupt signal listener.

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Info().Msgf("Shutting down HTTP server..")
		if err := server.Shutdown(ctx); err != nil {
			// eg. timeout
			log.Err(err).Msgf("HTTP server Shutdown")
		}
		log.Info().Msgf("Stopped serving new connections.")
		close(idleConnsClosed)
	}()

	log.Info().Msgf("Serving at %v..\n", address)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatal().Msgf("HTTP server ListendAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Info().Msgf("Bye bye")
}

func Empty(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
