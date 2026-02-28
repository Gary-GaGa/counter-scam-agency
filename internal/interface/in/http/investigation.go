package http

import (
	"log"
	"net/http"
)

// handleListMissions returns all available missions.
func (s *Server) handleListMissions(w http.ResponseWriter, r *http.Request) {
	missions, err := s.investigations.ListMissions(r.Context())
	if err != nil {
		internalError(w, "handleListMissions", err)
		return
	}
	writeJSON(w, http.StatusOK, missions)
}

// handleGetMission returns mission details by ID.
func (s *Server) handleGetMission(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validatePathID(w, id, "id") {
		return
	}
	detail, err := s.investigations.GetMission(r.Context(), id)
	if err != nil {
		log.Printf("handleGetMission: %v", err)
		writeError(w, http.StatusNotFound, "mission not found")
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

type startInvestigationRequest struct {
	InvestigationID string `json:"investigationId"`
	PlayerID        string `json:"playerId"`
	MissionID       string `json:"missionId"`
	StartNodeID     string `json:"startNodeId"`
}

// handleStartInvestigation creates a new investigation.
func (s *Server) handleStartInvestigation(w http.ResponseWriter, r *http.Request) {
	var req startInvestigationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.investigations.StartInvestigation(r.Context(), req.InvestigationID, req.PlayerID, req.MissionID, req.StartNodeID)
	if err != nil {
		log.Printf("handleStartInvestigation: %v", err)
		writeError(w, http.StatusBadRequest, "failed to start investigation")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

type advanceNodeRequest struct {
	NodeID   string `json:"nodeId"`
	OptionID string `json:"optionId"`
}

// handleAdvanceNode advances the investigation to the next node.
func (s *Server) handleAdvanceNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validatePathID(w, id, "id") {
		return
	}
	var req advanceNodeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.investigations.AdvanceNode(r.Context(), id, req.NodeID, req.OptionID)
	if err != nil {
		log.Printf("handleAdvanceNode: %v", err)
		writeError(w, http.StatusBadRequest, "failed to advance node")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type submitEvidenceRequest struct {
	EvidenceID string `json:"evidenceId"`
}

// handleSubmitEvidence submits evidence for the investigation.
func (s *Server) handleSubmitEvidence(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validatePathID(w, id, "id") {
		return
	}
	var req submitEvidenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.investigations.SubmitEvidence(r.Context(), id, req.EvidenceID)
	if err != nil {
		log.Printf("handleSubmitEvidence: %v", err)
		writeError(w, http.StatusBadRequest, "failed to submit evidence")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleCompleteInvestigation finalizes the investigation.
func (s *Server) handleCompleteInvestigation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validatePathID(w, id, "id") {
		return
	}
	result, err := s.investigations.CompleteInvestigation(r.Context(), id)
	if err != nil {
		log.Printf("handleCompleteInvestigation: %v", err)
		writeError(w, http.StatusBadRequest, "failed to complete investigation")
		return
	}
	writeJSON(w, http.StatusOK, result)
}
