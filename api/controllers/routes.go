package controllers

import "NemWebGoApi/api/middlewares"

func (s *Server) initializeRoutes() {
	s.Router.Use(middlewares.LoggingMW)
	s.Router.Use(middlewares.TimingMW)

	unitRouter := s.Router.PathPrefix("/units").Subrouter()
	unitRouter.HandleFunc("", s.GetAllUnits).Methods("GET")

	dataRouter := s.Router.PathPrefix("/data").Subrouter()
	dataRouter.HandleFunc("/demand", s.GetDemandData).Methods("GET")

	dataRouter.HandleFunc("/rooftop", s.GetRooftopData).Methods("GET")

	dataRouter.HandleFunc("/generation", s.GetGeneratingData).Methods("GET")
	dataRouter.HandleFunc("/generation/grouped", s.GetGenerationDataGrouped).Methods("GET")
}
