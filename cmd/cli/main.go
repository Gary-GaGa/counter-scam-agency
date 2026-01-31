package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"counter-scam-agency/internal/domain/personnel"
	mongoinfra "counter-scam-agency/internal/infrastructure/persistence/mongo"
	invrepo "counter-scam-agency/internal/infrastructure/persistence/mongo/investigation"
	"counter-scam-agency/internal/infrastructure/persistence/mongo/mission"
	"counter-scam-agency/internal/infrastructure/persistence/mongo/player"
	"counter-scam-agency/internal/usecase/dto"
	usecase "counter-scam-agency/internal/usecase/investigation"
)

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

	missionID, err := selectMission(ctx, svc)
	if err != nil {
		panic(err)
	}

	missionDetail, err := svc.GetMission(ctx, missionID)
	if err != nil {
		panic(err)
	}

	invID := fmt.Sprintf("inv-%d", time.Now().UnixNano())
	startResult, err := svc.StartInvestigation(ctx, invID, playerID, missionID, "")
	if err != nil {
		panic(err)
	}

	currentNodeID := startResult.CurrentNodeID
	for {
		node := findNode(missionDetail, currentNodeID)
		if node == nil {
			fmt.Println("找不到節點，流程終止。")
			return
		}

		renderNode(node)
		if node.IsTerminal || len(node.Options) == 0 {
			break
		}

		selectedOption := chooseOption(node)
		progress, err := svc.AdvanceNode(ctx, invID, node.ID, selectedOption.ID)
		if err != nil {
			panic(err)
		}

		renderEvidence(missionDetail, selectedOption.EvidenceIDs)
		currentNodeID = progress.NextNodeID
		if progress.Status != "Active" {
			break
		}
	}

	complete, err := svc.CompleteInvestigation(ctx, invID)
	if err != nil {
		panic(err)
	}

	renderSummary(complete)
}

func ensurePlayer(ctx context.Context, repo *player.MongoRepository, playerID string) error {
	found, err := repo.FindByID(ctx, playerID)
	if err == nil && found != nil {
		return nil
	}
	if err != nil && !strings.Contains(err.Error(), "no documents") {
		return err
	}
	newPlayer := personnel.NewPlayer(playerID)
	return repo.Save(ctx, newPlayer)
}

func selectMission(ctx context.Context, svc *usecase.Service) (string, error) {
	missions, err := svc.ListMissions(ctx)
	if err != nil {
		return "", err
	}
	if len(missions) == 0 {
		return "", fmt.Errorf("沒有可用的任務，請先 seed 資料")
	}

	fmt.Println("可用案件：")
	for i, m := range missions {
		fmt.Printf("%d) %s [%s] 難度:%d 聲望權重:%d\n", i+1, m.Title, m.Type, m.Difficulty, m.ReputationWeight)
	}

	choice := promptIndex("選擇案件編號", len(missions))
	return missions[choice-1].ID, nil
}

func findNode(mission *dto.MissionDetail, nodeID string) *dto.NarrativeNode {
	for i := range mission.Nodes {
		if mission.Nodes[i].ID == nodeID {
			return &mission.Nodes[i]
		}
	}
	return nil
}

func renderNode(node *dto.NarrativeNode) {
	fmt.Println("\n====================")
	fmt.Printf("%s\n\n", node.Title)
	fmt.Println(node.Body)
	fmt.Println("--------------------")
	for i, opt := range node.Options {
		fmt.Printf("%d) %s\n", i+1, opt.Label)
	}
}

func chooseOption(node *dto.NarrativeNode) dto.NarrativeOption {
	choice := promptIndex("選擇選項編號", len(node.Options))
	return node.Options[choice-1]
}

func renderEvidence(mission *dto.MissionDetail, evidenceIDs []string) {
	if len(evidenceIDs) == 0 {
		return
	}
	fmt.Println("\n取得證據：")
	for _, id := range evidenceIDs {
		if ev := findEvidence(mission, id); ev != nil {
			fmt.Printf("- %s (%s)\n", ev.Description, ev.Type)
		}
	}
}

func renderSummary(result *dto.CompleteResult) {
	fmt.Println("\n====================")
	if result.Success {
		fmt.Println("結案成功")
	} else {
		fmt.Println("結案失敗")
	}
	fmt.Printf("聲望變化：%d\n", result.ReputationGained)
	fmt.Println("====================")
}

func findEvidence(mission *dto.MissionDetail, evidenceID string) *dto.Evidence {
	for i := range mission.EvidenceList {
		if mission.EvidenceList[i].ID == evidenceID {
			return &mission.EvidenceList[i]
		}
	}
	return nil
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
