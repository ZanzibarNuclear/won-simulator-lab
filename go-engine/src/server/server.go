package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"won/sim-lab/go-engine/sim"
)

type Server struct {
	simulation *sim.Simulation
}

func NewServer(simulation *sim.Simulation) *Server {
	return &Server{simulation: simulation}
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/status", s.handleStatus)
	http.HandleFunc("/start", s.handleStart)
	http.HandleFunc("/stop", s.handleStop)
	http.HandleFunc("/advance", s.handleAdvance)

	fmt.Printf("Server listening on port %s\n", port)
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

func (s *Server) handleAdvance(w http.ResponseWriter, r *http.Request) {
	ticks := 10 // Default value
	if ticksParam := r.URL.Query().Get("ticks"); ticksParam != "" {
		parsedTicks, err := strconv.Atoi(ticksParam)
		if err != nil || parsedTicks <= 0 {
			http.Error(w, "Invalid 'ticks' parameter. Must be a positive integer.", http.StatusBadRequest)
			return
		}
		ticks = parsedTicks
	}
	s.simulation.Advance(ticks)
	// w.WriteHeader(http.StatusOK)

	status := s.simulation.Status()
	json.NewEncoder(w).Encode(status)
}
