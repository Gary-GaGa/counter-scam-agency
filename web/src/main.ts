import Phaser from 'phaser';
import { GAME_WIDTH, GAME_HEIGHT, Colors } from './ui/constants';
import { IntroScene } from './scenes/IntroScene';
import { MainMenuScene } from './scenes/MainMenuScene';
import { CaseListScene } from './scenes/CaseListScene';
import { InvestigationScene } from './scenes/InvestigationScene';
import { SkillTreeScene } from './scenes/SkillTreeScene';
import { ProfileScene } from './scenes/ProfileScene';
import { BaseScene } from './scenes/BaseScene';
import { ContradictionScene } from './scenes/ContradictionScene';
import { SignalTraceScene } from './scenes/SignalTraceScene';
import { NegotiationScene } from './scenes/NegotiationScene';
import { MentalRecoveryScene } from './scenes/MentalRecoveryScene';
import { gameState } from './services/gameState';

// 首次進入播放 Intro，之後直接進主選單
const introSeen = gameState.get().introSeen;

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
  scene: introSeen
    ? [MainMenuScene, IntroScene, CaseListScene, InvestigationScene, SkillTreeScene, ProfileScene, BaseScene, ContradictionScene, SignalTraceScene, NegotiationScene, MentalRecoveryScene]
    : [IntroScene, MainMenuScene, CaseListScene, InvestigationScene, SkillTreeScene, ProfileScene, BaseScene, ContradictionScene, SignalTraceScene, NegotiationScene, MentalRecoveryScene],
};

new Phaser.Game(config);
