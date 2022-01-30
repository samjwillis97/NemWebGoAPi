package controllers

import (
	"NemWebGoApi/api/models"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) GetAllUnits(w http.ResponseWriter, r *http.Request) {
	unit := models.Unit{}
	units, err := unit.ReadAll(s.SQLDb, models.ParseUnitFilterMap(r.URL.Query()))
	if err != nil {
		log.Debugln("Error Reading Units:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.respond(w, r, units, http.StatusOK)
	return
}
