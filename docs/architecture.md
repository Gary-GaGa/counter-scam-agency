# System Architecture

## Overview
Based on **Clean Architecture** and **Domain-Driven Design (DDD)** principles, the system is designed to be testable, maintainable, and loosely coupled.

## Directory Structure
The project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
internal/
  ├── domain/           # [Core] Enterprise Business Rules (Pure Go, No external deps)
  │   ├── personnel/    # Context: Player, AIPartner
  │   ├── intelligence/ # Context: Mission, Evidence (Static Content)
  │   └── operation/    # Context: Investigation (Runtime State)
  │
  ├── usecase/          # [Application] Application Business Rules
  │   ├── port/         # Interfaces (Input/Output Ports)
  │   └── dto/          # Data Transfer Objects
  │
  ├── interface/        # [Adapters] Interface Adapters
  │   ├── in/           # Input: CLI, HTTP, Game Loop
  │   └── out/          # Output: Presenters
  │
  └── infrastructure/   # [Frameworks] External Details
      ├── persistence/  # Database Implementations
      │   └── mongo/    # MongoDB Adapters (per domain)
      └── memory/       # In-Memory/Cache Implementations
```

## 1. Domain Layer (Core)
The heart of the software. It contains the business objects and high-level rules.

### Bounded Contexts
We split the domain into three distinct contexts based on the "Agency" metaphor:

#### A. Personnel Context (`internal/domain/personnel`)
- **Responsibility**: Manages the agents (Players) and their equipment.
- **Key Entities**:
  - `Player` (Aggregate Root): Stats, Inventory.
  - `AIPartner` (Entity): Loadout, Personality.
  - `Stats` (Value Object): Logic, Tech, Charisma, Resilience.

#### B. Intelligence Context (`internal/domain/intelligence`)
- **Responsibility**: Manages the static case files and knowledge base.
- **Key Entities**:
  - `Mission` (Aggregate Root): The script of a case with narrative nodes and options.
  - `Evidence` (Entity): Clues and contradictions.
  - `ScamType` (Value Object): Classification of scams.

#### C. Operation Context (`internal/domain/operation`)
- **Responsibility**: Manages the runtime execution of a mission.
- **Key Entities**:
  - `Investigation` (Aggregate Root): Tracks progress, current narrative node, and collected evidence for a specific play-through.

### Repository Interfaces
The Domain layer defines **interfaces** for persistence, following the Dependency Inversion Principle.
- `PlayerRepository`
- `MissionRepository`
- `InvestigationRepository`

## 2. Infrastructure Layer
Contains implementations of the interfaces defined in the inner layers.

- **Persistence**: MongoDB is used for storage.
  - `internal/infrastructure/persistence/mongo/player`
  - `internal/infrastructure/persistence/mongo/mission`
  - `internal/infrastructure/persistence/mongo/investigation`

## 3. Usecase Layer (Next Step)
Orchestrates the flow of data to and from the entities.
- **Services**: Will implement high-level flows like `StartInvestigation`, `AdvanceNode`, `SubmitEvidence`, `CompleteInvestigation`, and reputation accumulation.

## Key Design Decisions
1.  **Repository Interface in Domain**: Allows the domain to define its storage needs without knowing the implementation details.
2.  **Rich Domain Models**: Entities contain logic (e.g., `Player.GetTotalStats()`), preventing Anemic Domain Models.
3.  **Split Infrastructure**: MongoDB implementations are split by domain context for modularity.
