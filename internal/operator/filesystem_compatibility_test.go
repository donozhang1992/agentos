// Copyright 2026 Agenova contributors.
// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const fixtureExportLimit = 1 << 20

// TestFilesystemLocalCompatibility uses real git and Go commands inside a
// controlled temporary directory. It proves ordinary-tool compatibility, not
// isolation from a hostile process; the reusable model suite records boundary
// denials separately.
func TestFilesystemLocalCompatibility(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	collector := filepath.Join(root, "collector")
	outside := filepath.Join(root, "outside-sentinel.txt")
	for _, directory := range []string{workspace, collector, filepath.Join(workspace, ".home"), filepath.Join(workspace, ".cache"), filepath.Join(workspace, ".tmp")} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFixtureFile(t, outside, []byte("outside-unchanged\n"))
	writeFixtureFile(t, filepath.Join(workspace, "go.mod"), []byte("module example.local/fsfixture\n\ngo 1.24\n"))
	writeFixtureFile(t, filepath.Join(workspace, "value.go"), []byte("package fsfixture\n\nfunc Value() int { return 1 }\n"))
	writeFixtureFile(t, filepath.Join(workspace, "value_test.go"), []byte("package fsfixture\n\nimport \"testing\"\n\nfunc TestValue(t *testing.T) { if Value() != 2 { t.Fatalf(\"Value = %d\", Value()) } }\n"))

	environment := fixtureEnvironment(workspace)
	runFixtureCommand(t, workspace, environment, "git", "init", "--quiet")
	runFixtureCommand(t, workspace, environment, "git", "-c", "user.name=Agenova Fixture", "-c", "user.email=fixture@invalid.example", "add", ".")
	runFixtureCommand(t, workspace, environment, "git", "-c", "user.name=Agenova Fixture", "-c", "user.email=fixture@invalid.example", "commit", "--quiet", "-m", "fixture baseline")

	writeFixtureFile(t, filepath.Join(workspace, "value.go"), []byte("package fsfixture\n\nfunc Value() int { return 2 }\n"))
	diff := runFixtureCommand(t, workspace, environment, "git", "diff", "--", "value.go")
	if !bytes.Contains(diff, []byte("return 2")) {
		t.Fatalf("git diff did not observe workspace edit:\n%s", diff)
	}
	runFixtureCommand(t, workspace, environment, "go", "test", "./...")

	writeFixtureFile(t, filepath.Join(workspace, "result.patch"), diff)
	exported, err := collectFixtureOutput(workspace, collector, "result.patch", fixtureExportLimit)
	if err != nil {
		t.Fatal(err)
	}
	wantDigest := sha256.Sum256(diff)
	gotDigest := sha256.Sum256(exported)
	if gotDigest != wantDigest {
		t.Fatalf("export digest = %x, want %x", gotDigest, wantDigest)
	}
	if got := mustReadFile(t, outside); string(got) != "outside-unchanged\n" {
		t.Fatalf("outside sentinel changed: %q", got)
	}

	if err := os.RemoveAll(workspace); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(workspace); !os.IsNotExist(err) {
		t.Fatalf("workspace retained after cleanup: %v", err)
	}
	if got := mustReadFile(t, filepath.Join(collector, "result.patch")); !bytes.Equal(got, diff) {
		t.Fatal("pre-termination export did not survive workspace cleanup")
	}
	t.Logf("FS-P1/FS-P2 cwd=%s git_diff_bytes=%d go_test=pass export_sha256=%x workspace_retained=false evidence=local-compatibility-not-isolation", workspace, len(diff), gotDigest)
}

func collectFixtureOutput(workspace, collector, relative string, limit int64) ([]byte, error) {
	cleaned := filepath.Clean(relative)
	if filepath.IsAbs(cleaned) || cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("fixture output must be a relative path inside the workspace")
	}
	source := filepath.Join(workspace, cleaned)
	info, err := os.Lstat(source)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fmt.Errorf("fixture output must be a bounded regular non-link file")
	}
	if err := rejectFixtureAliases(workspace, cleaned); err != nil {
		return nil, err
	}
	resolvedRoot, err := filepath.Abs(workspace)
	if err != nil {
		return nil, err
	}
	resolvedSource, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(resolvedRoot, resolvedSource)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("fixture output escaped the workspace")
	}
	data, err := os.ReadFile(resolvedSource)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(collector, filepath.Base(cleaned)), data, 0o600); err != nil {
		return nil, err
	}
	return data, nil
}

func rejectFixtureAliases(workspace, relative string) error {
	current := workspace
	for _, component := range strings.Split(filepath.Clean(relative), string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("fixture output path contains a link")
		}
	}
	return nil
}

func fixtureEnvironment(workspace string) []string {
	home := filepath.Join(workspace, ".home")
	cache := filepath.Join(workspace, ".cache")
	temp := filepath.Join(workspace, ".tmp")
	environment := append([]string(nil), os.Environ()...)
	environment = append(environment,
		"HOME="+home,
		"USERPROFILE="+home,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL="+os.DevNull,
		"GOCACHE="+filepath.Join(cache, "go-build"),
		"GOMODCACHE="+filepath.Join(cache, "go-mod"),
		"TMPDIR="+temp,
		"TEMP="+temp,
		"TMP="+temp,
	)
	return environment
}

func runFixtureCommand(t *testing.T, directory string, environment []string, name string, arguments ...string) []byte {
	t.Helper()
	command := exec.Command(name, arguments...)
	command.Dir = directory
	command.Env = environment
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(arguments, " "), err, output)
	}
	return output
}

func writeFixtureFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestCollectFixtureOutputRejectsEscapeLinkAndOversize(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	collector := filepath.Join(root, "collector")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(collector, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, filepath.Join(root, "outside.txt"), []byte("outside"))
	writeFixtureFile(t, filepath.Join(workspace, "large.bin"), bytes.Repeat([]byte{'x'}, 9))

	assertCollectorRejects(t, "escape", func() error { _, err := collectFixtureOutput(workspace, collector, "../outside.txt", 8); return err })
	assertCollectorRejects(t, "oversize", func() error { _, err := collectFixtureOutput(workspace, collector, "large.bin", 8); return err })
	if runtime.GOOS != "windows" {
		if err := os.Symlink(filepath.Join(root, "outside.txt"), filepath.Join(workspace, "link.txt")); err != nil {
			t.Fatal(err)
		}
		assertCollectorRejects(t, "symlink", func() error { _, err := collectFixtureOutput(workspace, collector, "link.txt", 8); return err })
	}
	if entries, err := os.ReadDir(collector); err != nil || len(entries) != 0 {
		t.Fatalf("rejected export wrote collector data: entries=%v error=%v", entries, err)
	}
}

func assertCollectorRejects(t *testing.T, name string, fn func() error) {
	t.Helper()
	if err := fn(); err == nil {
		t.Fatalf("%s: expected rejection", name)
	}
}
