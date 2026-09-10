# Ticket #30 final repository evidence

- Task: [RuntimeBackend MVP](../../../../work/0030-runtime-backend-mvp/task.md)
- Date: 2026-09-10
- Branch: `tomtian/e3-t1-runtime-backend-mvp`
- Baseline: `91977b83e92faaa2cde39388e3f9b44684714151` plus the cleanup-binding fix, with all PR Go sources pinned by the [source manifest](../reference-contract/source-sha256.txt).

| Gate | Exact command | Result | Output |
| --- | --- | --- | --- |
| Build | `go build ./...` | exit 0, no output | — |
| Repository | `pwsh -File scripts/check.ps1 -All` | exit 0, 12 checks pass | [output.txt](output.txt) |
| Race | `go test -count=1 -race ./...` | exit 0 | [race-output.txt](race-output.txt) |
| Consumers | `go test -count=1 -v ./internal/app/... ./internal/cli/... ./internal/toolgateway/... ./internal/modelgateway/... ./harness/e2e/...` | exit 0 | [consumers-output.txt](consumers-output.txt) |
| Integration compile | included in `pwsh -File scripts/check.ps1 -All` | pass; no cluster tests run | [excerpt from repository output](integration-compile-output.txt) |
| Explicit-context guard (retained from Slice 2) | `go test -count=1 -tags integration -run '^TestRuntimeBackend_AllocateObserveCleanup$' ./harness/integration/agentsandbox/` | expected exit 1 before any cluster call; rejection verified | [guard output](integration-context-guard-output.txt) |

Build, focused groups, race and the repository gate were rerun after the cleanup fix on the final Go sources. Only review/evidence/task bookkeeping was finalized afterwards and checked by the documentation gate. The context-guard output is the earlier accepted run at `91977b8`: its harness preflight is unchanged and the rejection occurs before adapter construction. Contract output is recorded [separately](../reference-contract/summary.md).

A passing compile or context guard is not an integration pass. The current environment has no kubectl and no confirmed test context; [real-backend verification remains blocked](../agent-sandbox/summary.md). The [technical review and #31 handoff](../../../../work/0030-runtime-backend-mvp/review.md) are complete; independent human review and acceptance remain outstanding.
