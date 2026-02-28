/**
 * Mock API — 當後端不可用時，使用本地資料模擬所有 API 呼叫。
 * 回傳型別與真實 API 完全一致。
 */

import type {
  MissionSummary,
  MissionDetail,
  InvestigationStartResult,
  NodeProgressResult,
  SubmitEvidenceResult,
  CompleteResult,
  PlayerSummary,
  SkillSummary,
  SkillActionResult,
  BaseSummary,
} from './types';

import {
  MOCK_MISSIONS,
  MOCK_MISSION_DETAILS,
  MOCK_PLAYER,
  MOCK_SKILLS,
  MOCK_BASE,
} from './mockData';

import { gameState } from './gameState';

const delay = (ms: number) => new Promise(r => setTimeout(r, ms));

// 模擬網路延遲
const SIM_DELAY = 200;

// ─── 調查狀態追蹤 ───

interface MockInvestigation {
  id: string;
  playerId: string;
  missionId: string;
  currentNodeId: string;
  evidenceCollected: string[];
  status: 'Active' | 'Completed';
}

const investigations = new Map<string, MockInvestigation>();

// ─── Mission ───

export async function listMissions(): Promise<MissionSummary[]> {
  await delay(SIM_DELAY);
  return [...MOCK_MISSIONS];
}

export async function getMission(id: string): Promise<MissionDetail> {
  await delay(SIM_DELAY);
  const m = MOCK_MISSION_DETAILS[id];
  if (!m) throw new Error('找不到此案件');
  return { ...m };
}

// ─── Investigation ───

export async function startInvestigation(
  investigationId: string,
  playerId: string,
  missionId: string,
): Promise<InvestigationStartResult> {
  await delay(SIM_DELAY);
  const mission = MOCK_MISSION_DETAILS[missionId];
  if (!mission) throw new Error('找不到此案件');

  const inv: MockInvestigation = {
    id: investigationId,
    playerId,
    missionId,
    currentNodeId: mission.nodes[0].id,
    evidenceCollected: [],
    status: 'Active',
  };
  investigations.set(investigationId, inv);

  return {
    investigationId,
    playerId,
    missionId,
    status: 'Active',
    currentNodeId: inv.currentNodeId,
  };
}

export async function advanceNode(
  investigationId: string,
  nodeId: string,
  optionId: string,
): Promise<NodeProgressResult> {
  await delay(SIM_DELAY);
  const inv = investigations.get(investigationId);
  if (!inv) throw new Error('找不到此調查');

  const mission = MOCK_MISSION_DETAILS[inv.missionId];
  const node = mission.nodes.find(n => n.id === nodeId);
  const option = node?.options.find(o => o.id === optionId);

  if (option?.evidenceIds) {
    option.evidenceIds.forEach(eid => {
      if (!inv.evidenceCollected.includes(eid)) {
        inv.evidenceCollected.push(eid);
      }
    });
  }

  const nextNodeId = option?.nextNodeId || '';
  const nextNode = mission.nodes.find(n => n.id === nextNodeId);

  if (option?.leadsToEnd || nextNode?.isTerminal) {
    inv.status = 'Completed';
    inv.currentNodeId = nextNodeId;
  } else {
    inv.currentNodeId = nextNodeId;
  }

  return {
    investigationId,
    nodeId,
    optionId,
    nextNodeId,
    status: inv.status,
  };
}

export async function submitEvidence(
  investigationId: string,
  evidenceId: string,
): Promise<SubmitEvidenceResult> {
  await delay(SIM_DELAY);
  const inv = investigations.get(investigationId);
  if (inv && !inv.evidenceCollected.includes(evidenceId)) {
    inv.evidenceCollected.push(evidenceId);
  }
  return { investigationId, evidenceId };
}

export async function completeInvestigation(
  investigationId: string,
): Promise<CompleteResult> {
  await delay(SIM_DELAY);
  const inv = investigations.get(investigationId);
  if (!inv) throw new Error('找不到此調查');

  inv.status = 'Completed';
  const mission = MOCK_MISSION_DETAILS[inv.missionId];
  const keyEvidence = mission.evidenceList.filter(e => e.isKey);
  const collectedKey = keyEvidence.filter(e => inv.evidenceCollected.includes(e.id));
  const ratio = keyEvidence.length > 0 ? collectedKey.length / keyEvidence.length : 0;
  const success = ratio >= 0.3;
  const reputationGained = success ? Math.round(mission.reputationWeight * (5 + ratio * 15)) : 0;

  // 更新 game state
  if (reputationGained > 0) {
    gameState.addReputation(reputationGained);
  }
  gameState.completeMission(inv.missionId);

  return {
    investigationId,
    success,
    reputationGained,
    evidenceCollected: inv.evidenceCollected.length,
  };
}

// ─── Player ───

export async function createPlayer(_playerId: string): Promise<PlayerSummary> {
  await delay(SIM_DELAY);
  return getPlayerSummary();
}

export async function getPlayer(_playerId: string): Promise<PlayerSummary> {
  await delay(SIM_DELAY);
  return getPlayerSummary();
}

export async function updatePlayerStats(
  _playerId: string,
  _stats: { logic: number; tech: number; charisma: number; resilience: number },
): Promise<PlayerSummary> {
  await delay(SIM_DELAY);
  // 離線模式下，gameState 已在 recordMiniGame 中更新屬性，此處直接回傳
  return getPlayerSummary();
}

function getPlayerSummary(): PlayerSummary {
  const state = gameState.get();
  return {
    ...MOCK_PLAYER,
    reputation: state.reputation,
    stats: { ...state.stats },
    totalStats: { ...state.stats },
    equippedSkills: state.unlockedSkills,
  };
}

// ─── Skills ───

export async function listSkills(_playerId: string): Promise<SkillSummary[]> {
  await delay(SIM_DELAY);
  const state = gameState.get();
  return MOCK_SKILLS.map(s => ({
    ...s,
    unlocked: state.unlockedSkills.includes(s.id),
    equipped: state.unlockedSkills.includes(s.id),
  }));
}

export async function unlockSkill(_playerId: string, skillId: string): Promise<SkillActionResult> {
  await delay(SIM_DELAY);
  const state = gameState.get();
  const skill = MOCK_SKILLS.find(s => s.id === skillId);
  if (!skill) throw new Error('找不到此技能');
  if (state.reputation < skill.reputationRequired) {
    throw new Error(`需要 ${skill.reputationRequired} 聲望才能解鎖`);
  }
  gameState.unlockSkill(skillId);
  return { playerId: 'player-1', skillId, unlocked: true, equipped: false, cooldownRemaining: 0 };
}

export async function equipSkill(_playerId: string, skillId: string): Promise<SkillActionResult> {
  await delay(SIM_DELAY);
  return { playerId: 'player-1', skillId, unlocked: true, equipped: true, cooldownRemaining: 0 };
}

export async function activateSkill(_playerId: string, skillId: string): Promise<SkillActionResult> {
  await delay(SIM_DELAY);
  return { playerId: 'player-1', skillId, unlocked: true, equipped: true, cooldownRemaining: 3 };
}

// ─── Defense Base ───

export async function createBase(
  baseId: string,
  ownerId: string,
  slots: number,
): Promise<BaseSummary> {
  await delay(SIM_DELAY);
  return { ...MOCK_BASE, id: baseId, ownerId, facilitySlots: slots };
}

export async function getBase(_baseId: string): Promise<BaseSummary> {
  await delay(SIM_DELAY);
  return { ...MOCK_BASE };
}

export async function addFacility(
  _baseId: string,
  facility: { id: string; type: string; name: string; level: number; maxLevel: number; description: string },
): Promise<BaseSummary> {
  await delay(SIM_DELAY);
  return {
    ...MOCK_BASE,
    facilities: [...MOCK_BASE.facilities, facility],
  };
}

export async function upgradeSecurityLevel(
  _baseId: string,
  _maxLevel: number,
): Promise<BaseSummary> {
  await delay(SIM_DELAY);
  return { ...MOCK_BASE, securityLevel: MOCK_BASE.securityLevel + 1 };
}

export async function upgradeFacility(
  _baseId: string,
  facilityId: string,
): Promise<BaseSummary> {
  await delay(SIM_DELAY);
  return {
    ...MOCK_BASE,
    facilities: MOCK_BASE.facilities.map(f =>
      f.id === facilityId ? { ...f, level: f.level + 1 } : f,
    ),
  };
}
