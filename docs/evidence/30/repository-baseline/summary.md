# Ticket #30 Slice 2 repository evidence

- Task: [RuntimeBackend MVP](../../../../work/0030-runtime-backend-mvp/task.md)
- Date: 2026-09-10
- Branch: `tomtian/e3-t1-runtime-backend-mvp`
- Baseline: `a72baf9d1ee228fb88717e56d5f69c96847264f5` plus Slice 2 sources pinned by the [source manifest](../reference-contract/source-sha256.txt).

| Gate | Exact command | Result | Output |
| --- | --- | --- | --- |
| Build | `go build ./...` | exit 0, no output | — |
| Repository | `pwsh -File scripts/check.ps1 -All` | exit 0, 12 checks pass | [output.txt](output.txt) |
| Race | `go test -count=1 -race ./...` | exit 0 | [race-output.txt](race-output.txt) |
| Consumers | `go test -count=1 -v ./internal/app/... ./internal/cli/... ./internal/toolgateway/... ./internal/modelgateway/... ./harness/e2e/...` | exit 0 | [consumers-output.txt](consumers-output.txt) |
| Integration compile | `go test -count=1 -tags integration -run '^$' ./harness/integration/agentsandbox/` | exit 0; no cluster tests run | [compile output](integration-compile-output.txt) |
| Explicit-context guard | `go test -count=1 -tags integration -run '^TestRuntimeBackend_AllocateObserveCleanup$' ./harness/integration/agentsandbox/` | expected exit 1 before any cluster call; rejection verified | [guard output](integration-context-guard-output.txt) |

The full gates ran on the final Go sources. Only task completion and these evidence documents were added afterwards; the documentation gate checks those final changes. Contract output is recorded [separately](../reference-contract/summary.md).

A passing compile or context guard is not an integration pass. The current environment has no kubectl and no confirmed test context; [real-backend verification remains blocked](../agent-sandbox/summary.md). Final human acceptance and the #31 interface handoff remain Slice 3 work.
