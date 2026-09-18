
# Noryn

Noryn is an AI coding agent written in Go that provides an agentic software development experience from the terminal.

Noryn uses LLMs as its reasoning engine and provides controlled tools for inspecting, understanding, modifying, testing, and interacting with software projects.

Noryn is **not** an ML/DL framework and does not train or implement its own models.

## Core Principles

### Go First

* Use Go as the primary implementation language.
* Prefer the Go standard library when practical.
* Use external dependencies only when they provide clear value.
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
│       └── main.go
├── internal/
│   ├── agent/
│   │   ├── agent.go
│   │   └── agent_test.go
│   ├── config/
│   │   └── config.go
│   ├── instructions/
│   │   ├── builder.go
│   │   └── builder_test.go
│   ├── llm/
│   │   ├── client.go
│   │   ├── types.go
│   │   └── providers/
│   │       ├── fake/
│   │       ├── openai/
│   │       └── openrouter/
│   ├── project/
│   │   ├── project.go
│   │   └── scanner.go
│   └── tools/
│       ├── tool.go
│       ├── read_file.go
│       ├── write_file.go
│       ├── list_directory.go
│       ├── shell.go
│       ├── search.go
│       ├── git_status.go
│       ├── git_diff.go
│       └── git_log.go
├── AGENTS.md
├── go.mod
└── go.sum
```

---

## Architecture

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

### Project

The `project` package is responsible for discovering the project root and resolving paths safely inside that project.

The project root is currently discovered by searching for `go.mod`.

Tools that access project files should use the project's path resolution instead of directly trusting arbitrary paths.

### Instructions

The `instructions` package is responsible for loading project-specific instructions for the LLM.

Currently it reads:

```text
AGENTS.md
```

The instructions are loaded before the agent request and included in the initial LLM context.

### Agent

The `agent` package coordinates the interaction between the LLM and tools.

The agent:

1. Receives a request.
2. Provides available tool definitions to the LLM.
3. Sends the request to the LLM.
4. Checks whether the LLM requested tools.
5. Executes requested tools.
6. Adds tool results to the conversation.
7. Sends the updated conversation to the LLM.
8. Continues until the LLM returns a response without tool calls.

### LLM

The `llm` package provides a provider-independent abstraction for interacting with language models.

The main interface is:

```go
type Client interface {
    Chat(ctx context.Context, request Request) (Response, error)
}
```

Providers implement this interface.

Current providers:

* Fake
* OpenAI
* OpenRouter

The provider and model can be selected through configuration or CLI flags.

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

* `read_file`
* `write_file`
* `list_directory`
* `shell`
* `search`
* `git_status`
* `git_diff`
* `git_log`

Tool definitions describe the tool's purpose and parameters to the LLM.

Tool execution happens locally inside Noryn.

---

## Agent Loop

The core agent loop is:

```text
User Prompt
     ↓
Build Request
     ↓
Send Request to LLM
     ↓
Does the LLM request a tool?
     ├── No → Return Response
     │
     └── Yes
           ↓
       Execute Tool
           ↓
       Add Tool Result
           ↓
       Send Request Again
           ↓
         Repeat
```

The LLM decides when a tool is needed.

Noryn is responsible for:

* exposing available tools
* validating tool calls
* executing tools
* returning tool results to the LLM
* continuing the agent loop

---

## Tools

### Filesystem

Filesystem tools currently provide:

* file reading
* file writing
* directory listing

Filesystem paths are resolved relative to the discovered project root.

Tools must not access paths outside the project root.

### Search

The search tool:

* searches text recursively
* supports a file or directory as the search scope
* skips common irrelevant directories such as `.git`, `node_modules`, and `vendor`
* skips binary files
* skips files larger than the configured search limit
* limits the number of returned matches

Search is intended to help the agent locate relevant code before reading or modifying files.

### Shell

The shell tool executes commands from the project root.

Shell output and command errors are returned to the agent.

Shell execution must remain transparent.

Do not silently execute commands that were not requested or generated by the agent.

### Git

Current Git capabilities:

* `git_status`
* `git_diff`
* `git_log`

Current Git tools are read-oriented and do not modify Git history.

Git modification capabilities such as `add` and `commit` may be added later.

---

## Configuration

Configuration is loaded from `.env`.

Current configuration variables include:

```text
NORYN_PROVIDER
NORYN_MODEL
OPENAI_API_KEY
OPENROUTER_API_KEY
```

CLI flags can override the configured provider and model:

```bash
go run ./cmd/noryn --provider <provider> --model <model> "<prompt>"
```

API keys must never be committed to Git or included directly in source code.

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

* agent behavior
* tool registration
* tool execution
* tool errors
* filesystem tools
* search
* Git tools
* OpenAI provider behavior
* instructions loading

Provider-specific behavior should preferably be tested with deterministic HTTP test servers rather than relying on live API calls.

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

---

## Current State

Noryn currently has the foundation of a functional terminal AI coding agent.

Implemented:

* LLM abstraction
* Agent loop
* Fake LLM provider
* OpenAI provider
* OpenRouter provider
* Provider selection
* CLI provider and model flags
* Project discovery
* Safe project path resolution
* Project instructions loading through `AGENTS.md`
* File reading
* File writing
* Directory listing
* Shell execution
* Recursive text search
* Binary-file detection
* Search file-size limits
* Search result limits
* Git status
* Git diff
* Git log
* Tool definitions
* Tool error handling
* Automated tests

The current focus is improving context management and agent reliability.

---

## Roadmap

Potential future work:

1. Improve context management.
2. Add conversation history.
3. Improve tool-result handling.
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
