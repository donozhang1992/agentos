// Copyright 2026 Dapeng Zhang and Agenova contributors.
// SPDX-License-Identifier: Apache-2.0

// Package contracttest provides the reusable behavioral suite for the reduced
// RuntimeBackend contract. Import it from a backend's test file and call Run
// to verify parity with the in-memory reference implementation.
//
// The suite asserts operation semantics only: identity association,
// readiness versus explicit start, termination and cleanup evidence, and the
// negative cases from the Ticket #30 specification. Pool counters and
// application phases are implementation-specific and are not asserted here.
package contracttest

import (
	"errors"
	"testing"

	"github.com/wunderforge/agenova/api/v1alpha1"
	"github.com/wunderforge/agenova/internal/runtime"
)

// Fixture is one fresh, pre-configured backend plus the control hooks the
// suite needs to reach states a backend does not expose through the contract.
//
// HoldReadiness and ReleaseReadiness are REQUIRED for the reference suite:
// without them the "known but not ready" cases cannot be proven and Run fails
// instead of skipping. Started is an optional probe; when nil only the probe
// assertions are omitted, the behavioral start assertions always run.
type Fixture struct {
	Backend     runtime.RuntimeBackend
	TemplateRef string

	// HoldReadiness makes the next Allocate for claimID report Ready=false.
	HoldReadiness func(claimID string)
	// ReleaseReadiness flips a held allocation to ready.
	ReleaseReadiness func(id v1alpha1.SandboxClaimBackendIdentity)
	// Started reports whether the backend has acknowledged work start.
	Started func(id v1alpha1.SandboxClaimBackendIdentity) bool
}

// Run exercises the reduced RuntimeBackend contract. newFixture must return a
// fresh fixture for each subtest.
func Run(t *testing.T, newFixture func(t *testing.T) Fixture) {
	t.Helper()

	cases := []struct {
		name string
		fn   func(t *testing.T, f Fixture)
	}{
		{"allocate returns resolvable identity", testAllocateReturnsResolvableIdentity},
		{"golden path allocate observe start terminate cleanup", testGoldenPath},
		{"unknown identity is rejected by every operation", testUnknownIdentityRejected},
		{"duplicate claim id is rejected without disturbing first allocation", testDuplicateClaimRejected},
		{"start is rejected until readiness is released", testStartRejectedUntilReady},
		{"observe never implies start and start is explicit once", testStartIsExplicitOnce},
		{"terminate before start cancels idempotently and blocks start", testTerminateBeforeStart},
		{"cleanup is idempotent and blocks start and terminate", testCleanupIdempotent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			requireFixture(t, f)
			tc.fn(t, f)
		})
	}
}

func requireFixture(t *testing.T, f Fixture) {
	t.Helper()
	if f.Backend == nil {
		t.Fatal("fixture: Backend is required")
	}
	if f.TemplateRef == "" {
		t.Fatal("fixture: TemplateRef is required")
	}
	if f.HoldReadiness == nil || f.ReleaseReadiness == nil {
		t.Fatal("fixture: HoldReadiness and ReleaseReadiness are required to prove not-ready behavior")
	}
}

func testAllocateReturnsResolvableIdentity(t *testing.T, f Fixture) {
	t.Helper()
	alloc := allocate(t, f, "research-run")
	if alloc.ClaimID != "research-run" {
		t.Fatalf("allocation claim id = %q, want research-run", alloc.ClaimID)
	}
	if alloc.Identity.Backend == "" || alloc.Identity.WorkerID == "" {
		t.Fatalf("allocation identity incomplete: %+v", alloc.Identity)
	}
	obs := observe(t, f, alloc.Identity)
	if obs.ClaimID != "research-run" {
		t.Fatalf("observe claim id = %q, want research-run", obs.ClaimID)
	}
	if obs.Identity != alloc.Identity {
		t.Fatalf("observe identity = %+v, want %+v", obs.Identity, alloc.Identity)
	}
	if obs.Released {
		t.Fatal("fresh allocation must not report released")
	}
}

func testGoldenPath(t *testing.T, f Fixture) {
	t.Helper()
	alloc := allocate(t, f, "golden")
	id := alloc.Identity

	if obs := observe(t, f, id); !obs.Ready {
		t.Fatalf("default fixture allocation should be ready, got %+v", obs)
	}
	if f.Started != nil && f.Started(id) {
		t.Fatal("readiness must not imply started")
	}
	if err := f.Backend.Start(id); err != nil {
		t.Fatalf("start: %v", err)
	}
	if f.Started != nil && !f.Started(id) {
		t.Fatal("start was acknowledged but probe reports not started")
	}
	if err := f.Backend.Terminate(id); err != nil {
		t.Fatalf("terminate: %v", err)
	}
	res, err := f.Backend.Cleanup(id)
	if err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if !res.Released {
		t.Fatalf("cleanup should report released, got %+v", res)
	}
	if res.Identity != id {
		t.Fatalf("cleanup identity = %+v, want %+v", res.Identity, id)
	}
	obs := observe(t, f, id)
	if !obs.Released {
		t.Fatalf("identity must stay resolvable with Released=true after cleanup, got %+v", obs)
	}
	if obs.Ready {
		t.Fatal("released allocation must not report ready")
	}
	if obs.ClaimID != "golden" {
		t.Fatalf("claim correlation lost after cleanup: %+v", obs)
	}
}

func testUnknownIdentityRejected(t *testing.T, f Fixture) {
	t.Helper()
	alloc := allocate(t, f, "known")
	unknown := []struct {
		name string
		id   v1alpha1.SandboxClaimBackendIdentity
	}{
		{"unknown worker", v1alpha1.SandboxClaimBackendIdentity{Backend: alloc.Identity.Backend, WorkerID: "no-such-worker"}},
		{"known worker under other backend", v1alpha1.SandboxClaimBackendIdentity{Backend: "other-backend", WorkerID: alloc.Identity.WorkerID}},
		{"empty identity", v1alpha1.SandboxClaimBackendIdentity{}},
	}
	for _, u := range unknown {
		if _, err := f.Backend.Observe(u.id); !errors.Is(err, runtime.ErrUnknownIdentity) {
			t.Fatalf("%s: observe error = %v, want ErrUnknownIdentity", u.name, err)
		}
		if err := f.Backend.Start(u.id); !errors.Is(err, runtime.ErrUnknownIdentity) {
			t.Fatalf("%s: start error = %v, want ErrUnknownIdentity", u.name, err)
		}
		if err := f.Backend.Terminate(u.id); !errors.Is(err, runtime.ErrUnknownIdentity) {
			t.Fatalf("%s: terminate error = %v, want ErrUnknownIdentity", u.name, err)
		}
		if _, err := f.Backend.Cleanup(u.id); !errors.Is(err, runtime.ErrUnknownIdentity) {
			t.Fatalf("%s: cleanup error = %v, want ErrUnknownIdentity", u.name, err)
		}
	}
	// The known allocation is untouched by the rejected calls.
	if obs := observe(t, f, alloc.Identity); obs.Released || !obs.Ready {
		t.Fatalf("known allocation disturbed by unknown-identity calls: %+v", obs)
	}
}

func testDuplicateClaimRejected(t *testing.T, f Fixture) {
	t.Helper()
	first := allocate(t, f, "dup")
	before := observe(t, f, first.Identity)

	if _, err := f.Backend.Allocate(runtime.AllocateRequest{ClaimID: "dup", TemplateRef: f.TemplateRef}); err == nil {
		t.Fatal("duplicate claim id should be rejected")
	}
	after := observe(t, f, first.Identity)
	if after != before {
		t.Fatalf("first allocation changed by duplicate attempt: before %+v after %+v", before, after)
	}
}

func testStartRejectedUntilReady(t *testing.T, f Fixture) {
	t.Helper()
	f.HoldReadiness("held")
	alloc := allocate(t, f, "held")
	id := alloc.Identity

	if obs := observe(t, f, id); obs.Ready {
		t.Fatalf("held allocation should not be ready, got %+v", obs)
	}
	if err := f.Backend.Start(id); !errors.Is(err, runtime.ErrNotReady) {
		t.Fatalf("start before ready error = %v, want ErrNotReady", err)
	}
	for i := 0; i < 3; i++ {
		if obs := observe(t, f, id); obs.Ready {
			t.Fatal("repeated observe must not flip readiness")
		}
	}
	if f.Started != nil && f.Started(id) {
		t.Fatal("rejected start must not mark the allocation started")
	}
	f.ReleaseReadiness(id)
	if obs := observe(t, f, id); !obs.Ready {
		t.Fatalf("released readiness should be observable, got %+v", obs)
	}
	if err := f.Backend.Start(id); err != nil {
		t.Fatalf("start after readiness released: %v", err)
	}
}

func testStartIsExplicitOnce(t *testing.T, f Fixture) {
	t.Helper()
	alloc := allocate(t, f, "explicit")
	id := alloc.Identity

	for i := 0; i < 3; i++ {
		if obs := observe(t, f, id); !obs.Ready {
			t.Fatalf("allocation should stay ready, got %+v", obs)
		}
	}
	if f.Started != nil && f.Started(id) {
		t.Fatal("observe must not trigger start")
	}
	if err := f.Backend.Start(id); err != nil {
		t.Fatalf("first start: %v", err)
	}
	if err := f.Backend.Start(id); !errors.Is(err, runtime.ErrAlreadyStarted) {
		t.Fatalf("second start error = %v, want ErrAlreadyStarted", err)
	}
	if f.Started != nil && !f.Started(id) {
		t.Fatal("probe should report started after explicit start")
	}
}

func testTerminateBeforeStart(t *testing.T, f Fixture) {
	t.Helper()
	alloc := allocate(t, f, "cancel")
	id := alloc.Identity

	if err := f.Backend.Terminate(id); err != nil {
		t.Fatalf("terminate before start (cancel): %v", err)
	}
	if err := f.Backend.Terminate(id); err != nil {
		t.Fatalf("second terminate should be idempotent: %v", err)
	}
	if err := f.Backend.Start(id); !errors.Is(err, runtime.ErrTerminated) {
		t.Fatalf("start after terminate error = %v, want ErrTerminated", err)
	}
	if f.Started != nil && f.Started(id) {
		t.Fatal("cancelled allocation must never report started")
	}
	res, err := f.Backend.Cleanup(id)
	if err != nil {
		t.Fatalf("cleanup after cancel: %v", err)
	}
	if !res.Released {
		t.Fatalf("cleanup after cancel should release, got %+v", res)
	}
}

func testCleanupIdempotent(t *testing.T, f Fixture) {
	t.Helper()
	alloc := allocate(t, f, "release")
	id := alloc.Identity

	if err := f.Backend.Start(id); err != nil {
		t.Fatalf("start: %v", err)
	}
	first, err := f.Backend.Cleanup(id)
	if err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if !first.Released {
		t.Fatalf("cleanup should release, got %+v", first)
	}
	second, err := f.Backend.Cleanup(id)
	if err != nil {
		t.Fatalf("second cleanup should be idempotent: %v", err)
	}
	if second != first {
		t.Fatalf("second cleanup result %+v, want %+v", second, first)
	}
	if err := f.Backend.Start(id); !errors.Is(err, runtime.ErrReleased) {
		t.Fatalf("start after cleanup error = %v, want ErrReleased", err)
	}
	if err := f.Backend.Terminate(id); !errors.Is(err, runtime.ErrReleased) {
		t.Fatalf("terminate after cleanup error = %v, want ErrReleased", err)
	}
	if obs := observe(t, f, id); !obs.Released || obs.Ready {
		t.Fatalf("released allocation observation = %+v", obs)
	}
}

// --- helpers ---

func allocate(t *testing.T, f Fixture, claimID string) runtime.Allocation {
	t.Helper()
	alloc, err := f.Backend.Allocate(runtime.AllocateRequest{ClaimID: claimID, TemplateRef: f.TemplateRef})
	if err != nil {
		t.Fatalf("allocate %s: %v", claimID, err)
	}
	return alloc
}

func observe(t *testing.T, f Fixture, id v1alpha1.SandboxClaimBackendIdentity) runtime.Observation {
	t.Helper()
	obs, err := f.Backend.Observe(id)
	if err != nil {
		t.Fatalf("observe %+v: %v", id, err)
	}
	return obs
}
