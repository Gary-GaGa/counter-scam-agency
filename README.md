# Counter Scam Agency

Counter Scam Agency 是一款以防範詐騙為主題的角色扮演遊戲（RPG）。玩家扮演「反詐情報官」，以文字冒險節點推進案件，協助市民解析詐術、阻止匯款、追查詐團，並透過累計聲望解鎖調查工具與模組。

## 目標
- 以遊戲化方式傳遞防詐知識。
- 保持劇情與案件推理的沉浸感。
- 打造可延伸的系統架構，方便快速迭代內容。

## 特色
- **案件調查循環**：節點敘事 → 選項 → 蒐集資訊 → 判讀 → 處置 → 回饋。
- **技能與工具養成**：社交洞察、金融偵測、數據取證、心理韌性等技能樹。
- **AI 夥伴系統**：模組化裝備、主動技能與性格演化。
- **數位防禦基地**：建設防火牆、SIEM、訓練室等設施，提升資安能力。
- **受害者心理側寫**：根據焦慮、信任、急迫、孤立等維度分析風險，提供應對策略。
- **詐術圖鑑**：記錄已破解的詐騙類型與關鍵徵兆，作為教育資源。
- **聲望規則**：任務成功才累計聲望，失敗不給聲望。

## 技術棧
| 層級 | 技術 | 說明 |
|------|------|------|
| 後端 | Go (Clean Architecture + DDD) | 領域邏輯、Usecase、REST API |
| 資料庫 | MongoDB | 任務、調查、玩家、基地資料持久化 |
| 前端 | Phaser 3 + TypeScript | 像素風格 2D 遊戲畫面（網頁版） |
| 部署 | Web（瀏覽器） | 開瀏覽器即可遊玩，零安裝 |

### 前後端職責分工（混合模式）
- **Go 後端**：核心判定邏輯（任務推進、聲望計算、技能冷卻、基地規則），透過 REST API 提供服務。
- **Phaser 前端**：像素風格場景渲染、對話框 UI、小遊戲執行（彈幕/卡牌/節奏等），透過 HTTP 與後端溝通。

## 專案結構
```
cmd/
  adventure/        # 文字冒險 CLI 入口
  api/              # HTTP API 伺服器入口
  cli/              # 完整功能 CLI 入口
  seed-missions/    # 種子資料工具

internal/
  domain/           # 領域層：純 Go 商業邏輯
  usecase/          # 應用層：用例服務、Port、DTO
  interface/        # 轉接層：HTTP Handler、Mongo Repository
  infrastructure/   # 基礎設施：MongoDB Client

web/                # Phaser.js 前端（Vite + TypeScript）
  src/
    scenes/         # 遊戲場景（主選單、案件、調查、技能樹、狀態）
    services/       # API 客戶端與型別定義
    ui/             # 共用 UI 元件

docs/
  architecture.md   # 系統架構文件
  game-concept.md   # 遊戲企劃草案
  roadmap.md        # 開發路線圖
  ui-flow.md        # UI 流程規劃
```

## 快速開始
```bash
# 啟動 MongoDB
mongod --dbpath /path/to/data

# 匯入種子任務
go run cmd/seed-missions/main.go

# 啟動 API 伺服器
go run cmd/api/main.go

# 啟動 Web 前端（另開終端）
cd web && npm run dev
# 開啟瀏覽器 http://localhost:5173

# 或使用完整 CLI（含基地、技能、模組管理）
go run cmd/cli/main.go

# 或啟動精簡冒險模式
go run cmd/adventure/main.go
```

## 下一步
1. ~~實作 Go HTTP API，將現有 Usecase 包成 REST 端點。~~ ✅
2. ~~建立 Phaser.js 前端專案，像素風格場景與節點對話 UI。~~ ✅
3. ~~串接前後端，完成網頁版可遊玩原型。~~ ✅
4. 加入認證與授權機制。
5. 實作小遊戲（矛盾擊破、訊號追蹤、談判牌局、心靈調適）。
6. 完善像素風格美術資源。

如欲參與或提出建議，歡迎提交 Issue 或 Pull Request。
