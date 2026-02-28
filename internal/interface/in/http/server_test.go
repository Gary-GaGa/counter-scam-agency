package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "counter-scam-agency/internal/interface/in/http"
	"counter-scam-agency/internal/usecase/dto"

	"github.com/stretchr/testify/assert"
)

// stubInvestigation implements portin.InvestigationUsecase for testing.
type stubInvestigation struct {
	missions []dto.MissionSummary
	detail   *dto.MissionDetail
}

func (s *stubInvestigation) ListMissions(_ context.Context) ([]dto.MissionSummary, error) {
	return s.missions, nil
}

func (s *stubInvestigation) GetMission(_ context.Context, _ string) (*dto.MissionDetail, error) {
	if s.detail == nil {
		return nil, assert.AnError
	}
	return s.detail, nil
}

func (s *stubInvestigation) StartInvestigation(_ context.Context, investigationID, playerID, missionID, _ string) (*dto.InvestigationStartResult, error) {
	return &dto.InvestigationStartResult{
		InvestigationID: investigationID,
		PlayerID:        playerID,
		MissionID:       missionID,
		Status:          "active",
	}, nil
}

func (s *stubInvestigation) AdvanceNode(_ context.Context, investigationID, nodeID, optionID string) (*dto.NodeProgressResult, error) {
	return &dto.NodeProgressResult{
		InvestigationID: investigationID,
		NodeID:          nodeID,
		OptionID:        optionID,
		Status:          "active",
	}, nil
}

func (s *stubInvestigation) SubmitEvidence(_ context.Context, investigationID, evidenceID string) (*dto.SubmitEvidenceResult, error) {
	return &dto.SubmitEvidenceResult{
		InvestigationID: investigationID,
		EvidenceID:      evidenceID,
	}, nil
}

func (s *stubInvestigation) CompleteInvestigation(_ context.Context, investigationID string) (*dto.CompleteResult, error) {
	return &dto.CompleteResult{
		InvestigationID: investigationID,
		Success:         true,
	}, nil
}

// stubPersonnel implements portin.PersonnelUsecase for testing.
type stubPersonnel struct {
	skills []dto.SkillSummary
}

func (s *stubPersonnel) CreatePlayer(_ context.Context, playerID string) (*dto.PlayerSummary, error) {
	return &dto.PlayerSummary{
		ID:         playerID,
		Reputation: 0,
		Stats:      dto.StatsSummary{Logic: 10, Tech: 10, Charisma: 10, Resilience: 10},
	}, nil
}

func (s *stubPersonnel) GetPlayer(_ context.Context, playerID string) (*dto.PlayerSummary, error) {
	return &dto.PlayerSummary{
		ID:         playerID,
		Reputation: 50,
		Stats:      dto.StatsSummary{Logic: 12, Tech: 10, Charisma: 11, Resilience: 10},
	}, nil
}

func (s *stubPersonnel) ListSkills(_ context.Context, _ string) ([]dto.SkillSummary, error) {
	return s.skills, nil
}

func (s *stubPersonnel) UnlockSkill(_ context.Context, playerID, skillID string) (*dto.SkillActionResult, error) {
	return &dto.SkillActionResult{PlayerID: playerID, SkillID: skillID, Unlocked: true}, nil
}

func (s *stubPersonnel) EquipSkill(_ context.Context, playerID, skillID string) (*dto.SkillActionResult, error) {
	return &dto.SkillActionResult{PlayerID: playerID, SkillID: skillID, Equipped: true}, nil
}

func (s *stubPersonnel) ActivateSkill(_ context.Context, playerID, skillID string) (*dto.SkillActionResult, error) {
	return &dto.SkillActionResult{PlayerID: playerID, SkillID: skillID, CooldownRemaining: 5}, nil
}

func (s *stubPersonnel) TickSkillCooldowns(_ context.Context, playerID string, _ int) (*dto.SkillActionResult, error) {
	return &dto.SkillActionResult{PlayerID: playerID}, nil
}

// stubDefense implements portin.DefenseUsecase for testing.
type stubDefense struct {
	base *dto.BaseSummary
}

func (s *stubDefense) CreateBase(_ context.Context, baseID, ownerID string, slots int) (*dto.BaseSummary, error) {
	return &dto.BaseSummary{ID: baseID, OwnerID: ownerID, FacilitySlots: slots, SecurityLevel: 1}, nil
}

func (s *stubDefense) GetBase(_ context.Context, _ string) (*dto.BaseSummary, error) {
	if s.base == nil {
		return nil, assert.AnError
	}
	return s.base, nil
}

func (s *stubDefense) AddFacility(_ context.Context, baseID string, f dto.FacilityInput) (*dto.BaseSummary, error) {
	return &dto.BaseSummary{
		ID: baseID,
		Facilities: []dto.FacilitySummary{
			{ID: f.ID, Name: f.Name, Type: f.Type, Level: f.Level, MaxLevel: f.MaxLevel},
		},
	}, nil
}

func (s *stubDefense) UpgradeSecurity(_ context.Context, baseID string, _ int) (*dto.BaseSummary, error) {
	return &dto.BaseSummary{ID: baseID, SecurityLevel: 2}, nil
}

func (s *stubDefense) UpgradeFacility(_ context.Context, baseID, _ string) (*dto.BaseSummary, error) {
	return &dto.BaseSummary{ID: baseID}, nil
}

func newTestServer() (*handler.Server, *stubInvestigation, *stubPersonnel, *stubDefense) {
	inv := &stubInvestigation{
		missions: []dto.MissionSummary{{ID: "m1", Title: "Test Mission"}},
		detail:   &dto.MissionDetail{ID: "m1", Title: "Test Mission"},
	}
	per := &stubPersonnel{
		skills: []dto.SkillSummary{{ID: "s1", Name: "Alpha"}},
	}
	def := &stubDefense{
		base: &dto.BaseSummary{ID: "b1", OwnerID: "p1", SecurityLevel: 1},
	}
	srv := handler.NewServer(inv, per, def)
	return srv, inv, per, def
}

func TestListMissions(t *testing.T) {
	srv, _, _, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/missions", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestGetMission(t *testing.T) {
	srv, _, _, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/missions/m1", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStartInvestigation(t *testing.T) {
	srv, _, _, _ := newTestServer()
	body, _ := json.Marshal(map[string]string{
		"investigationId": "inv-1",
		"playerId":        "p1",
		"missionId":       "m1",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/investigations", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAdvanceNode(t *testing.T) {
	srv, _, _, _ := newTestServer()
	body, _ := json.Marshal(map[string]string{
		"nodeId":   "node-1",
		"optionId": "opt-1",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/investigations/inv-1/advance", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListSkills(t *testing.T) {
	srv, _, _, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/players/p1/skills", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBase(t *testing.T) {
	srv, _, _, _ := newTestServer()
	body, _ := json.Marshal(map[string]any{
		"baseId":  "b1",
		"ownerId": "p1",
		"slots":   2,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/bases", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBase(t *testing.T) {
	srv, _, _, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/bases/b1", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePlayer(t *testing.T) {
	srv, _, _, _ := newTestServer()
	body, _ := json.Marshal(map[string]string{"playerId": "p1"})
	req := httptest.NewRequest(http.MethodPost, "/api/players", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "p1")
}

func TestGetPlayer(t *testing.T) {
	srv, _, _, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/players/p1", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "p1")
}
