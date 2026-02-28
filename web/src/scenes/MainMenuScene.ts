import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT } from '../ui/constants';
import { createButton, createTitle } from '../ui/components';
import { gameState, type GameStateData } from '../services/gameState';
import { isUsingMock } from '../services/api';

export class MainMenuScene extends Phaser.Scene {
  private particles: Phaser.GameObjects.Rectangle[] = [];

  constructor() {
    super('MainMenuScene');
  }

  create(): void {
    const state = gameState.get();
    this.particles = [];

    // 深色漸層背景
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    // 頂部氛圍條
    this.add.rectangle(GAME_WIDTH / 2, 0, GAME_WIDTH, 4, Colors.accent).setOrigin(0.5, 0).setAlpha(0.6);

    // 背景浮動粒子
    for (let i = 0; i < 20; i++) {
      const x = Phaser.Math.Between(0, GAME_WIDTH);
      const y = Phaser.Math.Between(0, GAME_HEIGHT);
      const size = Phaser.Math.Between(1, 2);
      const p = this.add.rectangle(x, y, size, size, Colors.accent, Phaser.Math.FloatBetween(0.03, 0.12));
      this.particles.push(p);
      this.tweens.add({
        targets: p,
        y: p.y - Phaser.Math.Between(30, 80),
        alpha: 0,
        duration: Phaser.Math.Between(4000, 8000),
        repeat: -1,
        yoyo: true,
      });
    }

    // 淡入
    this.cameras.main.fadeIn(600, 0, 0, 0);

    // ─── 標題區 ───
    createTitle(this, GAME_WIDTH / 2, 52, '🛡️ Counter Scam Agency', '32px');
    this.add.text(GAME_WIDTH / 2, 88, '— 反詐情報局 —', {
      fontFamily: 'monospace', fontSize: '16px', color: '#8888aa',
    }).setOrigin(0.5);

    // ─── 玩家進度摘要 ───
    this.drawProgressBar(state);

    // ─── 主要功能按鈕 ───
    const cx = GAME_WIDTH / 2;
    const leftCol = cx - 155;
    const rightCol = cx + 155;

    createButton(this, leftCol, 225, '📋 開始調查', () => {
      this.transitionTo('CaseListScene');
    }, 270, 44);

    createButton(this, rightCol, 225, '🌳 技能樹', () => {
      this.transitionTo('SkillTreeScene');
    }, 270, 44);

    createButton(this, leftCol, 280, '👤 角色狀態', () => {
      this.transitionTo('ProfileScene');
    }, 270, 44);

    createButton(this, rightCol, 280, '🏰 防禦基地', () => {
      this.transitionTo('BaseScene');
    }, 270, 44);

    // ─── 訓練小遊戲區塊 ───
    this.add.line(0, 0, cx - 300, 320, cx + 300, 320, Colors.panelLight, 0.3).setOrigin(0);
    this.add.text(cx, 340, '⚔️ 訓練小遊戲', {
      fontFamily: 'monospace', fontSize: '15px', color: '#8888aa',
    }).setOrigin(0.5);

    const games: Array<{
      label: string;
      scene: string;
      key: keyof GameStateData['miniGameBest'];
      stat: string;
      color: number;
    }> = [
      { label: '⚡ 矛盾擊破', scene: 'ContradictionScene', key: 'contradiction', stat: '邏輯', color: Colors.statLogic },
      { label: '🔌 訊號追蹤', scene: 'SignalTraceScene', key: 'signalTrace', stat: '技術', color: Colors.statTech },
      { label: '🃏 談判牌局', scene: 'NegotiationScene', key: 'negotiation', stat: '交涉', color: Colors.statCharisma },
      { label: '🧘 心靈調適', scene: 'MentalRecoveryScene', key: 'mentalRecovery', stat: '韌性', color: Colors.statResilience },
    ];

    games.forEach((g, i) => {
      const col = i % 2 === 0 ? leftCol : rightCol;
      const row = Math.floor(i / 2);
      const y = 375 + row * 65;

      const btn = createButton(this, col, y, g.label, () => {
        this.transitionTo(g.scene);
      }, 270, 42);

      // 最高分標籤
      const best = state.miniGameBest[g.key];
      if (best) {
        const badge = this.add.text(col + 135, y - 21, `🏅 ${best.score}`, {
          fontFamily: 'monospace', fontSize: '11px', color: '#ffc107',
        }).setOrigin(1, 0.5);
      }

      // 對應屬性
      const statDot = this.add.rectangle(col - 135 + 8, y + 14, 6, 6, g.color).setAlpha(0.5);
      this.add.text(col - 135 + 16, y + 14, g.stat, {
        fontFamily: 'monospace', fontSize: '10px', color: '#666688',
      }).setOrigin(0, 0.5);
    });

    // ─── 底部資訊 ───
    this.add.text(GAME_WIDTH / 2, 538, '💡 完成小遊戲可提升屬性，累積聲望解鎖技能', {
      fontFamily: 'monospace', fontSize: '12px', color: '#555577',
    }).setOrigin(0.5);

    // 重播 Intro 按鈕
    const replayBtn = this.add.text(24, GAME_HEIGHT - 24, '🎬 重播序幕', {
      fontFamily: 'monospace', fontSize: '12px', color: '#444466',
    }).setOrigin(0, 1).setInteractive({ useHandCursor: true });
    replayBtn.on('pointerover', () => replayBtn.setColor('#e94560'));
    replayBtn.on('pointerout', () => replayBtn.setColor('#444466'));
    replayBtn.on('pointerdown', () => this.scene.start('IntroScene'));

    // 重置存檔
    const resetBtn = this.add.text(GAME_WIDTH - 24, GAME_HEIGHT - 24, '🗑️ 重置', {
      fontFamily: 'monospace', fontSize: '12px', color: '#444466',
    }).setOrigin(1, 1).setInteractive({ useHandCursor: true });
    resetBtn.on('pointerover', () => resetBtn.setColor('#e94560'));
    resetBtn.on('pointerout', () => resetBtn.setColor('#444466'));
    resetBtn.on('pointerdown', () => {
      gameState.reset();
      this.scene.restart();
    });

    // 版本 + 連線狀態
    const versionText = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT - 10, 'MVP v0.3 — 載入中...', {
      fontFamily: 'monospace', fontSize: '10px', color: '#333355',
    }).setOrigin(0.5, 1);

    isUsingMock().then(mock => {
      versionText.setText(mock ? 'MVP v0.3 — 離線模式 🔌' : 'MVP v0.3 — 已連線後端 🟢');
      versionText.setColor(mock ? '#666688' : '#53d769');
    });
  }

  private drawProgressBar(state: GameStateData): void {
    const cx = GAME_WIDTH / 2;
    const panelW = 620;
    const panelH = 60;
    const panelY = 132;

    // 進度面板背景
    this.add.rectangle(cx, panelY, panelW, panelH, Colors.panel, 0.85)
      .setStrokeStyle(1, Colors.panelLight, 0.4);

    // 聲望
    this.add.text(cx - 280, panelY - 14, `⭐ 聲望 ${state.reputation}`, {
      fontFamily: 'monospace', fontSize: '13px', color: '#ffc107',
    }).setOrigin(0, 0.5);

    // 已完成案件
    this.add.text(cx - 280, panelY + 10, `📁 案件 ${state.completedMissions.length}/2`, {
      fontFamily: 'monospace', fontSize: '12px', color: '#8888aa',
    }).setOrigin(0, 0.5);

    // 四維屬性
    const stats = [
      { label: '邏輯', value: state.stats.logic, color: Colors.statLogic },
      { label: '技術', value: state.stats.tech, color: Colors.statTech },
      { label: '交涉', value: state.stats.charisma, color: Colors.statCharisma },
      { label: '韌性', value: state.stats.resilience, color: Colors.statResilience },
    ];

    const statStartX = cx + 10;
    stats.forEach((s, i) => {
      const x = statStartX + i * 70;
      this.add.rectangle(x, panelY - 10, 8, 8, s.color).setAlpha(0.8);
      this.add.text(x + 8, panelY - 10, s.label, {
        fontFamily: 'monospace', fontSize: '11px', color: '#8888aa',
      }).setOrigin(0, 0.5);
      this.add.text(x + 8, panelY + 8, `${s.value}`, {
        fontFamily: 'monospace', fontSize: '14px', color: '#eaeaea',
      }).setOrigin(0, 0.5);
    });
  }

  private transitionTo(sceneName: string): void {
    this.cameras.main.fadeOut(300, 0, 0, 0);
    this.cameras.main.once('camerafadeoutcomplete', () => {
      this.scene.start(sceneName);
    });
  }
}
