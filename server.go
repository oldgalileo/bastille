package main

import (
	"net/http"

	"bytes"
	"io"
	"strings"

	"io/ioutil"

	"encoding/json"

	"github.com/kennygrant/sanitize"
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
	mux.HandleFunc("/upload", HandleUpload)
	mux.HandleFunc("/health", HandleHealth)
	mux.HandleFunc("/readiness", HandleReadiness)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "https://thebastille.co"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	srvLog.Fatal(http.ListenAndServe(":22101", handler))
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Limit filesize to 2MB
	r.Body = http.MaxBytesReader(w, r.Body, 2097152)
	var localBuff bytes.Buffer

	author := r.FormValue("author")
	name := r.FormValue("name")
	desc := r.FormValue("desc")
	file, header, err := r.FormFile("exec")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Body.Close()
		srvLog.WithError(err).Error("Could not handle file upload")
		return
	}
	defer file.Close()
	fileName := sanitize.BaseName(header.Filename)
	fileName = strings.Split(header.Filename, ".")[0] + ".bstl"
	srvLog.WithFields(log.Fields{
		"filename": header.Filename,
	}).Info("New strategy uploaded")
	io.Copy(&localBuff, file)

	writeErr := ioutil.WriteFile(STRATEGIES_DIR+fileName, localBuff.Bytes(), 0755)
	if writeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Body.Close()
		srvLog.WithError(err).Error("Issue opening file for writing")
		return
	}
	localBuff.Reset()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	type success struct {
		Message string `json:"message"`
	}
	response := &success{"Upload successful!"}
	responseRaw, responseErr := json.Marshal(response)
	if responseErr == nil {
		w.Write(responseRaw)
	} else {
		srvLog.WithError(responseErr).Error("Could not marshal strategy submission response")
	}
	r.Body.Close()
	strategy := &Strategy{
		ID:          getStrategyID(),
		Name:        name,
		Author:      author,
		Description: desc,
		Path:        STRATEGIES_DIR + fileName,
		Matches:     []MatchID{},
	}
	trn.add(strategy)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func HandleReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
