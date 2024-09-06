package server

import (
	"encoding/json"
	"net/http"
	"won/sim-lab/go-engine/sim"
)

type Server struct {
	simulation *sim.Simulation
}

func NewServer(simulation *sim.Simulation) *Server {
	return &Server{simulation: simulation}
}

func (s *Server) Start() error {
	http.HandleFunc("/status", s.handleStatus)
	http.HandleFunc("/start", s.handleStart)
	http.HandleFunc("/stop", s.handleStop)

	return http.ListenAndServe(":8080", nil)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := s.simulation.Status()

	json.NewEncoder(w).Encode(status)
}

func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	s.simulation.Start()
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	s.simulation.Stop()
	w.WriteHeader(http.StatusOK)
}
