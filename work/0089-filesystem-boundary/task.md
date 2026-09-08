# Task: Define and test the minimum agent filesystem boundary

- Ticket: [#89](https://github.com/wunderforge/agenova/issues/89)
- Mission: Make one claim's writable task directory explicit without turning Agenova into a filesystem proxy or workspace service.
- Target: This planning PR changes only this Task + Spec + Design; later implementation targets the merged #30 runtime seam, reference backend and reusable contract cases.
- User value: An agent can prepare a repository, edit, compile, test and produce an output within a documented boundary whose evidence limitations are visible.
- PRD outcome: [Backend-neutral execution](../../docs/product/prd.md#3-backend-neutral-execution).

## Context to Read

Always:

- [Agent routing](../../AGENTS.md)
- [PRD](../../docs/product/prd.md)
- this task packet

Additional task-specific context:

- [Specification](spec.md) and [design](design.md).
- [Architecture contract](../../docs/product/architecture-contract.md): Backend Neutrality, Claim Lifecycle, Authority and Credentials, Scope Discipline.
- [AIDLC](../../docs/development/AIDLC.md), [playbooks](../../docs/harness/playbooks.md): Start a GitHub Ticket, Change a Core Contract, Add or Change a Runtime Backend, Elaborate Parallel Work.
- [Quality gates](../../docs/harness/quality-gates.md) and [gotchas](../../docs/harness/gotchas.md).
- [Product tour](../../docs/project-design.md): worker ownership, lifecycle and backend explanation only; this is not a source of implemented schemas.
- [Runtime contract](../../internal/runtime/backend.go), [reusable cases](../../internal/runtime/contracttest/run.go), [reference runtime](../../internal/operator/runtime.go) and [tests](../../internal/operator/runtime_test.go).
- [Adapter](../../internal/runtime/agentsandbox/adapter.go), [backend note](../../docs/backends/agent-sandbox.md), [canonical template](../../api/v1alpha1/agent_template.go).
- Producer [#30](https://github.com/wunderforge/agenova/issues/30) and draft [#106](https://github.com/wunderforge/agenova/pull/106); consumers [#48](https://github.com/wunderforge/agenova/issues/48), [#51](https://github.com/wunderforge/agenova/issues/51), [#52](https://github.com/wunderforge/agenova/issues/52), [#53](https://github.com/wunderforge/agenova/issues/53).

## Scope

In scope:

- Propose a backend-neutral, claim-scoped writable directory, deterministic outside-boundary rule, lifecycle and bounded output handoff.
- Specify compatibility, positive/negative cases and simulated versus real evidence before implementation.
- After the implementation gate opens, add minimum shared semantics and reference evidence to the accepted #30 seam and document capability gaps for consumers.

Out of scope:

- This PR: Go changes, executable filesystem fixtures, runtime implementation or adapter changes.
- Host home, arbitrary host paths, unrelated workspaces, host credentials, provider types in shared contracts, syscall interception, FUSE and per-file audit.
- Managed or persistent Workspace, Web IDE, Portal editing, interactive human access, orchestration or new gateway authority.
- Real isolation evidence (#51), provider mapping (#48), direct-egress experiment (#52), and the engineer artifact implementation (#53).

## Acceptance Criteria

- One backend-selected working directory is writable for the claim, including task files, scratch and tool caches; callers cannot select a host mount.
- A prepared example repository supports ordinary git, compiler/test and agent-process filesystem operations directly, without per-syscall Tool Gateway calls.
- Outside it, runtime files are read-only and all other filesystem data is unavailable through the supported configuration; unsupported enforcement is explicit.
- Host home, unrelated claims/workspaces and long-lived provider credentials are absent from supported worker configuration.
- Shared semantics and evidence use claim/backend identity and neutral values, with no Kubernetes, container or provider types.
- Reusable reference cases cover directory behavior and an outside-boundary denial with an explicit simulation label; local process compatibility is separately demonstrated.
- Termination ends worker access, cleanup releases the directory without workspace retention, and only outputs exported before termination survive through the supported handoff.
- Tool/Model Gateway authorization and credential boundaries remain unchanged.

## Negative Case

The [spec case matrix](spec.md#negative-cases) covers traversal and alias escape, outside mutation, cross-claim access, forbidden configuration, missing capability, late output export and failed termination/cleanup. Each denial must leave the relevant synthetic sentinel unchanged and must not fabricate success or affect another claim.

## Execution Todo

- [x] Scout current code, requested issues, PR #106 and repository planning requirements.
- [x] Draft the canonical Task + Spec + Design with compatibility and negative cases.
- [ ] Record Owner and named independent Reviewer approval of this packet in #89.
- [ ] Verify #30's final implemented RuntimeBackend contract is merged; reread its accepted types/tests and reconcile this packet. Merging planning documents alone does not satisfy this gate.
- [ ] Only after both gates: implement the smallest filesystem semantic slice against that merged seam, with reference cases and explicit simulation labels.
- [ ] Add a controlled local repository/compile/test/output compatibility fixture; keep it separate from security evidence.
- [ ] Complete fault cases, capability/gap handoff, focused and repository gates; obtain independent acceptance. Do not merge under this task's current authorization.

## Quality Gates

Planning:

- `pwsh -NoProfile -File scripts/check.ps1 -Docs`
- `pwsh -NoProfile -File scripts/check-pr-body.ps1 -BodyPath .tmp/0089-pr-body.md`
- `git diff --cached --check`
- `pwsh -NoProfile -File scripts/check.ps1 -All`

After implementation is unblocked:

- `go test -count=1 -v ./internal/operator/... ./internal/runtime/...`
- `go test -count=1 -v ./internal/toolgateway/... ./internal/modelgateway/... ./harness/e2e/...`
- `pwsh -NoProfile -File scripts/check.ps1 -All`
- G2 reference contract evidence is the strongest #89 gate; #51 owns G5 real-backend evidence. Planning checks prove neither.

## Evidence Required

- PR records base/commit, exact commands, exit results and limitations; no new behavior is claimed by this packet.
- Implementation records case IDs from the spec, directory identity, command cwd/exit/output, synthetic sentinel comparisons, cleanup/error observations and output digest/size.
- Reference evidence labels policy simulation, real local command execution and unverified isolation separately. Existing test commands without the new named cases cannot count as #89 completion.
- #48 receives the accepted semantics and capability/gap table; #51 supplies real mount/layout, worker negative access and cleanup evidence using the same fixture; #52 retains the egress boundary.

## Constraints

- Preserve the [architecture contract](../../docs/product/architecture-contract.md); do not broaden the PRD.
- Do not modify or freeze shared RuntimeBackend/filesystem Go contracts before final #30 implementation merge and recorded Owner/independent planning approval.
- Do not reshape #30 or Agent Sandbox semantics for convenience. If the merged producer cannot carry an accepted requirement, record the conflict and request a maintainer decision.
- No real secrets or host credential locations may be read by fixtures; use synthetic data only.
- GitHub owns mutable assignment and delivery status; this packet records execution gates and provenance.

## Decisions and Blockers

- Planning baseline: main `3365cd0e37d181146dd5b6f8e65a58e03b6fb39e`, inspected 2026-09-08. PR #106 was draft/unmerged at head `016989ec238622a4b19de84d7f2fb0afcaf85e80`; its operation names are proposals, not accepted API.
- The Owner authorized the overall direction and publishing this draft planning PR. That does not constitute the independent planning approval required by AIDLC.
- **Implementation is blocked on the final #30 RuntimeBackend implementation being merged AND Owner/independent Reviewer approval recorded in #89.** Automatic Codex review is advisory and cannot supply human approval.
- Owner decisions: confirm the exact outside-boundary rule, pre-termination output export/no workspace retention, and the compatibility/evidence split described in the spec. Name an independent Reviewer and record the decision in #89.
- Planning review correction: use independent-claim isolation rather than parent/child lineage, and do not imply that Agenova already has a generic artifact-output API. The integration must select an approved governed operation; the collector remains test-only.
- The requested `docs/aidlc.md` does not exist at this baseline; `docs/development/AIDLC.md` is the canonical workflow linked by AGENTS.md.
