package tools

import (
	"testing"
)

func TestShellTool_Name(t *testing.T) {
	tool := ShellTool{}
	got := tool.Name()

	if got != "shell" {
		t.Errorf("Name() = %q, want %q", got, "shell")
	}
}

func TestShellTool_Execute(t *testing.T) {
	tests := []struct {
		name       string
		args       string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "returns the command output",
			args:       `{"command":"echo hello"}`,
			wantOutput: "hello\n",
		},
		{
			name:       "runs chained commands",
			args:       `{"command":"printf 'a\\nb' && printf c"}`,
			wantOutput: "a\nbc",
		},
		{
			name: "empty command succeeds",
			args: `{"command":""}`,
		},
		{
			name:       "failing command returns its combined output and error",
			args:       `{"command":"echo boom >&2; exit 1"}`,
			wantOutput: "boom\n",
			wantErr:    true,
		},
		{
			name:    "non-zero exit returns an error",
			args:    `{"command":"exit 3"}`,
			wantErr: true,
		},
		{
			name:    "malformed arguments return an error",
			args:    "not-json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := ShellTool{}
			got, err := tool.Execute(tt.args)

			if got != tt.wantOutput {
				t.Errorf("Execute() output = %q, want %q", got, tt.wantOutput)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want an error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
		})
	}
}
