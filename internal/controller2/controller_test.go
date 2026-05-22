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
	"testing"

	"github.com/google/ax/internal/controller/executor"
	"github.com/google/ax/internal/controller/executor/executortest"
	"github.com/google/ax/internal/harness"
	"github.com/google/ax/internal/harness/harnesstest"
	"github.com/google/ax/proto"
)

func TestController2_ExecHelloWorld(t *testing.T) {
	ctx := context.Background()
	cid := "test-conversation-id"

	log := &executortest.MemoryEventLog{}
	reg := NewRegistry()
	c, err := New(ctx, Config{
		Registry: reg,
		EventLogBuilder: func() (executor.EventLog, error) {
			return log, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var outputs []*proto.Message
	handler := ExecHandler(func(resp *proto.ExecResponse) error {
		outputs = append(outputs, resp.Outputs...)
		return nil
	})

	inputs := []*proto.Message{
		{
			Role: "user",
			Content: &proto.Content{
				Type: &proto.Content_Text{
					Text: &proto.TextContent{Text: "Trigger prompt"},
				},
			},
		},
	}

	err = c.Exec(ctx, &proto.ExecRequest{
		ConversationId: cid,
		Inputs:         inputs,
	}, handler)
	if err != nil {
		t.Fatalf("Controller2.Exec failed: %v", err)
	}

	if len(outputs) != 1 {
		t.Fatalf("expected exactly 1 output message, got %d", len(outputs))
	}

	gotText := outputs[0].GetContent().GetText().GetText()
	if gotText != "Hello world" {
		t.Errorf("expected 'Hello world' output text response, got %q", gotText)
	}
}

func TestController2_ExecAntigravityFallback(t *testing.T) {
	ctx := context.Background()
	cid := "test-conversation-id"

	log := &executortest.MemoryEventLog{}
	reg := NewRegistry()

	// Build and register harness with bad path to trigger build-time fallback
	badHarness := BuildHarness(ctx, "antigravity", harness.HarnessConfig{
		AntigravityScriptPath: "non-existent-script.py",
	})
	reg.RegisterHarness("antigravity", badHarness)

	c, err := New(ctx, Config{
		Registry: reg,
		EventLogBuilder: func() (executor.EventLog, error) {
			return log, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var outputs []*proto.Message
	handler := ExecHandler(func(resp *proto.ExecResponse) error {
		outputs = append(outputs, resp.Outputs...)
		return nil
	})

	inputs := []*proto.Message{
		{
			Role: "user",
			Content: &proto.Content{
				Type: &proto.Content_Text{
					Text: &proto.TextContent{Text: "Trigger prompt"},
				},
			},
		},
	}

	// Request "antigravity" agent
	err = c.Exec(ctx, &proto.ExecRequest{
		ConversationId: cid,
		Inputs:         inputs,
		AgentId:        "antigravity",
	}, handler)
	if err != nil {
		t.Fatalf("Controller2.Exec failed: %v", err)
	}

	if len(outputs) != 1 {
		t.Fatalf("expected exactly 1 output message, got %d", len(outputs))
	}

	gotText := outputs[0].GetContent().GetText().GetText()
	if gotText != "Hello world" {
		t.Errorf("expected 'Hello world' output text response due to fallback, got %q", gotText)
	}
}

func TestController2_ExecRuntimeFallback(t *testing.T) {
	ctx := context.Background()
	cid := "test-conversation-id"

	log := &executortest.MemoryEventLog{}
	reg := NewRegistry() // Empty registry, will force runtime fallback for any requested agent

	c, err := New(ctx, Config{
		Registry: reg,
		EventLogBuilder: func() (executor.EventLog, error) {
			return log, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var outputs []*proto.Message
	handler := ExecHandler(func(resp *proto.ExecResponse) error {
		outputs = append(outputs, resp.Outputs...)
		return nil
	})

	inputs := []*proto.Message{
		{
			Role: "user",
			Content: &proto.Content{
				Type: &proto.Content_Text{
					Text: &proto.TextContent{Text: "Trigger prompt"},
				},
			},
		},
	}

	// Request "antigravity" agent, which is NOT registered
	err = c.Exec(ctx, &proto.ExecRequest{
		ConversationId: cid,
		Inputs:         inputs,
		AgentId:        "antigravity",
	}, handler)
	if err != nil {
		t.Fatalf("Controller2.Exec failed: %v", err)
	}

	if len(outputs) != 1 {
		t.Fatalf("expected exactly 1 output message, got %d", len(outputs))
	}

	gotText := outputs[0].GetContent().GetText().GetText()
	if gotText != "Hello world" {
		t.Errorf("expected 'Hello world' output text response due to runtime fallback, got %q", gotText)
	}
}

func TestController2_Exec_ResumptionAndIDGeneration(t *testing.T) {
	t.Skip("Feature Gap: Resumption and Event Logging are not yet implemented in controller2")

	// This test is sketched out for when the features are implemented.
	ctx := context.Background()
	cid := "test-conv"

	inputs := []*proto.Message{
		{
			Role: "user",
			Content: &proto.Content{
				Type: &proto.Content_Text{
					Text: &proto.TextContent{Text: "hello"},
				},
			},
		},
	}

	log := &executortest.MemoryEventLog{}
	reg := NewRegistry()
	mockHarness := harnesstest.New()
	reg.RegisterHarness("mock-agent", mockHarness)

	c, err := New(ctx, Config{
		Registry: reg,
		EventLogBuilder: func() (executor.EventLog, error) {
			return log, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	err = c.Exec(ctx, &proto.ExecRequest{
		ConversationId: cid,
		Inputs:         inputs,
		AgentId:        "mock-agent",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that events were logged and exec ID was generated.
	// ...
}

func TestController2_Exec_LastSeq_Empty(t *testing.T) {
	t.Skip("Feature Gap: History playback and LastSeq are not yet implemented in controller2")
}

func TestController2_Exec_LastSeq(t *testing.T) {
	t.Skip("Feature Gap: History playback and LastSeq are not yet implemented in controller2")
}

func TestController2_Exec_LastSeq_NotFound(t *testing.T) {
	t.Skip("Feature Gap: History playback and LastSeq are not yet implemented in controller2")
}

func TestController2_Exec_WaitsForConfirmation(t *testing.T) {
	t.Skip("Feature Gap: Resumption and Confirmation handling are not yet implemented in controller2")

	// Sketch for when implemented:
	ctx := context.Background()
	cid := "test-conv-conf"

	log := &executortest.MemoryEventLog{}
	reg := NewRegistry()
	mockHarness := harnesstest.New()

	// Configure mock harness to return a confirmation question on first Run
	mockHarness.DefaultRunFunc = func(ctx context.Context, execID string, handler harness.Handler) error {
		questionMsg := &proto.Message{
			Role: "assistant",
			Content: &proto.Content{
				Type: &proto.Content_Confirmation{
					Confirmation: &proto.ConfirmationContent{
						Question: "Are you sure?",
					},
				},
			},
		}
		if err := handler.OnMessage(ctx, execID, questionMsg); err != nil {
			return err
		}
		return handler.OnComplete(ctx, execID)
	}
	reg.RegisterHarness("mock-agent", mockHarness)

	c, err := New(ctx, Config{
		Registry: reg,
		EventLogBuilder: func() (executor.EventLog, error) {
			return log, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	var msgs []*proto.Message
	handler := ExecHandler(func(resp *proto.ExecResponse) error {
		msgs = append(msgs, resp.Outputs...)
		return nil
	})

	err = c.Exec(ctx, &proto.ExecRequest{
		ConversationId: cid,
		AgentId:        "mock-agent",
	}, handler)
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	if msgs[0].GetContent().GetConfirmation().GetQuestion() != "Are you sure?" {
		t.Fatalf("expected 'Are you sure?', got %v", msgs[0].GetContent().GetConfirmation().GetQuestion())
	}
}

func TestController2_Exec_InternalOnly(t *testing.T) {
	t.Skip("Feature Gap: Event Logging and InternalOnly filtering are not yet implemented in controller2")
}
