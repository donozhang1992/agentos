# Technical Design: Minimum agent filesystem boundary

- Ticket: [#89](https://github.com/wunderforge/agenova/issues/89)
- Feature spec: [spec.md](spec.md)

Proposed approach; construction is blocked by the gates in [task.md](task.md).

## Current State and Constraints

At main `3365cd0e37d181146dd5b6f8e65a58e03b6fb39e`, RuntimeBackend mixes claim phases and pool methods. Its BackendClaim carries pool/input/status, and the reusable suite verifies pool replacement and lifecycle rather than directories or actual worker processes. The reference runtime stores in-memory objects; a successful existing suite is not filesystem evidence.

PR #106 at `016989ec238622a4b19de84d7f2fb0afcaf85e80` proposes allocation, observation, explicit start, termination and cleanup, retains concrete reference compatibility helpers, and narrows gateway claim lookup. It is a draft plan, not a merged contract. #89 must wait for the implemented producer and may not introduce a competing seam now.

The Agent Sandbox spike maps image/command to upstream objects and currently uses Ready as its start observation. #106 explicitly identifies that limitation. The current adapter code/backend note establishes no minimum writable-directory enforcement. This packet neither fixes that adapter nor infers isolation from its deletion/replacement behavior.

## Decision

Use one backend-selected working directory with the spec's fixed rule. Place scratch, synthetic HOME and tool caches below it, and leave ordinary local process/file APIs direct. The adapter implements the substrate boundary; shared Agenova code expresses the requirement, correlation and evidence only.

After #30 merges, inspect its actual types first. Add only the minimal neutral working-directory description and capability/evidence needed by its accepted launch/observation results. Reuse its allocation identity and errors. Do not add a general file API, a second lifecycle registry, a new public Workspace resource, or a host-path mapping. If no accepted extension point can carry these values, stop for a reviewed producer/consumer decision rather than editing #30 for convenience.

The minimum output demonstration is a file or bundle explicitly produced under W and exported while work is active. A test-only collector verifies bytes into a separate harness-owned destination. It simulates the external handoff, is never worker-mounted, and is not a proposed RuntimeBackend filesystem API. A real consumer must use its governed output/tool operation; this ticket does not invent one. Termination does not wait for a human or guarantee final output salvage.

## Ownership and Contract Boundaries

| Surface | Later #89 responsibility | Boundary retained |
| --- | --- | --- |
| internal/runtime | Minimal working-directory semantics attached to the final #30 seam | No provider types or application outcome ownership |
| internal/operator | Reference model correlated to claim/allocation | Simulated enforcement, no hostile isolation claim |
| internal/runtime/contracttest | Reusable success/fault cases from the spec | Factories hide implementation setup; no shared syscall proxy |
| Local compatibility fixture | Known repo, real git/compiler/test commands, output bytes | Test-only trusted workload; no execution of arbitrary user code |
| #48 / adapter note | Consume accepted capability requirements and list gaps | Provider layout/objects stay with adapter owner |
| #51 | Real working directory, mounts, negative access and cleanup | G5 evidence cannot be inferred from G2 |
| #52 / gateways | Preserve credential and bypass boundaries | Local file permission is not external-access authority |
| #53 | Consume cwd and export boundary during integration | Fixture/mock artifact can proceed independently |

## Alternatives Considered

- Send every file access through Tool Gateway: incompatible with ordinary git/compiler/process use and creates an unnecessary filesystem product.
- Expose configurable host mounts or a generic path allowlist: enlarges authority and requires a policy language, host semantics and additional credential risks.
- Enforce only cwd or string-prefix validation: compatible with tools but does not constrain process access, links or traversal; useful test setup cannot stand in for isolation.
- Persistent workspace with post-run retrieval: adds storage/access lifecycle and retention product scope. Prefer explicit pre-termination export and ephemeral cleanup.
- Force the current Agent Sandbox spike to define the rule: would reshape a shared contract for provider convenience. Keep unsupported gaps visible.

## Verification Strategy

Three evidence levels are deliberately separate:

1. **Reference model (simulated):** reusable cases exercise identity, supported operations, denial and cleanup with a modeled filesystem and synthetic outside/other-claim sentinels. Include injected termination/cleanup/export failures. Alias and permission checks here prove the model's contract decisions only. Do not assert that the in-memory backend can contain arbitrary native programs.
2. **Local process compatibility (real local commands, no isolation):** a harness-owned temporary root contains W, synthetic runtime/outside fixtures and a separate collector destination. Build a tiny deterministic repository with local git identity/config, set process cwd to W, redirect HOME/temp/caches under W, edit a source file, run git diff and a compiler/test command, and compare exported bytes/digest. Disable inherited git config/helpers and use no live credentials/network. The output collector uses explicit relative regular files, rejects aliases and special files, and enforces a documented fixture byte cap (proposed 1 MiB). Stop/quiesce the trusted writer before collecting so it cannot race path validation. No arbitrary code or native outside-write attack runs on the host in this fixture. Model denial is recorded separately.
3. **Real backend (#51, not executed by this ticket):** repeat FS-P1 and selected negative cases as the actual worker identity, capture native nonzero errors and unchanged synthetic sentinels, backend layout/mounts and worker configuration, then termination/cleanup observations. Include an alias escape attempt, fresh second claim and absence of host/credential mounts. Document supported OS/profile and residual gaps; do not replace real isolation checks with collector validation or model denials.

Future reproduction starts with the focused commands in task.md. Implementation must publish exact case names, fixture paths, command versions, exit codes and artifact identities; there is no new runnable fixture in this planning PR. The existing -All gate checks repository regression and compiles the integration package; it does not run a cluster or prove this feature.

## Risks and Compatibility

- Read-only runtime data requires tool cache/temp configuration; hard-coded writes outside W are unsupported, not a reason to loosen the rule. Runtime special facilities require explicit #48/#51 review.
- Link resolution alone has race hazards. Test collector validation for trusted local compatibility is not hostile isolation; the real backend must enforce the rule for worker native I/O independently.
- A source fixture may be prepared before start, but production acquisition must use already-authorized access. No global git helper, host home or provider secret is inherited to make clone/push work.
- Early or failed termination can lose unexported output. Failed cleanup must not silently retain a reusable workspace or change the work outcome; confirmed cleanup remains separate resource evidence.
- #30 may choose different names/shapes and error behavior. Reconcile this plan after its implementation merge without changing claim semantics; leave all Go contracts untouched in this PR.
- No changes are needed to PRD, architecture contract or project-status for a proposal. Once behavior is implemented and merged, update only the evidence/status sources whose owned facts actually changed.
