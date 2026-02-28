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

    createButton(this, cx, 460, '⚡ 矛盾擊破（小遊戲）', () => {
      this.scene.start('ContradictionScene');
    }, 300, 52);

    // Version tag
    this.add.text(GAME_WIDTH - 16, GAME_HEIGHT - 16, 'MVP v0.1', {
      fontFamily: 'monospace',
      fontSize: '12px',
      color: '#555577',
    }).setOrigin(1, 1);
  }
}
