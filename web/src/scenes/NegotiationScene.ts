import Phaser from 'phaser';
import { Colors, GAME_WIDTH, GAME_HEIGHT, PLAYER_ID } from '../ui/constants';
import { createButton, createTitle, createPanel } from '../ui/components';
import { gameState } from '../services/gameState';
import { updatePlayerStats } from '../services/api';

/**
 * 談判牌局（Negotiation Cards）
 * 卡牌對戰玩法，對應 Charisma 數值。
 *
 * 玩法：玩家有手牌（同理心、數據佐證、法律威攝、情感呼籲）。
 * 對手（詐騙者）打出「固執牌」（否認、威脅、情緒勒索、拖延）。
 * 每張牌有相剋關係，出正確的牌可以削減對方固執值。
 * 固執值歸零 → 勝利；玩家信任牌用完 → 失敗。
 */

interface Card {
  id: string;
  name: string;
  icon: string;
  type: 'player' | 'enemy';
  power: number;
  beats: string[]; // 剋制的敵方牌 type
}

const PLAYER_CARDS: Card[] = [
  { id: 'empathy', name: '同理心', icon: '💛', type: 'player', power: 15, beats: ['emotional'] },
  { id: 'evidence', name: '數據佐證', icon: '📊', type: 'player', power: 20, beats: ['denial'] },
  { id: 'legal', name: '法律威攝', icon: '⚖️', type: 'player', power: 18, beats: ['threat'] },
  { id: 'appeal', name: '情感呼籲', icon: '🤝', type: 'player', power: 12, beats: ['delay'] },
];

const ENEMY_CARDS: Card[] = [
  { id: 'denial', name: '否認一切', icon: '🙅', type: 'enemy', power: 10, beats: [] },
  { id: 'threat', name: '威脅恐嚇', icon: '😡', type: 'enemy', power: 12, beats: [] },
  { id: 'emotional', name: '情緒勒索', icon: '😭', type: 'enemy', power: 8, beats: [] },
  { id: 'delay', name: '拖延戰術', icon: '⏳', type: 'enemy', power: 6, beats: [] },
];

const ROUNDS_TOTAL = 5;

export class NegotiationScene extends Phaser.Scene {
  private stubbornness = 100; // 對手固執值
  private playerHP = 100;    // 玩家信任度
  private round = 0;
  private score = 0;
  private hand: Card[] = [];
  private currentEnemy: Card | null = null;
  private uiGroup!: Phaser.GameObjects.Group;
  private stubbornText!: Phaser.GameObjects.Text;
  private hpText!: Phaser.GameObjects.Text;
  private roundText!: Phaser.GameObjects.Text;
  private isGameOver = false;

  constructor() {
    super('NegotiationScene');
  }

  create(): void {
    this.stubbornness = 100;
    this.playerHP = 100;
    this.round = 0;
    this.score = 0;
    this.isGameOver = false;

    this.cameras.main.fadeIn(300, 0, 0, 0);
    this.add.rectangle(GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, Colors.bg);
    this.uiGroup = this.add.group();

    // HUD
    this.roundText = this.add.text(GAME_WIDTH / 2, 16, '', {
      fontFamily: 'monospace', fontSize: '18px', color: '#eaeaea',
    }).setOrigin(0.5, 0).setDepth(100);

    this.stubbornText = this.add.text(20, 16, '', {
      fontFamily: 'monospace', fontSize: '16px', color: '#e94560',
    }).setDepth(100);

    this.hpText = this.add.text(GAME_WIDTH - 20, 16, '', {
      fontFamily: 'monospace', fontSize: '16px', color: '#53d769',
    }).setOrigin(1, 0).setDepth(100);

    this.add.text(GAME_WIDTH / 2, 48, '🎴 選擇手牌對抗詐騙者的固執！', {
      fontFamily: 'monospace', fontSize: '13px', color: '#8888aa',
    }).setOrigin(0.5).setDepth(100);

    this.nextRound();
  }

  private updateHUD(): void {
    this.roundText.setText(`回合 ${this.round} / ${ROUNDS_TOTAL}`);
    this.stubbornText.setText(`🙅 固執值：${this.stubbornness}`);
    this.hpText.setText(`💚 信任度：${this.playerHP}`);
  }

  private nextRound(): void {
    if (this.isGameOver) return;
    this.round++;

    if (this.round > ROUNDS_TOTAL || this.stubbornness <= 0 || this.playerHP <= 0) {
      this.endGame();
      return;
    }

    this.uiGroup.clear(true, true);
    this.updateHUD();

    // 敵方出牌
    this.currentEnemy = ENEMY_CARDS[Phaser.Math.Between(0, ENEMY_CARDS.length - 1)];

    // 顯示敵方牌
    const enemyPanel = createPanel(this, GAME_WIDTH / 2, 160, 300, 120);
    this.uiGroup.add(enemyPanel);

    const enemyTitle = this.add.text(GAME_WIDTH / 2, 120, '詐騙者出牌：', {
      fontFamily: 'monospace', fontSize: '14px', color: '#8888aa',
    }).setOrigin(0.5);
    this.uiGroup.add(enemyTitle);

    const enemyCard = this.add.text(GAME_WIDTH / 2, 165,
      `${this.currentEnemy.icon} ${this.currentEnemy.name}`, {
        fontFamily: 'monospace', fontSize: '24px', color: '#e94560',
      }).setOrigin(0.5);
    this.uiGroup.add(enemyCard);

    const enemyPow = this.add.text(GAME_WIDTH / 2, 200,
      `威力：${this.currentEnemy.power}`, {
        fontFamily: 'monospace', fontSize: '14px', color: '#bbbbcc',
      }).setOrigin(0.5);
    this.uiGroup.add(enemyPow);

    // 洗牌給玩家
    this.hand = Phaser.Utils.Array.Shuffle([...PLAYER_CARDS]);

    // 顯示玩家手牌
    const handLabel = this.add.text(GAME_WIDTH / 2, 280, '選擇你的手牌：', {
      fontFamily: 'monospace', fontSize: '16px', color: '#eaeaea',
    }).setOrigin(0.5);
    this.uiGroup.add(handLabel);

    const cardWidth = 180;
    const gap = 20;
    const totalWidth = this.hand.length * cardWidth + (this.hand.length - 1) * gap;
    const startX = (GAME_WIDTH - totalWidth) / 2 + cardWidth / 2;

    this.hand.forEach((card, i) => {
      const x = startX + i * (cardWidth + gap);
      const y = 410;
      this.renderPlayerCard(x, y, card);
    });
  }

  private renderPlayerCard(x: number, y: number, card: Card): void {
    const w = 170;
    const h = 180;

    const bg = this.add.rectangle(x, y, w, h, Colors.panel, 0.95)
      .setStrokeStyle(2, Colors.panelLight)
      .setInteractive({ useHandCursor: true });
    this.uiGroup.add(bg);

    const icon = this.add.text(x, y - 50, card.icon, {
      fontSize: '36px',
    }).setOrigin(0.5);
    this.uiGroup.add(icon);

    const name = this.add.text(x, y - 5, card.name, {
      fontFamily: 'monospace', fontSize: '16px', color: '#eaeaea',
    }).setOrigin(0.5);
    this.uiGroup.add(name);

    const pow = this.add.text(x, y + 25, `威力：${card.power}`, {
      fontFamily: 'monospace', fontSize: '13px', color: '#ffc107',
    }).setOrigin(0.5);
    this.uiGroup.add(pow);

    // 剋制提示
    const beatsEnemy = this.currentEnemy && card.beats.includes(this.currentEnemy.id);
    if (beatsEnemy) {
      const hint = this.add.text(x, y + 50, '★ 剋制', {
        fontFamily: 'monospace', fontSize: '12px', color: '#53d769',
      }).setOrigin(0.5);
      this.uiGroup.add(hint);
    }

    bg.on('pointerover', () => bg.setStrokeStyle(3, Colors.accent));
    bg.on('pointerout', () => bg.setStrokeStyle(2, Colors.panelLight));
    bg.on('pointerdown', () => this.playCard(card));
  }

  private playCard(card: Card): void {
    if (this.isGameOver || !this.currentEnemy) return;

    const isCounter = card.beats.includes(this.currentEnemy.id);
    let damage = card.power;
    let enemyDamage = this.currentEnemy.power;

    if (isCounter) {
      damage = Math.floor(damage * 1.5);
      enemyDamage = Math.floor(enemyDamage * 0.5);
    }

    this.stubbornness = Math.max(0, this.stubbornness - damage);
    this.playerHP = Math.max(0, this.playerHP - enemyDamage);
    this.score += damage + (isCounter ? 15 : 0);

    // 顯示結果
    this.uiGroup.clear(true, true);
    this.updateHUD();

    const resultPanel = createPanel(this, GAME_WIDTH / 2, GAME_HEIGHT / 2, 500, 250);
    this.uiGroup.add(resultPanel);

    if (isCounter) {
      const counterMsg = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2 - 60,
        `🎯 剋制成功！`, {
          fontFamily: 'monospace', fontSize: '24px', color: '#53d769',
        }).setOrigin(0.5);
      this.uiGroup.add(counterMsg);
    }

    const yourPlay = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2 - 20,
      `${card.icon} ${card.name} → 造成 ${damage} 點動搖`, {
        fontFamily: 'monospace', fontSize: '16px', color: '#4fc3f7',
      }).setOrigin(0.5);
    this.uiGroup.add(yourPlay);

    const enemyPlay = this.add.text(GAME_WIDTH / 2, GAME_HEIGHT / 2 + 15,
      `${this.currentEnemy.icon} ${this.currentEnemy.name} → 造成 ${enemyDamage} 點信任損耗`, {
        fontFamily: 'monospace', fontSize: '16px', color: '#e94560',
      }).setOrigin(0.5);
    this.uiGroup.add(enemyPlay);

    const continueBtn = createButton(this, GAME_WIDTH / 2, GAME_HEIGHT / 2 + 80, '下一回合 →', () => {
      this.nextRound();
    }, 200, 44);
    this.uiGroup.add(continueBtn);
  }

  private endGame(): void {
    this.isGameOver = true;
    this.uiGroup.clear(true, true);

    const won = this.stubbornness <= 0 && this.playerHP > 0;

    const overlay = this.add.rectangle(
      GAME_WIDTH / 2, GAME_HEIGHT / 2, GAME_WIDTH, GAME_HEIGHT, 0x000000, 0.7,
    ).setDepth(300);

    createTitle(this, GAME_WIDTH / 2, 160, '🎴 談判牌局 — 結算', '28px').setDepth(301);

    const resultMsg = won ? '✅ 談判成功！詐騙者放棄固執！' : '❌ 談判失敗…信任度耗盡';
    this.add.text(GAME_WIDTH / 2, 230, resultMsg, {
      fontFamily: 'monospace', fontSize: '22px', color: won ? '#53d769' : '#e94560',
    }).setOrigin(0.5).setDepth(301);

    this.add.text(GAME_WIDTH / 2, 280, `最終分數：${this.score}`, {
      fontFamily: 'monospace', fontSize: '24px', color: '#ffc107',
    }).setOrigin(0.5).setDepth(301);

    const rank = this.getRank();
    this.add.text(GAME_WIDTH / 2, 320, `評價：${rank}`, {
      fontFamily: 'monospace', fontSize: '18px', color: '#bbbbcc',
    }).setOrigin(0.5).setDepth(301);

    const result = gameState.recordMiniGame('negotiation', this.score);
    const bonusColor = result.isNewBest ? '#53d769' : '#8888aa';
    this.add.text(GAME_WIDTH / 2, 355, result.isNewBest
      ? `🆕 新紀錄！ 交涉 +${result.bonus}`
      : `交涉 +${result.bonus}`, {
      fontFamily: 'monospace', fontSize: '16px', color: bonusColor,
    }).setOrigin(0.5).setDepth(301);

    // 同步至後端
    if (result.delta > 0) {
      updatePlayerStats(PLAYER_ID, { logic: 0, tech: 0, charisma: result.delta, resilience: 0 }).catch(() => {});
    }

    const retryBtn = createButton(this, GAME_WIDTH / 2 - 120, 420, '🔄 再挑戰', () => {
      this.scene.restart();
    }, 200, 48);
    retryBtn.setDepth(301);

    const backBtn = createButton(this, GAME_WIDTH / 2 + 120, 420, '← 返回選單', () => {
      this.scene.start('MainMenuScene');
    }, 200, 48);
    backBtn.setDepth(301);
  }

  private getRank(): string {
    if (this.score >= 150) return '🏆 S — 談判大師';
    if (this.score >= 100) return '🥇 A — 資深調解員';
    if (this.score >= 60) return '🥈 B — 合格談判者';
    if (this.score >= 30) return '🥉 C — 見習交涉官';
    return '📝 D — 需要更多訓練';
  }
}
