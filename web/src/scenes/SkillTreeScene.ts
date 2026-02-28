import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT, PLAYER_ID } from '../ui/constants';
import { createButton, createPanel, createTitle, createLabel, createNavBar } from '../ui/components';
import { listSkills, unlockSkill, equipSkill, activateSkill } from '../services/api';
import type { SkillSummary } from '../services/types';

const SKILL_TYPE_META: Record<string, { label: string; color: number; icon: string }> = {
  Analysis: { label: '分析', color: Colors.statLogic, icon: '🔍' },
  Negotiation: { label: '談判', color: Colors.statCharisma, icon: '🗣️' },
  Defense: { label: '防禦', color: Colors.statResilience, icon: '🛡️' },
  Forensics: { label: '鑑識', color: Colors.statTech, icon: '🔬' },
};

export class SkillTreeScene extends Phaser.Scene {
  private skills: SkillSummary[] = [];
  private uiGroup!: Phaser.GameObjects.Group;

  constructor() {
    super('SkillTreeScene');
  }

  create(): void {
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    this.uiGroup = this.add.group();
    createTitle(this, GAME_WIDTH / 2, 40, '🌳 技能樹', '26px');
    createNavBar(this);

    this.loadSkills();
  }

  private async loadSkills(): Promise<void> {
    const loadText = this.add.text(GAME_WIDTH / 2, 300, '載入中...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);

    try {
      this.skills = await listSkills(PLAYER_ID);
      loadText.destroy();
      this.renderSkillTree();
    } catch (err) {
      loadText.setText(`載入失敗：${(err as Error).message}`);
      loadText.setColor('#e94560');
    }
  }

  private renderSkillTree(): void {
    this.uiGroup.clear(true, true);

    if (this.skills.length === 0) {
      this.add.text(GAME_WIDTH / 2, 300, '目前沒有可用技能', {
        fontFamily: 'monospace', fontSize: '16px', color: '#8888aa',
      }).setOrigin(0.5);
      return;
    }

    // Group skills by type
    const grouped = new Map<string, SkillSummary[]>();
    for (const skill of this.skills) {
      if (!grouped.has(skill.type)) grouped.set(skill.type, []);
      grouped.get(skill.type)!.push(skill);
    }

    let colIndex = 0;
    const colWidth = 220;
    const startX = GAME_WIDTH / 2 - ((grouped.size - 1) * colWidth) / 2;

    for (const [type, skills] of grouped) {
      const x = startX + colIndex * colWidth;
      const meta = SKILL_TYPE_META[type] || { label: type, color: Colors.text, icon: '⚙️' };

      // Column header
      const header = createLabel(this, x, 90, `${meta.icon} ${meta.label}`, `#${meta.color.toString(16).padStart(6, '0')}`, '18px');
      this.uiGroup.add(header);

      // Connecting line
      const lineColor = meta.color;
      if (skills.length > 1) {
        const line = this.add.line(0, 0, x, 110, x, 110 + (skills.length - 1) * 130, lineColor, 0.3)
          .setLineWidth(2).setOrigin(0);
        this.uiGroup.add(line);
      }

      skills.forEach((skill, i) => {
        const y = 160 + i * 130;
        this.renderSkillNode(x, y, skill, meta.color);
      });

      colIndex++;
    }
  }

  private renderSkillNode(x: number, y: number, skill: SkillSummary, color: number): void {
    const nodeSize = 180;
    const nodeHeight = 100;

    // Background card
    const bg = this.add.rectangle(x, y, nodeSize, nodeHeight, Colors.panel, 0.9)
      .setStrokeStyle(2, skill.unlocked ? color : 0x444466);
    this.uiGroup.add(bg);

    // Locked overlay
    if (!skill.unlocked) {
      const overlay = this.add.rectangle(x, y, nodeSize, nodeHeight, 0x000000, 0.4);
      this.uiGroup.add(overlay);
    }

    // Name
    const name = this.add.text(x, y - 28, skill.name, {
      fontFamily: 'monospace', fontSize: '15px',
      color: skill.unlocked ? '#eaeaea' : '#666688',
      align: 'center',
    }).setOrigin(0.5);
    this.uiGroup.add(name);

    // Status badges
    let statusText = '';
    if (!skill.unlocked) {
      statusText = `🔒 需要聲望 ${skill.reputationRequired}`;
    } else if (skill.equipped) {
      statusText = skill.cooldownRemaining > 0
        ? `⏳ 冷卻 ${skill.cooldownRemaining}s`
        : '✅ 已裝備（就緒）';
    } else {
      statusText = '已解鎖';
    }

    const status = this.add.text(x, y + 2, statusText, {
      fontFamily: 'monospace', fontSize: '12px',
      color: skill.unlocked ? '#53d769' : '#e94560',
      align: 'center',
    }).setOrigin(0.5);
    this.uiGroup.add(status);

    // Action button
    if (!skill.unlocked) {
      const btn = this.createSmallButton(x, y + 32, '解鎖', async () => {
        await unlockSkill(PLAYER_ID, skill.id);
        this.loadSkills();
      });
      this.uiGroup.add(btn);
    } else if (!skill.equipped) {
      const btn = this.createSmallButton(x, y + 32, '裝備', async () => {
        await equipSkill(PLAYER_ID, skill.id);
        this.loadSkills();
      });
      this.uiGroup.add(btn);
    } else if (skill.cooldownRemaining === 0) {
      const btn = this.createSmallButton(x, y + 32, '啟動', async () => {
        await activateSkill(PLAYER_ID, skill.id);
        this.loadSkills();
      });
      this.uiGroup.add(btn);
    }
  }

  private createSmallButton(x: number, y: number, label: string, onClick: () => void): Phaser.GameObjects.Container {
    return createButton(this, x, y, label, onClick, 80, 28);
  }
}
