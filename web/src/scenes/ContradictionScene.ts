import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT } from '../ui/constants';
import { createButton, createTitle } from '../ui/components';
import { gameState } from '../services/gameState';

/**
 * 矛盾擊破（Contradiction Breaker）
 * 彈幕射擊式小遊戲，對應 Logic 數值。
 *
 * 玩法：畫面上方會飄過「證詞泡泡」，其中有些是「謊言」。
 * 玩家操作底部的十字準心，點擊擊破謊言泡泡。
 * 擊中謊言 +10 分，誤擊真話 -5 分，漏掉謊言 -3 分。
 */

interface Statement {
  text: string;
  isLie: boolean;
}

interface Bubble {
  container: Phaser.GameObjects.Container;
  statement: Statement;
  speed: number;
}

// 預設題庫：真話與謊言
const STATEMENT_POOL: Statement[] = [
  { text: '我是銀行客服，需要您的密碼', isLie: true },
  { text: '請撥打 165 反詐騙專線確認', isLie: false },
  { text: '限時三分鐘內匯款否則帳戶凍結', isLie: true },
  { text: '正規機構不會要求提供驗證碼', isLie: false },
  { text: '投資保證月報酬 30%', isLie: true },
  { text: '高報酬必定伴隨高風險', isLie: false },
  { text: '加入群組就能穩賺不賠', isLie: true },
  { text: '切勿點擊不明連結', isLie: false },
  { text: '中獎通知：請先繳納手續費', isLie: true },
  { text: '正規抽獎不會事先收費', isLie: false },
  { text: '您的包裹有違禁品需繳保證金', isLie: true },
  { text: '遇可疑電話應主動求證', isLie: false },
  { text: '交往對象急著借錢是警訊', isLie: false },
  { text: '只要轉帳到安全帳戶就沒事', isLie: true },
  { text: '監管機關不會透過電話辦案', isLie: false },
  { text: '提供身分證照片才能領獎', isLie: true },
];

const GAME_DURATION = 30; // 秒
const SPAWN_INTERVAL_MIN = 800;
const SPAWN_INTERVAL_MAX = 1800;

export class ContradictionScene extends Phaser.Scene {
  private score = 0;
  private timeLeft = GAME_DURATION;
  private bubbles: Bubble[] = [];
  private scoreText!: Phaser.GameObjects.Text;
  private timerText!: Phaser.GameObjects.Text;
  private comboCount = 0;
  private comboText!: Phaser.GameObjects.Text;
  private spawnTimer!: Phaser.Time.TimerEvent;
  private gameTimer!: Phaser.Time.TimerEvent;
  private isGameOver = false;
  private usedStatements: Set<number> = new Set();

  constructor() {
    super('ContradictionScene');
  }

  create(): void {
    this.score = 0;
    this.timeLeft = GAME_DURATION;
    this.bubbles = [];
    this.comboCount = 0;
    this.isGameOver = false;
    this.usedStatements.clear();

    this.cameras.main.fadeIn(300, 0, 0, 0);
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);

    // HUD
    this.scoreText = this.add.text(20, 16, '分數：0', {
      fontFamily: 'monospace', fontSize: '20px', color: '#eaeaea',
    }).setDepth(100);

    this.timerText = this.add.text(GAME_WIDTH - 20, 16, `⏱ ${GAME_DURATION}s`, {
      fontFamily: 'monospace', fontSize: '20px', color: '#ffc107',
    }).setOrigin(1, 0).setDepth(100);

    this.comboText = this.add.text(GAME_WIDTH / 2, 16, '', {
      fontFamily: 'monospace', fontSize: '18px', color: '#53d769',
    }).setOrigin(0.5, 0).setDepth(100);

    // 指示
    this.add.text(GAME_WIDTH / 2, 50, '👆 點擊謊言泡泡！放過真話！', {
      fontFamily: 'monospace', fontSize: '14px', color: '#8888aa',
    }).setOrigin(0.5).setDepth(100);

    // 分隔線（安全區）
    this.add.line(0, 0, 0, 580, GAME_WIDTH, 580, 0x333355, 0.5)
      .setOrigin(0).setDepth(50);

    // 開始生成泡泡
    this.spawnTimer = this.time.addEvent({
      delay: this.getSpawnDelay(),
      callback: this.spawnBubble,
      callbackScope: this,
      loop: false,
    });

    // 倒數計時
    this.gameTimer = this.time.addEvent({
      delay: 1000,
      callback: this.tick,
      callbackScope: this,
      loop: true,
    });
  }

  update(): void {
    if (this.isGameOver) return;

    // 移動泡泡
    for (let i = this.bubbles.length - 1; i >= 0; i--) {
      const bubble = this.bubbles[i];
      bubble.container.x += bubble.speed;

      // 超出畫面
      if (bubble.container.x > GAME_WIDTH + 200 || bubble.container.x < -200) {
        if (bubble.statement.isLie) {
          this.score -= 3;
          this.comboCount = 0;
          this.flashEffect(0xe94560);
        }
        bubble.container.destroy();
        this.bubbles.splice(i, 1);
      }
    }

    this.scoreText.setText(`分數：${this.score}`);
  }

  private getSpawnDelay(): number {
    // 隨時間加快
    const elapsed = GAME_DURATION - this.timeLeft;
    const factor = Math.max(0.4, 1 - elapsed / GAME_DURATION * 0.6);
    return Phaser.Math.Between(
      SPAWN_INTERVAL_MIN * factor,
      SPAWN_INTERVAL_MAX * factor,
    );
  }

  private getRandomStatement(): Statement {
    // 盡量不重複
    const available = STATEMENT_POOL
      .map((s, i) => ({ s, i }))
      .filter(({ i }) => !this.usedStatements.has(i));

    if (available.length === 0) {
      this.usedStatements.clear();
      return STATEMENT_POOL[Phaser.Math.Between(0, STATEMENT_POOL.length - 1)];
    }

    const pick = available[Phaser.Math.Between(0, available.length - 1)];
    this.usedStatements.add(pick.i);
    return pick.s;
  }

  private spawnBubble(): void {
    if (this.isGameOver) return;

    const statement = this.getRandomStatement();

    // 左側或右側進入
    const fromLeft = Math.random() > 0.5;
    const x = fromLeft ? -180 : GAME_WIDTH + 180;
    const y = Phaser.Math.Between(80, 520);
    const speed = (fromLeft ? 1 : -1) * Phaser.Math.FloatBetween(1.2, 2.8);

    // 建立泡泡
    const textW = Math.min(statement.text.length * 16, 360);
    const bg = this.add.rectangle(0, 0, textW + 30, 44, Colors.panel, 0.92)
      .setStrokeStyle(2, statement.isLie ? Colors.danger : Colors.panelLight);

    const label = this.add.text(0, 0, statement.text, {
      fontFamily: 'monospace', fontSize: '15px',
      color: '#eaeaea',
      wordWrap: { width: textW },
      align: 'center',
    }).setOrigin(0.5);

    const container = this.add.container(x, y, [bg, label]);
    container.setSize(textW + 30, 44);
    container.setInteractive({ useHandCursor: true });
    container.setDepth(10);

    container.on('pointerdown', () => {
      if (this.isGameOver) return;
      this.popBubble(bubble);
    });

    const bubble: Bubble = { container, statement, speed };
    this.bubbles.push(bubble);

    // 安排下一個泡泡
    this.spawnTimer = this.time.addEvent({
      delay: this.getSpawnDelay(),
      callback: this.spawnBubble,
      callbackScope: this,
      loop: false,
    });
  }

  private popBubble(bubble: Bubble): void {
    const idx = this.bubbles.indexOf(bubble);
    if (idx === -1) return;

    this.bubbles.splice(idx, 1);

    if (bubble.statement.isLie) {
      // 正確擊破謊言
      this.comboCount++;
      const comboBonus = Math.min(this.comboCount - 1, 5) * 2;
      this.score += 10 + comboBonus;
      this.showPopEffect(bubble.container.x, bubble.container.y, '✓ 擊破！', '#53d769');

      if (this.comboCount >= 2) {
        this.comboText.setText(`🔥 連擊 x${this.comboCount}`);
        this.time.delayedCall(1200, () => {
          if (this.comboText.active) this.comboText.setText('');
        });
      }
    } else {
      // 誤擊真話
      this.score -= 5;
      this.comboCount = 0;
      this.comboText.setText('');
      this.showPopEffect(bubble.container.x, bubble.container.y, '✗ 那是真話！', '#e94560');
      this.flashEffect(0xe94560);
    }

    bubble.container.destroy();
    this.scoreText.setText(`分數：${this.score}`);
  }

  private showPopEffect(x: number, y: number, msg: string, color: string): void {
    const text = this.add.text(x, y, msg, {
      fontFamily: 'monospace', fontSize: '16px', color,
    }).setOrigin(0.5).setDepth(50);

    this.tweens.add({
      targets: text,
      y: y - 40,
      alpha: 0,
      duration: 800,
      onComplete: () => text.destroy(),
    });
  }

  private flashEffect(color: number): void {
    const flash = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, color, 0.15,
    ).setDepth(200);
    this.tweens.add({
      targets: flash,
      alpha: 0,
      duration: 300,
      onComplete: () => flash.destroy(),
    });
  }

  private tick(): void {
    if (this.isGameOver) return;
    this.timeLeft--;
    this.timerText.setText(`⏱ ${this.timeLeft}s`);

    if (this.timeLeft <= 5) {
      this.timerText.setColor('#e94560');
    }

    if (this.timeLeft <= 0) {
      this.endGame();
    }
  }

  private endGame(): void {
    this.isGameOver = true;
    this.spawnTimer?.destroy();
    this.gameTimer?.destroy();

    // 清除剩餘泡泡
    this.bubbles.forEach(b => b.container.destroy());
    this.bubbles = [];

    // 結果畫面
    const overlay = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, 0x000000, 0.7,
    ).setDepth(300);

    const rank = this.getRank();

    const resultTitle = createTitle(this, GAME_WIDTH / 2, 180,
      '⚡ 矛盾擊破 — 結算', '28px');
    resultTitle.setDepth(301);

    const scoreResult = this.add.text(GAME_WIDTH / 2, 260, `最終分數：${this.score}`, {
      fontFamily: 'monospace', fontSize: '28px', color: '#53d769',
    }).setOrigin(0.5).setDepth(301);

    const rankText = this.add.text(GAME_WIDTH / 2, 310, `評價：${rank}`, {
      fontFamily: 'monospace', fontSize: '22px', color: '#ffc107',
    }).setOrigin(0.5).setDepth(301);

    // 記錄至遊戲狀態
    const result = gameState.recordMiniGame('contradiction', this.score);
    const statLabel = '邏輯';
    const bonusColor = result.isNewBest ? '#53d769' : '#8888aa';
    this.add.text(GAME_WIDTH / 2, 350, result.isNewBest
      ? `🆕 新紀錄！ ${statLabel} +${result.bonus}`
      : `${statLabel} +${result.bonus}`, {
      fontFamily: 'monospace', fontSize: '16px', color: bonusColor,
    }).setOrigin(0.5).setDepth(301);

    const retryBtn = createButton(this, GAME_WIDTH / 2 - 120, 420, '🔄 再挑戰', () => {
      this.scene.restart();
    }, 200, 48);
    retryBtn.setDepth(301);

    const backBtn = createButton(this, GAME_WIDTH / 2 + 120, 420, '← 返回選單', () => {
      this.scene.start('MainMenuScene');
    }, 200, 48);
    backBtn.setDepth(301);
  }

  private getRank(): string {
    if (this.score >= 120) return '🏆 S — 頂尖情報官';
    if (this.score >= 80) return '🥇 A — 資深偵查員';
    if (this.score >= 50) return '🥈 B — 合格探員';
    if (this.score >= 20) return '🥉 C — 見習生';
    return '📝 D — 需要更多訓練';
  }
}
