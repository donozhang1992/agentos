# Ticket #89 filesystem-boundary evidence

## Result

- The shared runtime seam now reports one backend-selected, claim-scoped working directory and an explicit filesystem evidence level.
- The reference backend passes reusable model cases for task read/write, traversal and outside-boundary denial, replacement freshness, and denial after termination or cleanup. These cases are labelled `Simulated`; they do not claim process isolation.
- A separate trusted local fixture successfully initialized and committed a small Git repository, edited it, ran `go test ./...`, exported one bounded regular file, deleted the workspace, and retained only the exported bytes and digest. This proves ordinary tool compatibility, not isolation.
- The Agent Sandbox adapter reports `Unsupported` until #48 maps the substrate and #51 supplies real worker isolation evidence.

## Commands

All commands completed successfully on 2026-09-10:

```text
go test -count=1 -v ./internal/operator/... ./internal/runtime/...
go test -count=1 ./internal/operator/... ./internal/runtime/...
go test -count=1 ./internal/toolgateway/... ./internal/modelgateway/... ./harness/e2e/...
pwsh -NoLogo -NoProfile -File scripts/check.ps1 -All
```

The local compatibility run reported `git_diff_bytes=192`, `go_test=pass`, exported SHA-256 `5501a550dfeaa6aa1cf476e3c06bcb5c3e05d038b4f40402dc43aa6c458b9186`, and `workspace_retained=false`.

## Evidence boundary

The reference model and local fixture are deterministic contract and compatibility evidence only. They do not prove containment against a hostile process. Ticket #51 owns that proof on the real Agent Sandbox backend using the same semantics and negative cases.
