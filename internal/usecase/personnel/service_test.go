package personnel

import (
	"context"
	"testing"

	domain "counter-scam-agency/internal/domain/personnel"

	"github.com/stretchr/testify/assert"
)

type playerRepo struct {
	players map[string]*domain.Player
}

func (r *playerRepo) Save(_ context.Context, player *domain.Player) error {
	r.players[player.ID] = player
	return nil
}

func (r *playerRepo) FindByID(_ context.Context, id string) (*domain.Player, error) {
	return r.players[id], nil
}

func TestSkillFlow(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	skill := domain.Skill{
		ID:                 "skill-1",
		Name:               "Rapid",
		CooldownSeconds:    5,
		ReputationRequired: 10,
	}
	svc := NewService(repo, []domain.Skill{skill})
	player := domain.NewPlayer("player-1")
	player.AddReputation(10)
	repo.players[player.ID] = player

	unlockRes, err := svc.UnlockSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)
	assert.True(t, unlockRes.Unlocked)

	equipRes, err := svc.EquipSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)
	assert.True(t, equipRes.Equipped)

	activateRes, err := svc.ActivateSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)
	assert.Greater(t, activateRes.CooldownRemaining, 0)
}
