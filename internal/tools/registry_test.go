package tools

import (
	"errors"
	"testing"

	"github.com/alexnakagama/noryn/internal/llm"
)

type stubTool struct {
	name         string
	lastArgs     string
	result       string
	executeError error
}

func (t *stubTool) Name() string {
	return t.name
}

func (t *stubTool) Description() string {
	return "stub tool " + t.name
}

func (t *stubTool) Execute(arguments string) (string, error) {
	t.lastArgs = arguments
	return t.result, t.executeError
}

func (t *stubTool) Definition() llm.ToolDefinition {
	return llm.ToolDefinition{Name: t.name, Description: "stub tool " + t.name}
}

func TestRegistry_Get_UnknownToolReturnsNotFound(t *testing.T) {
	registry := NewRegistry()

	_, ok := registry.Get("missing")

	if ok {
		t.Fatalf("Get() ok = true, want false")
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	tool := &stubTool{name: "one"}

	if err := registry.Register(tool); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	got, ok := registry.Get("one")
	if !ok {
		t.Fatalf("Get() ok = false, want true")
	}

	if got != tool {
		t.Errorf("Get() = %v, want %v", got, tool)
	}
}

func TestRegistry_RegisterDuplicateReturnsError(t *testing.T) {
	registry := NewRegistry(&stubTool{name: "one"})

	err := registry.Register(&stubTool{name: "one"})

	if err == nil {
		t.Fatal("Register() error = nil, want an error")
	}

	want := "tool already registered: one"
	if err.Error() != want {
		t.Errorf("Register() error = %q, want %q", err.Error(), want)
	}
}

func TestRegistry_NewRegistryRegistersAllTools(t *testing.T) {
	first := &stubTool{name: "first"}
	second := &stubTool{name: "second"}
	registry := NewRegistry(first, second)

	if _, ok := registry.Get("first"); !ok {
		t.Errorf("Get(first) ok = false, want true")
	}

	if _, ok := registry.Get("second"); !ok {
		t.Errorf("Get(second) ok = false, want true")
	}

	if _, ok := registry.Get("missing"); ok {
		t.Errorf("Get(missing) ok = true, want false")
	}
}

func TestRegistry_Execute(t *testing.T) {
	tests := []struct {
		name     string
		call     llm.ToolCall
		tool     *stubTool
		want     string
		wantArgs string
		wantErr  bool
	}{
		{
			name:     "executes a registered tool",
			call:     llm.ToolCall{Name: "one", Arguments: "hello"},
			tool:     &stubTool{name: "one", result: "world"},
			want:     "world",
			wantArgs: "hello",
		},
		{
			name: "passes empty arguments",
			call: llm.ToolCall{Name: "one"},
			tool: &stubTool{name: "one", result: "world"},
			want: "world",
		},
		{
			name:    "unknown tool returns an error",
			call:    llm.ToolCall{Name: "missing"},
			tool:    &stubTool{name: "one", result: "world"},
			wantErr: true,
		},
		{
			name:    "propagates the tool error",
			call:    llm.ToolCall{Name: "one", Arguments: "x"},
			tool:    &stubTool{name: "one", executeError: errors.New("boom")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry(tt.tool)

			got, err := registry.Execute(tt.call)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}

			if tt.tool.lastArgs != tt.wantArgs {
				t.Errorf("tool received arguments %q, want %q", tt.tool.lastArgs, tt.wantArgs)
			}
		})
	}
}

func TestRegistry_ExecuteErrorMessages(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.Execute(llm.ToolCall{Name: "missing"})

	if err == nil {
		t.Fatal("Execute() error = nil, want an error")
	}

	want := "tool not found: missing"
	if err.Error() != want {
		t.Errorf("Execute() error = %q, want %q", err.Error(), want)
	}
}

func TestRegistry_Definitions(t *testing.T) {
	registry := NewRegistry(
		&stubTool{name: "one"},
		&stubTool{name: "two"},
	)

	got := registry.Definitions()

	if len(got) != 2 {
		t.Fatalf("Definitions() returned %d definitions, want 2", len(got))
	}

	names := make(map[string]bool)
	for _, definition := range got {
		names[definition.Name] = true
	}

	if !names["one"] {
		t.Errorf("Definitions() missing tool %q", "one")
	}

	if !names["two"] {
		t.Errorf("Definitions() missing tool %q", "two")
	}
}

func TestRegistry_DefinitionsEmpty(t *testing.T) {
	registry := NewRegistry()

	got := registry.Definitions()

	if len(got) != 0 {
		t.Errorf("Definitions() returned %d definitions, want 0", len(got))
	}
}
