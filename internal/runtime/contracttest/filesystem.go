// Copyright 2026 Agenova contributors.
// SPDX-License-Identifier: Apache-2.0

package contracttest

import (
	"bytes"
	"errors"
	"path"
	"testing"

	"github.com/wunderforge/agenova/api/v1alpha1"
	"github.com/wunderforge/agenova/internal/runtime"
)

// FilesystemFixture supplies reference-model probes for the filesystem
// contract. These probes are test controls, not RuntimeBackend operations and
// not a filesystem gateway.
type FilesystemFixture struct {
	Backend         runtime.RuntimeBackend
	TemplateRef     string
	WriteTaskFile   func(v1alpha1.SandboxClaimBackendIdentity, string, []byte) error
	ReadTaskFile    func(v1alpha1.SandboxClaimBackendIdentity, string) ([]byte, error)
	OutsideSentinel func() []byte
}

// RunFilesystem exercises the backend-neutral filesystem description and
// reference semantics. It proves a simulated model only; hostile native
// process isolation requires separate real-backend evidence.
func RunFilesystem(t *testing.T, newFixture func(t *testing.T) FilesystemFixture) {
	t.Helper()
	cases := []struct {
		name string
		fn   func(*testing.T, FilesystemFixture)
	}{
		{"FS-P1 reports one simulated ephemeral task directory", testFilesystemDescription},
		{"FS-P2 reads and writes task data after explicit start", testFilesystemTaskData},
		{"FS-N1 rejects outside and traversal writes", testFilesystemOutsideBoundary},
		{"FS-N4 replacement claim receives fresh data", testFilesystemFreshReplacement},
		{"FS-N7 and FS-N8 deny access after termination and cleanup", testFilesystemLifecycle},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			requireFilesystemFixture(t, f)
			tc.fn(t, f)
		})
	}
}

func requireFilesystemFixture(t *testing.T, f FilesystemFixture) {
	t.Helper()
	if f.Backend == nil || f.TemplateRef == "" || f.WriteTaskFile == nil || f.ReadTaskFile == nil || f.OutsideSentinel == nil {
		t.Fatal("filesystem fixture requires backend, template and model read/write probes")
	}
}

func testFilesystemDescription(t *testing.T, f FilesystemFixture) {
	alloc := allocateFilesystem(t, f, "fs-description")
	want := runtime.FilesystemBoundary{
		WorkingDirectory: "/workspace",
		OutsideBoundary:  runtime.FilesystemOutsideRuntimeReadOnlyOtherUnavailable,
		Ephemeral:        true,
		EvidenceLevel:    runtime.FilesystemEvidenceSimulated,
	}
	if alloc.Filesystem != want {
		t.Fatalf("filesystem = %+v, want %+v", alloc.Filesystem, want)
	}
	obs, err := f.Backend.Observe(alloc.Identity)
	if err != nil || obs.Filesystem != alloc.Filesystem {
		t.Fatalf("observed filesystem = %+v, error = %v", obs.Filesystem, err)
	}
}

func testFilesystemTaskData(t *testing.T, f FilesystemFixture) {
	alloc := startFilesystem(t, f, "fs-task-data")
	want := []byte("claim-scoped output")
	if err := f.WriteTaskFile(alloc.Identity, "repo/result.txt", want); err != nil {
		t.Fatal(err)
	}
	got, err := f.ReadTaskFile(alloc.Identity, path.Join(alloc.Filesystem.WorkingDirectory, "repo/result.txt"))
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("read = %q, error = %v", got, err)
	}
	got[0] = 'X'
	again, _ := f.ReadTaskFile(alloc.Identity, "repo/result.txt")
	if !bytes.Equal(again, want) {
		t.Fatalf("caller mutated stored task data: %q", again)
	}
}

func testFilesystemOutsideBoundary(t *testing.T, f FilesystemFixture) {
	alloc := startFilesystem(t, f, "fs-outside")
	before := f.OutsideSentinel()
	for _, target := range []string{"/outside/sentinel", "../sentinel", "/workspace-sibling/sentinel", `C:\\host\\secret`} {
		if err := f.WriteTaskFile(alloc.Identity, target, []byte("mutated")); !errors.Is(err, runtime.ErrFilesystemBoundary) {
			t.Fatalf("target %q error = %v, want ErrFilesystemBoundary", target, err)
		}
	}
	if _, err := f.ReadTaskFile(alloc.Identity, "/outside/sentinel"); !errors.Is(err, runtime.ErrFilesystemBoundary) {
		t.Fatalf("outside read error = %v, want ErrFilesystemBoundary", err)
	}
	if after := f.OutsideSentinel(); !bytes.Equal(after, before) {
		t.Fatalf("outside sentinel changed: before %q after %q", before, after)
	}
}

func testFilesystemFreshReplacement(t *testing.T, f FilesystemFixture) {
	first := startFilesystem(t, f, "fs-first")
	if err := f.WriteTaskFile(first.Identity, "repo/private.txt", []byte("first claim")); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Backend.Cleanup(first.Identity); err != nil {
		t.Fatal(err)
	}
	second := startFilesystem(t, f, "fs-second")
	if _, err := f.ReadTaskFile(second.Identity, "repo/private.txt"); err == nil {
		t.Fatal("replacement claim inherited the previous claim's task file")
	}
}

func testFilesystemLifecycle(t *testing.T, f FilesystemFixture) {
	alloc := startFilesystem(t, f, "fs-lifecycle")
	if err := f.WriteTaskFile(alloc.Identity, "repo/output.txt", []byte("unexported")); err != nil {
		t.Fatal(err)
	}
	if err := f.Backend.Terminate(alloc.Identity); err != nil {
		t.Fatal(err)
	}
	if _, err := f.ReadTaskFile(alloc.Identity, "repo/output.txt"); !errors.Is(err, runtime.ErrTerminated) {
		t.Fatalf("read after termination = %v, want ErrTerminated", err)
	}
	if _, err := f.Backend.Cleanup(alloc.Identity); err != nil {
		t.Fatal(err)
	}
	if err := f.WriteTaskFile(alloc.Identity, "repo/late.txt", []byte("late")); !errors.Is(err, runtime.ErrReleased) {
		t.Fatalf("write after cleanup = %v, want ErrReleased", err)
	}
}

func allocateFilesystem(t *testing.T, f FilesystemFixture, claimID string) runtime.Allocation {
	t.Helper()
	alloc, err := f.Backend.Allocate(runtime.AllocateRequest{ClaimID: claimID, TemplateRef: f.TemplateRef})
	if err != nil {
		t.Fatalf("allocate %s: %v", claimID, err)
	}
	return alloc
}

func startFilesystem(t *testing.T, f FilesystemFixture, claimID string) runtime.Allocation {
	t.Helper()
	alloc := allocateFilesystem(t, f, claimID)
	if err := f.Backend.Start(alloc.Identity); err != nil {
		t.Fatalf("start %s: %v", claimID, err)
	}
	return alloc
}
