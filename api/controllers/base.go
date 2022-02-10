package controllers

import (
	"database/sql"
	"net/http"

	"NemWebGoApi/internal/config"
	"NemWebGoApi/internal/influxdb"
	"NemWebGoApi/internal/sqlite"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	SQLDb    *sql.DB
	InfluxDB influxdb2.Client
	Router   *mux.Router
	Config   *config.Config
}

func (s *Server) Init(cfg *config.Config) error {
	s.SQLDb = sqlite.New(cfg.SQLFilePath())
	s.InfluxDB = influxdb.New(cfg.InfluxHost(), cfg.InfluxToken())
	s.Router = mux.NewRouter()
	s.Config = cfg
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
			"http://127.0.0.1:3000",
			"http://127.0.0.1",
			"http://localhost:3005",
			"http://localhost:3000",
			"http://localhost",
		},
		AllowCredentials: true,
		// Debug:            app.Cfg.Debug(),
	})

	log.Infoln("Listening to Port ", port)
	log.Fatal(http.ListenAndServe(port, corsWrapper.Handler(s.Router)))
}
