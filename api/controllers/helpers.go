package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Warnln("Error Encoding JSON:", err)
		}
	}
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	//return json.NewDecoder(r.Body).Decode(v)
	decoder := json.NewDecoder(r.Body)
	for {
		err := decoder.Decode(v)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("decode: error decoding: %v", err)
		}
	}
}

// r.URL.Query().Get("limit")
