// DTO types matching the Go backend usecase/dto package.

export interface MissionSummary {
  id: string;
  title: string;
  description: string;
  type: string;
  difficulty: number;
  reputationWeight: number;
}

export interface MissionDetail {
  id: string;
  title: string;
  description: string;
  type: string;
  difficulty: number;
  reputationWeight: number;
  nodes: NarrativeNode[];
  evidenceList: Evidence[];
  victimProfile: VictimProfile | null;
}

export interface NarrativeNode {
  id: string;
  title: string;
  body: string;
  isTerminal: boolean;
  options: NarrativeOption[];
}

export interface NarrativeOption {
  id: string;
  label: string;
  nextNodeId: string;
  evidenceIds: string[];
  leadsToEnd: boolean;
  successEnd: boolean;
}

export interface Evidence {
  id: string;
  description: string;
  type: string;
  isKey: boolean;
}

export interface VictimProfile {
  anxiety: number;
  trust: number;
  urgency: number;
  isolation: number;
  riskScore: number;
  riskLevel: string;
}

export interface InvestigationStartResult {
  investigationId: string;
  playerId: string;
  missionId: string;
  status: string;
  currentNodeId: string;
}

export interface NodeProgressResult {
  investigationId: string;
  nodeId: string;
  optionId: string;
  nextNodeId: string;
  status: string;
}

export interface SubmitEvidenceResult {
  investigationId: string;
  evidenceId: string;
}

export interface CompleteResult {
  investigationId: string;
  success: boolean;
  reputationGained: number;
  evidenceCollected: number;
}

export interface PlayerSummary {
  id: string;
  reputation: number;
  stats: StatsSummary;
  totalStats: StatsSummary;
  partnerPersonality: string;
  equippedModules: string[];
  equippedSkills: string[];
}

export interface StatsSummary {
  logic: number;
  tech: number;
  charisma: number;
  resilience: number;
}

export interface SkillSummary {
  id: string;
  name: string;
  description: string;
  type: string;
  unlocked: boolean;
  equipped: boolean;
  cooldownRemaining: number;
  reputationRequired: number;
}

export interface SkillActionResult {
  playerId: string;
  skillId: string;
  unlocked: boolean;
  equipped: boolean;
  cooldownRemaining: number;
}

export interface BaseSummary {
  id: string;
  ownerId: string;
  securityLevel: number;
  facilitySlots: number;
  facilities: FacilitySummary[];
}

export interface FacilitySummary {
  id: string;
  type: string;
  name: string;
  level: number;
  maxLevel: number;
  description: string;
}
