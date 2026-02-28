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

func TestListSkills(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	skills := []domain.Skill{
		{ID: "s1", Name: "Alpha", CooldownSeconds: 3},
		{ID: "s2", Name: "Beta", CooldownSeconds: 5},
	}
	svc := NewService(repo, skills)
	player := domain.NewPlayer("player-1")
	repo.players[player.ID] = player

	result, err := svc.ListSkills(context.Background(), player.ID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListSkills_PlayerNotFound(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)

	_, err := svc.ListSkills(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrPlayerNotFound)
}

func TestTickSkillCooldowns(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	skill := domain.Skill{
		ID:                 "skill-1",
		Name:               "Rapid",
		CooldownSeconds:    10,
		ReputationRequired: 0,
	}
	svc := NewService(repo, []domain.Skill{skill})
	player := domain.NewPlayer("player-1")
	repo.players[player.ID] = player

	_, err := svc.UnlockSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)
	_, err = svc.EquipSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)
	_, err = svc.ActivateSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)

	result, err := svc.TickSkillCooldowns(context.Background(), player.ID, 5)
	assert.NoError(t, err)
	assert.Equal(t, player.ID, result.PlayerID)
}

func TestTickSkillCooldowns_PlayerNotFound(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)

	_, err := svc.TickSkillCooldowns(context.Background(), "missing", 5)
	assert.ErrorIs(t, err, ErrPlayerNotFound)
}

func TestEquipSkill_NotUnlocked(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	skill := domain.Skill{
		ID:                 "skill-1",
		Name:               "Rapid",
		CooldownSeconds:    5,
		ReputationRequired: 100,
	}
	svc := NewService(repo, []domain.Skill{skill})
	player := domain.NewPlayer("player-1")
	repo.players[player.ID] = player

	_, err := svc.EquipSkill(context.Background(), player.ID, "skill-1")
	assert.ErrorIs(t, err, ErrSkillLocked)
}

func TestActivateSkill_NotEquipped(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	skill := domain.Skill{
		ID:                 "skill-1",
		Name:               "Rapid",
		CooldownSeconds:    5,
		ReputationRequired: 0,
	}
	svc := NewService(repo, []domain.Skill{skill})
	player := domain.NewPlayer("player-1")
	repo.players[player.ID] = player

	_, err := svc.UnlockSkill(context.Background(), player.ID, "skill-1")
	assert.NoError(t, err)

	_, err = svc.ActivateSkill(context.Background(), player.ID, "skill-1")
	assert.ErrorIs(t, err, ErrSkillNotEquipped)
}

func TestUnlockSkill_SkillNotFound(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)
	player := domain.NewPlayer("player-1")
	repo.players[player.ID] = player

	_, err := svc.UnlockSkill(context.Background(), player.ID, "nonexistent")
	assert.ErrorIs(t, err, ErrSkillNotFound)
}

func TestCreatePlayer(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)

	result, err := svc.CreatePlayer(context.Background(), "player-1")
	assert.NoError(t, err)
	assert.Equal(t, "player-1", result.ID)
	assert.Equal(t, 0, result.Reputation)
	assert.Equal(t, 10, result.Stats.Logic)
}

func TestCreatePlayer_Duplicate(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)

	_, err := svc.CreatePlayer(context.Background(), "player-1")
	assert.NoError(t, err)

	_, err = svc.CreatePlayer(context.Background(), "player-1")
	assert.ErrorIs(t, err, ErrPlayerExists)
}

func TestGetPlayer(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)
	player := domain.NewPlayer("player-1")
	player.AddReputation(25)
	repo.players[player.ID] = player

	result, err := svc.GetPlayer(context.Background(), player.ID)
	assert.NoError(t, err)
	assert.Equal(t, "player-1", result.ID)
	assert.Equal(t, 25, result.Reputation)
}

func TestGetPlayer_NotFound(t *testing.T) {
	repo := &playerRepo{players: map[string]*domain.Player{}}
	svc := NewService(repo, nil)

	_, err := svc.GetPlayer(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrPlayerNotFound)
}
