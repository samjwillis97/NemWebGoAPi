package controllers

import "NemWebGoApi/api/middlewares"

func (s *Server) initializeRoutes() {
	s.Router.Use(middlewares.LoggingMW)
	s.Router.Use(middlewares.TimingMW)

	unitRouter := s.Router.PathPrefix("/units").Subrouter()
	unitRouter.HandleFunc("", s.GetAllUnits).Methods("GET")
}
