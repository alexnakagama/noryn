# Noryn

Noryn is an AI coding agent that runs in your terminal. Built in Go, it pairs an LLM with a set of curated, project-scoped tools so you can explore, understand, and modify a codebase from the command line.

Noryn is an agent orchestration layer — it does not train or implement models. It delegates reasoning to your chosen LLM provider and focuses on the parts that matter around it: tool calling, context assembly, and safe execution.

> **Status:** early-stage. The core agent loop and tooling work, but expect rough edges and ongoing change.

## Features

- Interactive REPL driven by an LLM
- Project discovery via `go.mod`
- Agent loop with automatic tool-call handling
- Tools for filesystem access, shell execution, search, and Git inspection
- All file operations confined to the project root
- Project-specific instructions loaded from `AGENTS.md`
- Swappable LLM providers (OpenAI, OpenRouter, Fake)
- Provider and model selection via `.env` or CLI flags

## Requirements

- Go 1.27 or newer
- An API key for your chosen LLM provider

## Installation

Run directly from the source tree:

```bash
go run ./cmd/noryn
```

Or build a binary:

```bash
go build -o noryn ./cmd/noryn
./noryn
```

## Configuration

Noryn reads configuration from a `.env` file in the directory where it is run. The file must exist; copy the template below and fill it in:

```bash
# .env
NORYN_PROVIDER=openai
NORYN_MODEL=gpt-5
OPENAI_API_KEY=sk-...
OPENROUTER_API_KEY=sk-or-...
```

| Variable             | Description                                   | Default   |
| -------------------- | --------------------------------------------- | --------- |
| `NORYN_PROVIDER`     | LLM provider (`openai`, `openrouter`, `fake`) | `openai`  |
| `NORYN_MODEL`        | Model to use for the current provider         | `gpt-5`   |
| `OPENAI_API_KEY`     | API key for the OpenAI provider               | —         |
| `OPENROUTER_API_KEY` | API key for the OpenRouter provider           | —         |

CLI flags override the configured provider and model per session:

```bash
go run ./cmd/noryn --provider openrouter --model anthropic/claude-sonnet-4
```

Never commit API keys or your `.env` file to a repository.

## Usage

Start Noryn from the root of a Go project:

```bash
noryn
```

You get an interactive prompt:

```
Noryn
Type 'exit' to quit.

noryn> What does this project do?
```

Type `exit` to quit. Each prompt is combined with the discovered project context and your `AGENTS.md` instructions before being sent to the model. The agent may invoke tools to inspect files, run commands, or check Git state, and it keeps calling tools until it has enough context to answer.

### Supported tools

| Tool             | Purpose                                            |
| ---------------- | -------------------------------------------------- |
| `read_file`      | Read a file inside the project                     |
| `write_file`     | Write a file inside the project                    |
| `list_directory` | List directory contents                            |
| `shell`          | Execute commands from the project root             |
| `search`         | Recursive text search (skips binaries, `.git`, …)  |
| `git_status`     | Short status of the working tree                   |
| `git_diff`       | Diff of the working tree, optionally per path      |
| `git_log`        | Recent commit history with configurable limit      |

## How it works

```text
User
  ↓
CLI
  ↓
Project Discovery
  ↓
Instructions Builder
  ↓
Agent
  ├── LLM Client
  └── Tools
       ├── Filesystem
       ├── Search
       ├── Shell
       └── Git
```

The agent loop:

```text
Prompt → Send to LLM → Tool requested?
  ├── No → Return response
  └── Yes → Execute tool → Feed result back → Repeat
```

Noryn exposes the available tools, validates and executes tool calls, and feeds results back to the model until it responds without requesting a tool.

### Safety model

- Filesystem paths are resolved relative to the discovered project root and cannot escape it.
- Shell commands run inside the project directory.
- Git tools are read-only; they never modify the working tree or history.
- Actions are driven by explicit tool calls from the model rather than hidden side effects.

## Project structure

```text
noryn/
├── cmd/
│   └── noryn/
│       └── main.go
├── internal/
│   ├── agent/          agent loop and orchestration
│   ├── config/         .env loading and defaults
│   ├── context/        project context assembly
│   ├── instructions/   AGENTS.md loading
│   ├── llm/            provider abstraction and types
│   │   └── providers/  fake, openai, openrouter
│   ├── project/        root discovery and safe path resolution
│   └── tools/          filesystem, search, shell, and git tools
├── AGENTS.md
├── go.mod
└── go.sum
```

## Development

```bash
gofmt -w .
go test ./...
go build ./...
```

Provider-specific behavior is tested with deterministic HTTP test servers rather than live API calls. Any meaningful change should include tests.

## Roadmap

- Improve context management and agent reliability
- Add conversation history
- Add safer file-editing capabilities
- Add tool-execution confirmations and permissions
- Improve CLI and TUI interaction
- Add additional LLM providers
- Add Git modification capabilities
