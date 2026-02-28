/**
 * Mock 資料 — 讓遊戲在沒有 Go 後端時也能完整遊玩。
 * 資料結構對應 types.ts 的 DTO 型別。
 */

import type {
  MissionSummary,
  MissionDetail,
  PlayerSummary,
  SkillSummary,
  BaseSummary,
} from './types';

// ─── 案件資料 ───

export const MOCK_MISSIONS: MissionSummary[] = [
  {
    id: 'mission-phishing',
    title: '假客服來電',
    description: '一名市民接到自稱電信客服的電話，聲稱帳號異常需要驗證。你必須協助辨識詐騙手法並保護受害者。',
    type: 'Phishing',
    difficulty: 2,
    reputationWeight: 1.0,
  },
  {
    id: 'mission-investment',
    title: '投資社群圈套',
    description: '社群平台出現高報酬投資群組，已有多名受害者匯款。深入調查這個詐騙集團的運作模式。',
    type: 'Investment',
    difficulty: 3,
    reputationWeight: 1.5,
  },
];

export const MOCK_MISSION_DETAILS: Record<string, MissionDetail> = {
  'mission-phishing': {
    id: 'mission-phishing',
    title: '假客服來電',
    description: '一名市民接到自稱電信客服的電話，聲稱帳號異常需要驗證。你必須協助辨識詐騙手法並保護受害者。',
    type: 'Phishing',
    difficulty: 2,
    reputationWeight: 1.0,
    victimProfile: {
      anxiety: 72,
      trust: 65,
      urgency: 80,
      isolation: 40,
      riskScore: 64,
      riskLevel: '高風險',
    },
    evidenceList: [
      { id: 'ev-caller-id', description: '來電顯示為非官方號碼', type: 'Digital', isKey: true },
      { id: 'ev-pressure', description: '對方施壓要求立即操作', type: 'Behavioral', isKey: true },
      { id: 'ev-ask-code', description: '要求提供簡訊驗證碼', type: 'Behavioral', isKey: true },
      { id: 'ev-official-check', description: '官方客服確認無此通話記錄', type: 'Document', isKey: false },
    ],
    nodes: [
      {
        id: 'node-1',
        title: '接獲通報',
        body: '情報中心收到一則緊急通報：市民王小姐接到自稱「中華電信客服」的電話，對方聲稱她的帳號遭到異常登入，需要立刻驗證身份。\n\n王小姐目前在線上等待你的指示。你會如何開始這次調查？',
        isTerminal: false,
        options: [
          { id: 'opt-1a', label: '先確認來電號碼是否為官方客服', nextNodeId: 'node-2a', evidenceIds: ['ev-caller-id'], leadsToEnd: false, successEnd: false },
          { id: 'opt-1b', label: '請王小姐描述對方說了什麼', nextNodeId: 'node-2b', evidenceIds: [], leadsToEnd: false, successEnd: false },
          { id: 'opt-1c', label: '直接建議她掛斷電話', nextNodeId: 'node-2c', evidenceIds: [], leadsToEnd: false, successEnd: false },
        ],
      },
      {
        id: 'node-2a',
        title: '號碼查證',
        body: '你查驗來電號碼，發現這是一個網路電話號碼，並非中華電信官方客服（0800-080-123）。\n\n這是一個重要的線索！正規電信公司有固定客服號碼。\n\n🔍 取得證據：「來電顯示為非官方號碼」',
        isTerminal: false,
        options: [
          { id: 'opt-2a-1', label: '繼續追問對方要求什麼操作', nextNodeId: 'node-3', evidenceIds: ['ev-pressure'], leadsToEnd: false, successEnd: false },
          { id: 'opt-2a-2', label: '立即通報 165 反詐騙專線', nextNodeId: 'node-end-good', evidenceIds: [], leadsToEnd: true, successEnd: true },
        ],
      },
      {
        id: 'node-2b',
        title: '蒐集情報',
        body: '王小姐表示，對方語氣非常急促，說：「您的帳號在 5 分鐘前從境外 IP 登入，如果不立刻驗證，帳號將被永久凍結。」\n\n對方接著要求她提供最近收到的簡訊驗證碼。\n\n💡 這種施壓手法是典型的詐騙話術特徵。',
        isTerminal: false,
        options: [
          { id: 'opt-2b-1', label: '確認來電號碼', nextNodeId: 'node-2a', evidenceIds: ['ev-caller-id'], leadsToEnd: false, successEnd: false },
          { id: 'opt-2b-2', label: '記錄施壓行為作為證據', nextNodeId: 'node-3', evidenceIds: ['ev-pressure', 'ev-ask-code'], leadsToEnd: false, successEnd: false },
        ],
      },
      {
        id: 'node-2c',
        title: '快速處置',
        body: '你建議王小姐立刻掛斷電話。她照做了，但你缺少了蒐集更多證據的機會。\n\n雖然受害者暫時安全了，但沒有足夠的情報來追蹤詐騙集團。',
        isTerminal: false,
        options: [
          { id: 'opt-2c-1', label: '結案報告（證據不足）', nextNodeId: 'node-end-partial', evidenceIds: [], leadsToEnd: true, successEnd: true },
        ],
      },
      {
        id: 'node-3',
        title: '深入調查',
        body: '你已掌握多項關鍵證據：\n\n📋 非官方來電號碼\n📋 施壓手法（限時威脅）\n📋 要求提供驗證碼\n\n這些都是典型的電信詐騙手法。你決定如何結案？',
        isTerminal: false,
        options: [
          { id: 'opt-3a', label: '通報 165 並協助王小姐凍結帳戶', nextNodeId: 'node-end-good', evidenceIds: ['ev-official-check'], leadsToEnd: true, successEnd: true },
          { id: 'opt-3b', label: '建議王小姐自行撥打客服確認', nextNodeId: 'node-end-good', evidenceIds: [], leadsToEnd: true, successEnd: true },
        ],
      },
      {
        id: 'node-end-good',
        title: '案件結案',
        body: '🎉 調查完成！\n\n你成功識破了電信詐騙手法，保護了受害者的資產安全。\n所蒐集的證據已提交至反詐情報資料庫，將有助於未來的偵辦工作。\n\n📊 結案評估中...',
        isTerminal: true,
        options: [],
      },
      {
        id: 'node-end-partial',
        title: '案件結案',
        body: '📝 調查結束。\n\n受害者已脫離險境，但本次調查蒐集的證據有限。\n建議未來在類似案件中多花時間蒐證，以利追蹤詐騙集團。\n\n📊 結案評估中...',
        isTerminal: true,
        options: [],
      },
    ],
  },
  'mission-investment': {
    id: 'mission-investment',
    title: '投資社群圈套',
    description: '社群平台出現高報酬投資群組，已有多名受害者匯款。深入調查這個詐騙集團的運作模式。',
    type: 'Investment',
    difficulty: 3,
    reputationWeight: 1.5,
    victimProfile: {
      anxiety: 55,
      trust: 82,
      urgency: 68,
      isolation: 72,
      riskScore: 69,
      riskLevel: '高風險',
    },
    evidenceList: [
      { id: 'ev-fake-profit', description: '群組內展示的獲利截圖疑似造假', type: 'Digital', isKey: true },
      { id: 'ev-no-license', description: '該投資平台未持有金管會核發之營業執照', type: 'Document', isKey: true },
      { id: 'ev-pressure-sell', description: '使用「限時加碼」話術催促匯款', type: 'Behavioral', isKey: true },
      { id: 'ev-victim-testimony', description: '受害者 A 證實無法提領獲利', type: 'Testimony', isKey: false },
      { id: 'ev-crypto-trace', description: '資金流向追蹤至境外虛擬錢包', type: 'Digital', isKey: false },
    ],
    nodes: [
      {
        id: 'inv-node-1',
        title: '案件來源',
        body: '情報中心截獲一個名為「財富自由俱樂部」的 LINE 投資群組，宣稱跟著「老師」操作即可月獲利 30%。\n\n群組已有 500+ 成員，近期多名受害者向 165 通報無法提款。\n\n你打算如何展開調查？',
        isTerminal: false,
        options: [
          { id: 'inv-opt-1a', label: '派遣探員臥底加入群組', nextNodeId: 'inv-node-2a', evidenceIds: [], leadsToEnd: false, successEnd: false },
          { id: 'inv-opt-1b', label: '先查證該投資平台的合法性', nextNodeId: 'inv-node-2b', evidenceIds: ['ev-no-license'], leadsToEnd: false, successEnd: false },
          { id: 'inv-opt-1c', label: '聯繫已通報的受害者取得證詞', nextNodeId: 'inv-node-2c', evidenceIds: ['ev-victim-testimony'], leadsToEnd: false, successEnd: false },
        ],
      },
      {
        id: 'inv-node-2a',
        title: '臥底蒐證',
        body: '你派遣 AI 夥伴「小盾」以假身份加入群組。群組內充斥著獲利截圖和「老師」的分析報告。\n\n小盾發現：截圖中的帳戶餘額數字有明顯的影像編輯痕跡，且所有「獲利見證」都來自同一批帳號。\n\n🔍 取得證據：「群組內獲利截圖疑似造假」',
        isTerminal: false,
        options: [
          { id: 'inv-opt-2a-1', label: '進一步追蹤資金流向', nextNodeId: 'inv-node-3', evidenceIds: ['ev-fake-profit', 'ev-crypto-trace'], leadsToEnd: false, successEnd: false },
          { id: 'inv-opt-2a-2', label: '蒐集足夠證據，準備結案', nextNodeId: 'inv-node-3', evidenceIds: ['ev-fake-profit'], leadsToEnd: false, successEnd: false },
        ],
      },
      {
        id: 'inv-node-2b',
        title: '平台查證',
        body: '你查詢金管會的合法業者名單，確認「財富自由投資平台」並未取得任何金融營業執照。\n\n此外，該平台的網域註冊僅 3 個月，伺服器位於東南亞。\n\n🔍 取得證據：「未持有金管會核發之營業執照」\n\n💡 正規的投資平台必須取得金管會核准才能營運。',
        isTerminal: false,
        options: [
          { id: 'inv-opt-2b-1', label: '加入群組蒐集更多內部證據', nextNodeId: 'inv-node-2a', evidenceIds: [], leadsToEnd: false, successEnd: false },
          { id: 'inv-opt-2b-2', label: '查看群組的銷售話術', nextNodeId: 'inv-node-2d', evidenceIds: ['ev-pressure-sell'], leadsToEnd: false, successEnd: false },
        ],
      },
      {
        id: 'inv-node-2c',
        title: '受害者證詞',
        body: '受害者陳先生表示：「一開始投入 5 萬，帳面上確實看到 30% 獲利。但當我想提領時，客服說要再繳 20% 手續費才能出金。繳了之後又說要稅金…」\n\n🔍 取得證據：「受害者證實無法提領」\n\n這是典型的「出金障礙」詐騙模式。',
        isTerminal: false,
        options: [
          { id: 'inv-opt-2c-1', label: '同時從平台合法性著手', nextNodeId: 'inv-node-2b', evidenceIds: ['ev-no-license'], leadsToEnd: false, successEnd: false },
          { id: 'inv-opt-2c-2', label: '直接結案通報', nextNodeId: 'inv-node-end-partial', evidenceIds: [], leadsToEnd: true, successEnd: true },
        ],
      },
      {
        id: 'inv-node-2d',
        title: '話術分析',
        body: '你截取了群組內的話術範例：\n\n🔴「名額只剩最後 10 位，錯過再等一年！」\n🔴「老師今天的策略已經驗證 5 次全勝！」\n🔴「現在加碼 10 萬，週末前獲利翻倍保證！」\n\n這些都是典型的「限時加碼」施壓話術。\n\n🔍 取得證據：「使用限時加碼話術催促匯款」',
        isTerminal: false,
        options: [
          { id: 'inv-opt-2d-1', label: '證據充足，進行最終分析', nextNodeId: 'inv-node-3', evidenceIds: [], leadsToEnd: false, successEnd: false },
        ],
      },
      {
        id: 'inv-node-3',
        title: '最終分析',
        body: '你已蒐集到多項關鍵證據，可以對案件做出最終判斷。\n\n目前掌握的證據指向：這是一個有組織的投資詐騙集團，利用社群平台招攬受害者。\n\n你決定如何處置？',
        isTerminal: false,
        options: [
          { id: 'inv-opt-3a', label: '提交完整調查報告至檢調單位', nextNodeId: 'inv-node-end-good', evidenceIds: [], leadsToEnd: true, successEnd: true },
          { id: 'inv-opt-3b', label: '先在群組發布警示訊息', nextNodeId: 'inv-node-end-good', evidenceIds: [], leadsToEnd: true, successEnd: true },
        ],
      },
      {
        id: 'inv-node-end-good',
        title: '案件結案',
        body: '🎉 出色的調查工作！\n\n你的調查報告已提交至相關單位。根據你蒐集的證據：\n• 該投資平台已被列入金管會黑名單\n• 群組管理員的 IP 已移交執法單位追查\n• 已通知所有已知受害者申請資產凍結\n\n你的努力保護了更多潛在受害者。',
        isTerminal: true,
        options: [],
      },
      {
        id: 'inv-node-end-partial',
        title: '案件結案',
        body: '📝 調查結束。\n\n雖然你掌握了部分證據，但資訊尚不夠完整。\n建議在未來的調查中，從多個角度蒐集證據，以建構更完整的案件檔案。',
        isTerminal: true,
        options: [],
      },
    ],
  },
};

// ─── 玩家資料 ───

export const MOCK_PLAYER: PlayerSummary = {
  id: 'player-1',
  reputation: 0,
  stats: { logic: 10, tech: 10, charisma: 10, resilience: 10 },
  totalStats: { logic: 10, tech: 10, charisma: 10, resilience: 10 },
  partnerPersonality: 'Balanced',
  equippedModules: [],
  equippedSkills: [],
};

// ─── 技能資料 ───

export const MOCK_SKILLS: SkillSummary[] = [
  { id: 'skill-pattern-recognition', name: '模式識別', description: '快速辨識詐騙話術的共同模式', type: 'Analysis', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 10 },
  { id: 'skill-data-mining', name: '資料探勘', description: '從大量資訊中提取關鍵線索', type: 'Analysis', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 25 },
  { id: 'skill-active-listening', name: '積極傾聽', description: '透過傾聽建立信任，引導受害者配合', type: 'Negotiation', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 10 },
  { id: 'skill-crisis-mediation', name: '危機調解', description: '在高壓情境下穩定受害者情緒', type: 'Negotiation', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 30 },
  { id: 'skill-firewall-deploy', name: '防火牆部署', description: '為受害者裝置設定臨時防護', type: 'Defense', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 15 },
  { id: 'skill-trace-route', name: '路徑追蹤', description: '追蹤詐騙訊號來源的 IP 位址', type: 'Forensics', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 15 },
  { id: 'skill-crypto-analysis', name: '加密分析', description: '解析詐騙集團使用的加密通訊', type: 'Forensics', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 35 },
  { id: 'skill-emotional-shield', name: '情緒護盾', description: '協助受害者建立心理防線', type: 'Defense', unlocked: false, equipped: false, cooldownRemaining: 0, reputationRequired: 20 },
];

// ─── 基地資料 ───

export const MOCK_BASE: BaseSummary = {
  id: 'base-1',
  ownerId: 'player-1',
  securityLevel: 1,
  facilitySlots: 4,
  facilities: [
    { id: 'fac-firewall', type: 'Firewall', name: '基礎防火牆', level: 1, maxLevel: 5, description: '過濾惡意流量' },
  ],
};
