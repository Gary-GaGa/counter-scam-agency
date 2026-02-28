import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT } from '../ui/constants';
import { createButton, createTitle } from '../ui/components';

export class MainMenuScene extends Phaser.Scene {
  constructor() {
    super('MainMenuScene');
  }

  create(): void {
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);

    createTitle(this, GAME_WIDTH / 2, 120, '🛡️ Counter Scam Agency', '36px');

    this.add.text(GAME_WIDTH / 2, 180, '— 反詐情報局 —', {
      fontFamily: 'monospace',
      fontSize: '20px',
      color: '#8888aa',
    }).setOrigin(0.5);

    this.add.text(GAME_WIDTH / 2, 230, '扮演反詐情報官，識破詐術、保護市民', {
      fontFamily: 'monospace',
      fontSize: '16px',
      color: '#eaeaea',
    }).setOrigin(0.5);

    const cx = GAME_WIDTH / 2;
    const leftCol = cx - 150;
    const rightCol = cx + 150;

    createButton(this, leftCol, 310, '📋 開始調查', () => {
      this.scene.start('CaseListScene');
    }, 260, 48);

    createButton(this, rightCol, 310, '🌳 技能樹', () => {
      this.scene.start('SkillTreeScene');
    }, 260, 48);

    createButton(this, leftCol, 380, '👤 角色狀態', () => {
      this.scene.start('ProfileScene');
    }, 260, 48);

    createButton(this, rightCol, 380, '🏰 防禦基地', () => {
      this.scene.start('BaseScene');
    }, 260, 48);

    // 小遊戲區塊
    this.add.text(cx, 440, '— 訓練小遊戲 —', {
      fontFamily: 'monospace', fontSize: '14px', color: '#8888aa',
    }).setOrigin(0.5);

    createButton(this, leftCol, 480, '⚡ 矛盾擊破', () => {
      this.scene.start('ContradictionScene');
    }, 260, 42);

    createButton(this, rightCol, 480, '🔌 訊號追蹤', () => {
      this.scene.start('SignalTraceScene');
    }, 260, 42);

    createButton(this, leftCol, 540, '🃏 談判牌局', () => {
      this.scene.start('NegotiationScene');
    }, 260, 42);

    createButton(this, rightCol, 540, '🧘 心靈調適', () => {
      this.scene.start('MentalRecoveryScene');
    }, 260, 42);

    // Version tag
    this.add.text(GAME_WIDTH - 16, GAME_HEIGHT - 16, 'MVP v0.1', {
      fontFamily: 'monospace',
      fontSize: '12px',
      color: '#555577',
    }).setOrigin(1, 1);
  }
}
