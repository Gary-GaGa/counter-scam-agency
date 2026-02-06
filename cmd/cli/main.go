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
	usecasepersonnel "counter-scam-agency/internal/usecase/personnel"
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
	skillSvc := usecasepersonnel.NewService(playersRepo, buildSkillCatalog())

	playerID := getenv("PLAYER_ID", "player-1")
	if err := ensurePlayer(ctx, playersRepo, playerID); err != nil {
		panic(err)
	}
	runSkillTreeMenu(ctx, skillSvc, playerID)

	missionID, err := selectMission(ctx, svc)
	if err != nil {
		panic(err)
	}

	missionDetail, err := svc.GetMission(ctx, missionID)
	if err != nil {
		panic(err)
	}
	renderMissionDetail(missionDetail)
	renderVictimProfile(missionDetail)

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
	newPlayer := domainpersonnel.NewPlayer(playerID)
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

	renderCasePanel(missions)

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

func renderMissionDetail(mission *dto.MissionDetail) {
	if mission == nil {
		return
	}
	fmt.Println("\n==================== 案件詳情 ====================")
	fmt.Printf("名稱:%s\n", mission.Title)
	fmt.Printf("類型:%s  難度:%d  聲望權重:%d\n", mission.Type, mission.Difficulty, mission.ReputationWeight)
	if strings.TrimSpace(mission.Description) != "" {
		fmt.Printf("描述:%s\n", mission.Description)
	}
}

func renderVictimProfile(mission *dto.MissionDetail) {
	if mission == nil || mission.VictimProfile == nil {
		return
	}
	profile := mission.VictimProfile
	fmt.Println("\n--- 受害者心理側寫 ---")
	fmt.Printf("焦慮:%d 信任:%d 急迫:%d 孤立:%d\n", profile.Anxiety, profile.Trust, profile.Urgency, profile.Isolation)
	fmt.Printf("風險分數:%d 風險等級:%s\n", profile.RiskScore, profile.RiskLevel)
}

func runSkillTreeMenu(ctx context.Context, svc *usecasepersonnel.Service, playerID string) {
	for {
		skills, err := svc.ListSkills(ctx, playerID)
		if err != nil {
			fmt.Printf("技能樹載入失敗：%v\n", err)
			return
		}
		renderSkillTree(skills)

		fmt.Println("1) 解鎖技能  2) 裝備技能  3) 啟動技能  4) 跳過")
		action := promptIndex("選擇操作", 4)
		if action == 4 {
			return
		}
		if len(skills) == 0 {
			fmt.Println("目前沒有可用技能。")
			continue
		}
		skillIndex := promptIndex("選擇技能編號", len(skills))
		skill := skills[skillIndex-1]

		switch action {
		case 1:
			res, err := svc.UnlockSkill(ctx, playerID, skill.ID)
			if err != nil {
				fmt.Printf("解鎖失敗：%v\n", err)
				continue
			}
			if res.Unlocked {
				fmt.Println("解鎖成功。")
			} else {
				fmt.Println("尚未達成解鎖條件。")
			}
		case 2:
			res, err := svc.EquipSkill(ctx, playerID, skill.ID)
			if err != nil {
				fmt.Printf("裝備失敗：%v\n", err)
				continue
			}
			if res.Equipped {
				fmt.Println("裝備成功。")
			} else {
				fmt.Println("裝備失敗。")
			}
		case 3:
			res, err := svc.ActivateSkill(ctx, playerID, skill.ID)
			if err != nil {
				fmt.Printf("啟動失敗：%v\n", err)
				continue
			}
			fmt.Printf("技能啟動，冷卻剩餘:%d\n", res.CooldownRemaining)
		}
	}
}

func renderSkillTree(skills []dto.SkillSummary) {
	fmt.Println("\n==================== 技能樹 ====================")
	if len(skills) == 0 {
		fmt.Println("尚無技能資料。")
		return
	}
	for i, skill := range skills {
		status := "未解鎖"
		if skill.Unlocked {
			status = "已解鎖"
		}
		if skill.Equipped {
			status += "/已裝備"
		}
		cooldown := "-"
		if skill.CooldownRemaining > 0 {
			cooldown = fmt.Sprintf("%d", skill.CooldownRemaining)
		}
		fmt.Printf("%d) %s [%s]  冷卻:%s  狀態:%s\n", i+1, skill.Name, skill.Type, cooldown, status)
		if strings.TrimSpace(skill.Description) != "" {
			fmt.Printf("   %s\n", skill.Description)
		}
	}
}

func buildSkillCatalog() []domainpersonnel.Skill {
	return []domainpersonnel.Skill{
		{
			ID:                 "skill-logic-1",
			Type:               domainpersonnel.SkillTypeAnalysis,
			Name:               "矛盾掃描",
			Description:        "短時間內分析對話矛盾，提高破案效率。",
			CooldownSeconds:    10,
			ReputationRequired: 0,
		},
		{
			ID:                 "skill-nego-1",
			Type:               domainpersonnel.SkillTypeNegotiation,
			Name:               "談判增幅",
			Description:        "提升談判成功率與獲得資訊品質。",
			CooldownSeconds:    12,
			ReputationRequired: 5,
		},
		{
			ID:                 "skill-defense-1",
			Type:               domainpersonnel.SkillTypeDefense,
			Name:               "防禦強化",
			Description:        "強化資安防線，降低詐騙滲透風險。",
			CooldownSeconds:    15,
			ReputationRequired: 10,
		},
	}
}

func renderCasePanel(missions []dto.MissionSummary) {
	fmt.Println("\n==================== 案件面板 ====================")
	for i, m := range missions {
		desc := trimText(m.Description, 48)
		fmt.Printf("%d) %s\n", i+1, m.Title)
		fmt.Printf("   類型:%s  難度:%d  聲望權重:%d\n", m.Type, m.Difficulty, m.ReputationWeight)
		if desc != "" {
			fmt.Printf("   摘要:%s\n", desc)
		}
		fmt.Println("------------------------------------------------")
	}
}

func trimText(text string, max int) string {
	clean := strings.TrimSpace(text)
	if clean == "" {
		return ""
	}
	if max <= 0 || len([]rune(clean)) <= max {
		return clean
	}
	runes := []rune(clean)
	return string(runes[:max]) + "…"
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
