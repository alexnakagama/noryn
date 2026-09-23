# Noryn

Noryn is an AI coding agent written in Go that provides an agentic software development experience from the terminal.

Noryn uses LLMs as its reasoning engine and provides controlled tools for inspecting, understanding, modifying, testing, and interacting with software projects.

Noryn is **not** an ML/DL framework and does not train or implement its own models.

---

## Core Principles

### Go First

* Use Go as the primary implementation language.
* Prefer the Go standard library when practical.
* Use external dependencies only when they provide clear value (Bubble Tea ecosystem and `godotenv` are the current exceptions).
* Write idiomatic Go.
* Prefer small interfaces and explicit error handling.
* Prefer composition over inheritance.
* Keep data structures and control flow simple.

### Keep Architecture Simple

* Separate responsibilities clearly.
* Create abstractions only when they solve an actual problem.
* Avoid unnecessary interfaces, factories, and generic abstractions.
* Avoid deep package hierarchies.
* Do not rewrite working code without a reason.
* Prefer simple implementations over premature abstractions.

### Preserve User Work

* Never silently destroy user files or Git changes.
* Keep modifications focused on the requested task.
* Inspect existing code before modifying it.
* Do not overwrite unrelated changes.
* Prefer small, understandable changes.

---

## Project Structure

```text
noryn/
├── cmd/
│   └── noryn/
│       └── main.go                     entry point: config, discovery, wiring, flags
├── internal/
│   ├── agent/
│   │   ├── agent.go                    agent loop (Chat, ChatStream)
│   │   ├── agent_test.go
│   │   ├── context.go                  ContextManager, token-budgeted history
│   │   ├── context_test.go
│   │   ├── event.go                    streaming event types
│   │   ├── event_test.go
│   │   ├── token_counter.go            token estimation
│   │   └── token_counter_test.go
│   ├── config/
│   │   └── config.go                   .env loading and defaults
│   ├── instructions/
│   │   ├── builder.go                  loads AGENTS.md
│   │   ├── builder_test.go
│   │   ├── prompt.go                   assembles the initial prompt
│   │   ├── prompt_test.go
│   │   └── system.go                   base system prompt
│   ├── llm/
│   │   ├── client.go                   Client and StreamingClient interfaces
│   │   ├── types.go                    request/response/tool/stream types
│   │   └── providers/
│   │       ├── providers.go            provider factory
│   │       ├── fake/                   deterministic test client
│   │       ├── openai/                 OpenAI Responses API
│   │       └── openrouter/             OpenRouter chat completions + streaming
│   ├── project/
│   │   ├── context.go                  project file snapshot (text files only)
│   │   ├── context_test.go
│   │   ├── project.go                  safe path resolution
│   │   ├── project_test.go
│   │   ├── scanner.go                  root discovery via go.mod
│   │   └── scanner_test.go
│   ├── tools/
│   │   ├── tool.go                     Tool interface
│   │   ├── registry.go                 tool registration, lookup, definitions
│   │   ├── registry_test.go
│   │   ├── read_file.go / _test.go
│   │   ├── write_file.go / _test.go
│   │   ├── list_directory.go / _test.go
│   │   ├── shell.go / _test.go
│   │   ├── search.go / _test.go
│   │   ├── git_status.go / _test.go
│   │   ├── git_diff.go / _test.go
│   │   └── git_log.go / _test.go
│   └── tui/
│       ├── model.go                    top-level Bubble Tea model
│       ├── router.go                   screen routing
│       ├── tui.go                      program bootstrap (alternate screen)
│       ├── screens/
│       │   ├── chat/                   chat screen (model + view)
│       │   └── welcome/                welcome screen (model + view)
│       └── styles/                     Tokyo Night lipgloss palette
├── .gitignore                          ignores .env
├── AGENTS.md
├── go.mod
└── go.sum
```

---

## Architecture

```text
User
  ↓
TUI (Bubble Tea: welcome → chat)
  ↓
Project Discovery         FindRoot walks up to the nearest go.mod
  ↓
Project Context           BuildContext snapshots the repo's text files
  ↓
Instructions Builder      SystemPrompt + AGENTS.md + project context + user prompt
  ↓
Agent
  ├── LLM Client           provider abstraction (Chat / ChatStream)
  ├── ContextManager       token-budgeted history trimming
  ├── Registry             tool registration, validation, execution
  └── Event stream         text / tool calls / tool results / done / error
```

### Project

The `project` package is responsible for discovering the project root and resolving paths safely inside that project.

* `FindRoot` walks upward from the start directory until it finds a `go.mod`; it returns an error when no root exists.
* `Discover` returns a `*Project` pointing at the root.
* `ResolvePath` resolves a relative or absolute path against the root and **rejects any path that escapes it**.
* `BuildContext` walks the project and returns a sorted list of text files (skipping `.git`, `node_modules`, `vendor`, files larger than 1 MB, and binaries — detected via `http.DetectContentType`).

Tools that access project files must use the project's path resolution instead of trusting arbitrary paths.

### Instructions

The `instructions` package is responsible for loading and assembling the LLM's initial context.

* `system.go` holds the base `SystemPrompt` (agent identity and tool-use guidance).
* `builder.go` loads the project's `AGENTS.md` (resolved through the project).
* `prompt.go` assembles the final prompt as: system prompt + project instructions + project context + user prompt.

The instructions and project context are included in the initial LLM context, before the first user turn.

### Agent

The `agent` package coordinates the interaction between the LLM and tools.

The agent:

1. Receives a request.
2. Appends it to conversation history.
3. Builds the request through the `ContextManager` and attaches tool definitions.
4. Sends the request to the LLM.
5. If the LLM requested tools, executes each call through the registry, truncates and appends the results, and repeats.
6. Continues until the LLM returns a response without tool calls.

The agent supports two execution paths:

* `Chat` — non-streaming call/response loop.
* `ChatStream` — streaming loop that requires a `llm.StreamingClient` and emits an `agent.Event` channel with `EventText`, `EventToolCall`, `EventToolResult`, `EventDone`, and `EventError`.

Both paths expose `SetToolCallHandler` / `SetToolResultHandler` callbacks for UI integration.

### Context management

`ContextManager` keeps prompts bounded:

* Trims history to a fixed token budget (`maxContextTokens`).
* Estimates tokens with `TokenCounter` (an `EstimateTokenCounter` heuristic: 4 characters ≈ 1 token).
* Never trims into the middle of a tool-call turn.
* Individual tool results are truncated to `maxToolResultLength` (10,000 chars).
* The agent aborts after `maxToolIterations` (20) consecutive tool-turn iterations.

### LLM

The `llm` package provides provider-independent abstractions.

```go
type Client interface {
    Chat(ctx context.Context, request Request) (Response, error)
}

type StreamingClient interface {
    Client
    ChatStream(ctx context.Context, request Request) (<-chan StreamChunk, error)
}
```

Providers:

* **Fake** — deterministic, no network; useful for tests and local development.
* **OpenAI** — uses the Responses API (`/v1/responses`), non-streaming.
* **OpenRouter** — uses chat completions (`/v1/chat/completions`) with SSE streaming and incremental tool-call accumulation.

The provider is selected through configuration or the `--provider` CLI flag. Missing API keys for OpenAI/OpenRouter are reported at startup.

### Tools

Tools implement the following interface:

```go
type Tool interface {
    Name() string
    Description() string
    Definition() llm.ToolDefinition
    Execute(arguments string) (string, error)
}
```

Current tools:

* `read_file`   — 1 MB file limit
* `write_file`
* `list_directory`
* `shell`       — runs from the project root, 30-second timeout
* `search`      — 1 MB per-file limit, 200 result cap, skips binaries and noise dirs
* `git_status`
* `git_diff`    — optional path scope
* `git_log`     — limit 1–50

`Registry` holds registered tools by name, reports duplicate registrations, exposes `Definitions()` for the LLM, and dispatches `Execute` for validated tool calls. Tool definitions describe the tool's purpose and parameters to the LLM. Tool execution happens locally inside Noryn.

### TUI

The `tui` package is the Bubble Tea front end:

* `Router` switches between `ScreenWelcome` and `ScreenChat`.
* The welcome screen renders the Noryn logo and a "Press Enter to start" hint.
* The chat screen renders a conversation area, an input bar, and a status sidebar (status, project, branch, model, usage). Prompt dispatch to the streaming agent is still work in progress.
* `styles` defines the Tokyo Night color palette and shared lipgloss styles.

---

## Agent Loop

The core agent loop is:

```text
User Prompt
     ↓
Build Request (ContextManager + tool definitions)
     ↓
Send Request to LLM
     ↓
Does the LLM request a tool?
     ├── No → Return Response
     │
     └── Yes
           ↓
       Execute Tool (Registry)
           ↓
       Truncate + Add Tool Result
           ↓
       Send Request Again
           ↓
         Repeat
```

The streaming variant applies the same loop but surfaces progress (text deltas, tool calls, tool results) as `agent.Event`s.

The LLM decides when a tool is needed. Noryn is responsible for:

* exposing available tools
* validating tool calls
* executing tools
* returning tool results to the LLM
* continuing the agent loop
* keeping history within the token budget

---

## Tools

### Filesystem

Filesystem tools require a `*project.Project` and resolve every path through `ResolvePath`; tools must not access paths outside the project root.

* `read_file` — reads text files up to 1 MB, rejects larger files.
* `write_file` — writes content (mode `0644`) to a project path.
* `list_directory` — lists immediate entries of a directory.

### Search

The search tool:

* searches text recursively
* requires both `query` and `path` (a file or directory)
* skips common irrelevant directories (`.git`, `node_modules`, `vendor`)
* skips binary files (detected by embedded null bytes)
* skips files larger than 1 MB
* caps results at 200 matches of `path:line:content`

Search is intended to help the agent locate relevant code before reading or modifying files.

### Shell

The shell tool executes a command via `sh -c` from the project root with a hard 30-second timeout.

Shell output and command errors are returned to the agent.

Shell execution must remain transparent: do not silently execute commands that were not requested or generated by the agent.

### Git

Current Git capabilities: `git_status` (short), `git_diff` (optional path), `git_log` (1–50 commits).

Current Git tools are read-oriented and do not modify Git history. Git modification capabilities (e.g. `add`, `commit`) may be added later, carefully.

---

## Configuration

Configuration is loaded from `.env` (must exist) using `godotenv`.

Current configuration variables:

```text
NORYN_PROVIDER          default: openai
NORYN_MODEL             default: gpt-5
OPENAI_API_KEY
OPENROUTER_API_KEY
```

CLI flags can override the configured provider and model:

```bash
go run ./cmd/noryn --provider <provider> --model <model>
```

`providers.New` resolves the configured provider name to a client and validates API keys. API keys must never be committed to Git or included directly in source code.

---

## Testing

Every meaningful change should have appropriate tests.

Useful development commands:

```bash
gofmt -w .
go test ./...
go build ./...
```

Tests currently cover:

* agent loop behavior (including streaming events)
* context management and token budgeting
* tool registration, lookup, and duplicate detection
* tool execution and tool errors
* filesystem tools (`read_file`, `write_file`, `list_directory`)
* search behavior
* Git tools (`git_status`, `git_diff`, `git_log`)
* OpenAI and OpenRouter provider behavior
* project discovery and path confinement
* instructions loading and prompt assembly

Provider-specific behavior should preferably be tested with deterministic HTTP test servers rather than relying on live API calls. The Fake provider is preferred for testing the agent loop.

---

## Development Rules

* Inspect existing code before modifying it.
* Make focused changes.
* Preserve existing behavior unless the task requires changing it.
* Add tests for new behavior.
* Run tests after meaningful changes.
* Format Go code.
* Keep error messages useful.
* Prefer explicit and simple implementations.
* Avoid premature abstractions.
* Do not introduce ML/DL functionality unless explicitly required.
* Do not add Python as a runtime dependency unless explicitly required.
* Do not silently modify or destroy user data.
* Do not silently modify Git history.
* Do not rewrite working code unnecessarily.
* Keep the docs (README.md and AGENTS.md) in sync with the code.

---

## Current State

Noryn has the foundation of a functional terminal AI coding agent.

Implemented:

* LLM abstraction with `Chat` and `ChatStream`
* Agent loop with tool-call handling
* Streaming agent events
* Context manager with token-budgeted history trimming
* Tool-result truncation and iteration guards
* Fake, OpenAI, and OpenRouter providers
* Provider selection via `.env` and CLI flags
* Project discovery and safe path resolution
* Project context snapshot generation
* Instructions loading through `AGENTS.md` and prompt assembly
* Filesystem tools (read, write, list)
* Shell execution
* Recursive text search
* Git inspection tools
* Tool registry with validation
* Bubble Tea TUI with welcome and chat screens
* Automated tests

The current focus is improving context management and agent reliability, and wiring the chat screen to the streaming agent.

---

## Roadmap

Potential future work:

1. Wire the chat screen to the streaming agent (prompt dispatch, live responses).
2. Add conversation history persistence.
3. Improve context management and tool-result handling.
4. Add safer file editing capabilities.
5. Add tool execution permissions and confirmations where appropriate.
6. Improve CLI and TUI interaction.
7. Add additional LLM providers.
8. Improve agent reliability and testing.
9. Add more advanced project understanding.
10. Add Git modification capabilities carefully.

---

## Important

Noryn is an AI coding agent, not an AI model.

The project should focus on:

* agent orchestration
* context management
* LLM integration
* tool calling
* project understanding
* safe code modification
* testing
* developer experience
* reliability

Do not turn Noryn into an ML/DL training framework unless that becomes an explicit project requirement.
