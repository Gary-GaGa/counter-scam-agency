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

const BASE_URL = '/api';

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

export function listMissions(): Promise<MissionSummary[]> {
  return request<MissionSummary[]>('/missions');
}

export function getMission(id: string): Promise<MissionDetail> {
  return request<MissionDetail>(`/missions/${id}`);
}

// --- Investigation ---

export function startInvestigation(
  investigationId: string,
  playerId: string,
  missionId: string,
): Promise<InvestigationStartResult> {
  return request<InvestigationStartResult>('/investigations', {
    method: 'POST',
    body: JSON.stringify({ investigationId, playerId, missionId }),
  });
}

export function advanceNode(
  investigationId: string,
  nodeId: string,
  optionId: string,
): Promise<NodeProgressResult> {
  return request<NodeProgressResult>(`/investigations/${investigationId}/advance`, {
    method: 'POST',
    body: JSON.stringify({ nodeId, optionId }),
  });
}

export function submitEvidence(
  investigationId: string,
  evidenceId: string,
): Promise<SubmitEvidenceResult> {
  return request<SubmitEvidenceResult>(`/investigations/${investigationId}/evidence`, {
    method: 'POST',
    body: JSON.stringify({ evidenceId }),
  });
}

export function completeInvestigation(
  investigationId: string,
): Promise<CompleteResult> {
  return request<CompleteResult>(`/investigations/${investigationId}/complete`, {
    method: 'POST',
  });
}

// --- Player ---

export function createPlayer(playerId: string): Promise<PlayerSummary> {
  return request<PlayerSummary>('/players', {
    method: 'POST',
    body: JSON.stringify({ playerId }),
  });
}

export function getPlayer(playerId: string): Promise<PlayerSummary> {
  return request<PlayerSummary>(`/players/${playerId}`);
}

// --- Skills ---

export function listSkills(playerId: string): Promise<SkillSummary[]> {
  return request<SkillSummary[]>(`/players/${playerId}/skills`);
}

export function unlockSkill(playerId: string, skillId: string): Promise<SkillActionResult> {
  return request<SkillActionResult>(`/players/${playerId}/skills/${skillId}/unlock`, {
    method: 'POST',
  });
}

export function equipSkill(playerId: string, skillId: string): Promise<SkillActionResult> {
  return request<SkillActionResult>(`/players/${playerId}/skills/${skillId}/equip`, {
    method: 'POST',
  });
}

export function activateSkill(playerId: string, skillId: string): Promise<SkillActionResult> {
  return request<SkillActionResult>(`/players/${playerId}/skills/${skillId}/activate`, {
    method: 'POST',
  });
}

// --- Defense Base ---

export function createBase(
  baseId: string,
  ownerId: string,
  slots: number,
): Promise<BaseSummary> {
  return request<BaseSummary>('/bases', {
    method: 'POST',
    body: JSON.stringify({ baseId, ownerId, slots }),
  });
}

export function getBase(baseId: string): Promise<BaseSummary> {
  return request<BaseSummary>(`/bases/${baseId}`);
}

export function addFacility(
  baseId: string,
  facility: { id: string; type: string; name: string; level: number; maxLevel: number; description: string },
): Promise<BaseSummary> {
  return request<BaseSummary>(`/bases/${baseId}/facilities`, {
    method: 'POST',
    body: JSON.stringify(facility),
  });
}

export function upgradeSecurityLevel(
  baseId: string,
  maxLevel: number,
): Promise<BaseSummary> {
  return request<BaseSummary>(`/bases/${baseId}/security/upgrade`, {
    method: 'POST',
    body: JSON.stringify({ maxLevel }),
  });
}

export function upgradeFacility(
  baseId: string,
  facilityId: string,
): Promise<BaseSummary> {
  return request<BaseSummary>(`/bases/${baseId}/facilities/${facilityId}/upgrade`, {
    method: 'POST',
  });
}
