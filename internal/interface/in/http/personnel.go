package http

import "net/http"

// handleListSkills returns available skills for a player.
func (s *Server) handleListSkills(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skills, err := s.personnel.ListSkills(r.Context(), playerID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, skills)
}

// handleUnlockSkill unlocks a skill for a player.
func (s *Server) handleUnlockSkill(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skillID := r.PathValue("skillID")
	result, err := s.personnel.UnlockSkill(r.Context(), playerID, skillID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleEquipSkill equips a skill to the AI partner.
func (s *Server) handleEquipSkill(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skillID := r.PathValue("skillID")
	result, err := s.personnel.EquipSkill(r.Context(), playerID, skillID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// handleActivateSkill activates a skill if equipped and not on cooldown.
func (s *Server) handleActivateSkill(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("id")
	skillID := r.PathValue("skillID")
	result, err := s.personnel.ActivateSkill(r.Context(), playerID, skillID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
	var req tickCooldownsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := s.personnel.TickSkillCooldowns(r.Context(), playerID, req.Seconds)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
