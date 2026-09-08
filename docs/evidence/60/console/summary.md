# #60 staged console evidence

- Date: 2026-09-08
- Branch: `codex/0060-claim-console`, stacked on #108 at `4954efe`.
- Approved planning packet: `0eb707b`; [independent approval](https://github.com/wunderforge/agenova/issues/60#issuecomment-5582660832).
- Implementation: `729cd90`; subsequent evidence commit contains documentation/artifacts only.
- Gate: `./scripts/check.ps1 -All`, PASS (exit 0). [Raw output](output.txt).
- Environment: Windows, Go 1.26.4, Node 24.11.1, npm 11.6.2, Playwright 1.63.0 / Chromium 153.0.8010.12. No package/lockfile changes from #59.
- Scope: fixture-driven portion only. #60 is not complete; #38/#68 remain blockers.

## Exact commands

| Command | Result |
| --- | --- |
| `npm --prefix ui test -- --run src/console.test.tsx` | PASS; 31 console route/source/component tests |
| `npm --prefix ui run test:smoke -- console.spec.ts` | PASS; 16 console browser checks |
| `npm --prefix ui run contracts:check` | PASS; canonical 12-case parity, stale-binding rejection, enum drift, new Go narrowing validation and 22 base fixture/shape tests |
| `npm --prefix ui run typecheck` | PASS |
| `npm --prefix ui test -- --run` | PASS; 71 tests in 3 files, including all 40 #59 tests |
| `npm --prefix ui run build` | PASS; 26 modules, JS 228.58 kB / 67.10 kB gzip |
| `npm --prefix ui run test:smoke` | PASS; 23 tests, including all 7 #59 smoke checks |
| `./scripts/check.ps1 -All` | PASS; 13 Go packages, vet/module/gofmt/docs/boundaries, integration-package compilation, complete frontend gate |
| `git diff --exit-code 0eb707b -- api/v1alpha1 harness/fixtures/contract/v0 ui/src/contracts.generated.ts` | PASS; canonical types, fixture payloads and generated bindings unchanged |
| `git diff --cached --check` | PASS |

Focused runs preceded the final full baseline; `output.txt` records the accepted full run after the mobile column readability adjustment. The temporary stacked PR base does not trigger the main-targeted PR workflow. Any manually dispatched Linux run is reported in the PR separately, without changing the stack or claiming a live-API test.

## Reproduce

Follow [UI prerequisites and commands](../../../../ui/README.md), then open these **client presentation** routes at Vite's printed local address:

```text
/console/requests/fix-payment-timeout
/console/claims/claim%3Afix-payment-timeout%3A1
/console/requests/fix-payment-timeout?scenario=narrowed
/console/requests/fix-payment-timeout-team-b
```

The remaining source states use `?scenario=loading`, `not-found`, `malformed`, or `unavailable` on the Team A request route. Loading deliberately stays pending; there is no timer/polling. Unknown references return not-found; malformed paths/encodings/kinds and unsupported scenario queries do not fall back to another claim.

## Source provenance

- `claim-request.valid.team-a-engineer-json` and `claim-request.valid.team-a-engineer-yaml`: canonical request intent, consumed directly from manifest inputs and normalized by Go.
- `issued-state.valid.team-a-engineer`: Allow / Running with request, template and claim references checked before pairing. Canonical requested/effective access is equal.
- `derived.console.narrowed`: starts from that Team A issued state, removes only `github.pull-request` from `effectiveAuthority.tools`, then runs `ParseSystemIssuedState` again. The Go test proves all other state and canonical rows remain unchanged. The page visibly labels the derivative; it is not an executed policy decision.
- `issued-state.valid.team-b-denial`: Deny without claim/authority/backend. No matching request document exists; Team A request intent is never reused.
- Malformed: deletes `evidence` from an in-memory clone of the selected issued state; generated shape validation emits `missing-source-field: evidence`.
- Loading/not-found/unavailable: deterministic EvidenceSource transport responses, not canonical governance records.
- All #59 invalid canonical cases retain their existing test coverage. No duplicated fixture payload or new API/evidence schema was introduced.

## Deterministic screenshots

Browser tests assert each state before capturing full-page images, disable animations/carets, wait for fonts, and verify no horizontal viewport overflow.

| State | 1100x1000 | 390x844 |
| --- | --- | --- |
| Allowed | [desktop](console-1100-allowed.png) | [mobile](console-390-allowed.png) |
| Derived narrowed | [desktop](console-1100-narrowed.png) | [mobile](console-390-narrowed.png) |
| Denied before claim | [desktop](console-1100-denied.png) | [mobile](console-390-denied.png) |
| Loading | [desktop](console-1100-loading.png) | [mobile](console-390-loading.png) |
| Not found | [desktop](console-1100-not-found.png) | [mobile](console-390-not-found.png) |
| Malformed | [desktop](console-1100-malformed.png) | [mobile](console-390-malformed.png) |
| Unavailable | [desktop](console-1100-unavailable.png) | [mobile](console-390-unavailable.png) |

The narrowing, denial and malformed views were visually inspected. Executable checks cover all 14 images plus direct claim route/reload, browser back/forward, keyboard links and scenario selection, visible focus, skip link, one main landmark, named navigation/table and polite busy announcements. This is an accessibility smoke check, not a comprehensive accessibility audit.

## Remaining completion blockers

- #38 must provide the accepted assembled evidence, then #68 the read-only API.
- HttpEvidenceSource, bounded polling, real running-to-terminal/revocation evidence, real-API G6 E2E and automated CLI/API/UI equality remain absent.
- Current fixtures provide no detailed per-call decisions/results, agent outcome or revocation proof. The console labels these as not supplied; it does not synthesize them from Running or backend identity.
- No HTTP contract, live endpoint, lineage, multi-agent UI or mutation control was added. Existing Go integration compilation is not real backend execution evidence.
- PR #111 remains stacked on #108 and unmerged. The passing fixture gate does not satisfy the full #60 Definition of Done.
