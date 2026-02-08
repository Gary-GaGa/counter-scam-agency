package http

import (
	"net/http"

	"counter-scam-agency/internal/usecase/dto"
)

type createBaseRequest struct {
	BaseID  string `json:"baseId"`
	OwnerID string `json:"ownerId"`
	Slots   int    `json:"slots"`
}

// handleCreateBase creates a new defense base.
func (s *Server) handleCreateBase(w http.ResponseWriter, r *http.Request) {
	var req createBaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.defense.CreateBase(r.Context(), req.BaseID, req.OwnerID, req.Slots)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// handleGetBase returns a base by ID.
func (s *Server) handleGetBase(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.defense.GetBase(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type addFacilityRequest struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	MaxLevel    int    `json:"maxLevel"`
	Description string `json:"description"`
}

// handleAddFacility adds a facility to a base.
func (s *Server) handleAddFacility(w http.ResponseWriter, r *http.Request) {
	baseID := r.PathValue("id")
	var req addFacilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.defense.AddFacility(r.Context(), baseID, dtoFacilityInput(req))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

type upgradeSecurityRequest struct {
	MaxLevel int `json:"maxLevel"`
}

// handleUpgradeSecurity increases the base security level.
func (s *Server) handleUpgradeSecurity(w http.ResponseWriter, r *http.Request) {
	baseID := r.PathValue("id")
	var req upgradeSecurityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.defense.UpgradeSecurity(r.Context(), baseID, req.MaxLevel)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleUpgradeFacility upgrades a specific facility.
func (s *Server) handleUpgradeFacility(w http.ResponseWriter, r *http.Request) {
	baseID := r.PathValue("id")
	facilityID := r.PathValue("facilityID")
	result, err := s.defense.UpgradeFacility(r.Context(), baseID, facilityID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func dtoFacilityInput(req addFacilityRequest) dto.FacilityInput {
	return dto.FacilityInput{
		ID:          req.ID,
		Type:        req.Type,
		Name:        req.Name,
		Level:       req.Level,
		MaxLevel:    req.MaxLevel,
		Description: req.Description,
	}
}
