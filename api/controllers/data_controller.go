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

	// TODO: Think of better method to filter, very confusing already caught me out twice
	// Currently if there are no DUID filters given it will then search
	filter := models.FilterMapToGenerationFilter(r.URL.Query())
	if len(filter.DuID.GetEq()) == 0 {
		units, err := unit.ReadAll(
			s.SQLDb,
			models.ParseUnitFilterMap(r.URL.Query()),
		)

		if err != nil {
			log.Debugln("Error Reading Units:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		duids := []string{}
		for _, val := range *units {
			duids = append(duids, val.DuID)
		}

		duids = append(duids, filter.DuID.GetEq()[:]...)
		filter.DuID.SetEq(duids)
	}

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

func (s *Server) GetGenerationDataGrouped(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterMapToGenerationGroupedFilter(r.URL.Query())

	units, _, err := filter.GetAllGroupUnitCombinations(s.SQLDb, r.URL.Query())
	if err != nil {
		log.Debugln("Error Getting Grouped Generation Data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := models.ReadGroupedGenerationData(
		s.InfluxDB.QueryAPI(s.Config.InfluxOrg()),
		s.Config.InfluxBucket(),
		filter,
		units,
	)

	if err != nil {
		log.Debugln("Error Getting Grouped Generation Data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.respond(w, r, data, http.StatusOK)
	return
}
