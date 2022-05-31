package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vijeyash1/gitevents/models"
)

const (
	streamName     = "GITMETRICS"
	streamSubjects = "GITMETRICS.*"
	eventSubject   = "GITMETRICS.event"
	allSubject     = "GITMETRICS.all"
	version        = "1.0.0"
)

//to read the token from env variables
var (
	gituser  = os.Getenv("GIT_USER")
	gittoken = os.Getenv("GIT_TOKEN")
	token    = os.Getenv("NATS_TOKEN")   //"UfmrJOYwYCCsgQvxvcfJ3BdI6c8WBbnD"
	natsurl  = os.Getenv("NATS_ADDRESS") //"nats://localhost:4222"
)

type config struct {
	port      int
	nats      string
	natstoken string
}

type application struct {
	config config
	models models.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8000, "Server port to listen on")
	flag.StringVar(&cfg.nats, "natsurl", natsurl, "nats connection url")
	flag.StringVar(&cfg.natstoken, token, "UfmrJOYwYCCsgQvxvcfJ3BdI6c8WBbnD", "nats token")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	js := openJS(cfg)

	app := &application{
		config: cfg,
		models: models.NewModels(js, gituser, gittoken),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Println("Starting server on port", cfg.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openJS(cfg config) nats.JetStreamContext {
	// Connect to NATS
	nc, err := nats.Connect(cfg.nats, nats.Name("Github metrics"), nats.Token(cfg.natstoken))
	if err != nil {
		log.Fatal(err)
	}
	// Creates JetStreamContext
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	// Creates stream
	err = createStream(js)
	if err != nil {
		log.Fatal(err)
	}
	return js

}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext) error {
	// Check if the METRICS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	log.Printf("Retrieved stream %s", fmt.Sprintf("%v", stream))
	if err != nil {
		log.Printf("Error getting stream %s", err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
