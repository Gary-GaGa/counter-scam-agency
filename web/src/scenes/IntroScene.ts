import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT } from '../ui/constants';
import { gameState } from '../services/gameState';

/**
 * 開場引導場景 — 打字機敘事 + 遊戲世界觀介紹。
 * 首次進入自動播放，之後可跳過。
 */

const INTRO_LINES = [
  '2026 年，臺灣。',
  '',
  '詐騙集團的手法日益精進，',
  '從假客服、投資群組到感情圈套…',
  '',
  '每天都有人失去積蓄、信任，甚至尊嚴。',
  '',
  '政府成立了一個秘密單位——',
  '',
  '「反詐情報局」',
  'Counter Scam Agency',
  '',
  '你，被選為新一代的情報官。',
  '',
  '你的任務：識破詐術，保護市民，',
  '讓每一個人都能免於恐懼地生活。',
  '',
  '準備好了嗎？',
];

export class IntroScene extends Phaser.Scene {
  private lineIndex = 0;
  private charIndex = 0;
  private textObjects: Phaser.GameObjects.Text[] = [];
  private currentText!: Phaser.GameObjects.Text;
  private typeTimer!: Phaser.Time.TimerEvent;
  private skipText!: Phaser.GameObjects.Text;
  private particles: Phaser.GameObjects.Rectangle[] = [];

  constructor() {
    super('IntroScene');
  }

  create(): void {
    this.lineIndex = 0;
    this.charIndex = 0;
    this.textObjects = [];
    this.particles = [];

    // 全黑背景
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, 0x0a0a16);

    // 動態背景粒子
    for (let i = 0; i < 30; i++) {
      const x = Phaser.Math.Between(0, GAME_WIDTH);
      const y = Phaser.Math.Between(0, GAME_HEIGHT);
      const size = Phaser.Math.Between(1, 3);
      const p = this.add.rectangle(x, y, size, size, Colors.accent, Phaser.Math.FloatBetween(0.05, 0.2));
      this.particles.push(p);
      this.tweens.add({
        targets: p,
        y: p.y - Phaser.Math.Between(20, 60),
        alpha: 0,
        duration: Phaser.Math.Between(3000, 6000),
        repeat: -1,
        yoyo: true,
      });
    }

    // 跳過按鈕
    this.skipText = this.add.text(GAME_WIDTH - 30, GAME_HEIGHT - 30, '按任意鍵跳過 →', {
      fontFamily: 'monospace', fontSize: '14px', color: '#555577',
    }).setOrigin(1, 1).setAlpha(0);

    this.tweens.add({
      targets: this.skipText, alpha: 0.6, delay: 2000, duration: 1000,
    });

    // 開始打字
    this.startTyping();

    // 點擊或按鍵跳過
    this.input.on('pointerdown', () => this.skipIntro());
    this.input.keyboard?.on('keydown', () => this.skipIntro());
  }

  private startTyping(): void {
    this.typeNextLine();
  }

  private typeNextLine(): void {
    if (this.lineIndex >= INTRO_LINES.length) {
      this.finishIntro();
      return;
    }

    const line = INTRO_LINES[this.lineIndex];
    const y = 120 + this.lineIndex * 28;

    if (line === '') {
      this.lineIndex++;
      this.time.delayedCall(300, () => this.typeNextLine());
      return;
    }

    // 特殊行：標題
    const isTitle = line === '「反詐情報局」' || line === 'Counter Scam Agency';
    const fontSize = isTitle ? '28px' : '17px';
    const color = isTitle ? '#e94560' : '#ccccdd';

    this.currentText = this.add.text(GAME_WIDTH / 2, y, '', {
      fontFamily: 'monospace', fontSize, color,
    }).setOrigin(0.5, 0).setAlpha(0);

    this.textObjects.push(this.currentText);
    this.charIndex = 0;

    this.tweens.add({ targets: this.currentText, alpha: 1, duration: 200 });

    const speed = isTitle ? 80 : 45;
    this.typeTimer = this.time.addEvent({
      delay: speed,
      callback: () => {
        if (this.charIndex < line.length) {
          this.currentText.setText(line.substring(0, this.charIndex + 1));
          this.charIndex++;
        } else {
          this.typeTimer.destroy();
          this.lineIndex++;
          this.time.delayedCall(isTitle ? 800 : 400, () => this.typeNextLine());
        }
      },
      loop: true,
    });
  }

  private finishIntro(): void {
    this.time.delayedCall(1200, () => {
      this.cameras.main.fadeOut(1000, 0, 0, 0);
      this.cameras.main.once('camerafadeoutcomplete', () => {
        gameState.markIntroSeen();
        this.scene.start('MainMenuScene');
      });
    });
  }

  private skipIntro(): void {
    this.typeTimer?.destroy();
    gameState.markIntroSeen();
    this.scene.start('MainMenuScene');
  }
}
