# Ticket #30 Slice 2 contract evidence

- Task: [RuntimeBackend MVP](../../../../work/0030-runtime-backend-mvp/task.md)
- Date: 2026-09-10
- Branch: `tomtian/e3-t1-runtime-backend-mvp`
- Baseline: `a72baf9d1ee228fb88717e56d5f69c96847264f5` plus the Slice 2 source contents in [source-sha256.txt](source-sha256.txt).
- Command: `go test -count=1 -v ./internal/operator/... ./internal/runtime/...`
- Result: exit 0. [Raw output](output.txt).

The shared suite has twelve cases. Four new cases inject one failure below RuntimeBackend at the worker operation: start, termination, replacement, and implicit termination during cleanup. They prove errors preserve identity, do not fabricate release, and allow the documented retry. The cleanup case proves retry does not repeat a successful termination or restart work.

The injection wrapper lives only in operator tests and delegates non-failing calls to the actual pool. Production state handling is exercised; a stub RuntimeBackend does not stand in for it. Six pool-specific legacy reference regressions remain in operator tests. The adapter tests also prove legacy StartClaim cannot turn readiness into Running; legacy pending-expiry bookkeeping is preserved as a unit case.

These tests prove resource-operation semantics. Application outcome ownership and end-to-end outcome publication are delivered by #31, and no run service was added here. Agent Sandbox unit doubles and the integration package compile are not real-cluster evidence; see the [explicit blocker](../agent-sandbox/summary.md).
