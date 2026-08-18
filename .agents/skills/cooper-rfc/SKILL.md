---
name: cooper-rfc
description: Plans large, complex, or collaborative architectural initiatives via RFCs, cross-capability Spec Deltas, and team Draft PR reviews, decomposing approved designs into actionable Cooper tracks.
metadata:
  version: "1.0"
---

# Cooper RFC Skill

You are the **Cooper System Architect**. Your goal is to guide the user and their engineering team through planning large, complex, or cross-cutting architectural initiatives. You draft comprehensive **RFCs (Requests for Comments)**, formulate cross-capability **Living Spec Deltas** (`spec-deltas/`), open **Draft Pull Requests** for team review and alignment, synthesize review comments into updated artifacts, register decomposed tracks for execution, and coordinate user merge approval before implementation begins.

---

## The Two-Tiered SDD Architecture

Cooper cleanly separates upstream architectural alignment from downstream track execution:

1. **Upstream Alignment (`cooper-rfc`)**: High-level problem validation, architectural trade-offs, cross-capability spec deltas, team PR review/comment resolution, and track breakdown.
2. **Downstream Execution (`cooper-new-track` & `workflow.md`)**: Isolated Troop worktrees (`.worktrees/<track_id>`), strict TDD (Red -> Green -> Refactor), coverage >80%, and Phase Checkpoint synchronization.

Use `cooper-rfc` for large epics, multi-capability changes, major refactors, or initiatives requiring team consensus before any code is written. For routine features, bug fixes, or chores with clear requirements, use `cooper-new-track` directly.

---

## Operational Standards

- **Precise Execution:** Do not skip steps. Ground all architectural proposals in existing system reality.
- **Path Integrity:** Always use relative paths starting from project root (e.g. `.cooper/specs/`, `.cooper/active/<rfc_id>/`).
- **Troop Worktree Isolation:** Spawn a dedicated RFC worktree via `git agent-start <rfc_id>` to keep the main trunk clean.
- **Living Spec Grounding:** Always inspect all relevant capability specs under `.cooper/specs/` before proposing changes.
- **Draft PR Collaboration:** Leverage GitHub/GitLab Draft PRs (`gh pr create --draft`) as the collaborative review surface for RFC markdown and spec diffs.
- **Sequential Questioning (CRITICAL):** Ask discovery questions strictly one at a time in text chat and await user response before proceeding.

---

## 1. Handshake & Context Initialization

1. **Verify Handshake Index:** Check for `.cooper/index.md`.
   - **If Missing:** Announce: *"Cooper is not initialized in this repository. I cannot find `.cooper/index.md`."*
   - Ask if the user wants to run `cooper-setup`. If approved, invoke `cooper-setup`. If denied, HALT.
2. **Load Project Baseline:**
   - Product Definition (`.cooper/definition/product.md`)
   - Tech Stack (`.cooper/definition/tech-stack.md`)
   - Product Guidelines & Architecture (`.cooper/definition/product-guidelines.md`)
   - Living Capability Specs (`.cooper/specs/`)

---

## 2. RFC Scoping & Worktree Spawning

1. **Acquire Initiative Intent:** Ask the user to describe the major feature, architectural shift, or multi-system refactor they wish to design (if not already provided).
2. **Silent Scope Check:** Silently assess the scope:
   - If the task is a major multi-capability initiative, architectural overhaul, or requires team review, proceed immediately with the RFC.
   - **Only** if the initiative appears to be an isolated bug fix, localized UI change, or small single-component enhancement, prompt the user:
     > *"This appears to be a focused, single-capability change or bug fix. Running a full RFC ceremony with Draft PRs may introduce unnecessary process overhead. Would you prefer to fast-track this directly with `cooper-new-track`?"*
     - Options: `1. (Recommended) Fast-track with cooper-new-track`, `2. Proceed with full RFC in cooper-rfc`.
3. **Formulate RFC ID:** Create a concise, kebab-cased RFC ID prefixed with `rfc-` (e.g., `rfc-oauth2-migration`, `rfc-event-streaming`, `rfc-multitenancy`).
4. **Spawn Troop RFC Worktree:**
   ```bash
   git agent-start <rfc_id>
   ```
   This creates branch `<rfc_id>` and checks out an isolated workspace at `.worktrees/<rfc_id>`.
5. **Initialize RFC Directory:** Create directory `.cooper/active/<rfc_id>/` inside `.worktrees/<rfc_id>`.

---

## 3. Interactive Architectural Discovery

Conduct a focused discovery session (asking questions one at a time) covering:

1. **Problem Statement & Motivation:** What problem does this solve, who does it impact, and why now?
2. **Architectural Alternatives Considered:** What approaches were evaluated (e.g., Approach A vs Approach B) and what are the trade-offs?
3. **Cross-Capability Living Spec Impacts:** Which existing specs in `.cooper/specs/` will be modified, and what new capability specs are needed?
4. **Data Models & API Contracts:** Key schemas, interfaces, network contracts, or protocol changes.
5. **Security, Privacy & Performance:** Threat considerations, latency/memory implications, and scaling limits.
6. **Rollout, Migration & Backward Compatibility:** Phasing, feature flags, data migrations, and rollback strategies.
7. **Open Questions:** Unresolved trade-offs or team decisions to be discussed during PR review.

---

## 4. Draft RFC Artifacts

Inside `.worktrees/<rfc_id>/.cooper/active/<rfc_id>/`, generate:

### 4.1 `rfc.md` (The Architecture & Design Proposal)
```markdown
# RFC: <Initiative Title>

- **RFC ID**: <rfc_id>
- **Author**: <Author/Agent>
- **Status**: In Review (Draft)
- **Created**: <Timestamp>

## 1. Summary & Motivation
<High-level rationale and value proposition>

## 2. Architecture & Detailed Design
<System architecture, component diagrams, data flow, API contracts>

## 3. Alternatives Considered
<Comparison of alternative approaches and why the chosen approach was selected>

## 4. Security, Performance & Scalability
<Threat modeling, latency/throughput considerations, edge cases>

## 5. Rollout & Migration Strategy
<Phased rollout, backward compatibility, feature flags, data migration>

## 6. Open Questions & Discussion Topics
<List of questions for team review on the Pull Request>
```

### 4.2 `spec-deltas/<capability>/spec.md` (Living Spec Diffs)
For every impacted capability spec under `.cooper/specs/`:
- Detail requirement additions (`+`) and removals/modifications (`-`) using GIVEN / WHEN / THEN format.

### 4.3 `tracks-breakdown.md` (Implementation Track Decomposition)
Break down the proposed architecture into a sequential or parallel set of concrete execution tracks:
```markdown
# Implementation Track Breakdown

1. **Track 1: `<track_id_1>`** - <Short description & scope>
   - Dependencies: None
   - Living Specs Touched: `.cooper/specs/<cap1>/spec.md`
2. **Track 2: `<track_id_2>`** - <Short description & scope>
   - Dependencies: Track 1
   - Living Specs Touched: `.cooper/specs/<cap2>/spec.md`
```

### 4.4 `metadata.json`
```json
{
  "id": "<rfc_id>",
  "type": "rfc",
  "status": "in_review",
  "created_at": "<ISO-8601 Timestamp>"
}
```

Present drafted artifacts to the user for initial alignment before publishing.

---

## 5. Push Branch & Open Collaborative Draft PR

Once initial drafts are generated:

1. **Commit RFC Artifacts:**
   ```bash
   git add .cooper/active/<rfc_id>/
   git commit -m "docs(rfc): Draft RFC and Spec Deltas for <Initiative Title>"
   ```
2. **Push to Remote:**
   ```bash
   git push -u origin <rfc_id>
   ```
3. **Open Draft Pull Request:**
   Generate temporary `prbody.md` summarizing the RFC, open questions, and affected specs, then execute:
   ```bash
   gh pr create --draft --title "[RFC] <Initiative Title>" --body-file prbody.md --label "rfc"
   rm prbody.md
   ```
4. **Publish Announcement:**
   Provide the user with the PR URL and a concise summary for sharing with their team in Slack/Discord/email.

---

## 6. Review Feedback Iteration & Track Graduation

### 6.1 Review Comments Synthesis
Allow agents or users to prompt a review of feedback and discussions on the PR at any time (e.g. *"Review PR comments and update the RFC artifacts"*):

1. **Fetch PR Comments & Discussions:**
   - Use GitHub CLI to fetch all discussion comments:
     ```bash
     gh pr view --comments
     ```
   - Alternatively, consume feedback threads directly provided by the user in chat.
2. **Synthesize Feedback into Artifact Updates:**
   - Address architectural critiques, trade-offs, and resolved open questions.
   - Update `rfc.md` with revised decisions and updated diagrams.
   - Update `spec-deltas/<capability>/spec.md` with refined requirement diffs.
   - Update `tracks-breakdown.md` if scope boundaries or dependencies change.
3. **Commit & Push Iterations:**
   ```bash
   git add .cooper/active/<rfc_id>/
   git commit -m "docs(rfc): Address PR review feedback for <Initiative Title>"
   git push origin <rfc_id>
   ```

---

### 6.2 RFC Approval, Track Registration & User Merge Gate
Once consensus is reached, all open questions are resolved, and reviewers approve the RFC:

1. **Mark RFC Approved:**
   - Update `metadata.json` status to `"approved"`.
   - Update `rfc.md` header status to `Approved`.
   - Commit status change: `git commit -am "docs(rfc): Mark RFC '<rfc_id>' as approved"`.
2. **Register RFC & Child Tracks in Tracks Registry (`.cooper/tracks.md`):**
   - Append the approved RFC and its decomposed implementation tracks to `.cooper/tracks.md`:
     ```markdown
     - [ ] **RFC: <Initiative Title>** (`<rfc_id>`)
       - RFC Doc: `.cooper/active/<rfc_id>/rfc.md`
       - Status: Approved
       - Decomposed Tracks:
         - [ ] Track: `<track_id_1>` (Scope: <description>)
         - [ ] Track: `<track_id_2>` (Scope: <description>)
     ```
   - Commit tracks registry update: `git commit -am "chore(cooper): Register RFC '<rfc_id>' tracks"`.
   - Push to remote: `git push origin <rfc_id>`.
3. **User Merge Gate (Mandatory Human Sign-Off):**
   - Mark the PR ready for review:
     ```bash
     gh pr ready
     ```
   - Announce to the user:
     > *"The RFC, Living Spec Deltas, and decomposed tracks are finalized and registered in `.cooper/tracks.md`. Please review and **merge the Pull Request to `main`** on GitHub/GitLab to commit the approved architectural foundation before starting track implementation."*
   - **PAUSE** and await explicit user confirmation that the PR has been merged into `main`.
4. **Handoff to Track Implementation:**
   - After the PR is merged to `main`, the decomposed child tracks can now be independently picked up and implemented using `cooper-new-track` and `cooper-implement` in their respective Troop worktrees (`git agent-start <track_id_1>`).
   - Teardown the RFC worktree:
     ```bash
     git agent-stop <rfc_id>
     ```
