// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller2

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/google/ax/internal/harness"
	"github.com/google/ax/internal/harness/harnesstest"
)

// BuildHarness builds a harness based on the requested type, with fallback to test harness.
func BuildHarness(ctx context.Context, harnessType string, cfg harness.HarnessConfig) harness.Harness {
	switch harnessType {
	case "antigravity":
		// Check if python3 is available
		if _, err := exec.LookPath("python3"); err != nil {
			log.Printf("WARNING: python3 not found in PATH, falling back to test harness: %v", err)
			return harnesstest.New()
		}
		// Check if script exists
		scriptPath := cfg.AntigravityScriptPath
		if scriptPath == "" {
			scriptPath = "examples/antigravity_agent/agent.py"
		}
		if _, err := os.Stat(scriptPath); err != nil {
			log.Printf("WARNING: Antigravity agent script not found at %s, falling back to test harness: %v", scriptPath, err)
			return harnesstest.New()
		}
		log.Printf("Using Antigravity harness with script: %s", scriptPath)

		builder := &harness.AntigravityHarnessBuilder{
			Config: harness.HarnessConfig{
				AntigravityScriptPath: scriptPath,
			},
		}
		return builder.Build()
	default:
		log.Printf("Using default test harness")
		return harnesstest.New()
	}
}
