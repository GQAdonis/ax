// Gar is a CLI tool for managing agent orchestrator sessions.
// It provides commands to trigger sessions, resume from checkpoints,
// inspect session state, register agents, and run the controller server.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
