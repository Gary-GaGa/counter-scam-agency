package http

import "net/http"

// handleListMissions returns all available missions.
func (s *Server) handleListMissions(w http.ResponseWriter, r *http.Request) {
	missions, err := s.investigations.ListMissions(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, missions)
}

// handleGetMission returns mission details by ID.
func (s *Server) handleGetMission(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	detail, err := s.investigations.GetMission(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
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
		writeError(w, http.StatusBadRequest, err.Error())
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
	var req advanceNodeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.investigations.AdvanceNode(r.Context(), id, req.NodeID, req.OptionID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
	var req submitEvidenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.investigations.SubmitEvidence(r.Context(), id, req.EvidenceID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleCompleteInvestigation finalizes the investigation.
func (s *Server) handleCompleteInvestigation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.investigations.CompleteInvestigation(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
