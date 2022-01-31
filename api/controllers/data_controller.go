package controllers

import (
	"NemWebGoApi/api/models"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) GetDemandData(w http.ResponseWriter, r *http.Request) {
	data, err := models.ReadDemandData(s.InfluxDB.QueryAPI(s.Config.InfluxOrg()), s.Config.InfluxBucket())

	if err != nil {
		log.Debugln("Error Getting Demand Data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respond(w, r, data, http.StatusOK)
	return
}
