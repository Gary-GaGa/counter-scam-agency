import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT, PLAYER_ID } from '../ui/constants';
import { createPanel, createTitle, createLabel, createNavBar } from '../ui/components';
import { getPlayer } from '../services/api';
import type { PlayerSummary } from '../services/types';

export class ProfileScene extends Phaser.Scene {
  private uiGroup!: Phaser.GameObjects.Group;

  constructor() {
    super('ProfileScene');
  }

  create(): void {
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    this.uiGroup = this.add.group();
    createTitle(this, GAME_WIDTH / 2, 40, '👤 角色狀態', '26px');
    createNavBar(this);

    this.loadProfile();
  }

  private async loadProfile(): Promise<void> {
    const loadText = this.add.text(GAME_WIDTH / 2, 300, '載入中...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);

    try {
      const player = await getPlayer(PLAYER_ID);
      loadText.destroy();
      this.renderProfile(player);
    } catch (err) {
      loadText.setText(`載入失敗：${(err as Error).message}`);
      loadText.setColor('#e94560');
    }
  }

  private renderProfile(player: PlayerSummary): void {
    this.uiGroup.clear(true, true);

    // Main info panel
    createPanel(this, GAME_WIDTH / 2, 200, 500, 200);

    const idLabel = createLabel(this, GAME_WIDTH / 2, 120, `情報官 ${player.id}`, '#eaeaea', '20px');
    this.uiGroup.add(idLabel);

    const repLabel = createLabel(this, GAME_WIDTH / 2, 155,
      `🏆 聲望：${player.reputation}`, '#ffc107', '18px');
    this.uiGroup.add(repLabel);

    if (player.partnerPersonality) {
      const personality = createLabel(this, GAME_WIDTH / 2, 185,
        `🤖 AI 性格：${player.partnerPersonality}`, '#8888aa', '15px');
      this.uiGroup.add(personality);
    }

    // Stats radar-like display
    createPanel(this, GAME_WIDTH / 2, 420, 600, 250);
    const statsTitle = createLabel(this, GAME_WIDTH / 2, 320, '— 能力值 —', '#e94560', '18px');
    this.uiGroup.add(statsTitle);

    const stats = player.stats;
    const totalStats = player.totalStats || stats;

    const statEntries = [
      { label: '邏輯 (Logic)', base: stats.logic, total: totalStats.logic, color: Colors.statLogic },
      { label: '技術 (Tech)', base: stats.tech, total: totalStats.tech, color: Colors.statTech },
      { label: '交涉 (Charisma)', base: stats.charisma, total: totalStats.charisma, color: Colors.statCharisma },
      { label: '韌性 (Resilience)', base: stats.resilience, total: totalStats.resilience, color: Colors.statResilience },
    ];

    statEntries.forEach((stat, i) => {
      const y = 360 + i * 40;
      const labelX = 250;
      const barX = 430;
      const barWidth = 200;
      const maxVal = 30;

      // Label
      const label = this.add.text(labelX, y, stat.label, {
        fontFamily: 'monospace', fontSize: '14px', color: '#eaeaea',
      }).setOrigin(0, 0.5);
      this.uiGroup.add(label);

      // Bar background
      const barBg = this.add.rectangle(barX + barWidth / 2, y, barWidth, 18, 0x333355);
      this.uiGroup.add(barBg);

      // Base stat bar
      const baseWidth = Math.min(stat.base / maxVal, 1) * barWidth;
      const baseFill = this.add.rectangle(barX, y, baseWidth, 18, stat.color)
        .setOrigin(0, 0.5).setAlpha(0.5);
      this.uiGroup.add(baseFill);

      // Total stat bar (with bonus)
      const totalWidth = Math.min(stat.total / maxVal, 1) * barWidth;
      const totalFill = this.add.rectangle(barX, y, totalWidth, 18, stat.color)
        .setOrigin(0, 0.5).setAlpha(0.8);
      this.uiGroup.add(totalFill);

      // Value text
      const bonus = stat.total - stat.base;
      const valStr = bonus > 0 ? `${stat.base} (+${bonus})` : `${stat.base}`;
      const valText = this.add.text(barX + barWidth + 12, y, valStr, {
        fontFamily: 'monospace', fontSize: '13px', color: '#eaeaea',
      }).setOrigin(0, 0.5);
      this.uiGroup.add(valText);
    });

    // Equipped info
    if (player.equippedModules && player.equippedModules.length > 0) {
      const modLabel = createLabel(this, GAME_WIDTH / 2, 540,
        `裝備模組：${player.equippedModules.join(', ')}`, '#8888aa', '13px');
      this.uiGroup.add(modLabel);
    }
  }
}
