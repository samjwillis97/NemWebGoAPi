package controllers

import (
	"database/sql"
	"net/http"

	"NemWebGoApi/internal/sqlite"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	SQLDb  *sql.DB
	Router *mux.Router
}

func (s *Server) Init(SQLFilePath string) error {
	s.SQLDb = sqlite.New(SQLFilePath)
	s.Router = mux.NewRouter()
	s.initializeRoutes()
	return nil
}

func (s *Server) Run(port string) {
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "OPTIONS", "POST", "DELETE", "PUT", "PATCH"},
		AllowedHeaders: []string{"Content-type", "Origin", "Accept", "*"},
		// AllowedHeaders: []string{"Content-type", "Origin", "Accept", "Access-Control-Allow-Origin"},
		AllowedOrigins: []string{
			"http://127.0.0.1:3005",
			"http://127.0.0.1",
			"http://localhost:3005",
			"http://localhost",
		},
		AllowCredentials: true,
		// Debug:            app.Cfg.Debug(),
	})

	log.Infoln("Listening to Port ", port)
	log.Fatal(http.ListenAndServe(port, corsWrapper.Handler(s.Router)))
}
