package main

import (
	"net/http"

	"bytes"
	"io"
	"strings"

	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

type Server struct{}

func (s *Server) init() {
	http.HandleFunc("/upload", HandlerUpload)
	log.Fatal(http.ListenAndServe(":22101", nil))
}

func HandlerR

func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	var localBuff bytes.Buffer

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Body.Close()
		log.WithError(err).Error("Could not handle file upload")
		return
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	log.WithFields(log.Fields{
		"filename": header.Filename,
	}).Info("New strategy uploaded")
	io.Copy(&localBuff, file)

	writeErr := ioutil.WriteFile(name[0]+".ipd", localBuff.Bytes(), 0755)
	if writeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.Body.Close()
		log.WithError(err).Error("Issue opening file for writing")
		return
	}

	localBuff.Reset()
	w.WriteHeader(http.StatusOK)
	r.Body.Close()
}
