package http

import (
	"log"
	"net/http"
)

type createPlayerRequest struct {
	PlayerID string `json:"playerId"`
}

// handleCreatePlayer creates a new player.
func (s *Server) handleCreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req createPlayerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.personnel.CreatePlayer(r.Context(), req.PlayerID)
	if err != nil {
		log.Printf("handleCreatePlayer: %v", err)
		writeError(w, http.StatusBadRequest, "failed to create player")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// handleGetPlayer returns a player by ID.
func (s *Server) handleGetPlayer(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	if !validatePathID(w, playerID, "id") {
		return
	}
	result, err := s.personnel.GetPlayer(r.Context(), playerID)
	if err != nil {
		log.Printf("handleGetPlayer: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleListSkills returns available skills for a player.
func (s *Server) handleListSkills(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	if !validatePathID(w, playerID, "id") {
		return
	}
	skills, err := s.personnel.ListSkills(r.Context(), playerID)
	if err != nil {
		log.Printf("handleListSkills: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}
	writeJSON(w, http.StatusOK, skills)
}

// handleUnlockSkill unlocks a skill for a player.
func (s *Server) handleUnlockSkill(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skillID := r.PathValue("skillID")
	if !validatePathID(w, playerID, "id") || !validatePathID(w, skillID, "skillID") {
		return
	}
	result, err := s.personnel.UnlockSkill(r.Context(), playerID, skillID)
	if err != nil {
		log.Printf("handleUnlockSkill: %v", err)
		writeError(w, http.StatusBadRequest, "failed to unlock skill")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleEquipSkill equips a skill to the AI partner.
func (s *Server) handleEquipSkill(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skillID := r.PathValue("skillID")
	if !validatePathID(w, playerID, "id") || !validatePathID(w, skillID, "skillID") {
		return
	}
	result, err := s.personnel.EquipSkill(r.Context(), playerID, skillID)
	if err != nil {
		log.Printf("handleEquipSkill: %v", err)
		writeError(w, http.StatusBadRequest, "failed to equip skill")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleActivateSkill activates a skill if equipped and not on cooldown.
func (s *Server) handleActivateSkill(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skillID := r.PathValue("skillID")
	if !validatePathID(w, playerID, "id") || !validatePathID(w, skillID, "skillID") {
		return
	}
	result, err := s.personnel.ActivateSkill(r.Context(), playerID, skillID)
	if err != nil {
		log.Printf("handleActivateSkill: %v", err)
		writeError(w, http.StatusBadRequest, "failed to activate skill")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type tickCooldownsRequest struct {
	Seconds int `json:"seconds"`
}

// handleTickCooldowns reduces cooldown timers for all skills.
func (s *Server) handleTickCooldowns(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	if !validatePathID(w, playerID, "id") {
		return
	}
	var req tickCooldownsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.personnel.TickSkillCooldowns(r.Context(), playerID, req.Seconds)
	if err != nil {
		log.Printf("handleTickCooldowns: %v", err)
		writeError(w, http.StatusBadRequest, "failed to tick cooldowns")
		return
	}
	writeJSON(w, http.StatusOK, result)
}
