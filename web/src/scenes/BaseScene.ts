import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT, PLAYER_ID } from '../ui/constants';
import { createButton, createPanel, createTitle, createLabel, createNavBar } from '../ui/components';
import { getBase, createBase, addFacility, upgradeFacility, upgradeSecurityLevel } from '../services/api';
import type { BaseSummary, FacilitySummary } from '../services/types';

const BASE_ID = `base-${PLAYER_ID}`;

const FACILITY_TEMPLATES = [
  { type: 'Firewall', name: '防火牆', icon: '🔥', description: '阻擋惡意流量與攻擊嘗試' },
  { type: 'SIEM', name: '安全監控中心', icon: '📡', description: '即時偵測異常行為與入侵事件' },
  { type: 'Training', name: '訓練中心', icon: '🎓', description: '提升團隊防詐意識與反應能力' },
];

export class BaseScene extends Phaser.Scene {
  private base: BaseSummary | null = null;
  private uiGroup!: Phaser.GameObjects.Group;

  constructor() {
    super('BaseScene');
  }

  create(): void {
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    this.uiGroup = this.add.group();
    createTitle(this, GAME_WIDTH / 2, 40, '🏰 數位防禦基地', '26px');
    createNavBar(this);

    this.loadBase();
  }

  private async loadBase(): Promise<void> {
    const loadText = this.add.text(GAME_WIDTH / 2, 300, '載入基地...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);

    try {
      this.base = await getBase(BASE_ID);
      loadText.destroy();
      this.renderBase();
    } catch {
      // 基地不存在，自動建立
      try {
        this.base = await createBase(BASE_ID, PLAYER_ID, 6);
        loadText.destroy();
        this.renderBase();
      } catch (err) {
        loadText.setText(`基地載入失敗：${(err as Error).message}`);
        loadText.setColor('#e94560');
      }
    }
  }

  private clearUI(): void {
    this.uiGroup.clear(true, true);
  }

  private renderBase(): void {
    this.clearUI();
    if (!this.base) return;

    // 安全等級面板
    const secPanel = createPanel(this, GAME_WIDTH / 2, 120, 600, 80);
    this.uiGroup.add(secPanel);

    const secLabel = createLabel(this, GAME_WIDTH / 2 - 100, 110,
      `🛡️ 安全等級`, '#eaeaea', '18px');
    this.uiGroup.add(secLabel);

    const secLevel = this.add.text(GAME_WIDTH / 2 + 60, 110,
      `Lv.${this.base.securityLevel}`, {
        fontFamily: 'monospace', fontSize: '24px', color: '#53d769',
      }).setOrigin(0.5);
    this.uiGroup.add(secLevel);

    const secBar = this.add.rectangle(GAME_WIDTH / 2, 140, 400, 12, 0x333355);
    this.uiGroup.add(secBar);
    const secFill = this.add.rectangle(
      GAME_WIDTH / 2 - 200, 140,
      Math.min(this.base.securityLevel / 10, 1) * 400, 12,
      Colors.success,
    ).setOrigin(0, 0.5);
    this.uiGroup.add(secFill);

    const upgradeSecBtn = createButton(this, GAME_WIDTH / 2 + 240, 120, '⬆ 升級', async () => {
      try {
        this.base = await upgradeSecurityLevel(BASE_ID, 10);
        this.renderBase();
      } catch (err) {
        this.showToast(`升級失敗：${(err as Error).message}`);
      }
    }, 100, 36);
    this.uiGroup.add(upgradeSecBtn);

    // 設施區域
    const slotsLabel = createLabel(this, GAME_WIDTH / 2, 180,
      `設施欄位：${this.base.facilities.length} / ${this.base.facilitySlots}`,
      '#8888aa', '15px');
    this.uiGroup.add(slotsLabel);

    // 渲染現有設施
    const facilities = this.base.facilities;
    const colWidth = 280;
    const startX = GAME_WIDTH / 2 - ((Math.min(facilities.length, 3) - 1) * colWidth) / 2;

    facilities.forEach((facility, i) => {
      const col = i % 3;
      const row = Math.floor(i / 3);
      const x = startX + col * colWidth;
      const y = 270 + row * 160;
      this.renderFacilityCard(x, y, facility);
    });

    // 新增設施按鈕
    if (this.base.facilities.length < this.base.facilitySlots) {
      const addY = 270 + Math.floor(facilities.length / 3) * 160 +
        (facilities.length % 3 === 0 ? 0 : 160);
      const addBtnY = Math.min(addY, 530);

      const addBtn = createButton(this, GAME_WIDTH / 2, addBtnY, '➕ 新增設施', () => {
        this.showAddFacilityMenu();
      }, 200, 44);
      this.uiGroup.add(addBtn);
    }
  }

  private renderFacilityCard(x: number, y: number, facility: FacilitySummary): void {
    const cardW = 250;
    const cardH = 130;
    const template = FACILITY_TEMPLATES.find(t => t.type === facility.type);
    const icon = template?.icon || '⚙️';

    const bg = this.add.rectangle(x, y, cardW, cardH, Colors.panel, 0.9)
      .setStrokeStyle(2, Colors.panelLight);
    this.uiGroup.add(bg);

    const name = this.add.text(x, y - 40, `${icon} ${facility.name}`, {
      fontFamily: 'monospace', fontSize: '16px', color: '#eaeaea',
    }).setOrigin(0.5);
    this.uiGroup.add(name);

    const levelText = this.add.text(x, y - 15, `Lv.${facility.level} / ${facility.maxLevel}`, {
      fontFamily: 'monospace', fontSize: '14px', color: '#53d769',
    }).setOrigin(0.5);
    this.uiGroup.add(levelText);

    // 等級進度條
    const barW = 160;
    const barBg = this.add.rectangle(x, y + 8, barW, 10, 0x333355);
    this.uiGroup.add(barBg);
    const fillW = facility.maxLevel > 0 ? (facility.level / facility.maxLevel) * barW : 0;
    const barFill = this.add.rectangle(x - barW / 2, y + 8, fillW, 10, Colors.accent)
      .setOrigin(0, 0.5);
    this.uiGroup.add(barFill);

    // 升級按鈕
    if (facility.level < facility.maxLevel) {
      const upBtn = createButton(this, x, y + 40, '升級', async () => {
        try {
          this.base = await upgradeFacility(BASE_ID, facility.id);
          this.renderBase();
        } catch (err) {
          this.showToast(`升級失敗：${(err as Error).message}`);
        }
      }, 80, 30);
      this.uiGroup.add(upBtn);
    } else {
      const maxLabel = this.add.text(x, y + 40, '已滿級', {
        fontFamily: 'monospace', fontSize: '13px', color: '#ffc107',
      }).setOrigin(0.5);
      this.uiGroup.add(maxLabel);
    }
  }

  private showAddFacilityMenu(): void {
    this.clearUI();
    const existingTypes = this.base?.facilities.map(f => f.type) || [];

    const title = createTitle(this, GAME_WIDTH / 2, 100, '選擇要建設的設施', '22px');
    this.uiGroup.add(title);

    const available = FACILITY_TEMPLATES.filter(t => !existingTypes.includes(t.type));

    if (available.length === 0) {
      const msg = this.add.text(GAME_WIDTH / 2, 250, '所有設施類型已建設完成', {
        fontFamily: 'monospace', fontSize: '16px', color: '#8888aa',
      }).setOrigin(0.5);
      this.uiGroup.add(msg);
    } else {
      available.forEach((tmpl, i) => {
        const y = 180 + i * 120;
        const panel = createPanel(this, GAME_WIDTH / 2, y, 500, 90);
        this.uiGroup.add(panel);

        const name = this.add.text(GAME_WIDTH / 2 - 140, y - 20,
          `${tmpl.icon} ${tmpl.name}`, {
            fontFamily: 'monospace', fontSize: '18px', color: '#eaeaea',
          }).setOrigin(0, 0.5);
        this.uiGroup.add(name);

        const desc = this.add.text(GAME_WIDTH / 2 - 140, y + 10, tmpl.description, {
          fontFamily: 'monospace', fontSize: '13px', color: '#8888aa',
        }).setOrigin(0, 0.5);
        this.uiGroup.add(desc);

        const btn = createButton(this, GAME_WIDTH / 2 + 180, y, '建設', async () => {
          try {
            const facilityId = `${tmpl.type.toLowerCase()}-1`;
            this.base = await addFacility(BASE_ID, {
              id: facilityId,
              type: tmpl.type,
              name: tmpl.name,
              level: 1,
              maxLevel: 5,
              description: tmpl.description,
            });
            this.renderBase();
          } catch (err) {
            this.showToast(`建設失敗：${(err as Error).message}`);
          }
        }, 100, 36);
        this.uiGroup.add(btn);
      });
    }

    const backBtn = createButton(this, GAME_WIDTH / 2, 530, '← 返回基地', () => {
      this.renderBase();
    }, 160, 40);
    this.uiGroup.add(backBtn);
  }

  private showToast(msg: string): void {
    const toast = this.add.text(GAME_WIDTH / 2, 570, msg, {
      fontFamily: 'monospace', fontSize: '14px', color: '#e94560',
      backgroundColor: '#1a1a2e',
      padding: { x: 12, y: 6 },
    }).setOrigin(0.5).setDepth(200);

    this.time.delayedCall(2500, () => toast.destroy());
  }
}
