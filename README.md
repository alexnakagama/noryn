# Noryn

Noryn is an AI coding agent that runs in your terminal. Built in Go, it pairs an LLM with a set of curated, project-scoped tools so you can explore, understand, and modify a codebase from the command line.

Noryn is an **agent orchestration layer**, not a model. It delegates reasoning to your chosen LLM provider and focuses on everything around it: tool calling, context assembly, history management, and safe execution.

> **Status:** early-stage. The agent loop, tooling, and context management work, and a Bubble Tea TUI is under active development. Expect rough edges and ongoing change.

---

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Supported tools](#supported-tools)
- [How it works](#how-it-works)
- [Context management](#context-management)
- [Safety model](#safety-model)
- [Project structure](#project-structure)
- [Development](#development)
- [Roadmap](#roadmap)

---

## Features

- **Terminal UI** built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), featuring a welcome screen and an interactive chat screen
- **Project discovery** by walking up to the nearest `go.mod`
- **Project context assembly** — a snapshot of the project's text files is loaded into the initial prompt
- **Project-specific instructions** loaded from `AGENTS.md`
- **Agent loop** with automatic, validated tool calls
- **Streaming and non-streaming** agent execution with a unified event stream
- **Context management** — conversation history is trimmed to a token budget, and tool results are truncated to keep prompts bounded
- **8 built-in tools** covering filesystem, search, shell, and Git inspection
- **Path confinement** — file operations cannot escape the project root
- **Swappable LLM providers** — OpenAI, OpenRouter, and a deterministic Fake provider for testing
- **Provider and model selection** via `.env` or CLI flags

## Requirements

- **Go 1.27 or newer**
- **An API key** for your chosen LLM provider (OpenAI or OpenRouter)
- **Git** on your `PATH` (for the Git tools)

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

You may want to move the binary onto your `PATH`:

```bash
mv noryn /usr/local/bin/noryn
```

## Configuration

Noryn reads configuration from a `.env` file in the directory where it is run (via `godotenv`). The file must exist; copy the template below and fill it in:

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
go run ./cmd/noryn --provider fake
```

> **Security:** never commit API keys or your `.env` file to a repository. `.env` is already listed in `.gitignore`.

### Providers

| Provider      | Value        | Requires             | Streaming |
| ------------- | ------------ | -------------------- | --------- |
| OpenAI        | `openai`     | `OPENAI_API_KEY`     | No        |
| OpenRouter    | `openrouter` | `OPENROUTER_API_KEY` | Yes       |
| Fake          | `fake`       | —                    | No        |

The Fake provider is deterministic and used for local testing of the agent loop without any external calls.

## Usage

Start Noryn from the root of a Go project:

```bash
noryn
```

A full-screen welcome screen appears:

```text
   ███╗   ██╗ ██████╗ ██████╗ ██╗   ██╗███╗   ██╗
   ████╗  ██║██╔═══██╗██╔══██╗╚██╗ ██╔╝████╗  ██║
   ██╔██╗ ██║██║   ██║██████╔╝ ╚████╔╝ ██╔██╗ ██║
   ██║╚██╗██║██║   ██║██╔══██╗  ╚██╔╝  ██║╚██╗██║
   ██║ ╚████║╚██████╔╝██║  ██║   ██║   ██║ ╚████║
   ╚═╝  ╚═══╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═══╝

Noryn v0.1.0 — AI coding assistant

                      Press Enter to start
```

Press `Enter` to open the chat screen, which shows a conversation area, an input bar, and a sidebar with status, project, branch, model, and token usage.

Each prompt is combined with the discovered project context and your `AGENTS.md` instructions before being sent to the model. The agent calls tools to inspect files, run commands, or check Git state, and keeps iterating until it has enough context to answer.

## Supported tools

| Tool             | Purpose                                                        | Limits                                          |
| ---------------- | -------------------------------------------------------------- | ----------------------------------------------- |
| `read_file`      | Read a file inside the project                                 | Files up to 1 MB                                |
| `write_file`     | Write content to a file inside the project                     | —                                               |
| `list_directory` | List the contents of a project directory                       | —                                               |
| `shell`          | Execute a shell command from the project root                  | 30-second timeout                               |
| `search`         | Recursive text search with `path:line:content` results         | Skips binaries, files > 1 MB, up to 200 matches |
| `git_status`     | Short status of the working tree                               | Read-only                                       |
| `git_diff`       | Diff of the working tree, optionally scoped to a path          | Read-only                                       |
| `git_log`        | Recent commit history with a configurable limit (1–50)         | Read-only                                       |

Search and project context both skip common irrelevant directories (`.git`, `node_modules`, `vendor`), binary files, and oversized files.

## How it works

```text
User
  ↓
TUI (Bubble Tea)
  ↓
Project Discovery        → find nearest go.mod
  ↓
Project Context          → snapshot of repo text files
  ↓
Instructions Builder     → load AGENTS.md + system prompt
  ↓
Agent
  ├── LLM Client          → provider abstraction
  ├── Context Manager     → token-budgeted history
  ├── Tool Registry       → validates and executes tool calls
  └── Event Stream        → text, tool calls, tool results, done, error
```

### The agent loop

Non-streaming (`Chat`):

```text
Prompt → Send to LLM → Tool requested?
  ├── No → Return response
  └── Yes → Execute tool → Feed result back → Repeat
```

Streaming (`ChatStream`):

```text
Prompt → Stream chunks from LLM
          ├── text chunk  → emit EventText
          └── tool call   → emit EventToolCall, execute, emit EventToolResult
        → emit EventDone / EventError when finished
```

In both paths Noryn exposes the available tools, validates tool calls through the registry, executes them, and feeds results back to the model until it responds without requesting a tool.

### Startup flow

`cmd/noryn/main.go` wires everything together:

1. Load configuration from `.env` (overridable with `--provider` and `--model`).
2. Discover the project root by searching upward for `go.mod`.
3. Build the project context (a structured snapshot of the repo's text files).
4. Build the LLM instructions from the system prompt and the project's `AGENTS.md`.
5. Pick the LLM provider based on configuration.
6. Register all tools in the `tools.Registry`.
7. Create the `agent.Agent` and hand it to the TUI.

## Context management

Noryn keeps the conversation bounded so prompts stay within the model's context window:

- **Token budget** — the `ContextManager` keeps the most recent messages that fit within a 32,000-token budget, using a rough estimate (4 characters ≈ 1 token).
- **Turn integrity** — trimming never starts mid-way through a tool-call turn, so the model never sees orphaned `tool` messages.
- **Tool-result truncation** — individual tool results are capped at 10,000 characters and flagged as truncated.
- **Iteration guard** — the agent stops after 20 consecutive tool iterations to prevent runaway loops.

## Safety model

- **Path confinement** — filesystem paths are resolved relative to the discovered project root and cannot escape it.
- **Binary detection** — project context and search skip binary files automatically.
- **Shell scope** — shell commands run inside the project directory with a hard timeout.
- **Read-only Git** — Git tools only inspect state; they never modify the working tree or history.
- **Explicit actions** — all operations are driven by explicit tool calls from the model rather than hidden side effects.

## Project structure

```text
noryn/
├── cmd/
│   └── noryn/
│       └── main.go               entry point: wiring, flags, startup
├── internal/
│   ├── agent/                    agent loop, context manager, streaming events
│   ├── config/                   .env loading and defaults
│   ├── instructions/             system prompt, AGENTS.md loading, prompt assembly
│   ├── llm/                      provider abstraction and shared types
│   │   └── providers/            fake, openai, openrouter, provider factory
│   ├── project/                  root discovery, safe path resolution, context
│   ├── tools/                    tool interface, registry, and 8 tools
│   └── tui/                      Bubble Tea UI, screen router, styles
├── AGENTS.md
├── go.mod
├── go.sum
└── .gitignore
```

### Package overview

| Package       | Responsibility                                                            |
| ------------- | ------------------------------------------------------------------------- |
| `agent`       | Runs the agent loop (`Chat`/`ChatStream`), manages history and token budget |
| `config`      | Loads configuration from `.env` and applies defaults                      |
| `instructions`| Builds the system prompt, loads `AGENTS.md`, assembles the initial prompt  |
| `llm`         | Defines `Client`/`StreamingClient` and the request/response types          |
| `providers`   | Factory for Fake, OpenAI, and OpenRouter clients                           |
| `project`     | Finds the project root, resolves paths safely, builds project context      |
| `tools`       | `Tool` interface, `Registry`, and the concrete tools                       |
| `tui`         | Bubble Tea application with welcome and chat screens                       |

## Development

```bash
gofmt -w .
go test ./...
go build ./...
```

- Every meaningful change should include tests.
- Provider-specific behavior is tested with deterministic HTTP test servers rather than live API calls.
- Tests currently cover the agent loop and context manager, tool registration and execution, filesystem and search tools, Git tools, project discovery and path confinement, instructions loading, and provider behavior.

## Roadmap

- Wire the chat screen to the streaming agent (prompt dispatch, live responses)
- Add conversation history persistence
- Add safer file-editing capabilities
- Add tool-execution confirmations and permissions
- Improve context management and agent reliability
- Add additional LLM providers
- Add Git modification capabilities carefully
- Improve TUI interaction (a real live-conversation view, multiple sessions)
