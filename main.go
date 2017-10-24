package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	flagServerPort string
	srv            *Server
	trn            *TournamentManager
)

func init() {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&flagServerPort, "port", "6024", "Port for the server to listen on")
}

func main() {
	flag.Parse()
	log.Info("Bastille is opening the rec. yard...")

	srv = &Server{}
	trn = &TournamentManager{}

	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitChan
		cleanup()
		os.Exit(1)
	}()

	go trn.init()
	srv.init()
}

func cleanup() {
	trn.cleanup()
	log.Info("Shutting down... Goodbye!")
}
