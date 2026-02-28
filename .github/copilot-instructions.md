# Copilot Instructions

## 語言

文件、commit message 與程式碼註釋使用**繁體中文**。Go 識別符號與 CLI flag 使用英文。

## 建置與測試

```bash
# 執行所有測試
go test ./...

# 執行單一套件測試
go test ./internal/domain/personnel/...

# 執行指定測試
go test ./internal/domain/personnel/... -run TestPlayerStatsCalculation

# 啟動 API 伺服器（需先啟動 MongoDB）
go run cmd/api/main.go

# 啟動完整 CLI（調查、技能樹、模組、基地）
go run cmd/cli/main.go

# 啟動精簡冒險模式
go run cmd/adventure/main.go

# 匯入種子任務資料
go run cmd/seed-missions/main.go
```

測試使用 `github.com/stretchr/testify`。Domain 與 Usecase 測試使用記憶體內 mock repository，不需 MongoDB。

## 遊戲核心循環（MVP）

遊戲以「案件調查 → 角色成長 → 解鎖更難案件」為核心循環：

1. **案件調查**：接收案件 → 文字冒險節點選擇 → 蒐證 → 處置 → 結案
2. **聲望累計**：成功結案才累計聲望（失敗不給），聲望加權由任務 `ReputationWeight` 決定
3. **技能樹養成**：四大技能類型（Analysis / Negotiation / Defense / Forensics），需聲望門檻解鎖 → 裝備 → 啟動（有冷卻）
4. **AI 模組裝備**：四種模組（VoiceAnalyzer / CryptoTracer / EmpathyEngine / MentalFirewall），各加成不同屬性
5. **數位防禦基地**：建設設施（Firewall / SIEM / Training）、升級安全等級
6. **受害者心理側寫**：焦慮、信任、急迫、孤立四維度分析，風險分數影響應對策略

### MVP 包含的案件

- **假客服來電**（Phishing, 難度 2）
- **投資社群圈套**（Investment, 難度 3）

### 角色基礎數值

四項數值影響小遊戲與案件表現：邏輯 (Logic)、技術 (Tech)、交涉 (Charisma)、韌性 (Resilience)。基礎值皆為 10，模組裝備可加成。

## 架構

**Clean Architecture + DDD**，四個限界上下文：

```
domain（純商業邏輯，零外部依賴）
  ↑
usecase（應用服務，跨上下文協調）
  ↑
interface/in（HTTP handlers）   interface/out（MongoDB repos）
  ↑                               ↑
infrastructure（DB client 設定）
```

**依賴規則**：import 只能向內。Domain 不 import usecase，usecase 不 import interface/infrastructure。

### 限界上下文

| 上下文 | 套件 | 聚合根 | 職責 |
|--------|------|--------|------|
| Personnel | `internal/domain/personnel` | `Player` | 玩家屬性、AI 夥伴、技能、模組 |
| Intelligence | `internal/domain/intelligence` | `Mission` | 案件敘事、證據、受害者側寫、詐騙類型 |
| Operation | `internal/domain/operation` | `Investigation` | 調查執行狀態、節點決策 |
| Defense | `internal/domain/defense` | `Base` | 防禦設施、升級 |

### 各層職責

- **Domain**：Rich domain model（非貧血模型）。Repository 介面定義在此層。
- **Usecase**：Input port 介面在 `usecase/port/in/`，DTO 在 `usecase/dto/`。Service 透過 constructor injection 接收 repository。
- **Interface In**：HTTP handler 使用 Go 1.22+ `http.ServeMux` routing，依賴 input port 介面。
- **Interface Out**：MongoDB repository 實作。Persistence Object (PO) 在 `po/` 下，雙向轉換器在 `po/convert/`。
- **Infrastructure**：共用 MongoDB client 設定。

### 跨上下文協調

Usecase 層負責協調多個上下文。例如 `InvestigationService` 讀取 Intelligence（任務），變更 Operation（調查），更新 Personnel（玩家聲望）。Domain 上下文之間不直接 import。

## REST API 端點

| 方法 | 路徑 | 說明 |
|------|------|------|
| GET | `/api/missions` | 案件列表 |
| GET | `/api/missions/{id}` | 案件詳情（含節點、證據、受害者側寫） |
| POST | `/api/investigations` | 開始調查 |
| POST | `/api/investigations/{id}/advance` | 推進節點 |
| POST | `/api/investigations/{id}/evidence` | 提交證據 |
| POST | `/api/investigations/{id}/complete` | 結案 |
| POST | `/api/players` | 建立玩家 |
| GET | `/api/players/{id}` | 玩家狀態 |
| GET | `/api/players/{id}/skills` | 技能樹（含解鎖/裝備/冷卻狀態） |
| POST | `/api/players/{id}/skills/{skillID}/unlock` | 解鎖技能 |
| POST | `/api/players/{id}/skills/{skillID}/equip` | 裝備技能 |
| POST | `/api/players/{id}/skills/{skillID}/activate` | 啟動技能 |
| POST | `/api/players/{id}/skills/tick` | 推進冷卻時間 |
| POST | `/api/bases` | 建立基地 |
| GET | `/api/bases/{id}` | 基地狀態 |
| POST | `/api/bases/{id}/facilities` | 新增設施 |
| POST | `/api/bases/{id}/security/upgrade` | 升級安全等級 |
| POST | `/api/bases/{id}/facilities/{facilityID}/upgrade` | 升級設施 |

## 關鍵慣例

### Constructor Injection

所有 Service 與 Repository 使用 `NewService(...)`、`NewMongoRepository(...)` 建構函式注入依賴。無 service locator 或全域變數。

### 編譯期介面檢查

每個 repository 與 service 實作都包含：

```go
var _ personnel.PlayerRepository = (*MongoRepository)(nil)
```

### 錯誤處理

- 使用 `fmt.Errorf("description: %w", err)` 包裝錯誤
- Domain 層定義 sentinel error
- **HTTP handler 不得將 `err.Error()` 回傳給客戶端**——使用 `internalError()` 記錄伺服器端日誌並回傳通用訊息

### Persistence Object (PO) 模式

Domain entity 不直接存入 MongoDB。使用專用 PO struct 與轉換器：
- `po/convert/po/` — Domain → PO（寫入）
- `po/convert/model/` — PO → Domain（讀取）

### 測試模式

- **Domain 測試**：直接建構 entity 並斷言
- **Usecase 測試**：實作 domain repository 介面的記憶體 struct（不使用 mock 框架）
- **HTTP 測試**：透過 stub usecase 測試 handler 行為

### 命名慣例

- Domain entity：`Player`, `Mission`, `Investigation`, `Base`
- DTO：`PlayerSummary`, `MissionDetail`, `InvestigationStartResult`
- Persistence object：`PlayerDocPo`, `MissionDocPo`
- Service：各 usecase 子套件內命名為 `Service`（如 `investigation.Service`）

## 安全性規範

### HTTP 安全

- 所有回應加上安全標頭（`X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`）
- Request body 限制 1 MB（`http.MaxBytesReader`）
- JSON decoder 啟用 `DisallowUnknownFields()`
- Path parameter 驗證：僅允許 `[a-zA-Z0-9_-]{1,128}`
- CORS origin 透過 `CORS_ORIGIN` 環境變數設定（開發用 `*`，正式環境指定 origin）
- HTTP server 設定完整 timeout（ReadTimeout / WriteTimeout / IdleTimeout）

### 待實作安全項目

- [ ] 認證與授權中介層（JWT 或 session-based）
- [ ] 請求速率限制
- [ ] TLS 或反向代理文件化

## 環境變數

| 變數 | 預設值 | 說明 |
|------|--------|------|
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB 連線字串 |
| `MONGO_DB` | `counter_scam_agency` | 資料庫名稱 |
| `ADDR` | `:8080` | HTTP 伺服器監聽位址 |
| `CORS_ORIGIN` | `*` | 允許的 CORS origin |
| `PLAYER_ID` | `player-1` | CLI 模式預設玩家 ID |

## 前端（Phaser.js Web UI）

```bash
# 安裝依賴
cd web && npm install

# 開發模式（含 API proxy 到 localhost:8080）
npm run dev

# 正式建置
npm run build

# TypeScript 型別檢查
npx tsc --noEmit
```

### 前端架構

- **Vite + TypeScript + Phaser 3**：像素風格 2D 遊戲引擎
- 所有 API 呼叫透過 `src/services/api.ts`，Vite dev server 自動 proxy `/api` 到 Go 後端
- 場景（Scene）為 Phaser 的頁面單位，放在 `src/scenes/`
- 共用 UI 元件（按鈕、面板、導航列）在 `src/ui/`
- DTO 型別定義在 `src/services/types.ts`，對應 Go 後端 `usecase/dto`

### 場景列表

| 場景 | 檔案 | 功能 |
|------|------|------|
| MainMenuScene | 主選單 | 遊戲入口，導航至各功能 |
| CaseListScene | 案件列表 | 任務卡片選擇 |
| InvestigationScene | 調查場景 | 文字冒險核心玩法 |
| SkillTreeScene | 技能樹 | 技能視覺化 + 解鎖/裝備/啟動 |
| ProfileScene | 角色狀態 | 屬性值 + 聲望 + AI 狀態 |
