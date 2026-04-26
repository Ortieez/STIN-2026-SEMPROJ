# STIN-2026 Semestral Project Implementation Plan

## Objective
Finish the STIN-2026 Semestral Project based on the requirements in `DSP - STIN 2026.md`. The development follows a **Feature-Centric (Vertical Slices)** branching strategy from `develop`, where each phase is delivered end-to-end, tested, and reviewed via a GitHub Pull Request before moving to the next.

Maintain KISS, DRY, YAGNI, and SOLID principles. Avoid unnecessary comments.

## Key Files & Context
- `backend/`: Go backend (Gin). Endpoints (`/latest`, `/strongest`, `/weakest`, `/average`), caching (`cache.go`), storage, and auth.
- `frontend/`: Needs full initialization in React.
- `DSP - STIN 2026.md`: Source of truth for requirements.

## Implementation Steps

### Phase 1: Backend Testing & CI/CD Pipeline (`feature/tests-cicd`) - COMPLETE
- **Backend:** Write unit tests for `main.go`, `api.go`, and `cache.go` (>80% coverage).
- **CI/CD:** Setup GitHub Actions for Build, Test, and Coverage.
- **Verification:** Pipeline passes on PR.

### Phase 2: Backend Persistence & Authentication (`feature/backend-persistence-auth`) - IN PROGRESS
- **Backend:** 
  - Implement persistent JSON storage for `UserSettings` and Logs.
  - Create hashed login logic (client sends hash, server compares against locally hashed `.env` credentials).
  - Secure endpoints with token middleware.
- **Verification:** 
  - Endpoints return 401 without token.
  - Login returns 200 with valid hashed credentials.
  - **New Requirement:** Every new endpoint MUST be tested using yaak (the MCP server).

### Phase 3: Frontend Setup & Authentication UI (`feature/frontend-auth`)
- **Frontend:** 
  - Initialize React framework.
  - Implement Login screen with client-side hashing (SHA256).
  - Secure frontend routes.
- **Verification:** User can log in and store token in localStorage.

### Phase 4: Frontend Settings & Dashboard UI (`feature/frontend-ui`)
- **Frontend:** 
  - Implement Settings UI (Base currency, tracked currencies).
  - Implement Dashboard (Latest Rates, Analytics charts/tables).
  - Handle API error states.
- **Verification:** End-to-end integration works; UI updates based on backend API responses.

## Verification & Workflow
1. For each phase, work will be committed to its respective `feature/*` branch.
2. For every new endpoint created, test it using yaak (the MCP server).
3. The agent will commit (multiple times if needed) and push to the respective `feature/*` branch and pause, allowing you to create a PR to `develop`.
4. You will review the code, test the feature, and merge it.
5. Once merged, the agent will move to the next phase.
