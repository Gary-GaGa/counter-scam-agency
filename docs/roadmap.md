# Roadmap

## Phase 0 — Foundations
- [x] 建立公開 repo，設定 README / LICENSE / .gitignore
- [x] 撰寫遊戲企劃草案（`docs/game-concept.md`）
- [x] 定義領域模型（Mission, Evidence, ScamType, Player, AIPartner, Investigation）
- [x] 規劃 DDD/Clean Architecture 分層與 Repository 介面
- [x] 定義文字冒險節點資料結構（節點、選項、結果）
- [x] 定義聲望累計加權規則（失敗不給聲望）
- [x] 草擬文字版 UI Flow（案件閱讀、選項、結算、成長，見 `docs/ui-flow.md`）

## Phase 1 — MVP Prototype
- [x] **CA + DDD 核心流程落地**（近期 Sprint）
    - [x] Domain：定義文字冒險節點模型（Node/Option/Outcome）與聲望規則（Value Object）
    - [x] Domain：Investigation 聚合內實作節點推進、蒐證、結案的狀態轉移
    - [x] Usecase：建立 `StartInvestigation` / `AdvanceNode` / `SubmitEvidence` / `CompleteInvestigation`
    - [x] Usecase：定義 Input/Output Port 與 DTO（避免 Domain 直出 UI 結構）
    - [x] Interface：CLI/遊戲流程 Adapter（讀取任務、輸入選項、輸出結算）
    - [x] Infrastructure：Mongo Repository 串接（Mission/Investigation/Player）
    - [x] 測試：Domain 狀態轉移與 Usecase 流程測試
- [x] 實作文字冒險原型：兩種案件（假客服、投資社群）
- [x] 節點流程：讀取任務 → 選項 → 蒐證 → 處置 → 結案
- [x] 聲望累計加權計算（成功才加分）
- [x] 模組解鎖規則：累計聲望門檻
- [x] 建立任務與圖鑑資料基礎（MongoDB 儲存）
- [x] 撰寫自動化測試（Domain/Usecase）

## Phase 2 — Visual & UX Upgrade
- [x] Domain：AI 技能模型（Skill）
- [x] Domain：受害者心理側寫（VictimProfile）
- [x] Domain：數位防禦基地（Base / Facility）
- [x] Usecase：基地管理流程（建立 / 升級 / 設施）
- [x] Infrastructure：基地 Mongo Repository
- [x] Seed：任務心理側寫資料
- [x] Usecase：AI 技能流程（解鎖 / 裝備 / 啟動 / 冷卻）
- [ ] Ebiten UI：案件面板、推理棋盤、技能樹介面
- [ ] **AI 模組與技能系統**：裝備介面與主動技能實作
- [ ] 直播宣導 / 情報擴散玩法
- [ ] **數位防禦基地**：基地建設與資安升級系統
- [ ] **小遊戲開發**：
    - [ ] 矛盾擊破 (Logic) - 彈幕射擊玩法
    - [ ] 訊號追蹤 (Tech) - 接水管/駭客玩法
    - [ ] 談判牌局 (Charisma) - 卡牌對戰玩法
    - [ ] 心靈調適 (Resilience) - 節奏/休閒玩法
- [ ] **受害者心理側寫**：心理狀態分析與應對機制

## Phase 3 — Content Expansion
- [ ] 新案件章節：情感攻防、高科技對決、跨國詐團
- [ ] **AI 性格演化系統**：性格參數影響回饋與風格描述
- [ ] 多語言在地化、教育資源連結
- [ ] 評測與教案整合（可選）
- [ ] **反向腳本模擬**：詐騙話術組裝沙盒
- [ ] **社群共創系統**：玩家投稿與案件審核機制

> 本文件會隨進度調整，可透過 Issue / PR 更新。