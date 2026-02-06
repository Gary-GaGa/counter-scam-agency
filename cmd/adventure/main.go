package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domainpersonnel "counter-scam-agency/internal/domain/personnel"
	mongoinfra "counter-scam-agency/internal/infrastructure/persistence/mongo"
	invrepo "counter-scam-agency/internal/interface/out/persistence/mongo/investigation"
	"counter-scam-agency/internal/interface/out/persistence/mongo/mission"
	"counter-scam-agency/internal/interface/out/persistence/mongo/player"
	"counter-scam-agency/internal/usecase/dto"
	usecase "counter-scam-agency/internal/usecase/investigation"
)

// main provides a streamlined text-adventure mode focused on narrative immersion.
// Unlike the full CLI (cmd/cli), this entry point skips skill/base menus
// and jumps directly into a mission investigation loop.
func main() {
	ctx := context.Background()
	cfg := mongoinfra.Config{
		URI:      getenv("MONGO_URI", "mongodb://localhost:27017"),
		Database: getenv("MONGO_DB", "counter_scam_agency"),
	}

	client, err := mongoinfra.NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = mongoinfra.Disconnect(ctx, client)
	}()

	db, err := mongoinfra.NewDatabase(client, cfg)
	if err != nil {
		panic(err)
	}

	missionsRepo := mission.NewMongoRepository(db)
	investigationsRepo := invrepo.NewMongoRepository(db)
	playersRepo := player.NewMongoRepository(db)
	svc := usecase.NewService(missionsRepo, investigationsRepo, playersRepo)

	playerID := getenv("PLAYER_ID", "player-1")
	if err := ensurePlayer(ctx, playersRepo, playerID); err != nil {
		panic(err)
	}

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║     Counter Scam Agency — 冒險模式    ║")
	fmt.Println("╚══════════════════════════════════════╝")

	for {
		missionID, err := selectMission(ctx, svc)
		if err != nil {
			fmt.Printf("錯誤：%v\n", err)
			return
		}

		missionDetail, err := svc.GetMission(ctx, missionID)
		if err != nil {
			fmt.Printf("錯誤：%v\n", err)
			return
		}

		renderMissionBrief(missionDetail)
		renderVictimProfile(missionDetail)

		invID := fmt.Sprintf("inv-%d", time.Now().UnixNano())
		startResult, err := svc.StartInvestigation(ctx, invID, playerID, missionID, "")
		if err != nil {
			fmt.Printf("開始調查失敗：%v\n", err)
			return
		}

		runAdventure(ctx, svc, invID, missionDetail, startResult.CurrentNodeID)

		complete, err := svc.CompleteInvestigation(ctx, invID)
		if err != nil {
			fmt.Printf("結案失敗：%v\n", err)
			return
		}

		renderSummary(complete)

		fmt.Print("\n繼續下一個案件？(y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) != "y" {
			fmt.Println("感謝遊玩，下次見！")
			return
		}
	}
}

func runAdventure(ctx context.Context, svc *usecase.Service, invID string, missionDetail *dto.MissionDetail, startNodeID string) {
	currentNodeID := startNodeID
	for {
		node := findNode(missionDetail, currentNodeID)
		if node == nil {
			fmt.Println("找不到節點，流程終止。")
			return
		}

		renderNode(node)
		if node.IsTerminal || len(node.Options) == 0 {
			return
		}

		selectedOption := chooseOption(node)
		progress, err := svc.AdvanceNode(ctx, invID, node.ID, selectedOption.ID)
		if err != nil {
			fmt.Printf("推進失敗：%v\n", err)
			return
		}

		renderEvidence(missionDetail, selectedOption.EvidenceIDs)
		currentNodeID = progress.NextNodeID
		if progress.Status != "Active" {
			return
		}
	}
}

func ensurePlayer(ctx context.Context, repo *player.MongoRepository, playerID string) error {
	found, err := repo.FindByID(ctx, playerID)
	if err == nil && found != nil {
		return nil
	}
	if err != nil && !strings.Contains(err.Error(), "no documents") {
		return err
	}
	newPlayer := domainpersonnel.NewPlayer(playerID)
	return repo.Save(ctx, newPlayer)
}

func selectMission(ctx context.Context, svc *usecase.Service) (string, error) {
	missions, err := svc.ListMissions(ctx)
	if err != nil {
		return "", err
	}
	if len(missions) == 0 {
		return "", fmt.Errorf("沒有可用的任務，請先執行 seed-missions")
	}

	fmt.Println("\n--- 可用案件 ---")
	for i, m := range missions {
		fmt.Printf("%d) %s [%s] 難度:%d\n", i+1, m.Title, m.Type, m.Difficulty)
	}

	choice := promptIndex("選擇案件", len(missions))
	return missions[choice-1].ID, nil
}

func findNode(m *dto.MissionDetail, nodeID string) *dto.NarrativeNode {
	for i := range m.Nodes {
		if m.Nodes[i].ID == nodeID {
			return &m.Nodes[i]
		}
	}
	return nil
}

func renderMissionBrief(m *dto.MissionDetail) {
	if m == nil {
		return
	}
	fmt.Println("\n══════════════════════════════════════")
	fmt.Printf("  案件：%s\n", m.Title)
	fmt.Printf("  類型：%s  難度：%d  聲望權重：%d\n", m.Type, m.Difficulty, m.ReputationWeight)
	if strings.TrimSpace(m.Description) != "" {
		fmt.Printf("  描述：%s\n", m.Description)
	}
	fmt.Println("══════════════════════════════════════")
}

func renderVictimProfile(m *dto.MissionDetail) {
	if m == nil || m.VictimProfile == nil {
		return
	}
	p := m.VictimProfile
	fmt.Println("\n--- 受害者心理側寫 ---")
	fmt.Printf("焦慮:%d  信任:%d  急迫:%d  孤立:%d\n", p.Anxiety, p.Trust, p.Urgency, p.Isolation)
	fmt.Printf("風險分數:%d  風險等級:%s\n", p.RiskScore, p.RiskLevel)
}

func renderNode(node *dto.NarrativeNode) {
	fmt.Println("\n──────────────────────────────────────")
	fmt.Printf("【%s】\n\n", node.Title)
	fmt.Println(node.Body)
	fmt.Println("──────────────────────────────────────")
	for i, opt := range node.Options {
		fmt.Printf("  %d) %s\n", i+1, opt.Label)
	}
}

func chooseOption(node *dto.NarrativeNode) dto.NarrativeOption {
	choice := promptIndex("選擇行動", len(node.Options))
	return node.Options[choice-1]
}

func renderEvidence(m *dto.MissionDetail, evidenceIDs []string) {
	if len(evidenceIDs) == 0 {
		return
	}
	fmt.Println("\n  📋 取得證據：")
	for _, id := range evidenceIDs {
		for i := range m.EvidenceList {
			if m.EvidenceList[i].ID == id {
				fmt.Printf("    - %s (%s)\n", m.EvidenceList[i].Description, m.EvidenceList[i].Type)
				break
			}
		}
	}
}

func renderSummary(result *dto.CompleteResult) {
	fmt.Println("\n══════════════════════════════════════")
	if result.Success {
		fmt.Println("  ✅ 結案成功！")
	} else {
		fmt.Println("  ❌ 結案失敗")
	}
	fmt.Printf("  聲望變化：%d\n", result.ReputationGained)
	fmt.Println("══════════════════════════════════════")
}

func promptIndex(label string, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (1-%d): ", label, max)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		value, err := strconv.Atoi(input)
		if err == nil && value >= 1 && value <= max {
			return value
		}
		fmt.Println("輸入無效，請重新輸入。")
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
