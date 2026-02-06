package main

import (
	"context"
	"fmt"
	"os"

	"counter-scam-agency/internal/domain/intelligence"
	mongoinfra "counter-scam-agency/internal/infrastructure/persistence/mongo"
	"counter-scam-agency/internal/interface/out/persistence/mongo/mission"
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

	repo := mission.NewMongoRepository(db)
	missions := seedMissions()
	for _, m := range missions {
		if err := repo.Save(ctx, m); err != nil {
			panic(err)
		}
	}

	fmt.Printf("seeded %d missions\n", len(missions))
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
