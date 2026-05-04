package main

import (
	"fleetglance/internal/version"
	"fmt"
)

func main() {
	// prints greeting to the console with version information
	fmt.Printf("Hello, Agent!\nVersion: %s\nCommit: %s\nBuiltAt: %s\n", version.Version, version.Commit, version.BuiltAt)
}
