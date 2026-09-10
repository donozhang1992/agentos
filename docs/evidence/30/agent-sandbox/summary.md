# Ticket #30 real-backend evidence blocker

- Task: [Reduce RuntimeBackend to the MVP contract](../../../../work/0030-runtime-backend-mvp/task.md)
- Gate: Agent Sandbox integration
- Date: 2026-09-10
- Branch: `tomtian/e3-t1-runtime-backend-mvp`
- Implementation baseline: `a72baf9d1ee228fb88717e56d5f69c96847264f5`; Slice 2 changes migrate the integration gate to the reduced contract.
- Command to run on an available test environment: `pwsh -File scripts/check.ps1 -Integration -KubeContext <explicitly-confirmed-test-context>`
- Result: **blocked, not run**. The executable probe returned no kubectl binary, and the Owner has not confirmed a test context.
- Probe output: [output.txt](output.txt).

The package can be compiled without contacting a cluster. Compilation and fake-controller tests do not verify Agent Sandbox allocation, readiness, or cleanup on a real cluster.

A future run must use an explicitly confirmed test context with compatible CRDs/controller and test-resource permissions. It must capture the identity-matched readiness and claim/sandbox absence asserted by the updated harness. Unknown allocation identity leaves resources for manual recovery; it does not permit blind cleanup. Replace this blocker with actual integration output only after that gate runs.
