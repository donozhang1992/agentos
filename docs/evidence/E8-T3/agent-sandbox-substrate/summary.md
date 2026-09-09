# Evidence Summary

- Ticket: E8-T3 (#50), PR #99 remediation
- Gate: agent-sandbox-substrate
- Date: 2026-09-09
- Branch: codex/e8-t3-pr99-unblock
- Base commit: b031f6a10a450f71aea48a3931ab582dc8a6f4bb
- Verified script commit: f76a2b05c3f761b90afa638b1cd8f7eea1fb6fbf (clean working tree)
- Script SHA256: 9f53fe05c0a491ad07caecf9b58efc18d411c59f0be590826f30953c078c4f91
- Command: `bash harness/spike/agent-sandbox-substrate/reproduce.sh up --capture`
- Result: **blocked** — Docker is unavailable on the Windows repair machine.
- Upstream target: Agent Sandbox v0.4.6 / extensions.agents.x-k8s.io/v1alpha1

The old Darwin capture at adf4d76 predates the final script and is superseded.
Its summary/output are recoverable from Git history; they do not prove the repair.

The isolated command-double gate and repository baseline are separate from
real-backend acceptance. The baseline passed (including integration compilation);
no actual integration or race run was performed on this machine. Go emitted a
telemetry-cache permission warning, but the baseline exited 0.

| Check | Command | Result |
| --- | --- | --- |
| Bash syntax | `bash -n harness/spike/agent-sandbox-substrate/reproduce.sh` | pass |
| Isolated regressions | `bash harness/spike/agent-sandbox-substrate/test-reproduce.sh` | pass, 15 scenarios |
| Repository baseline | `pwsh -NoProfile -File scripts/check.ps1 -All` | pass |
| Actual prerequisites from committed script | `bash harness/spike/agent-sandbox-substrate/reproduce.sh up --capture` | exit 1, Docker unavailable |

Before #99 can be merged as complete, run the final committed script twice with
`all --capture` on Docker/kind and replace this blocker with one accepted capture.
The capture must identify the actual commit, script hash, tools/server/CRD versions,
controller readiness, claim Ready, pod/namespace cleanup and cluster deletion.
No Agenova claim-governance or adapter proof is asserted here.

Raw prerequisite result: [output.txt](output.txt).
