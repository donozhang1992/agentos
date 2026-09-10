# Evidence Summary

- Ticket: E8-T3 (#50)
- Gate: agent-sandbox-substrate
- Date: 2026-09-10
- Branch / commit: neo/e8-t3-kind-agent-sandbox / dab25fd (run below used identical script content)
- Command: `bash harness/spike/agent-sandbox-substrate/reproduce.sh all --capture`
- Host: Darwin arm64, Docker Desktop, existing Homebrew `kind v0.32.0` + `kubectl v1.36.2` (both resolved as `existing`)
- Upstream target: Agent Sandbox **v0.4.6** / `extensions.agents.x-k8s.io/v1alpha1`
- Result: **partial — blocked at controller image pull by this network's TLS interception, not by the substrate or the script**

## What the real run proved (`output.txt`)

- `kind` and `kubectl` resolved from `PATH` as `existing` and recorded; nothing downloaded.
- Disposable `kind` cluster `agenova-k8s-lab` created (`kindest/node:v1.36.1`), control-plane Ready in ~19s.
- Kube context verified as `kind-agenova-k8s-lab` before any mutation (re-checked before each apply).
- Ownership receipt (Docker control-plane container ID + `kube-system` UID) written to `.tmp/agenova-k8s-lab-owner/identity` and re-verified.
- Pinned Agent Sandbox v0.4.6 `manifest.yaml` downloaded and applied: `agent-sandbox-system` namespace, controller ServiceAccount / ClusterRole / ClusterRoleBinding / Service / Deployment, and CRD `sandboxes.agents.x-k8s.io` all created.
- Scoped teardown proven separately (`reproduce.sh down`): re-verified context + the ownership receipt, deleted **only** `agenova-k8s-lab` (unrelated local `kind` clusters untouched), and removed the receipt.

## Blocker

The `agent-sandbox-controller` Deployment never reached Ready. The `kind` node's
containerd cannot pull `registry.k8s.io/agent-sandbox/agent-sandbox-controller:v0.4.6`:

```
Failed to pull image ".../agent-sandbox-controller:v0.4.6": ... failed to do request:
Head "https://registry.k8s.io/v2/.../manifests/v0.4.6":
tls: failed to verify certificate: x509: certificate signed by unknown authority
```

This machine's network runs TLS inspection (Zscaler). The host trusts the
inspection CA — `docker pull` of the same image on the host succeeds — but the
`kind` node's containerd ships its own CA bundle and does not. This is an
environment property, not a defect in the substrate, the manifests, or
`reproduce.sh`, which correctly reported it as a non-zero failure rather than a
silent pass.

## Required before #99 is accepted

Run `reproduce.sh all --capture` **twice** on a network without TLS interception
(or with the inspection root CA added to the `kind` node's trust store), and
replace this file and `output.txt` with a capture that also shows: each CRD's
served/storage versions, controller image + readiness, `SandboxClaim`
`Ready=True`, sandbox pods -> 0 after claim/pool/template teardown, the smoke
namespace deleted, and the owned cluster deleted. No Agenova adapter or
claim-governance path is exercised here (that is E8-T4 / #51).

Raw output: [output.txt](output.txt).
