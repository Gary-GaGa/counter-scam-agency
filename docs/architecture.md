# System Architecture

## Overview
Based on **Clean Architecture** and **Domain-Driven Design (DDD)** principles, the system is designed to be testable, maintainable, and loosely coupled.

## Directory Structure
The project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
cmd/
  ├── adventure/          # Text-adventure CLI entry point
  ├── cli/                # General CLI entry point
  └── seed-missions/      # Seed data tooling

internal/
  ├── domain/             # [Core] Enterprise Business Rules (Pure Go, No external deps)
  │   ├── personnel/      # Context: Player, AIPartner, Module, Skill, Stats
  │   ├── intelligence/   # Context: Mission, Evidence, ScamType, VictimProfile
  │   ├── operation/      # Context: Investigation (Runtime State)
  │   └── defense/        # Context: Base, Facility, Upgrade (Digital Defense)
  │
  ├── usecase/            # [Application] Application Business Rules
  │   ├── port/
  │   │   └── in/         # Input Ports (Usecase Interfaces)
  │   ├── dto/            # Data Transfer Objects (UI/API boundary)
  │   ├── investigation/  # Investigation Service
  │   ├── personnel/      # Personnel Service (Player/AI Skill flows)
  │   └── defense/        # Defense Service (Base management)
  │
  ├── interface/          # [Adapters] Interface Adapters
  │   ├── in/
  │   │   └── http/       # HTTP Handlers (planned)
  │   └── out/
  │       └── persistence/
  │           └── mongo/  # MongoDB Repository Implementations
  │               ├── player/
  │               ├── mission/
  │               ├── investigation/
  │               ├── defense/
  │               └── po/           # Persistence Objects & Converters
  │                   ├── convert/
  │                   ├── defense/
  │                   ├── intelligence/
  │                   ├── operation/
  │                   └── personnel/
  │
  └── infrastructure/     # [Frameworks] External Details
      ├── persistence/
      │   └── mongo/      # MongoDB Client/Database setup (shared)
      ├── memory/
      │   └── redis/      # Redis cache (planned)
      └── logger/         # Logging (planned)
```

## 1. Domain Layer (Core)
The heart of the software. It contains the business objects and high-level rules. Domain layer has **zero external dependencies** — only pure Go types and logic.

### Bounded Contexts
We split the domain into four distinct contexts based on the "Agency" metaphor:

#### A. Personnel Context (`internal/domain/personnel`)
- **Responsibility**: Manages the agents (Players), AI partner equipment, and skill systems.
- **Key Entities**:
  - `Player` (Aggregate Root): Owns Stats, AI Partner, reputation, unlocked modules/skills.
  - `AIPartner` (Entity): Loadout (modules), Skills, SkillCooldowns, Personality evolution.
  - `Module` (Entity): Upgradeable AI component (VoiceAnalyzer, CryptoTracer, EmpathyEngine, MentalFirewall) with stat bonuses.
  - `Skill` (Entity): Active AI ability with cooldown and module dependencies (Analysis, Negotiation, Defense, Forensics).
  - `Stats` (Value Object): Logic, Tech, Charisma, Resilience — supports arithmetic operations.

#### B. Intelligence Context (`internal/domain/intelligence`)
- **Responsibility**: Manages the static case files, knowledge base, and victim psychology.
- **Key Entities**:
  - `Mission` (Aggregate Root): The script of a case with narrative nodes, options, evidence, and reputation weight.
  - `NarrativeNode` (Entity): A node in the text adventure flow with options and terminal flags.
  - `NarrativeOption` (Value Object): A player choice linking to the next node and associated evidence.
  - `Evidence` (Entity): Clues and contradictions (Dialogue, Transaction, Image, Document).
  - `VictimProfile` (Entity): Psychological profile (Anxiety, Trust, Urgency, Isolation) with risk scoring.
  - `ScamType` (Value Object): Classification of scams (Phishing, Investment, Romance, Impersonation).

#### C. Operation Context (`internal/domain/operation`)
- **Responsibility**: Manages the runtime execution of a mission play-through.
- **Key Entities**:
  - `Investigation` (Aggregate Root): Tracks status (Active/Completed/Failed), current node, decision history, collected evidence, and suspicion level.
  - `NodeDecision` (Value Object): Records a player choice (nodeID, optionID, nextNodeID).

#### D. Defense Context (`internal/domain/defense`)
- **Responsibility**: Manages the player's digital defense base and facility upgrades.
- **Key Entities**:
  - `Base` (Aggregate Root): Owns facilities, upgrade history, security level, and slot capacity.
  - `Facility` (Entity): A buildable structure (Firewall, SIEM, Training) with level progression.
  - `Upgrade` (Value Object): Base-level upgrade with level cap.

### Repository Interfaces
The Domain layer defines **interfaces** for persistence, following the Dependency Inversion Principle.
- `PlayerRepository` — Save, FindByID
- `MissionRepository` — FindByID, FindAll, Save
- `InvestigationRepository` — Save, FindByID, FindByPlayerID
- `BaseRepository` — Save, FindByID, FindByOwnerID

## 2. Usecase Layer (Application)
Orchestrates the flow of data to and from the domain entities. Each service depends only on Domain Repository interfaces (injected via constructor).

### Input Ports (`internal/usecase/port/in`)
Define the application's public API as Go interfaces, consumed by adapters (CLI, HTTP, etc.):
- `InvestigationUsecase` — ListMissions, GetMission, StartInvestigation, AdvanceNode, SubmitEvidence, CompleteInvestigation
- `PersonnelUsecase` — ListSkills, UnlockSkill, EquipSkill, ActivateSkill, TickSkillCooldowns
- `DefenseUsecase` — CreateBase, GetBase, AddFacility, UpgradeSecurity, UpgradeFacility

### DTOs (`internal/usecase/dto`)
Data Transfer Objects decouple domain models from UI/API representation:
- `MissionSummary`, `MissionDetail`, `NarrativeNode`, `NarrativeOption`, `Evidence`, `VictimProfile`
- `InvestigationStartResult`, `NodeProgressResult`, `SubmitEvidenceResult`, `CompleteResult`
- `SkillSummary`, `SkillActionResult`
- `BaseSummary`, `FacilitySummary`, `FacilityInput`

### Services
- **Investigation Service** (`internal/usecase/investigation`): Orchestrates mission browsing, investigation lifecycle (start → advance → evidence → complete), and reputation accumulation. Cross-context: reads Mission (Intelligence), mutates Investigation (Operation) and Player (Personnel).
- **Personnel Service** (`internal/usecase/personnel`): Manages AI skill catalog, unlock/equip/activate flows, and cooldown ticking.
- **Defense Service** (`internal/usecase/defense`): Manages base creation, facility installation, and security/facility upgrades.

## 3. Interface Layer (Adapters)

### Input Adapters (`internal/interface/in`)
- **HTTP** (`internal/interface/in/http`): REST API handlers (planned).
- **CLI** (`cmd/cli`, `cmd/adventure`): Command-line game loop and adventure mode.

### Output Adapters (`internal/interface/out`)
Implement the repository interfaces defined by the Domain layer.

- **Persistence Adapters**: MongoDB repositories per domain context.
  - `internal/interface/out/persistence/mongo/player` — Implements `PlayerRepository`
  - `internal/interface/out/persistence/mongo/mission` — Implements `MissionRepository`
  - `internal/interface/out/persistence/mongo/investigation` — Implements `InvestigationRepository`
  - `internal/interface/out/persistence/mongo/defense` — Implements `BaseRepository`
- **Persistence Objects** (`internal/interface/out/persistence/mongo/po`): MongoDB document structs and bidirectional converters (Domain ↔ PO) per context.

## 4. Infrastructure Layer (Frameworks)
Contains external setup details shared across adapters.

- **Persistence**: MongoDB connection/database configuration (`internal/infrastructure/persistence/mongo`)
- **Memory**: Redis cache setup (`internal/infrastructure/memory/redis`, planned)
- **Logger**: Structured logging (`internal/infrastructure/logger`, planned)

## Key Design Decisions
1.  **Repository Interface in Domain**: Allows the domain to define its storage needs without knowing the implementation details (Dependency Inversion Principle).
2.  **Rich Domain Models**: Entities contain logic (e.g., `Player.GetTotalStats()`, `Investigation.RecordDecision()`, `Base.AddFacility()`), preventing Anemic Domain Models.
3.  **Split Infrastructure vs Interface Out**: Shared infrastructure (DB client) lives in `infrastructure/`, while per-context repository implementations live in `interface/out/`.
4.  **Persistence Objects (PO)**: Dedicated `po/` layer with converters decouples MongoDB document structure from domain models, allowing independent evolution.
5.  **Input Ports as Interfaces**: Usecase interfaces in `port/in/` enable adapters (CLI, HTTP) to depend on abstractions, not concrete services.
6.  **DTO Boundary**: Domain entities never leak to UI/API; all external communication uses DTOs defined in `usecase/dto/`.
7.  **Cross-Context Orchestration in Usecase**: Services like `InvestigationService` coordinate multiple bounded contexts (Intelligence + Operation + Personnel) without coupling domain layers to each other.
8.  **Compile-time Interface Guards**: Each repository and usecase service includes `var _ Interface = (*Impl)(nil)` to ensure interface satisfaction at compile time.

## 5. Frontend Layer (Planned)

### Technology
- **Phaser 3** + **TypeScript**: Pixel-art 2D game engine running in the browser.
- **Architecture**: Mixed mode — core game logic stays on the Go backend (REST API), frontend handles rendering and mini-game execution.

### Responsibilities
| Aspect | Frontend (Phaser.js) | Backend (Go) |
|---|---|---|
| Investigation | Node rendering, option UI, animations | State transitions, reputation calculation, save |
| Skills/Modules | Skill tree panel, equipment drag & drop | Unlock validation, cooldown timing |
| Defense Base | Building visuals, facility placement | Build rules, upgrade logic |
| Mini-games | Full game execution (bullet hell, cards…) | Result validation, score recording |
| Victim Profile | Radar chart visualization, strategy tips | Risk score calculation |

### Communication
- Frontend ↔ Backend via **REST API** (`internal/interface/in/http`)
- JSON request/response using existing DTO structures
- Stateless HTTP for most operations; WebSocket reserved for future real-time features
