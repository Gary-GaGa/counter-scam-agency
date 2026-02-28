import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT } from '../ui/constants';
import { createButton, createPanel, createTitle, createBodyText, createLabel } from '../ui/components';
import { getMission, startInvestigation, advanceNode, completeInvestigation } from '../services/api';
import type { MissionDetail, NarrativeNode, NarrativeOption, CompleteResult } from '../services/types';

interface SceneData {
  missionId: string;
  playerId: string;
}

export class InvestigationScene extends Phaser.Scene {
  private missionId!: string;
  private playerId!: string;
  private investigationId!: string;
  private currentNodeId!: string;
  private mission!: MissionDetail;
  private uiGroup!: Phaser.GameObjects.Group;

  constructor() {
    super('InvestigationScene');
  }

  init(data: SceneData): void {
    this.missionId = data.missionId;
    this.playerId = data.playerId;
  }

  create(): void {
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    this.uiGroup = this.add.group();

    const loadingText = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2, '載入案件中...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);

    getMission(this.missionId)
      .then(async (mission) => {
        this.mission = mission;
        loadingText.destroy();
        this.showBriefing();
      })
      .catch((err) => {
        loadingText.setText(`載入失敗：${err.message}`);
        loadingText.setColor('#e94560');
      });
  }

  private clearUI(): void {
    this.uiGroup.clear(true, true);
  }

  private showBriefing(): void {
    this.clearUI();

    createPanel(this, GAME_WIDTH / 2, GAME_HEIGHT / 2, 800, 500);

    const title = createTitle(this, GAME_WIDTH / 2, 100, `📁 ${this.mission.title}`, '24px');
    this.uiGroup.add(title);

    const typeLabel: Record<string, string> = {
      Phishing: '網路釣魚', Investment: '投資詐騙', Romance: '感情詐騙', Impersonation: '冒名詐騙',
    };
    const meta = createLabel(this, GAME_WIDTH / 2, 140,
      `類型：${typeLabel[this.mission.type] || this.mission.type}  難度：${'★'.repeat(this.mission.difficulty)}  聲望權重：${this.mission.reputationWeight}`);
    this.uiGroup.add(meta);

    if (this.mission.description) {
      const desc = createBodyText(this, GAME_WIDTH / 2, 170, this.mission.description, 650);
      this.uiGroup.add(desc);
    }

    // Victim profile
    if (this.mission.victimProfile) {
      const vp = this.mission.victimProfile;
      const vpTitle = createLabel(this, GAME_WIDTH / 2, 230, '— 受害者心理側寫 —', '#e94560', '16px');
      this.uiGroup.add(vpTitle);

      const bars = [
        { label: '焦慮', value: vp.anxiety, color: Colors.danger },
        { label: '信任', value: vp.trust, color: Colors.statCharisma },
        { label: '急迫', value: vp.urgency, color: Colors.warning },
        { label: '孤立', value: vp.isolation, color: Colors.statResilience },
      ];

      bars.forEach((bar, i) => {
        const y = 265 + i * 30;
        const labelText = this.add.text(300, y, bar.label, {
          fontFamily: 'monospace', fontSize: '14px', color: '#8888aa',
        }).setOrigin(1, 0.5);
        this.uiGroup.add(labelText);

        const barBg = this.add.rectangle(310 + 150, y, 300, 16, 0x333355).setOrigin(0.5);
        this.uiGroup.add(barBg);

        const barFill = this.add.rectangle(310, y, 3 * bar.value, 16, bar.color)
          .setOrigin(0, 0.5);
        this.uiGroup.add(barFill);

        const valText = this.add.text(620, y, `${bar.value}`, {
          fontFamily: 'monospace', fontSize: '13px', color: '#eaeaea',
        }).setOrigin(0, 0.5);
        this.uiGroup.add(valText);
      });

      const risk = createLabel(this, GAME_WIDTH / 2, 390,
        `風險等級：${vp.riskLevel}（${vp.riskScore} 分）`, '#e94560', '15px');
      this.uiGroup.add(risk);
    }

    const btn = createButton(this, GAME_WIDTH / 2, 460, '▶ 開始調查', () => {
      this.beginInvestigation();
    }, 200, 48);
    this.uiGroup.add(btn);
  }

  private async beginInvestigation(): Promise<void> {
    this.clearUI();
    const loadText = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2, '開始調查...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);
    this.uiGroup.add(loadText);

    try {
      this.investigationId = `inv-${Date.now()}`;
      const result = await startInvestigation(this.investigationId, this.playerId, this.missionId);
      loadText.destroy();
      const startNodeId = result.currentNodeId || this.mission.nodes[0]?.id;
      this.showNode(startNodeId);
    } catch (err) {
      loadText.setText(`開始調查失敗：${(err as Error).message}`);
      loadText.setColor('#e94560');
    }
  }

  private showNode(nodeId: string): void {
    this.clearUI();
    this.currentNodeId = nodeId;

    const node = this.mission.nodes.find((n) => n.id === nodeId);
    if (!node) {
      this.showError('找不到節點');
      return;
    }

    createPanel(this, GAME_WIDTH / 2, GAME_HEIGHT / 2, 820, 540);

    const title = createTitle(this, GAME_WIDTH / 2, 80, `【${node.title}】`, '22px');
    this.uiGroup.add(title);

    const body = createBodyText(this, GAME_WIDTH / 2, 120, node.body, 700);
    this.uiGroup.add(body);

    if (node.isTerminal || node.options.length === 0) {
      this.showTerminalActions(node);
      return;
    }

    // Options
    const optStartY = Math.min(body.y + body.displayHeight + 40, 350);
    node.options.forEach((option, i) => {
      const btn = createButton(this, GAME_WIDTH / 2, optStartY + i * 60, `${i + 1}. ${option.label}`, () => {
        this.selectOption(option);
      }, 650, 46);
      this.uiGroup.add(btn);
    });
  }

  private async selectOption(option: NarrativeOption): Promise<void> {
    this.clearUI();
    const loadText = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2, '推進中...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);
    this.uiGroup.add(loadText);

    // Show evidence gained
    if (option.evidenceIds && option.evidenceIds.length > 0) {
      const evidence = this.mission.evidenceList.filter(
        (e) => option.evidenceIds.includes(e.id),
      );
      if (evidence.length > 0) {
        loadText.destroy();
        this.showEvidenceGained(evidence, option);
        return;
      }
    }

    try {
      const result = await advanceNode(this.investigationId, this.currentNodeId, option.id);
      loadText.destroy();
      if (result.status !== 'Active' && result.status !== 'active') {
        this.finishInvestigation();
      } else {
        this.showNode(result.nextNodeId || option.nextNodeId);
      }
    } catch {
      loadText.destroy();
      this.showNode(option.nextNodeId);
    }
  }

  private showEvidenceGained(evidence: { id: string; description: string; type: string }[], option: NarrativeOption): void {
    this.clearUI();
    createPanel(this, GAME_WIDTH / 2, GAME_HEIGHT / 2, 600, 300);

    const title = createTitle(this, GAME_WIDTH / 2, GAME_HEIGHT / 2 - 90, '📋 取得證據', '22px');
    this.uiGroup.add(title);

    evidence.forEach((e, i) => {
      const y = GAME_HEIGHT / 2 - 40 + i * 35;
      const text = this.add.text(GAME_WIDTH / 2, y, `• ${e.description} (${e.type})`, {
        fontFamily: 'monospace', fontSize: '15px', color: '#53d769',
      }).setOrigin(0.5);
      this.uiGroup.add(text);
    });

    const btn = createButton(this, GAME_WIDTH / 2, GAME_HEIGHT / 2 + 90, '繼續', async () => {
      try {
        const result = await advanceNode(this.investigationId, this.currentNodeId, option.id);
        if (result.status !== 'Active' && result.status !== 'active') {
          this.finishInvestigation();
        } else {
          this.showNode(result.nextNodeId || option.nextNodeId);
        }
      } catch {
        this.showNode(option.nextNodeId);
      }
    }, 160, 44);
    this.uiGroup.add(btn);
  }

  private showTerminalActions(_node: NarrativeNode): void {
    const btn = createButton(this, GAME_WIDTH / 2, 480, '📝 結案', () => {
      this.finishInvestigation();
    }, 200, 48);
    this.uiGroup.add(btn);
  }

  private async finishInvestigation(): Promise<void> {
    this.clearUI();
    const loadText = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2, '結案中...', {
      fontFamily: 'monospace', fontSize: '18px', color: '#8888aa',
    }).setOrigin(0.5);
    this.uiGroup.add(loadText);

    try {
      const result = await completeInvestigation(this.investigationId);
      loadText.destroy();
      this.showSummary(result);
    } catch (err) {
      loadText.setText(`結案失敗：${(err as Error).message}`);
      loadText.setColor('#e94560');
      const btn = createButton(this, GAME_WIDTH / 2, GAME_HEIGHT / 2 + 60, '返回案件列表', () => {
        this.scene.start('CaseListScene');
      }, 200, 44);
      this.uiGroup.add(btn);
    }
  }

  private showSummary(result: CompleteResult): void {
    this.clearUI();
    createPanel(this, GAME_WIDTH / 2, GAME_HEIGHT / 2, 600, 350);

    if (result.success) {
      const icon = createTitle(this, GAME_WIDTH / 2, 200, '✅ 結案成功！', '28px');
      this.uiGroup.add(icon);
    } else {
      const icon = this.add.text(GAME_WIDTH / 2, 200, '❌ 結案失敗', {
        fontFamily: 'monospace', fontSize: '28px', color: '#e94560',
      }).setOrigin(0.5);
      this.uiGroup.add(icon);
    }

    const repText = this.add.text(GAME_WIDTH / 2, 260, `聲望變化：+${result.reputationGained}`, {
      fontFamily: 'monospace', fontSize: '20px',
      color: result.reputationGained > 0 ? '#53d769' : '#8888aa',
    }).setOrigin(0.5);
    this.uiGroup.add(repText);

    const evText = this.add.text(GAME_WIDTH / 2, 300, `蒐集證據：${result.evidenceCollected} 件`, {
      fontFamily: 'monospace', fontSize: '16px', color: '#bbbbcc',
    }).setOrigin(0.5);
    this.uiGroup.add(evText);

    const backBtn = createButton(this, GAME_WIDTH / 2 - 120, 400, '返回案件列表', () => {
      this.scene.start('CaseListScene');
    }, 200, 44);
    this.uiGroup.add(backBtn);

    const skillBtn = createButton(this, GAME_WIDTH / 2 + 120, 400, '查看技能樹', () => {
      this.scene.start('SkillTreeScene');
    }, 200, 44);
    this.uiGroup.add(skillBtn);
  }

  private showError(msg: string): void {
    this.clearUI();
    this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2, msg, {
      fontFamily: 'monospace', fontSize: '18px', color: '#e94560',
    }).setOrigin(0.5);
    createButton(this, GAME_WIDTH / 2, GAME_HEIGHT / 2 + 60, '返回', () => {
      this.scene.start('CaseListScene');
    }, 160, 44);
  }
}
