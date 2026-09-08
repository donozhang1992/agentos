# Feature Specification: Single-claim console: staged fixture-driven screen

- Ticket: [#60](https://github.com/wunderforge/agenova/issues/60)
- PRD outcome: [Read-only claim console](../../docs/product/prd.md#7-read-only-claim-console), fixture stage only under [#107](https://github.com/wunderforge/agenova/issues/107).

## Intent

Present one governed assignment as a readable story using the already accepted canonical types and fixture source. Make data availability and fixture provenance visible while the assembled evidence and live API are still pending.

## In Scope

- One client route and one read-only page, canonical-object correlation, requested/effective comparison, explicit source states, responsive/accessibility smoke and deterministic evidence.
- Minimal, labeled in-memory scenario derivation; no copied payload files.

## Out of Scope

- HTTP contract or client, polling, real API/CLI equality, live terminal progression, lineage, multi-agent UI, mutation or policy evaluation.

## Requirements

- Given a request-reference or claim-ID client route, load that identity through the existing EvidenceSource, including direct navigation/reload and browser history. Reject malformed encoding and unknown reference kinds visibly; never select the first fixture as a silent fallback.
- Given request and issued results, retain their generated types and compare only when `request.metadata.name == issued.requestRef`, `claim.requestRef` (if present) matches, and template references agree. A mismatch withholds the comparison and identifies the gap; do not combine unrelated workers.
- Given the canonical Team A request and issued snapshot, show request task/runtime, trusted principal from issued state, policy decision/reason, requested/effective access, Running phase and recorded ClaimRunning event, reference backend/worker identity and known-empty invocation lists.
- Given Team A's equal requested/effective access, label equality honestly. For the derived narrowing demonstration, remove only `github.pull-request` from the issued effective tools before canonical Go system validation. Show that tool as requested but not granted and label the entire scenario as derived. The browser compares recorded sets; it neither calculates policy nor changes the grant.
- Given Team B's denial, show its principal, policy and denial reason, request reference, no claim issued, no authority issued, no backend allocated in this snapshot, and no matching request document supplied. Do not attach Team A request intent to it.
- Given a Running snapshot, show that no terminal agent outcome is supplied. Given missing outcome, detailed invocation decisions, timestamps or revocation evidence, show a specific unavailable/contract-gap label. Do not fabricate lifecycle steps, events, success, cleanup or revocation from backend readiness.
- Given known-empty invocation arrays, show zero recorded calls. Given missing arrays or malformed/unknown fields, show the source diagnostic, not zero. Future detailed invocation behavior remains owned by canonical upstream evidence.
- Given delayed, missing, malformed or failed fixture-source responses, visibly render loading, not-found, malformed and unavailable. These are controlled transport scenarios, not invented governance results. Clear previous claim data while loading; ignore stale completions after navigation/unmount.
- Given keyboard-only use, route links and scenario controls have accessible names, focus remains visible, the page has meaningful landmarks/headings, and loading/errors have appropriate announcements. Screenshots are captured deterministically at 1100x1000 and 390x844 after state assertions; no horizontal page overflow.
- Given a source replacement with equivalent canonical objects, the page semantics remain unchanged. No component imports fixture files or uses scenario labels as governance truth.

## Negative Cases

- Cross-request/template mismatch, malformed encoding, absent/unknown IDs, stale Allow after a newer Deny, rejected source promises, and missing/unknown fields.
- Denial cannot display claim/grant/backend information from the previous case.
- Missing evidence cannot become successful outcome, proof of revocation, or zero observed calls.
- Derived narrowed data cannot be mistaken for canonical or live evidence; canonical fixture files must remain unchanged.

## Compatibility

- Reuse #59 EvidenceSource and generated ClaimRequest/IssuedState types. Any UI-local pairing/provenance structure is non-serialized composition only; it is not an API/evidence contract.
- The old fixture storyboard and its tests remain useful regression coverage. Console presentation may reuse components without changing their underlying semantics.
- Canonical parsers validate positive derived domain data before the fixture source exposes it. Invalid display cases remain sanitized diagnostics.
- Route syntax and fixture scenario transport are presentation details only. #38 and #68 determine the later real evidence/HTTP semantics.
- No completion claim for #60: its HTTP source, bounded polling, real-API G6 E2E and CLI equality are intentionally pending.

## Open Decisions

- Independent human Reviewer approval of this packet is required before construction, particularly the correlation rules, declared narrowing derivative and honest outcome/invocation gaps.
- No live schema decision is delegated to this ticket; #38/#68 remain the producing authorities.
