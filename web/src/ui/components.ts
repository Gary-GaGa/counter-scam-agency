import Phaser from 'phaser';
import { Colors } from './constants';

// Creates a styled interactive button.
export function createButton(
  scene: Phaser.Scene,
  x: number,
  y: number,
  label: string,
  onClick: () => void,
  width = 260,
  height = 48,
): Phaser.GameObjects.Container {
  const bg = scene.add.rectangle(0, 0, width, height, Colors.panelLight, 1)
    .setStrokeStyle(2, Colors.accent);
  const text = scene.add.text(0, 0, label, {
    fontFamily: 'monospace',
    fontSize: '18px',
    color: '#eaeaea',
    align: 'center',
  }).setOrigin(0.5);

  const container = scene.add.container(x, y, [bg, text]);
  container.setSize(width, height);
  container.setInteractive({ useHandCursor: true });

  container.on('pointerover', () => {
    bg.setFillStyle(Colors.accent);
  });
  container.on('pointerout', () => {
    bg.setFillStyle(Colors.panelLight);
  });
  container.on('pointerdown', onClick);

  return container;
}

// Creates a panel background.
export function createPanel(
  scene: Phaser.Scene,
  x: number,
  y: number,
  width: number,
  height: number,
): Phaser.GameObjects.Rectangle {
  return scene.add.rectangle(x, y, width, height, Colors.panel, 0.9)
    .setStrokeStyle(2, Colors.panelLight)
    .setOrigin(0.5);
}

// Creates a styled title text.
export function createTitle(
  scene: Phaser.Scene,
  x: number,
  y: number,
  text: string,
  fontSize = '28px',
): Phaser.GameObjects.Text {
  return scene.add.text(x, y, text, {
    fontFamily: 'monospace',
    fontSize,
    color: '#e94560',
    align: 'center',
  }).setOrigin(0.5);
}

// Creates body text with word wrap.
export function createBodyText(
  scene: Phaser.Scene,
  x: number,
  y: number,
  content: string,
  maxWidth = 700,
): Phaser.GameObjects.Text {
  return scene.add.text(x, y, content, {
    fontFamily: 'monospace',
    fontSize: '16px',
    color: '#eaeaea',
    wordWrap: { width: maxWidth },
    lineSpacing: 6,
  }).setOrigin(0.5, 0);
}

// Creates a small label.
export function createLabel(
  scene: Phaser.Scene,
  x: number,
  y: number,
  content: string,
  color = '#8888aa',
  fontSize = '14px',
): Phaser.GameObjects.Text {
  return scene.add.text(x, y, content, {
    fontFamily: 'monospace',
    fontSize,
    color,
  }).setOrigin(0.5);
}

// Bottom navigation bar.
export function createNavBar(scene: Phaser.Scene): void {
  const y = 610;
  const barBg = scene.add.rectangle(480, y, 960, 60, Colors.panel, 0.95);
  barBg.setDepth(100);

  const navItems = [
    { label: '🏠 主選單', scene: 'MainMenuScene' },
    { label: '📋 案件', scene: 'CaseListScene' },
    { label: '🌳 技能樹', scene: 'SkillTreeScene' },
    { label: '👤 狀態', scene: 'ProfileScene' },
  ];

  const startX = 140;
  const gap = 220;

  navItems.forEach((item, i) => {
    const btn = scene.add.text(startX + i * gap, y, item.label, {
      fontFamily: 'monospace',
      fontSize: '16px',
      color: '#8888aa',
      align: 'center',
    }).setOrigin(0.5).setInteractive({ useHandCursor: true }).setDepth(101);

    btn.on('pointerover', () => btn.setColor('#e94560'));
    btn.on('pointerout', () => btn.setColor('#8888aa'));
    btn.on('pointerdown', () => {
      if (scene.scene.key !== item.scene) {
        scene.scene.start(item.scene);
      }
    });
  });
}
