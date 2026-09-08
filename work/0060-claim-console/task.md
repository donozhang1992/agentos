# Task: Single-claim console: staged fixture-driven screen

- Ticket: [#60](https://github.com/wunderforge/agenova/issues/60)
- Mission: Make one governed assignment understandable on a routed read-only console while retaining the live requirements as explicitly unfinished work.
- Target: `ui/src/` console/route/composition, fixture adapter scenario setup, focused tests and `ui/smoke/`; contributor and task evidence documentation.
- User value: A reviewer can inspect request intent, trusted principal, decision, granted access and recorded execution evidence for one assignment without confusing fixture demonstrations with live governance proof.
- PRD outcome: [Read-only claim console](../../docs/product/prd.md#7-read-only-claim-console), staged under [#107](https://github.com/wunderforge/agenova/issues/107). This packet does not replace #60's live Definition of Done.

## Context to Read

Always:

- `AGENTS.md`
- `docs/product/prd.md`
- this task packet

Additional task-specific context:

- [Console presentation specification](spec.md).
- [#60 stable scope](https://github.com/wunderforge/agenova/issues/60), [#107 execution plan](https://github.com/wunderforge/agenova/issues/107#issuecomment-5580468473), and [single-claim scope decision](https://github.com/wunderforge/agenova/issues/107#issuecomment-5580872705).
- [AIDLC planning review](../../docs/development/AIDLC.md#from-existing-ticket-to-task-packet) and [Start a GitHub Ticket](../../docs/harness/playbooks.md#start-a-github-ticket).
- [Architecture](../../docs/product/architecture-contract.md): Evidence Surfaces, Submission and Resolution, Claim Lifecycle, Authority and Credentials, Backend Neutrality.
- [#59 task packet](../0059-react-contract-foundation/task.md), [spec](../0059-react-contract-foundation/spec.md), and [fixture foundation documentation](../../ui/README.md).
- `ui/src/evidence-source.ts`, `fixture-source.ts`, `shape-check.ts`, `contracts.generated.ts`, `App.tsx`, `main.tsx`, their current tests, `ui/contractgen/fixtures.go`, `ui/smoke/storyboard.spec.ts` and `ui/playwright.config.ts`.
- [Canonical fixture manifest](../../harness/fixtures/contract/v0/manifest.json), its selected ClaimRequest/IssuedState inputs, `api/v1alpha1/claim_request.go` and `sandbox_claim.go` for validation of derived scenarios.
- [#38 evidence assembly](https://github.com/wunderforge/agenova/issues/38), [#68 read-only API](https://github.com/wunderforge/agenova/issues/68), and [quality gates](../../docs/harness/quality-gates.md).

## Scope

In scope after independent planning approval:

- One read-only Claim Console and client-side route keyed by request reference or claim ID, using #59's EvidenceSource and fixture-source implementation.
- One correlated assignment on screen: request/reference and principal, decision/policy, requested vs effective access, recorded lifecycle, invocation evidence, backend identity and outcome availability.
- Fixture-driven loading, not-found, malformed, unavailable, allowed, explicitly derived narrowed, and denied-before-claim states.
- Responsive layout, semantic structure, keyboard focus/navigation, accessible loading/error announcements and deterministic browser screenshots.
- Focused tests and shared baseline evidence for the staged fixture portion; preserve #59 regression tests.

Out of scope / completion deferred:

- HttpEvidenceSource, live endpoint or JSON contract design, bounded polling/timers, terminal progression over HTTP, real-API G6 E2E, and CLI/API/UI equality. These remain #60 completion blockers until #38 then #68 land.
- Claim search/listing, dashboards, mutation controls, policy editing, lineage, child claims, multi-agent UI, invented facts or claim authority.
- Merging #108 or this stacked PR, or marking #60 Done.

## Acceptance Criteria

Staged acceptance only:

- A direct client route and reload select one request or claim through the injected source; unknown/malformed route references produce explicit states. Route syntax is a presentation concern, not a proposed API endpoint.
- The page displays canonical domain objects without a new governance model. Request-to-issued pairing must match request name/reference, claim/request reference and template; missing or inconsistent pairs are visible gaps.
- Requested and effective tools/scopes are distinct. Set differences describe recorded input/output only; the UI never evaluates policy or grants access.
- Canonical Team A shows Allow/Running; a labeled derived scenario removes `github.pull-request` from the effective tools and passes canonical system-state validation before rendering narrowing. Canonical files remain unchanged.
- Team B shows Deny with no claim, grant or backend. Its absent matching request document is explicitly unavailable; do not silently reuse Team A's request.
- Lifecycle shows only recorded phase/events. Detailed invocation outcome, agent result and revocation evidence absent from v0 remain visibly unavailable; backend identity remains `reference` when that is what the fixture records.
- Loading, not-found, malformed, unavailable, allowed/narrowed and denial have component and browser evidence, with no stale evidence shown during a selection/source change.
- Keyboard and 1100x1000 / 390x844 viewport checks pass; screenshots supplement assertions. The baseline passes without claiming live acceptance.

Full #60 acceptance still blocked:

- #38 supplies accepted assembled evidence and #68 exposes it.
- A later reviewed live slice implements the real HTTP source and bounded polling against that contract.
- Real API E2E proves success/denial/terminal progression; automated CLI equality and live screenshots pass. Only then may #60 be marked complete.

## Negative Case

- Invalid route encoding, wrong reference kind, unknown reference and unrelated request/claim pairing do not fall back to another claim.
- A late Allow response cannot overwrite a newly selected denial or missing/error state.
- Missing required fields, unknown enums and unavailable source states expose sanitized diagnostics and withhold validated grants.
- Empty known invocation arrays remain distinct from missing observations. Running does not imply agent success; backend cleanup/readiness does not imply outcome or revocation.
- A derived narrowing scenario is always labeled as a fixture demonstration, never as an executed policy evaluation or live result.

## Execution Todo

- [x] Inspect #60, #107, architecture, #59 and current UI; verify latest #108 head.
- [x] Create `codex/0060-claim-console` from #108 head `4954efe91c2bd081e045cc1ad37244467849ea41`.
- [x] Record the user's bounded staged-direction approval on #60 and publish this packet for independent planning review.
- [x] Obtain and record required independent human planning approval on #60 before implementation.
- [x] Add the single-claim route and fixture reference lookup through the existing EvidenceSource.
- [x] Compose the correlated single-claim presentation and visible missing-data sections.
- [x] Add canonical-parser-validated derived narrowing and deterministic source-status scenarios without copied fixture payloads.
- [x] Add component, route, reference-correlation, keyboard, viewport and browser screenshot evidence.
- [x] Run focused commands and `./scripts/check.ps1 -All`; record exact outputs and fixture/derivation identities.
- [ ] Push only #60 changes to the stacked PR, request `@codex review`, and report the implemented subset plus live blockers. Do not merge or close #60.

## Quality Gates

Planning:

- `./scripts/check.ps1 -Docs`
- `git diff --check`

After approval, staged implementation:

- `npm --prefix ui run contracts:check`
- `npm --prefix ui run typecheck`
- `npm --prefix ui test -- --run`
- `npm --prefix ui run build`
- `npm --prefix ui run test:smoke`
- `./scripts/check.ps1 -All`
- CI PR profile when retargeted to main; the current workflow only triggers PR checks for base main, so a stacked base does not automatically supply that CI evidence.

## Evidence Required

- Exact command results and source/evidence commit IDs, plus a diff against #108 proving #59 changes remain outside this PR's delta.
- Screenshots at 1100x1000 and 390x844 for single-claim allowed, derived narrowed, pre-claim denied, loading, not-found, malformed and unavailable states.
- Base IDs: `claim-request.valid.team-a-engineer-json`, `claim-request.valid.team-a-engineer-yaml`, `issued-state.valid.team-a-engineer`, `issued-state.valid.team-b-denial`; preserve all #59 invalid cases.
- Derivation log: Team A issued snapshot with `effectiveAuthority.tools` excluding `github.pull-request`, revalidated through canonical Go; malformed data based on the same snapshot with explicit missing/unknown path; injected deferred/not-found/unavailable source responses are transport fixtures, not governance facts.
- Keyboard-only navigation, visible focus, landmarks/headings, announcement assertions, horizontal-overflow checks, and direct-route reload evidence. No claim of a full accessibility audit.
- Explicit completion blockers for live API, polling, terminal revocation/progression, real-API E2E and CLI equality.

## Constraints

- Preserve `docs/product/architecture-contract.md`; do not broaden the Ticket or PRD without a recorded human decision.
- Stay single-claim per the user's instruction, #60 and #107; the older PRD's lineage ambitions do not enlarge this staged scope.
- Keep Go semantic validation and generated bindings authoritative. No frontend policy evaluation or inferred authority.
- Preserve #108 unchanged: temporary PR base is `codex/0059-react-contract-foundation`, head `codex/0060-claim-console`. After #108 is human-merged, retarget/rebase with normal conflict review, not by merging it from this task.
- Keep dependency/assignment/status truth on GitHub; #60 remains unfinished and its #38/#68 dependencies are not removed by fixture progress.

## Decisions and Blockers

- The user explicitly approved starting this bounded fixture-driven direction before #68 on 2026-09-08. That direction approval does not substitute for the independent planning approval explicitly requested in the same instruction.
- Task + Spec is the minimum depth: correlation, presentation gaps and derived narrowing require reviewed behavior. No new design/API document is justified.
- Proposed route: `/console/requests/:reference` or `/console/claims/:id`, with percent-encoded opaque values; fixture scenario selection, if needed, uses a clearly labeled presentation query parameter. This does not define HTTP resource paths or polling semantics.
- #59 returns separate request and issued results. The composition layer may load both through the existing interface, retaining their canonical types and checking correlation; it must not invent a serialized combined evidence contract. Unmatched request data stays unavailable.
- Canonical Team A requested/effective access is equal. Narrowing therefore needs a declared, canonically validated in-memory derivative. Current v0 snapshots have no detailed invocation decisions, agent outcome or terminal revocation evidence; these are display gaps owned upstream, not fields to invent.
- #38 and #68 are open. Full #60 remains blocked on their accepted live contract and evidence, regardless of staged frontend progress.
- Independent planning approval at 0eb707b was recorded on #60: https://github.com/wunderforge/agenova/issues/60#issuecomment-5582660832. The approved fixture portion is authorized; full live completion remains blocked.

### Planning evidence (2026-09-08)

- Stack base verified with `gh pr view 108 --json headRefOid` and `git rev-parse HEAD`: `4954efe91c2bd081e045cc1ad37244467849ea41`.
- `./scripts/check.ps1 -Docs`: PASS, all documentation, metadata, architecture, links, boundaries and delivery-contract checks.
- `git diff --check`: PASS. Planning delta contains only this Task + Spec; no executable source changes.
- #60-specific frontend gates and screenshots: not run/not produced; construction awaits independent planning approval. #59 evidence remains provenance for the base, not proof of this staged console.
