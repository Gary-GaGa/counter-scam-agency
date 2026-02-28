import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT, PLAYER_ID } from '../ui/constants';
import { createButton, createTitle } from '../ui/components';
import { gameState } from '../services/gameState';
import { updatePlayerStats } from '../services/api';

/**
 * 心靈調適（Mental Recovery）
 * 節奏/休閒玩法，對應 Resilience 數值。
 *
 * 玩法：畫面上會落下「情緒泡泡」（正面/負面）。
 * 正面泡泡（綠色）落到底線時按空白鍵或點擊接住 → 恢復能量。
 * 負面泡泡（紅色）要避開（不按任何鍵讓它自然消失）。
 * 連續正確接住正面泡泡可產生 combo。
 */

interface EmoBubble {
  sprite: Phaser.GameObjects.Container;
  positive: boolean;
  lane: number;
  speed: number;
}

const POSITIVE_LABELS = ['希望', '冷靜', '信任', '勇氣', '陪伴', '支持', '理解', '安全'];
const NEGATIVE_LABELS = ['恐懼', '焦慮', '憤怒', '羞恥', '孤立', '急迫', '懷疑', '絕望'];

const LANE_COUNT = 4;
const CATCH_ZONE_Y = 540;
const GAME_DURATION = 35;
const SPAWN_INTERVAL_BASE = 900;

export class MentalRecoveryScene extends Phaser.Scene {
  private bubbles: EmoBubble[] = [];
  private score = 0;
  private energy = 50;       // 0~100
  private combo = 0;
  private timeLeft = GAME_DURATION;
  private isGameOver = false;

  private scoreText!: Phaser.GameObjects.Text;
  private timerText!: Phaser.GameObjects.Text;
  private comboText!: Phaser.GameObjects.Text;
  private energyBar!: Phaser.GameObjects.Rectangle;
  private energyBarBg!: Phaser.GameObjects.Rectangle;
  private energyLabel!: Phaser.GameObjects.Text;
  private spawnTimer!: Phaser.Time.TimerEvent;
  private gameTimer!: Phaser.Time.TimerEvent;
  private laneButtons: Phaser.GameObjects.Rectangle[] = [];

  constructor() {
    super('MentalRecoveryScene');
  }

  create(): void {
    this.score = 0;
    this.energy = 50;
    this.combo = 0;
    this.timeLeft = GAME_DURATION;
    this.isGameOver = false;
    this.bubbles = [];

    this.cameras.main.fadeIn(300, 0, 0, 0);
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);

    // HUD
    this.scoreText = this.add.text(20, 16, '分數：0', {
      fontFamily: 'monospace', fontSize: '18px', color: '#eaeaea',
    }).setDepth(100);

    this.timerText = this.add.text(GAME_WIDTH - 20, 16, `⏱ ${GAME_DURATION}s`, {
      fontFamily: 'monospace', fontSize: '18px', color: '#ffc107',
    }).setOrigin(1, 0).setDepth(100);

    this.comboText = this.add.text(GAME_WIDTH / 2, 16, '', {
      fontFamily: 'monospace', fontSize: '16px', color: '#ba68c8',
    }).setOrigin(0.5, 0).setDepth(100);

    // 能量條
    this.energyLabel = this.add.text(GAME_WIDTH / 2, 46, '🧠 心靈能量', {
      fontFamily: 'monospace', fontSize: '13px', color: '#8888aa',
    }).setOrigin(0.5).setDepth(100);

    this.energyBarBg = this.add.rectangle(GAME_WIDTH / 2, 68, 400, 16, 0x333355).setDepth(100);
    this.energyBar = this.add.rectangle(GAME_WIDTH / 2 - 200, 68, 200, 16, Colors.statResilience)
      .setOrigin(0, 0.5).setDepth(101);

    // 落點區域（4 車道）
    const laneWidth = 160;
    const lanesStart = (GAME_WIDTH - LANE_COUNT * laneWidth) / 2;

    for (let i = 0; i < LANE_COUNT; i++) {
      const x = lanesStart + i * laneWidth + laneWidth / 2;

      // 車道分隔線
      if (i > 0) {
        this.add.line(0, 0, lanesStart + i * laneWidth, 90, lanesStart + i * laneWidth, CATCH_ZONE_Y + 30, 0x333355, 0.3)
          .setOrigin(0).setDepth(1);
      }

      // 接住區
      const catchZone = this.add.rectangle(x, CATCH_ZONE_Y, laneWidth - 10, 40, Colors.panelLight, 0.4)
        .setStrokeStyle(2, Colors.statResilience, 0.5)
        .setInteractive({ useHandCursor: true })
        .setDepth(5);
      this.laneButtons.push(catchZone);

      catchZone.on('pointerdown', () => this.catchBubbleInLane(i));

      // 車道標籤
      this.add.text(x, CATCH_ZONE_Y + 30, `${i + 1}`, {
        fontFamily: 'monospace', fontSize: '14px', color: '#555577',
      }).setOrigin(0.5).setDepth(5);
    }

    // 接住線
    this.add.line(0, 0, lanesStart, CATCH_ZONE_Y - 20, lanesStart + LANE_COUNT * laneWidth, CATCH_ZONE_Y - 20, Colors.statResilience, 0.3)
      .setOrigin(0).setDepth(3);

    // 鍵盤控制
    this.input.keyboard?.on('keydown-ONE', () => this.catchBubbleInLane(0));
    this.input.keyboard?.on('keydown-TWO', () => this.catchBubbleInLane(1));
    this.input.keyboard?.on('keydown-THREE', () => this.catchBubbleInLane(2));
    this.input.keyboard?.on('keydown-FOUR', () => this.catchBubbleInLane(3));

    // 開始生成
    this.scheduleSpawn();

    this.gameTimer = this.time.addEvent({
      delay: 1000,
      callback: this.tick,
      callbackScope: this,
      loop: true,
    });
  }

  update(): void {
    if (this.isGameOver) return;

    for (let i = this.bubbles.length - 1; i >= 0; i--) {
      const bubble = this.bubbles[i];
      bubble.sprite.y += bubble.speed;

      // 超過接住區
      if (bubble.sprite.y > CATCH_ZONE_Y + 60) {
        if (bubble.positive) {
          // 漏掉正面泡泡
          this.energy = Math.max(0, this.energy - 3);
          this.combo = 0;
          this.comboText.setText('');
        }
        // 負面泡泡自然消失 = 正確
        if (!bubble.positive) {
          this.score += 3;
        }
        bubble.sprite.destroy();
        this.bubbles.splice(i, 1);
      }
    }

    // 更新能量條
    const barWidth = Math.max(0, (this.energy / 100) * 400);
    this.energyBar.setSize(barWidth, 16);

    const barColor = this.energy > 60 ? Colors.statResilience
      : this.energy > 30 ? Colors.warning : Colors.danger;
    this.energyBar.setFillStyle(barColor);

    this.scoreText.setText(`分數：${this.score}`);

    if (this.energy <= 0) {
      this.endGame();
    }
  }

  private scheduleSpawn(): void {
    const elapsed = GAME_DURATION - this.timeLeft;
    const factor = Math.max(0.5, 1 - elapsed / GAME_DURATION * 0.4);
    const delay = SPAWN_INTERVAL_BASE * factor;

    this.spawnTimer = this.time.addEvent({
      delay,
      callback: () => {
        if (!this.isGameOver) {
          this.spawnBubble();
          this.scheduleSpawn();
        }
      },
      callbackScope: this,
    });
  }

  private spawnBubble(): void {
    const lane = Phaser.Math.Between(0, LANE_COUNT - 1);
    const positive = Math.random() > 0.35; // 65% 正面

    const laneWidth = 160;
    const lanesStart = (GAME_WIDTH - LANE_COUNT * laneWidth) / 2;
    const x = lanesStart + lane * laneWidth + laneWidth / 2;

    const label = positive
      ? POSITIVE_LABELS[Phaser.Math.Between(0, POSITIVE_LABELS.length - 1)]
      : NEGATIVE_LABELS[Phaser.Math.Between(0, NEGATIVE_LABELS.length - 1)];

    const color = positive ? 0x53d769 : 0xe94560;
    const textColor = positive ? '#53d769' : '#e94560';

    const circle = this.add.circle(0, 0, 28, color, 0.25)
      .setStrokeStyle(2, color, 0.8);
    const text = this.add.text(0, 0, label, {
      fontFamily: 'monospace', fontSize: '14px', color: textColor,
    }).setOrigin(0.5);

    const container = this.add.container(x, 90, [circle, text]).setDepth(10);

    const speed = Phaser.Math.FloatBetween(1.5, 3.0);
    this.bubbles.push({ sprite: container, positive, lane, speed });
  }

  private catchBubbleInLane(lane: number): void {
    if (this.isGameOver) return;

    // 閃光效果
    const btn = this.laneButtons[lane];
    if (btn) {
      btn.setFillStyle(Colors.statResilience, 0.6);
      this.time.delayedCall(150, () => btn.setFillStyle(Colors.panelLight, 0.4));
    }

    // 找最接近接住區的泡泡
    const nearBubbles = this.bubbles
      .filter(b => b.lane === lane && b.sprite.y > CATCH_ZONE_Y - 60 && b.sprite.y < CATCH_ZONE_Y + 40)
      .sort((a, b) => Math.abs(a.sprite.y - CATCH_ZONE_Y) - Math.abs(b.sprite.y - CATCH_ZONE_Y));

    if (nearBubbles.length === 0) return;

    const bubble = nearBubbles[0];
    const idx = this.bubbles.indexOf(bubble);
    if (idx === -1) return;

    this.bubbles.splice(idx, 1);

    if (bubble.positive) {
      // 正確接住正面泡泡
      this.combo++;
      const comboBonus = Math.min(this.combo - 1, 5) * 2;
      this.score += 8 + comboBonus;
      this.energy = Math.min(100, this.energy + 5);

      this.showEffect(bubble.sprite.x, bubble.sprite.y, '✓', '#53d769');

      if (this.combo >= 2) {
        this.comboText.setText(`✨ 連續 x${this.combo}`);
        this.time.delayedCall(1500, () => {
          if (this.comboText.active) this.comboText.setText('');
        });
      }
    } else {
      // 接到負面泡泡
      this.energy = Math.max(0, this.energy - 10);
      this.combo = 0;
      this.comboText.setText('');
      this.score -= 5;

      this.showEffect(bubble.sprite.x, bubble.sprite.y, '✗', '#e94560');
      this.flashScreen(0xe94560);
    }

    bubble.sprite.destroy();
  }

  private showEffect(x: number, y: number, msg: string, color: string): void {
    const text = this.add.text(x, y, msg, {
      fontFamily: 'monospace', fontSize: '20px', color,
    }).setOrigin(0.5).setDepth(50);
    this.tweens.add({
      targets: text, y: y - 40, alpha: 0, duration: 600,
      onComplete: () => text.destroy(),
    });
  }

  private flashScreen(color: number): void {
    const flash = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, color, 0.1,
    ).setDepth(200);
    this.tweens.add({
      targets: flash, alpha: 0, duration: 250,
      onComplete: () => flash.destroy(),
    });
  }

  private tick(): void {
    if (this.isGameOver) return;
    this.timeLeft--;
    this.timerText.setText(`⏱ ${this.timeLeft}s`);
    if (this.timeLeft <= 8) this.timerText.setColor('#e94560');
    if (this.timeLeft <= 0) this.endGame();
  }

  private endGame(): void {
    this.isGameOver = true;
    this.spawnTimer?.destroy();
    this.gameTimer?.destroy();
    this.bubbles.forEach(b => b.sprite.destroy());
    this.bubbles = [];

    const overlay = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, 0x000000, 0.7,
    ).setDepth(300);

    const rank = this.getRank();
    const survived = this.energy > 0;

    createTitle(this, GAME_WIDTH / 2, 170, '🧘 心靈調適 — 結算', '28px').setDepth(301);

    this.add.text(GAME_WIDTH / 2, 240, survived ? '✅ 心靈能量維持穩定！' : '❌ 心靈能量耗盡…', {
      fontFamily: 'monospace', fontSize: '20px', color: survived ? '#53d769' : '#e94560',
    }).setOrigin(0.5).setDepth(301);

    this.add.text(GAME_WIDTH / 2, 280, `最終分數：${this.score}`, {
      fontFamily: 'monospace', fontSize: '24px', color: '#ba68c8',
    }).setOrigin(0.5).setDepth(301);

    this.add.text(GAME_WIDTH / 2, 315, `剩餘能量：${this.energy}%`, {
      fontFamily: 'monospace', fontSize: '16px', color: '#bbbbcc',
    }).setOrigin(0.5).setDepth(301);

    this.add.text(GAME_WIDTH / 2, 350, `評價：${rank}`, {
      fontFamily: 'monospace', fontSize: '20px', color: '#ffc107',
    }).setOrigin(0.5).setDepth(301);

    const result = gameState.recordMiniGame('mentalRecovery', this.score);
    const bonusColor = result.isNewBest ? '#53d769' : '#8888aa';
    this.add.text(GAME_WIDTH / 2, 385, result.isNewBest
      ? `🆕 新紀錄！ 韌性 +${result.bonus}`
      : `韌性 +${result.bonus}`, {
      fontFamily: 'monospace', fontSize: '16px', color: bonusColor,
    }).setOrigin(0.5).setDepth(301);

    // 同步至後端
    if (result.delta > 0) {
      updatePlayerStats(PLAYER_ID, { logic: 0, tech: 0, charisma: 0, resilience: result.delta }).catch(() => {});
    }

    const retryBtn = createButton(this, GAME_WIDTH / 2 - 120, 450, '🔄 再挑戰', () => {
      this.scene.restart();
    }, 200, 48);
    retryBtn.setDepth(301);

    const backBtn = createButton(this, GAME_WIDTH / 2 + 120, 450, '← 返回選單', () => {
      this.scene.start('MainMenuScene');
    }, 200, 48);
    backBtn.setDepth(301);
  }

  private getRank(): string {
    if (this.score >= 120) return '🏆 S — 心靈導師';
    if (this.score >= 80) return '🥇 A — 堅定意志';
    if (this.score >= 50) return '🥈 B — 穩定心態';
    if (this.score >= 20) return '🥉 C — 初學者';
    return '📝 D — 需要更多修練';
  }
}
