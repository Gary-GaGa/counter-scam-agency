import Phaser from 'phaser';
import { GAME_WIDTH, GAME_HEIGHT, Colors } from './ui/constants';
import { MainMenuScene } from './scenes/MainMenuScene';
import { CaseListScene } from './scenes/CaseListScene';
import { InvestigationScene } from './scenes/InvestigationScene';
import { SkillTreeScene } from './scenes/SkillTreeScene';
import { ProfileScene } from './scenes/ProfileScene';
import { BaseScene } from './scenes/BaseScene';

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  width: GAME_WIDTH,
  height: GAME_HEIGHT,
  parent: 'game-container',
  backgroundColor: `#${Colors.bg.toString(16)}`,
  scale: {
    mode: Phaser.Scale.FIT,
    autoCenter: Phaser.Scale.CENTER_BOTH,
  },
  scene: [MainMenuScene, CaseListScene, InvestigationScene, SkillTreeScene, ProfileScene, BaseScene],
};

new Phaser.Game(config);
