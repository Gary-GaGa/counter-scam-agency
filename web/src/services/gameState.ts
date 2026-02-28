/**
 * 遊戲狀態管理 — 使用 localStorage 持久化玩家進度。
 * 所有場景共用此狀態，小遊戲結果會反映在玩家屬性上。
 */

export interface MiniGameScore {
  score: number;
  rank: string;
  playedAt: number;
}

export interface GameStateData {
  reputation: number;
  stats: {
    logic: number;
    tech: number;
    charisma: number;
    resilience: number;
  };
  miniGameBest: {
    contradiction: MiniGameScore | null;
    signalTrace: MiniGameScore | null;
    negotiation: MiniGameScore | null;
    mentalRecovery: MiniGameScore | null;
  };
  completedMissions: string[];
  unlockedSkills: string[];
  introSeen: boolean;
  totalGamesPlayed: number;
}

const STORAGE_KEY = 'counter-scam-agency-state';

const DEFAULT_STATE: GameStateData = {
  reputation: 0,
  stats: { logic: 10, tech: 10, charisma: 10, resilience: 10 },
  miniGameBest: {
    contradiction: null,
    signalTrace: null,
    negotiation: null,
    mentalRecovery: null,
  },
  completedMissions: [],
  unlockedSkills: [],
  introSeen: false,
  totalGamesPlayed: 0,
};

type MiniGameKey = keyof GameStateData['miniGameBest'];
type StatKey = keyof GameStateData['stats'];

const GAME_STAT_MAP: Record<MiniGameKey, StatKey> = {
  contradiction: 'logic',
  signalTrace: 'tech',
  negotiation: 'charisma',
  mentalRecovery: 'resilience',
};

// 根據分數計算屬性加成
function calcStatBonus(score: number): number {
  if (score >= 120) return 5;
  if (score >= 80) return 3;
  if (score >= 50) return 2;
  if (score >= 20) return 1;
  return 0;
}

function getRankFromScore(score: number): string {
  if (score >= 120) return 'S';
  if (score >= 80) return 'A';
  if (score >= 50) return 'B';
  if (score >= 20) return 'C';
  return 'D';
}

class GameState {
  private data: GameStateData;

  constructor() {
    this.data = this.load();
  }

  private load(): GameStateData {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw);
        // 用 spread 確保新增的欄位有預設值
        return { ...DEFAULT_STATE, ...parsed, stats: { ...DEFAULT_STATE.stats, ...parsed.stats }, miniGameBest: { ...DEFAULT_STATE.miniGameBest, ...parsed.miniGameBest } };
      }
    } catch { /* 忽略解析錯誤 */ }
    return { ...DEFAULT_STATE, stats: { ...DEFAULT_STATE.stats }, miniGameBest: { ...DEFAULT_STATE.miniGameBest } };
  }

  private save(): void {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.data));
    } catch { /* localStorage 不可用時靜默失敗 */ }
  }

  get(): Readonly<GameStateData> {
    return this.data;
  }

  /**
   * 記錄小遊戲結果。回傳屬性提升量與增量（可能為 0）。
   */
  recordMiniGame(game: MiniGameKey, score: number): { statName: StatKey; bonus: number; delta: number; isNewBest: boolean } {
    const current = this.data.miniGameBest[game];
    const oldBonus = current ? calcStatBonus(current.score) : 0;
    const isNewBest = !current || score > current.score;

    if (isNewBest) {
      this.data.miniGameBest[game] = {
        score,
        rank: getRankFromScore(score),
        playedAt: Date.now(),
      };
    }

    this.data.totalGamesPlayed++;

    // 根據最高分重新計算屬性（避免重複加成）
    const statName = GAME_STAT_MAP[game];
    const bestScore = this.data.miniGameBest[game]!.score;
    const bonus = calcStatBonus(bestScore);
    const delta = bonus - oldBonus;
    this.data.stats[statName] = 10 + bonus; // 基礎 10 + 小遊戲加成

    this.save();
    return { statName, bonus, delta, isNewBest };
  }

  addReputation(amount: number): void {
    this.data.reputation += amount;
    this.save();
  }

  completeMission(missionId: string): void {
    if (!this.data.completedMissions.includes(missionId)) {
      this.data.completedMissions.push(missionId);
      this.save();
    }
  }

  unlockSkill(skillId: string): void {
    if (!this.data.unlockedSkills.includes(skillId)) {
      this.data.unlockedSkills.push(skillId);
      this.save();
    }
  }

  markIntroSeen(): void {
    this.data.introSeen = true;
    this.save();
  }

  reset(): void {
    this.data = { ...DEFAULT_STATE, stats: { ...DEFAULT_STATE.stats }, miniGameBest: { ...DEFAULT_STATE.miniGameBest } };
    this.save();
  }
}

export const gameState = new GameState();
