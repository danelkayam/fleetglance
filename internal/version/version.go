package version

import "fmt"

var (
	Version = "dev"
	Commit  = "unknown"
	BuiltAt = "unknown"
)

func Format(component string) string {
	return fmt.Sprintf("Fleetglance %s\nversion=%s\ncommit=%s\nbuilt_at=%s\n", component, Version, Commit, BuiltAt)
}
