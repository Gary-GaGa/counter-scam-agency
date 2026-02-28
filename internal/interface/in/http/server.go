package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	portin "counter-scam-agency/internal/usecase/port/in"
)

const maxBodySize = 1 << 20 // 1 MB

var validID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

// Server is the HTTP REST API server that wraps usecase ports.
type Server struct {
	mux            *http.ServeMux
	investigations portin.InvestigationUsecase
	personnel      portin.PersonnelUsecase
	defense        portin.DefenseUsecase
}

// NewServer creates a new HTTP server with all routes registered.
func NewServer(
	investigations portin.InvestigationUsecase,
	personnel portin.PersonnelUsecase,
	defense portin.DefenseUsecase,
) *Server {
	s := &Server{
		mux:            http.NewServeMux(),
		investigations: investigations,
		personnel:      personnel,
		defense:        defense,
	}
	s.registerRoutes()
	return s
}

// Handler returns the underlying http.Handler for use with http.ListenAndServe.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	// Investigation / Mission endpoints.
	s.mux.HandleFunc("GET /api/missions", s.handleListMissions)
	s.mux.HandleFunc("GET /api/missions/{id}", s.handleGetMission)
	s.mux.HandleFunc("POST /api/investigations", s.handleStartInvestigation)
	s.mux.HandleFunc("POST /api/investigations/{id}/advance", s.handleAdvanceNode)
	s.mux.HandleFunc("POST /api/investigations/{id}/evidence", s.handleSubmitEvidence)
	s.mux.HandleFunc("POST /api/investigations/{id}/complete", s.handleCompleteInvestigation)

	// Personnel / Player & Skill endpoints.
	s.mux.HandleFunc("POST /api/players", s.handleCreatePlayer)
	s.mux.HandleFunc("GET /api/players/{id}", s.handleGetPlayer)
	s.mux.HandleFunc("POST /api/players/{id}/stats", s.handleUpdateStats)
	s.mux.HandleFunc("GET /api/players/{id}/skills", s.handleListSkills)
	s.mux.HandleFunc("POST /api/players/{id}/skills/{skillID}/unlock", s.handleUnlockSkill)
	s.mux.HandleFunc("POST /api/players/{id}/skills/{skillID}/equip", s.handleEquipSkill)
	s.mux.HandleFunc("POST /api/players/{id}/skills/{skillID}/activate", s.handleActivateSkill)
	s.mux.HandleFunc("POST /api/players/{id}/skills/tick", s.handleTickCooldowns)

	// Defense / Base endpoints.
	s.mux.HandleFunc("POST /api/bases", s.handleCreateBase)
	s.mux.HandleFunc("GET /api/bases/{id}", s.handleGetBase)
	s.mux.HandleFunc("POST /api/bases/{id}/facilities", s.handleAddFacility)
	s.mux.HandleFunc("POST /api/bases/{id}/security/upgrade", s.handleUpgradeSecurity)
	s.mux.HandleFunc("POST /api/bases/{id}/facilities/{facilityID}/upgrade", s.handleUpgradeFacility)
}

// writeJSON encodes v as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("http: encode response: %v", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// decodeJSON reads and decodes a JSON request body into dst.
func decodeJSON(r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// internalError logs the real error and returns a generic message to the client.
func internalError(w http.ResponseWriter, handler string, err error) {
	log.Printf("%s: %v", handler, err)
	writeError(w, http.StatusInternalServerError, "internal error")
}

// validatePathID checks that a path parameter is a safe identifier.
func validatePathID(w http.ResponseWriter, id, name string) bool {
	if !validID.MatchString(id) {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid %s", name))
		return false
	}
	return true
}
