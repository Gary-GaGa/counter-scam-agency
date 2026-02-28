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

import * as mockApi from './mockApi';

const BASE_URL = '/api';

// 自動偵測後端是否可用，不可用時降級為 mock
let _useMock: boolean | null = null;

async function shouldUseMock(): Promise<boolean> {
  if (_useMock !== null) return _useMock;
  try {
    const res = await fetch(`${BASE_URL}/missions`, { signal: AbortSignal.timeout(1500) });
    _useMock = !res.ok;
  } catch {
    _useMock = true;
  }
  if (_useMock) {
    console.log('%c[CSA] 後端不可用，使用離線模式 🎮', 'color:#ffc107;font-weight:bold');
  }
  return _useMock;
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

// --- Mission ---

export async function listMissions(): Promise<MissionSummary[]> {
  if (await shouldUseMock()) return mockApi.listMissions();
  return request<MissionSummary[]>('/missions');
}

export async function getMission(id: string): Promise<MissionDetail> {
  if (await shouldUseMock()) return mockApi.getMission(id);
  return request<MissionDetail>(`/missions/${id}`);
}

// --- Investigation ---

export async function startInvestigation(
  investigationId: string,
  playerId: string,
  missionId: string,
): Promise<InvestigationStartResult> {
  if (await shouldUseMock()) return mockApi.startInvestigation(investigationId, playerId, missionId);
  return request<InvestigationStartResult>('/investigations', {
    method: 'POST',
    body: JSON.stringify({ investigationId, playerId, missionId }),
  });
}

export async function advanceNode(
  investigationId: string,
  nodeId: string,
  optionId: string,
): Promise<NodeProgressResult> {
  if (await shouldUseMock()) return mockApi.advanceNode(investigationId, nodeId, optionId);
  return request<NodeProgressResult>(`/investigations/${investigationId}/advance`, {
    method: 'POST',
    body: JSON.stringify({ nodeId, optionId }),
  });
}

export async function submitEvidence(
  investigationId: string,
  evidenceId: string,
): Promise<SubmitEvidenceResult> {
  if (await shouldUseMock()) return mockApi.submitEvidence(investigationId, evidenceId);
  return request<SubmitEvidenceResult>(`/investigations/${investigationId}/evidence`, {
    method: 'POST',
    body: JSON.stringify({ evidenceId }),
  });
}

export async function completeInvestigation(
  investigationId: string,
): Promise<CompleteResult> {
  if (await shouldUseMock()) return mockApi.completeInvestigation(investigationId);
  return request<CompleteResult>(`/investigations/${investigationId}/complete`, {
    method: 'POST',
  });
}

// --- Player ---

export async function createPlayer(playerId: string): Promise<PlayerSummary> {
  if (await shouldUseMock()) return mockApi.createPlayer(playerId);
  return request<PlayerSummary>('/players', {
    method: 'POST',
    body: JSON.stringify({ playerId }),
  });
}

export async function getPlayer(playerId: string): Promise<PlayerSummary> {
  if (await shouldUseMock()) return mockApi.getPlayer(playerId);
  return request<PlayerSummary>(`/players/${playerId}`);
}

// --- Skills ---

export async function listSkills(playerId: string): Promise<SkillSummary[]> {
  if (await shouldUseMock()) return mockApi.listSkills(playerId);
  return request<SkillSummary[]>(`/players/${playerId}/skills`);
}

export async function unlockSkill(playerId: string, skillId: string): Promise<SkillActionResult> {
  if (await shouldUseMock()) return mockApi.unlockSkill(playerId, skillId);
  return request<SkillActionResult>(`/players/${playerId}/skills/${skillId}/unlock`, {
    method: 'POST',
  });
}

export async function equipSkill(playerId: string, skillId: string): Promise<SkillActionResult> {
  if (await shouldUseMock()) return mockApi.equipSkill(playerId, skillId);
  return request<SkillActionResult>(`/players/${playerId}/skills/${skillId}/equip`, {
    method: 'POST',
  });
}

export async function activateSkill(playerId: string, skillId: string): Promise<SkillActionResult> {
  if (await shouldUseMock()) return mockApi.activateSkill(playerId, skillId);
  return request<SkillActionResult>(`/players/${playerId}/skills/${skillId}/activate`, {
    method: 'POST',
  });
}

// --- Defense Base ---

export async function createBase(
  baseId: string,
  ownerId: string,
  slots: number,
): Promise<BaseSummary> {
  if (await shouldUseMock()) return mockApi.createBase(baseId, ownerId, slots);
  return request<BaseSummary>('/bases', {
    method: 'POST',
    body: JSON.stringify({ baseId, ownerId, slots }),
  });
}

export async function getBase(baseId: string): Promise<BaseSummary> {
  if (await shouldUseMock()) return mockApi.getBase(baseId);
  return request<BaseSummary>(`/bases/${baseId}`);
}

export async function addFacility(
  baseId: string,
  facility: { id: string; type: string; name: string; level: number; maxLevel: number; description: string },
): Promise<BaseSummary> {
  if (await shouldUseMock()) return mockApi.addFacility(baseId, facility);
  return request<BaseSummary>(`/bases/${baseId}/facilities`, {
    method: 'POST',
    body: JSON.stringify(facility),
  });
}

export async function upgradeSecurityLevel(
  baseId: string,
  maxLevel: number,
): Promise<BaseSummary> {
  if (await shouldUseMock()) return mockApi.upgradeSecurityLevel(baseId, maxLevel);
  return request<BaseSummary>(`/bases/${baseId}/security/upgrade`, {
    method: 'POST',
    body: JSON.stringify({ maxLevel }),
  });
}

export async function upgradeFacility(
  baseId: string,
  facilityId: string,
): Promise<BaseSummary> {
  if (await shouldUseMock()) return mockApi.upgradeFacility(baseId, facilityId);
  return request<BaseSummary>(`/bases/${baseId}/facilities/${facilityId}/upgrade`, {
    method: 'POST',
  });
}
