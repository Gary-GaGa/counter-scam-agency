# Roadmap

## Phase 0 — Foundations
- [x] 建立公開 repo，設定 README / LICENSE / .gitignore
- [x] 撰寫遊戲企劃草案（`docs/game-concept.md`）
- [x] 定義領域模型（Mission, Evidence, ScamType, Player, AIPartner, Investigation）
- [x] 規劃 DDD/Clean Architecture 分層與 Repository 介面
- [ ] 定義文字冒險節點資料結構（節點、選項、結果）
- [ ] 定義聲望累計加權規則（失敗不給聲望）
- [ ] 草擬文字版 UI Flow（案件閱讀、選項、結算、成長）

## Phase 1 — MVP Prototype
- [ ] 實作文字冒險原型：兩種案件（假客服、投資社群）
- [ ] 節點流程：讀取任務 → 選項 → 蒐證 → 處置 → 結案
- [ ] 聲望累計加權計算（成功才加分）
- [ ] 模組解鎖規則：累計聲望門檻
- [ ] 建立任務與圖鑑資料基礎（MongoDB 儲存）
- [ ] 撰寫自動化測試（Domain/Usecase）

## Phase 2 — Visual & UX Upgrade
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