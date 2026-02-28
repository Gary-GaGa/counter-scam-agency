import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT, PLAYER_ID } from '../ui/constants';
import { createButton, createTitle, createPanel, createNavBar } from '../ui/components';
import { listMissions } from '../services/api';
import type { MissionSummary } from '../services/types';

export class CaseListScene extends Phaser.Scene {
  constructor() {
    super('CaseListScene');
  }

  create(): void {
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    createTitle(this, GAME_WIDTH / 2, 40, '📋 案件面板', '26px');
    createNavBar(this);

    const loadingText = this.add.text(GAME_WIDTH / 2, 300, '載入中...', {
      fontFamily: 'monospace',
      fontSize: '18px',
      color: '#8888aa',
    }).setOrigin(0.5);

    listMissions()
      .then((missions) => {
        loadingText.destroy();
        this.renderMissions(missions);
      })
      .catch((err) => {
        loadingText.setText(`載入失敗：${err.message}`);
        loadingText.setColor('#e94560');
      });
  }

  private renderMissions(missions: MissionSummary[]): void {
    if (missions.length === 0) {
      this.add.text(GAME_WIDTH / 2, 300, '沒有可用案件\n請先執行 go run cmd/seed-missions/main.go', {
        fontFamily: 'monospace',
        fontSize: '16px',
        color: '#8888aa',
        align: 'center',
      }).setOrigin(0.5);
      return;
    }

    const startY = 90;
    const cardHeight = 120;
    const cardWidth = 700;
    const gap = 16;

    missions.forEach((mission, i) => {
      const y = startY + i * (cardHeight + gap) + cardHeight / 2;

      createPanel(this, GAME_WIDTH / 2, y, cardWidth, cardHeight);

      // Title
      this.add.text(GAME_WIDTH / 2 - cardWidth / 2 + 24, y - 36, mission.title, {
        fontFamily: 'monospace',
        fontSize: '20px',
        color: '#e94560',
      });

      // Metadata
      const typeLabel = this.scamTypeLabel(mission.type);
      const stars = '★'.repeat(mission.difficulty) + '☆'.repeat(5 - mission.difficulty);
      this.add.text(GAME_WIDTH / 2 - cardWidth / 2 + 24, y - 8, `類型：${typeLabel}  難度：${stars}  聲望權重：${mission.reputationWeight}`, {
        fontFamily: 'monospace',
        fontSize: '14px',
        color: '#8888aa',
      });

      // Description
      if (mission.description) {
        const desc = mission.description.length > 60
          ? mission.description.substring(0, 60) + '…'
          : mission.description;
        this.add.text(GAME_WIDTH / 2 - cardWidth / 2 + 24, y + 16, desc, {
          fontFamily: 'monospace',
          fontSize: '14px',
          color: '#bbbbcc',
        });
      }

      // Start button
      createButton(this, GAME_WIDTH / 2 + cardWidth / 2 - 80, y, '開始調查', () => {
        this.scene.start('InvestigationScene', {
          missionId: mission.id,
          playerId: PLAYER_ID,
        });
      }, 120, 40);
    });
  }

  private scamTypeLabel(type: string): string {
    const map: Record<string, string> = {
      Phishing: '網路釣魚',
      Investment: '投資詐騙',
      Romance: '感情詐騙',
      Impersonation: '冒名詐騙',
    };
    return map[type] || type;
  }
}
