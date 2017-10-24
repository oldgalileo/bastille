package main

import (
	"net/http"

	"bytes"
	"io"
	"strings"

	"io/ioutil"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var (
	srvLog = log.WithFields(log.Fields{
		"prefix": "api",
	})
)

type Server struct{}

func (s *Server) init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", HandlerUpload)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "https://thebastille.co"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	srvLog.Fatal(http.ListenAndServe(":22101", handler))
}

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	var localBuff bytes.Buffer

	file, header, err := r.FormFile("strategy")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Body.Close()
		srvLog.WithError(err).Error("Could not handle file upload")
		return
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	srvLog.WithFields(log.Fields{
		"filename": header.Filename,
	}).Info("New strategy uploaded")
	io.Copy(&localBuff, file)

	writeErr := ioutil.WriteFile(STRATEGIES_DIR+name[0]+".ipd", localBuff.Bytes(), 0755)
	if writeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Body.Close()
		srvLog.WithError(err).Error("Issue opening file for writing")
		return
	}

	localBuff.Reset()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{'message': 'Upload successful'}"))
	r.Body.Close()
	strategy := &Strategy{
		ID:      getStrategyID(),
		Name:    name[0] + ".ipd",
		Author:  "tbd",
		Path:    STRATEGIES_DIR + name[0],
		Matches: []MatchID{},
	}
	trn.add(strategy)
}
