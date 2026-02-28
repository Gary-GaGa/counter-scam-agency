import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT } from '../ui/constants';
import { createButton, createTitle } from '../ui/components';
import { gameState } from '../services/gameState';

/**
 * 訊號追蹤（Signal Trace）
 * 接水管/駭客玩法，對應 Tech 數值。
 *
 * 玩法：5x5 網格，每格是一段管線（直線、彎管、T 字管、十字管）。
 * 點擊格子旋轉管線，讓訊號從起點（左側）連通到終點（右側）。
 * 限時 45 秒完成越多關卡越好。
 */

// 管線類型：每個方向 (上右下左) 為 boolean
type PipeType = 'straight' | 'bend' | 'tee' | 'cross';

interface PipeCell {
  type: PipeType;
  rotation: number; // 0=0° 1=90° 2=180° 3=270°
  row: number;
  col: number;
  connected: boolean;
}

const GRID_SIZE = 5;
const CELL_SIZE = 72;
const GRID_OFFSET_X = (GAME_WIDTH - GRID_SIZE * CELL_SIZE) / 2;
const GRID_OFFSET_Y = 100;
const GAME_DURATION = 45;

// 各管線類型的基本連接方向 (rotation=0): [上, 右, 下, 左]
const PIPE_CONNECTIONS: Record<PipeType, boolean[]> = {
  straight: [true, false, true, false],   // ┃ 上下
  bend:     [true, true, false, false],    // ┗ 上右
  tee:      [true, true, true, false],     // ┣ 上右下
  cross:    [true, true, true, true],      // ╋ 全通
};

function getConnections(cell: PipeCell): boolean[] {
  const base = PIPE_CONNECTIONS[cell.type];
  const r = cell.rotation % 4;
  // 順時針旋轉 r 次
  const result = [...base];
  for (let i = 0; i < r; i++) {
    const last = result.pop()!;
    result.unshift(last);
  }
  return result;
}

export class SignalTraceScene extends Phaser.Scene {
  private grid: PipeCell[][] = [];
  private cellGraphics: Phaser.GameObjects.Graphics[][] = [];
  private score = 0;
  private level = 1;
  private timeLeft = GAME_DURATION;
  private scoreText!: Phaser.GameObjects.Text;
  private timerText!: Phaser.GameObjects.Text;
  private levelText!: Phaser.GameObjects.Text;
  private gameTimer!: Phaser.Time.TimerEvent;
  private isGameOver = false;

  constructor() {
    super('SignalTraceScene');
  }

  create(): void {
    this.score = 0;
    this.level = 1;
    this.timeLeft = GAME_DURATION;
    this.isGameOver = false;

    this.cameras.main.fadeIn(300, 0, 0, 0);
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);

    // HUD
    this.scoreText = this.add.text(20, 16, '分數：0', {
      fontFamily: 'monospace', fontSize: '20px', color: '#eaeaea',
    }).setDepth(100);

    this.timerText = this.add.text(GAME_WIDTH - 20, 16, `⏱ ${GAME_DURATION}s`, {
      fontFamily: 'monospace', fontSize: '20px', color: '#ffc107',
    }).setOrigin(1, 0).setDepth(100);

    this.levelText = this.add.text(GAME_WIDTH / 2, 16, `關卡 ${this.level}`, {
      fontFamily: 'monospace', fontSize: '20px', color: '#81c784',
    }).setOrigin(0.5, 0).setDepth(100);

    this.add.text(GAME_WIDTH / 2, 55, '🖱️ 點擊格子旋轉管線，連通左側 → 右側訊號！', {
      fontFamily: 'monospace', fontSize: '13px', color: '#8888aa',
    }).setOrigin(0.5).setDepth(100);

    // 起點/終點標記
    const startY = GRID_OFFSET_Y + 2 * CELL_SIZE + CELL_SIZE / 2;
    this.add.text(GRID_OFFSET_X - 40, startY, '📡', {
      fontSize: '28px',
    }).setOrigin(0.5).setDepth(50);
    this.add.text(GRID_OFFSET_X + GRID_SIZE * CELL_SIZE + 40, startY, '🖥️', {
      fontSize: '28px',
    }).setOrigin(0.5).setDepth(50);

    this.generateLevel();

    this.gameTimer = this.time.addEvent({
      delay: 1000,
      callback: this.tick,
      callbackScope: this,
      loop: true,
    });
  }

  private generateLevel(): void {
    // 清除舊格子
    this.cellGraphics.forEach(row => row.forEach(g => g.destroy()));
    this.cellGraphics = [];
    this.grid = [];

    // 先生成一條通路，再隨機打亂旋轉
    const path = this.generatePath();

    for (let r = 0; r < GRID_SIZE; r++) {
      this.grid[r] = [];
      this.cellGraphics[r] = [];
      for (let c = 0; c < GRID_SIZE; c++) {
        const type = path[r][c];
        // 隨機旋轉（讓玩家需要解謎）
        const rotation = Phaser.Math.Between(0, 3);
        const cell: PipeCell = { type, rotation, row: r, col: c, connected: false };
        this.grid[r][c] = cell;

        const g = this.add.graphics();
        this.cellGraphics[r][c] = g;

        // 互動區域
        const x = GRID_OFFSET_X + c * CELL_SIZE;
        const y = GRID_OFFSET_Y + r * CELL_SIZE;
        const zone = this.add.zone(x + CELL_SIZE / 2, y + CELL_SIZE / 2, CELL_SIZE, CELL_SIZE)
          .setInteractive({ useHandCursor: true });

        zone.on('pointerdown', () => {
          if (this.isGameOver) return;
          cell.rotation = (cell.rotation + 1) % 4;
          this.drawGrid();
          this.checkConnection();
        });
      }
    }
    this.drawGrid();
  }

  private generatePath(): PipeType[][] {
    // 建立基本棋盤（全部 straight）
    const board: PipeType[][] = Array.from({ length: GRID_SIZE }, () =>
      Array.from({ length: GRID_SIZE }, () => 'straight' as PipeType),
    );

    // 生成隨機路徑：從 (2, 0) 到 (2, 4)
    const visited = new Set<string>();
    const path: { r: number; c: number }[] = [];

    const buildPath = (r: number, c: number): boolean => {
      if (c === GRID_SIZE - 1) {
        path.push({ r, c });
        return true;
      }
      if (r < 0 || r >= GRID_SIZE || c < 0 || c >= GRID_SIZE) return false;
      const key = `${r},${c}`;
      if (visited.has(key)) return false;
      visited.add(key);
      path.push({ r, c });

      // 優先往右走
      const dirs = [
        { dr: 0, dc: 1 },
        { dr: -1, dc: 0 },
        { dr: 1, dc: 0 },
      ];
      Phaser.Utils.Array.Shuffle(dirs);
      // 保證往右的方向排前面
      dirs.sort((a, b) => b.dc - a.dc);

      for (const d of dirs) {
        if (buildPath(r + d.dr, c + d.dc)) return true;
      }

      path.pop();
      visited.delete(key);
      return false;
    };

    buildPath(2, 0);

    // 根據路徑設定管線類型
    for (let i = 0; i < path.length; i++) {
      const { r, c } = path[i];
      const prev = i > 0 ? path[i - 1] : { r: 2, c: -1 };
      const next = i < path.length - 1 ? path[i + 1] : { r: path[i].r, c: GRID_SIZE };

      const fromDir = this.getDirection(r, c, prev.r, prev.c);
      const toDir = this.getDirection(r, c, next.r, next.c);

      board[r][c] = this.getPipeForDirections(fromDir, toDir);
    }

    // 非路徑格子隨機類型
    const types: PipeType[] = ['straight', 'bend', 'tee', 'cross'];
    for (let r = 0; r < GRID_SIZE; r++) {
      for (let c = 0; c < GRID_SIZE; c++) {
        if (!path.some(p => p.r === r && p.c === c)) {
          board[r][c] = types[Phaser.Math.Between(0, types.length - 1)];
        }
      }
    }

    return board;
  }

  private getDirection(fromR: number, fromC: number, toR: number, toC: number): number {
    // 0=上 1=右 2=下 3=左
    if (toR < fromR) return 0;
    if (toC > fromC) return 1;
    if (toR > fromR) return 2;
    return 3;
  }

  private getPipeForDirections(d1: number, d2: number): PipeType {
    const dirs = [d1, d2].sort();
    // 對向 = straight
    if ((dirs[0] === 0 && dirs[1] === 2) || (dirs[0] === 1 && dirs[1] === 3)) return 'straight';
    // 其他 = bend
    return 'bend';
  }

  private drawGrid(): void {
    // 先檢查連通性來上色
    this.updateConnections();

    for (let r = 0; r < GRID_SIZE; r++) {
      for (let c = 0; c < GRID_SIZE; c++) {
        const cell = this.grid[r][c];
        const g = this.cellGraphics[r][c];
        g.clear();

        const x = GRID_OFFSET_X + c * CELL_SIZE;
        const y = GRID_OFFSET_Y + r * CELL_SIZE;
        const cx = x + CELL_SIZE / 2;
        const cy = y + CELL_SIZE / 2;

        // 背景
        g.fillStyle(Colors.panel, 0.9);
        g.fillRect(x + 2, y + 2, CELL_SIZE - 4, CELL_SIZE - 4);
        g.lineStyle(1, Colors.panelLight);
        g.strokeRect(x + 2, y + 2, CELL_SIZE - 4, CELL_SIZE - 4);

        // 繪製管線
        const pipeColor = cell.connected ? 0x53d769 : 0x4fc3f7;
        const pipeWidth = 10;
        g.lineStyle(pipeWidth, pipeColor, 0.9);

        const connections = getConnections(cell);
        const halfCell = CELL_SIZE / 2;

        // 從中心到各方向
        if (connections[0]) { // 上
          g.lineBetween(cx, cy, cx, y + 2);
        }
        if (connections[1]) { // 右
          g.lineBetween(cx, cy, x + CELL_SIZE - 2, cy);
        }
        if (connections[2]) { // 下
          g.lineBetween(cx, cy, cx, y + CELL_SIZE - 2);
        }
        if (connections[3]) { // 左
          g.lineBetween(cx, cy, x + 2, cy);
        }

        // 中心圓點
        g.fillStyle(pipeColor, 1);
        g.fillCircle(cx, cy, 6);
      }
    }
  }

  private updateConnections(): void {
    // BFS 從左側中間行開始
    for (let r = 0; r < GRID_SIZE; r++) {
      for (let c = 0; c < GRID_SIZE; c++) {
        this.grid[r][c].connected = false;
      }
    }

    const queue: { r: number; c: number }[] = [];
    // 起點：左側邊界，中間行，需要有左連接
    const startRow = 2;
    const startConn = getConnections(this.grid[startRow][0]);
    if (startConn[3]) { // 有左連接（通往起點）
      this.grid[startRow][0].connected = true;
      queue.push({ r: startRow, c: 0 });
    }

    const dirs = [
      { dr: -1, dc: 0, from: 0, to: 2 }, // 上：自己的上 = 鄰居的下
      { dr: 0, dc: 1, from: 1, to: 3 },   // 右
      { dr: 1, dc: 0, from: 2, to: 0 },   // 下
      { dr: 0, dc: -1, from: 3, to: 1 },  // 左
    ];

    while (queue.length > 0) {
      const { r, c } = queue.shift()!;
      const conn = getConnections(this.grid[r][c]);

      for (const d of dirs) {
        if (!conn[d.from]) continue;
        const nr = r + d.dr;
        const nc = c + d.dc;
        if (nr < 0 || nr >= GRID_SIZE || nc < 0 || nc >= GRID_SIZE) continue;
        if (this.grid[nr][nc].connected) continue;

        const neighborConn = getConnections(this.grid[nr][nc]);
        if (neighborConn[d.to]) {
          this.grid[nr][nc].connected = true;
          queue.push({ r: nr, c: nc });
        }
      }
    }
  }

  private checkConnection(): void {
    // 檢查右側是否有連通格子且有右向連接
    for (let r = 0; r < GRID_SIZE; r++) {
      const cell = this.grid[r][GRID_SIZE - 1];
      if (cell.connected) {
        const conn = getConnections(cell);
        if (conn[1]) { // 有右連接（通往終點）
          this.levelComplete();
          return;
        }
      }
    }
  }

  private levelComplete(): void {
    const timeBonus = this.timeLeft * 2;
    const levelPoints = this.level * 20;
    this.score += levelPoints + timeBonus;
    this.level++;

    // 閃光效果
    const flash = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, 0x53d769, 0.2,
    ).setDepth(200);
    this.tweens.add({
      targets: flash, alpha: 0, duration: 400,
      onComplete: () => flash.destroy(),
    });

    const msg = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2, `✓ 連通！+${levelPoints + timeBonus}`, {
      fontFamily: 'monospace', fontSize: '24px', color: '#53d769',
    }).setOrigin(0.5).setDepth(201);
    this.tweens.add({
      targets: msg, y: msg.y - 50, alpha: 0, duration: 1000,
      onComplete: () => msg.destroy(),
    });

    this.scoreText.setText(`分數：${this.score}`);
    this.levelText.setText(`關卡 ${this.level}`);

    this.time.delayedCall(600, () => {
      if (!this.isGameOver) this.generateLevel();
    });
  }

  private tick(): void {
    if (this.isGameOver) return;
    this.timeLeft--;
    this.timerText.setText(`⏱ ${this.timeLeft}s`);
    if (this.timeLeft <= 10) this.timerText.setColor('#e94560');
    if (this.timeLeft <= 0) this.endGame();
  }

  private endGame(): void {
    this.isGameOver = true;
    this.gameTimer?.destroy();

    const overlay = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, 0x000000, 0.7,
    ).setDepth(300);

    const rank = this.getRank();

    createTitle(this, GAME_WIDTH / 2, 180, '📡 訊號追蹤 — 結算', '28px').setDepth(301);
    this.add.text(GAME_WIDTH / 2, 250, `最終分數：${this.score}`, {
      fontFamily: 'monospace', fontSize: '28px', color: '#81c784',
    }).setOrigin(0.5).setDepth(301);
    this.add.text(GAME_WIDTH / 2, 290, `完成關卡：${this.level - 1}`, {
      fontFamily: 'monospace', fontSize: '18px', color: '#bbbbcc',
    }).setOrigin(0.5).setDepth(301);
    this.add.text(GAME_WIDTH / 2, 330, `評價：${rank}`, {
      fontFamily: 'monospace', fontSize: '22px', color: '#ffc107',
    }).setOrigin(0.5).setDepth(301);

    const result = gameState.recordMiniGame('signalTrace', this.score);
    const bonusColor = result.isNewBest ? '#53d769' : '#8888aa';
    this.add.text(GAME_WIDTH / 2, 365, result.isNewBest
      ? `🆕 新紀錄！ 技術 +${result.bonus}`
      : `技術 +${result.bonus}`, {
      fontFamily: 'monospace', fontSize: '16px', color: bonusColor,
    }).setOrigin(0.5).setDepth(301);

    const retryBtn = createButton(this, GAME_WIDTH / 2 - 120, 430, '🔄 再挑戰', () => {
      this.scene.restart();
    }, 200, 48);
    retryBtn.setDepth(301);

    const backBtn = createButton(this, GAME_WIDTH / 2 + 120, 430, '← 返回選單', () => {
      this.scene.start('MainMenuScene');
    }, 200, 48);
    backBtn.setDepth(301);
  }

  private getRank(): string {
    if (this.score >= 200) return '🏆 S — 駭客大師';
    if (this.score >= 120) return '🥇 A — 資深工程師';
    if (this.score >= 60) return '🥈 B — 合格技術員';
    if (this.score >= 30) return '🥉 C — 實習生';
    return '📝 D — 需要更多練習';
  }
}
