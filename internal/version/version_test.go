package version

import "testing"

func TestFormat(t *testing.T) {
	oldVersion := Version
	oldCommit := Commit
	oldBuiltAt := BuiltAt
	t.Cleanup(func() {
		Version = oldVersion
		Commit = oldCommit
		BuiltAt = oldBuiltAt
	})

	Version = "v1.2.3"
	Commit = "abc1234"
	BuiltAt = "2026-05-11T10:00:00Z"

	got := Format("Agent")
	want := "Fleetglance Agent\nversion=v1.2.3\ncommit=abc1234\nbuilt_at=2026-05-11T10:00:00Z\n"
	if got != want {
		t.Fatalf("expected version format %q, got %q", want, got)
	}
}
