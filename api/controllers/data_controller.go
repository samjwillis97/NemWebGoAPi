package controllers

import (
	"NemWebGoApi/api/models"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) GetDemandData(w http.ResponseWriter, r *http.Request) {
	data, err := models.ReadDemandData(
		s.InfluxDB.QueryAPI(s.Config.InfluxOrg()),
		s.Config.InfluxBucket(),
		models.FilterMaptoDemandFilter(r.URL.Query()),
	)

	if err != nil {
		log.Debugln("Error Getting Demand Data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respond(w, r, data, http.StatusOK)
	return
}

func (s *Server) GetRooftopData(w http.ResponseWriter, r *http.Request) {
	data, err := models.ReadRooftapData(
		s.InfluxDB.QueryAPI(s.Config.InfluxOrg()),
		s.Config.InfluxBucket(),
		models.FilterMaptoRooftopFilter(r.URL.Query()),
	)

	if err != nil {
		log.Debugln("Error Getting Demand Data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respond(w, r, data, http.StatusOK)
	return
}

func (s *Server) GetGeneratingData(w http.ResponseWriter, r *http.Request) {
	unit := models.Unit{}
	units, err := unit.ReadAll(
		s.SQLDb,
		models.ParseUnitFilterMap(r.URL.Query()),
	)

	duids := []string{}
	for _, val := range *units {
		duids = append(duids, val.DuID)
	}
	log.Warnln(duids)

	filter := models.FilterMapToGenerationFilter(r.URL.Query())
	filter.DuID.SetEq(duids)

	data, err := models.ReadGenerationData(
		s.InfluxDB.QueryAPI(s.Config.InfluxOrg()),
		s.Config.InfluxBucket(),
		filter,
	)

	if err != nil {
		log.Debugln("Error Getting Demand Data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respond(w, r, data, http.StatusOK)
	return
}
