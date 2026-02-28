package main

import (
	"context"
	"fmt"
	"os"

	"counter-scam-agency/internal/domain/defense"
	"counter-scam-agency/internal/domain/intelligence"
	"counter-scam-agency/internal/domain/personnel"
	mongoinfra "counter-scam-agency/internal/infrastructure/persistence/mongo"
	defenserepo "counter-scam-agency/internal/interface/out/persistence/mongo/defense"
	"counter-scam-agency/internal/interface/out/persistence/mongo/mission"
	playerrepo "counter-scam-agency/internal/interface/out/persistence/mongo/player"
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

	repo := mission.NewMongoRepository(db)
	missions := seedMissions()
	for _, m := range missions {
		if err := repo.Save(ctx, m); err != nil {
			fmt.Fprintf(os.Stderr, "seed mission: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("seeded %d missions\n", len(missions))

	// 建立預設玩家
	playerRepo := playerrepo.NewMongoRepository(db)
	existing, _ := playerRepo.FindByID(ctx, "player-1")
	if existing == nil {
		player := personnel.NewPlayer("player-1")
		if err := playerRepo.Save(ctx, player); err != nil {
			fmt.Fprintf(os.Stderr, "seed player: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("seeded player: player-1")
	} else {
		fmt.Println("player-1 already exists, skipped")
	}

	// 建立預設基地
	baseRepo := defenserepo.NewMongoRepository(db)
	existingBase, _ := baseRepo.FindByID(ctx, "base-1")
	if existingBase == nil {
		base := defense.NewBase("base-1", "player-1", 4)
		base.AddFacility(defense.Facility{
			ID:          "facility-firewall",
			Type:        defense.FacilityTypeFirewall,
			Name:        "基礎防火牆",
			Level:       1,
			MaxLevel:    5,
			Description: "阻擋基本惡意流量",
		})
		if err := baseRepo.Save(ctx, base); err != nil {
			fmt.Fprintf(os.Stderr, "seed base: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("seeded base: base-1")
	} else {
		fmt.Println("base-1 already exists, skipped")
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func seedMissions() []*intelligence.Mission {
	return []*intelligence.Mission{
		seedFakeSupportMission(),
		seedInvestmentGroupMission(),
	}
}

func seedFakeSupportMission() *intelligence.Mission {
	mission := intelligence.NewMission(
		"mission-fake-support",
		"假客服來電",
		"冒充官方客服引導受害者匯款。",
		intelligence.ScamTypePhishing,
		2,
		2,
	)
	mission.VictimProfile = &intelligence.VictimProfile{
		Anxiety:   70,
		Trust:     60,
		Urgency:   65,
		Isolation: 40,
	}

	mission.AddEvidence(*intelligence.NewEvidence("ev-support-1", "對方要求提供驗證碼", intelligence.EvidenceTypeDialogue, true))
	mission.AddEvidence(*intelligence.NewEvidence("ev-support-2", "來電號碼與官方不符", intelligence.EvidenceTypeDocument, true))

	mission.AddNode(intelligence.NarrativeNode{
		ID:    "start",
		Title: "陌生客服來電",
		Body:  "對方自稱客服，表示帳戶異常並要求配合處理。",
		Options: []intelligence.NarrativeOption{
			{
				ID:          "opt-verify",
				Label:       "要求對方提供官方管道並自行查證",
				NextNodeID:  "verify",
				EvidenceIDs: []string{"ev-support-2"},
			},
			{
				ID:         "opt-transfer",
				Label:      "依照指示立即轉帳",
				NextNodeID: "end-fail",
				LeadsToEnd: true,
				SuccessEnd: false,
			},
		},
	})

	mission.AddNode(intelligence.NarrativeNode{
		ID:    "verify",
		Title: "查證號碼",
		Body:  "你改用官方網站聯絡客服，確認對方為詐騙。",
		Options: []intelligence.NarrativeOption{
			{
				ID:          "opt-report",
				Label:       "回報並收集證據",
				NextNodeID:  "end-success",
				EvidenceIDs: []string{"ev-support-1"},
				LeadsToEnd:  true,
				SuccessEnd:  true,
			},
		},
	})

	mission.AddNode(intelligence.NarrativeNode{
		ID:         "end-success",
		Title:      "成功阻止",
		Body:       "你成功識破假客服，阻止匯款並收集證據。",
		IsTerminal: true,
	})

	mission.AddNode(intelligence.NarrativeNode{
		ID:         "end-fail",
		Title:      "失敗結案",
		Body:       "你依照指示轉帳，損失擴大。",
		IsTerminal: true,
	})

	return mission
}

func seedInvestmentGroupMission() *intelligence.Mission {
	mission := intelligence.NewMission(
		"mission-investment-group",
		"投資社群圈套",
		"投資群組以高報酬話術誘導加入。",
		intelligence.ScamTypeInvestment,
		3,
		2,
	)
	mission.VictimProfile = &intelligence.VictimProfile{
		Anxiety:   55,
		Trust:     80,
		Urgency:   75,
		Isolation: 60,
	}

	mission.AddEvidence(*intelligence.NewEvidence("ev-invest-1", "群組要求先付入會費", intelligence.EvidenceTypeTransaction, true))
	mission.AddEvidence(*intelligence.NewEvidence("ev-invest-2", "收益截圖無法驗證", intelligence.EvidenceTypeImage, true))

	mission.AddNode(intelligence.NarrativeNode{
		ID:    "start",
		Title: "加入群組",
		Body:  "群組內成員分享獲利截圖，邀你入會。",
		Options: []intelligence.NarrativeOption{
			{
				ID:          "opt-question",
				Label:       "詢問報酬來源並要求公開證明",
				NextNodeID:  "probe",
				EvidenceIDs: []string{"ev-invest-2"},
			},
			{
				ID:         "opt-pay",
				Label:      "先支付入會費",
				NextNodeID: "end-fail",
				LeadsToEnd: true,
				SuccessEnd: false,
			},
		},
	})

	mission.AddNode(intelligence.NarrativeNode{
		ID:    "probe",
		Title: "深入追問",
		Body:  "對方開始迴避問題，並要求你私下轉帳。",
		Options: []intelligence.NarrativeOption{
			{
				ID:          "opt-collect",
				Label:       "蒐集證據並退出群組",
				NextNodeID:  "end-success",
				EvidenceIDs: []string{"ev-invest-1"},
				LeadsToEnd:  true,
				SuccessEnd:  true,
			},
		},
	})

	mission.AddNode(intelligence.NarrativeNode{
		ID:         "end-success",
		Title:      "成功阻止",
		Body:       "你保留證據並成功避免損失。",
		IsTerminal: true,
	})

	mission.AddNode(intelligence.NarrativeNode{
		ID:         "end-fail",
		Title:      "失敗結案",
		Body:       "你支付入會費後被踢出群組。",
		IsTerminal: true,
	})

	return mission
}
