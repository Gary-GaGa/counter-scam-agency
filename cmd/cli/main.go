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
	defenserepo "counter-scam-agency/internal/interface/out/persistence/mongo/defense"
	invrepo "counter-scam-agency/internal/interface/out/persistence/mongo/investigation"
	"counter-scam-agency/internal/interface/out/persistence/mongo/mission"
	"counter-scam-agency/internal/interface/out/persistence/mongo/player"
	usecasedefense "counter-scam-agency/internal/usecase/defense"
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
		fmt.Fprintf(os.Stderr, "mongo connect: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = mongoinfra.Disconnect(ctx, client)
	}()

	db, err := mongoinfra.NewDatabase(client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mongo database: %v\n", err)
		os.Exit(1)
	}

	missionsRepo := mission.NewMongoRepository(db)
	investigationsRepo := invrepo.NewMongoRepository(db)
	playersRepo := player.NewMongoRepository(db)
	basesRepo := defenserepo.NewMongoRepository(db)

	svc := usecase.NewService(missionsRepo, investigationsRepo, playersRepo)
	skillSvc := usecasepersonnel.NewService(playersRepo, buildSkillCatalog())
	defenseSvc := usecasedefense.NewService(basesRepo)

	playerID := getenv("PLAYER_ID", "player-1")
	if err := ensurePlayer(ctx, playersRepo, playerID); err != nil {
		fmt.Fprintf(os.Stderr, "ensure player: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║    Counter Scam Agency — 指揮中心      ║")
	fmt.Println("╚══════════════════════════════════════╝")

	runMainMenu(ctx, svc, skillSvc, defenseSvc, playersRepo, playerID)
}

// ---------------------------------------------------------------------------
// Main Menu
// ---------------------------------------------------------------------------

func runMainMenu(ctx context.Context, svc *usecase.Service, skillSvc *usecasepersonnel.Service, defenseSvc *usecasedefense.Service, playersRepo *player.MongoRepository, playerID string) {
	for {
		fmt.Println("\n==================== 主選單 ====================")
		fmt.Println("1) 調查案件")
		fmt.Println("2) 技能樹管理")
		fmt.Println("3) AI 模組裝備")
		fmt.Println("4) 防禦基地管理")
		fmt.Println("5) 查看角色狀態")
		fmt.Println("6) 離開")

		choice := promptIndex("選擇操作", 6)
		switch choice {
		case 1:
			runInvestigation(ctx, svc, playerID)
		case 2:
			runSkillTreeMenu(ctx, skillSvc, playerID)
		case 3:
			runModuleMenu(ctx, playersRepo, playerID)
		case 4:
			runBaseMenu(ctx, defenseSvc, playerID)
		case 5:
			renderPlayerStatus(ctx, playersRepo, playerID)
		case 6:
			fmt.Println("感謝遊玩，下次見！")
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Investigation Flow
// ---------------------------------------------------------------------------

func runInvestigation(ctx context.Context, svc *usecase.Service, playerID string) {
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
	renderMissionDetail(missionDetail)
	renderVictimProfile(missionDetail)
	renderVictimStrategy(missionDetail)

	invID := fmt.Sprintf("inv-%d", time.Now().UnixNano())
	startResult, err := svc.StartInvestigation(ctx, invID, playerID, missionID, "")
	if err != nil {
		fmt.Printf("開始調查失敗：%v\n", err)
		return
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
			fmt.Printf("推進失敗：%v\n", err)
			return
		}

		renderEvidence(missionDetail, selectedOption.EvidenceIDs)
		currentNodeID = progress.NextNodeID
		if progress.Status != "Active" {
			break
		}
	}

	complete, err := svc.CompleteInvestigation(ctx, invID)
	if err != nil {
		fmt.Printf("結案失敗：%v\n", err)
		return
	}

	renderSummary(complete)
}

// ---------------------------------------------------------------------------
// AI Module Equipment
// ---------------------------------------------------------------------------

func runModuleMenu(ctx context.Context, playersRepo *player.MongoRepository, playerID string) {
	catalog := buildModuleCatalog()

	for {
		p, err := playersRepo.FindByID(ctx, playerID)
		if err != nil || p == nil {
			fmt.Println("無法載入角色資料。")
			return
		}

		fmt.Println("\n==================== AI 模組裝備 ====================")
		for i, mod := range catalog {
			status := "未解鎖"
			if p.IsModuleUnlocked(mod.ID) {
				status = "已解鎖"
			}
			equipped := false
			if p.Partner != nil {
				for _, installed := range p.Partner.Loadout {
					if installed.ID == mod.ID {
						equipped = true
						break
					}
				}
			}
			if equipped {
				status += "/已裝備"
			}
			fmt.Printf("%d) %s [%s]  需求聲望:%d  狀態:%s\n", i+1, mod.Name, mod.Type, mod.ReputationRequired, status)
			fmt.Printf("   %s\n", mod.Description)
			fmt.Printf("   加成：邏輯+%d 技術+%d 交涉+%d 韌性+%d\n",
				mod.StatBonus.Logic, mod.StatBonus.Tech, mod.StatBonus.Charisma, mod.StatBonus.Resilience)
		}

		fmt.Println("\n1) 解鎖模組  2) 裝備模組  3) 返回")
		action := promptIndex("選擇操作", 3)
		if action == 3 {
			return
		}

		modIndex := promptIndex("選擇模組編號", len(catalog))
		mod := catalog[modIndex-1]

		switch action {
		case 1:
			if p.UnlockModule(mod) {
				if err := playersRepo.Save(ctx, p); err != nil {
					fmt.Printf("儲存失敗：%v\n", err)
					continue
				}
				fmt.Printf("模組「%s」解鎖成功！\n", mod.Name)
			} else {
				fmt.Println("解鎖失敗：可能已解鎖或聲望不足。")
			}
		case 2:
			if p.EquipPartnerModule(mod) {
				if err := playersRepo.Save(ctx, p); err != nil {
					fmt.Printf("儲存失敗：%v\n", err)
					continue
				}
				fmt.Printf("模組「%s」裝備成功！\n", mod.Name)
			} else {
				fmt.Println("裝備失敗：可能尚未解鎖。")
			}
		}
	}
}

func buildModuleCatalog() []domainpersonnel.Module {
	return []domainpersonnel.Module{
		{
			ID:                 "mod-voice",
			Type:               domainpersonnel.ModuleVoiceAnalyzer,
			Name:               "語音分析模組",
			Description:        "自動標記語氣異常的關鍵字，在矛盾擊破中提供輔助。",
			Level:              1,
			ReputationRequired: 0,
			StatBonus:          domainpersonnel.Stats{Logic: 3},
		},
		{
			ID:                 "mod-crypto",
			Type:               domainpersonnel.ModuleCryptoTracer,
			Name:               "金流追蹤模組",
			Description:        "自動鎖定金流路徑，在訊號追蹤中提供輔助。",
			Level:              1,
			ReputationRequired: 5,
			StatBonus:          domainpersonnel.Stats{Tech: 3},
		},
		{
			ID:                 "mod-empathy",
			Type:               domainpersonnel.ModuleEmpathyEngine,
			Name:               "情感模擬模組",
			Description:        "預測對方情緒屬性，在談判牌局中提供輔助。",
			Level:              1,
			ReputationRequired: 8,
			StatBonus:          domainpersonnel.Stats{Charisma: 3},
		},
		{
			ID:                 "mod-firewall",
			Type:               domainpersonnel.ModuleMentalFirewall,
			Name:               "防火牆模組",
			Description:        "減少調查失敗時的壓力值扣除。",
			Level:              1,
			ReputationRequired: 10,
			StatBonus:          domainpersonnel.Stats{Resilience: 3},
		},
	}
}

// ---------------------------------------------------------------------------
// Defense Base Management
// ---------------------------------------------------------------------------

func runBaseMenu(ctx context.Context, defenseSvc *usecasedefense.Service, playerID string) {
	baseID := fmt.Sprintf("base-%s", playerID)

	for {
		fmt.Println("\n==================== 防禦基地 ====================")
		fmt.Println("1) 查看基地")
		fmt.Println("2) 建立基地")
		fmt.Println("3) 新增設施")
		fmt.Println("4) 升級設施")
		fmt.Println("5) 升級安全等級")
		fmt.Println("6) 返回")

		action := promptIndex("選擇操作", 6)
		if action == 6 {
			return
		}

		switch action {
		case 1:
			base, err := defenseSvc.GetBase(ctx, baseID)
			if err != nil {
				fmt.Printf("尚未建立基地或查詢失敗：%v\n", err)
				continue
			}
			renderBase(base)
		case 2:
			base, err := defenseSvc.CreateBase(ctx, baseID, playerID, 4)
			if err != nil {
				fmt.Printf("建立失敗：%v\n", err)
				continue
			}
			fmt.Println("基地建立成功！")
			renderBase(base)
		case 3:
			facilityTemplates := buildFacilityTemplates()
			fmt.Println("\n可用設施：")
			for i, f := range facilityTemplates {
				fmt.Printf("%d) %s [%s] — %s\n", i+1, f.Name, f.Type, f.Description)
			}
			idx := promptIndex("選擇設施", len(facilityTemplates))
			tmpl := facilityTemplates[idx-1]
			facID := fmt.Sprintf("fac-%s-%d", strings.ToLower(tmpl.Type), time.Now().UnixNano())

			base, err := defenseSvc.AddFacility(ctx, baseID, dto.FacilityInput{
				ID:          facID,
				Type:        tmpl.Type,
				Name:        tmpl.Name,
				Level:       1,
				MaxLevel:    tmpl.MaxLevel,
				Description: tmpl.Description,
			})
			if err != nil {
				fmt.Printf("新增設施失敗：%v\n", err)
				continue
			}
			fmt.Printf("設施「%s」已安裝！\n", tmpl.Name)
			renderBase(base)
		case 4:
			base, err := defenseSvc.GetBase(ctx, baseID)
			if err != nil {
				fmt.Printf("查詢失敗：%v\n", err)
				continue
			}
			if len(base.Facilities) == 0 {
				fmt.Println("基地尚無設施。")
				continue
			}
			fmt.Println("\n設施列表：")
			for i, f := range base.Facilities {
				fmt.Printf("%d) %s [%s] 等級:%d/%d\n", i+1, f.Name, f.Type, f.Level, f.MaxLevel)
			}
			idx := promptIndex("選擇升級設施", len(base.Facilities))
			updated, err := defenseSvc.UpgradeFacility(ctx, baseID, base.Facilities[idx-1].ID)
			if err != nil {
				fmt.Printf("升級失敗：%v\n", err)
				continue
			}
			fmt.Println("設施升級成功！")
			renderBase(updated)
		case 5:
			base, err := defenseSvc.UpgradeSecurity(ctx, baseID, 10)
			if err != nil {
				fmt.Printf("升級失敗：%v\n", err)
				continue
			}
			fmt.Println("安全等級提升！")
			renderBase(base)
		}
	}
}

func renderBase(base *dto.BaseSummary) {
	if base == nil {
		fmt.Println("基地資料為空。")
		return
	}
	fmt.Printf("\n  基地ID:    %s\n", base.ID)
	fmt.Printf("  擁有者:    %s\n", base.OwnerID)
	fmt.Printf("  安全等級:  %d\n", base.SecurityLevel)
	fmt.Printf("  設施欄位:  %d/%d\n", len(base.Facilities), base.FacilitySlots)
	if len(base.Facilities) > 0 {
		fmt.Println("  設施清單:")
		for _, f := range base.Facilities {
			fmt.Printf("    - %s [%s] 等級:%d/%d\n", f.Name, f.Type, f.Level, f.MaxLevel)
		}
	}
}

type facilityTemplate struct {
	Type        string
	Name        string
	MaxLevel    int
	Description string
}

func buildFacilityTemplates() []facilityTemplate {
	return []facilityTemplate{
		{Type: "Firewall", Name: "防火牆", MaxLevel: 5, Description: "阻擋外部入侵與惡意連線"},
		{Type: "SIEM", Name: "資安事件監控中心", MaxLevel: 5, Description: "即時偵測並分析威脅情報"},
		{Type: "Training", Name: "人員訓練室", MaxLevel: 3, Description: "提升團隊防詐意識與應變能力"},
	}
}

// ---------------------------------------------------------------------------
// Victim Profile Strategy
// ---------------------------------------------------------------------------

func renderVictimStrategy(m *dto.MissionDetail) {
	if m == nil || m.VictimProfile == nil {
		return
	}
	p := m.VictimProfile

	fmt.Println("\n--- 應對策略建議 ---")

	strategies := make([]string, 0, 4)

	if p.Anxiety >= 60 {
		strategies = append(strategies, "• 焦慮偏高：先穩定對方情緒，避免直接質問，以同理心引導。")
	}
	if p.Trust >= 60 {
		strategies = append(strategies, "• 信任偏高：對方容易輕信他人，需用具體數據與證據打破迷思。")
	}
	if p.Urgency >= 60 {
		strategies = append(strategies, "• 急迫偏高：對方傾向衝動行事，應先拖延時間以創造冷靜空間。")
	}
	if p.Isolation >= 60 {
		strategies = append(strategies, "• 孤立偏高：對方缺乏社會支持，建議引入親友或專業人員協助。")
	}

	if len(strategies) == 0 {
		fmt.Println("  此案件受害者心理狀態相對穩定，依正常流程調查即可。")
	} else {
		for _, s := range strategies {
			fmt.Println("  " + s)
		}
	}

	fmt.Printf("  綜合風險等級：%s（分數：%d）\n", p.RiskLevel, p.RiskScore)
}

// ---------------------------------------------------------------------------
// Player Status
// ---------------------------------------------------------------------------

func renderPlayerStatus(ctx context.Context, playersRepo *player.MongoRepository, playerID string) {
	p, err := playersRepo.FindByID(ctx, playerID)
	if err != nil || p == nil {
		fmt.Println("無法載入角色資料。")
		return
	}

	totalStats := p.GetTotalStats()
	fmt.Println("\n==================== 角色狀態 ====================")
	fmt.Printf("  玩家ID:    %s\n", p.ID)
	fmt.Printf("  聲望:      %d\n", p.Reputation)
	fmt.Printf("  基礎數值:  邏輯:%d 技術:%d 交涉:%d 韌性:%d\n",
		p.Stats.Logic, p.Stats.Tech, p.Stats.Charisma, p.Stats.Resilience)
	fmt.Printf("  總合數值:  邏輯:%d 技術:%d 交涉:%d 韌性:%d\n",
		totalStats.Logic, totalStats.Tech, totalStats.Charisma, totalStats.Resilience)

	if p.Partner != nil {
		fmt.Printf("  AI 性格:   %s\n", p.Partner.Personality)
		fmt.Printf("  已裝模組:  %d 個\n", len(p.Partner.Loadout))
		for _, mod := range p.Partner.Loadout {
			fmt.Printf("    - %s [%s]\n", mod.Name, mod.Type)
		}
		fmt.Printf("  已學技能:  %d 個\n", len(p.Partner.Skills))
		for _, skill := range p.Partner.Skills {
			cd := p.Partner.CooldownRemaining(skill.ID)
			cdStr := "就緒"
			if cd > 0 {
				cdStr = fmt.Sprintf("冷卻:%d", cd)
			}
			fmt.Printf("    - %s [%s] %s\n", skill.Name, skill.Type, cdStr)
		}
	}
}

// ---------------------------------------------------------------------------
// Skill Tree (existing, refined)
// ---------------------------------------------------------------------------

func runSkillTreeMenu(ctx context.Context, svc *usecasepersonnel.Service, playerID string) {
	for {
		skills, err := svc.ListSkills(ctx, playerID)
		if err != nil {
			fmt.Printf("技能樹載入失敗：%v\n", err)
			return
		}
		renderSkillTree(skills)

		fmt.Println("1) 解鎖技能  2) 裝備技能  3) 啟動技能  4) 返回")
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
		{
			ID:                 "skill-forensics-1",
			Type:               domainpersonnel.SkillTypeForensics,
			Name:               "數位鑑識",
			Description:        "深度取證分析，揭露隱藏的數位足跡。",
			CooldownSeconds:    18,
			ReputationRequired: 15,
		},
	}
}

// ---------------------------------------------------------------------------
// Mission Selection & Rendering
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

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
