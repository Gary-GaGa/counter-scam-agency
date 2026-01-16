package investigation_test

import (
	"context"
	"testing"

	"counter-scam-agency/internal/domain/intelligence"
	"counter-scam-agency/internal/domain/operation"
	"counter-scam-agency/internal/domain/personnel"
	"counter-scam-agency/internal/usecase/investigation"

	"github.com/stretchr/testify/assert"
)

type missionRepo struct {
	missions map[string]*intelligence.Mission
}

func (r *missionRepo) FindByID(_ context.Context, id string) (*intelligence.Mission, error) {
	return r.missions[id], nil
}

func (r *missionRepo) FindAll(_ context.Context) ([]*intelligence.Mission, error) {
	results := make([]*intelligence.Mission, 0, len(r.missions))
	for _, mission := range r.missions {
		results = append(results, mission)
	}
	return results, nil
}

func (r *missionRepo) Save(_ context.Context, mission *intelligence.Mission) error {
	r.missions[mission.ID] = mission
	return nil
}

type investigationRepo struct {
	investigations map[string]*operation.Investigation
}

func (r *investigationRepo) Save(_ context.Context, inv *operation.Investigation) error {
	r.investigations[inv.ID] = inv
	return nil
}

func (r *investigationRepo) FindByID(_ context.Context, id string) (*operation.Investigation, error) {
	return r.investigations[id], nil
}

func (r *investigationRepo) FindByPlayerID(_ context.Context, playerID string) ([]*operation.Investigation, error) {
	results := make([]*operation.Investigation, 0)
	for _, inv := range r.investigations {
		if inv.PlayerID == playerID {
			results = append(results, inv)
		}
	}
	return results, nil
}

type playerRepo struct {
	players   map[string]*personnel.Player
	saveCount int
}

func (r *playerRepo) Save(_ context.Context, player *personnel.Player) error {
	r.players[player.ID] = player
	r.saveCount++
	return nil
}

func (r *playerRepo) FindByID(_ context.Context, id string) (*personnel.Player, error) {
	return r.players[id], nil
}

func TestCompleteInvestigationSuccess(t *testing.T) {
	missions := &missionRepo{missions: map[string]*intelligence.Mission{}}
	investigations := &investigationRepo{investigations: map[string]*operation.Investigation{}}
	players := &playerRepo{players: map[string]*personnel.Player{}}

	mission := intelligence.NewMission("mission-1", "Case A", "desc", intelligence.ScamTypePhishing, 3, 2)
	missions.missions[mission.ID] = mission

	player := personnel.NewPlayer("player-1")
	players.players[player.ID] = player

	inv := operation.NewInvestigation("inv-1", player.ID, mission.ID)
	inv.Complete()
	investigations.investigations[inv.ID] = inv

	svc := investigation.NewService(missions, investigations, players)
	result, err := svc.CompleteInvestigation(context.Background(), inv.ID)

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 6, result.ReputationGained)
	assert.Equal(t, 6, players.players[player.ID].Reputation)
	assert.Equal(t, 1, players.saveCount)
}

func TestCompleteInvestigationFailure(t *testing.T) {
	missions := &missionRepo{missions: map[string]*intelligence.Mission{}}
	investigations := &investigationRepo{investigations: map[string]*operation.Investigation{}}
	players := &playerRepo{players: map[string]*personnel.Player{}}

	mission := intelligence.NewMission("mission-1", "Case A", "desc", intelligence.ScamTypePhishing, 3, 2)
	missions.missions[mission.ID] = mission

	player := personnel.NewPlayer("player-1")
	players.players[player.ID] = player

	inv := operation.NewInvestigation("inv-1", player.ID, mission.ID)
	inv.Fail()
	investigations.investigations[inv.ID] = inv

	svc := investigation.NewService(missions, investigations, players)
	result, err := svc.CompleteInvestigation(context.Background(), inv.ID)

	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, 0, result.ReputationGained)
	assert.Equal(t, 0, players.players[player.ID].Reputation)
	assert.Equal(t, 0, players.saveCount)
}

func TestCompleteInvestigationActive(t *testing.T) {
	missions := &missionRepo{missions: map[string]*intelligence.Mission{}}
	investigations := &investigationRepo{investigations: map[string]*operation.Investigation{}}
	players := &playerRepo{players: map[string]*personnel.Player{}}

	mission := intelligence.NewMission("mission-1", "Case A", "desc", intelligence.ScamTypePhishing, 1, 1)
	missions.missions[mission.ID] = mission

	player := personnel.NewPlayer("player-1")
	players.players[player.ID] = player

	inv := operation.NewInvestigation("inv-1", player.ID, mission.ID)
	investigations.investigations[inv.ID] = inv

	svc := investigation.NewService(missions, investigations, players)
	_, err := svc.CompleteInvestigation(context.Background(), inv.ID)

	assert.ErrorIs(t, err, investigation.ErrInvestigationNotFinished)
}
