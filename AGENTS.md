# Noryn — Agent Instructions

## Project Overview

Noryn is an AI coding agent for the terminal, built entirely around Go.

The goal of Noryn is to provide an agentic software development experience directly from the command line. Noryn uses LLMs as its reasoning engine and provides them with controlled tools to inspect, understand, modify, and test software projects.

Noryn is **not** an ML/DL framework and does not train or implement machine learning models.

The project focuses on **AI Engineering and Software Engineering**:

* LLM integration
* Agent loops
* Tool calling
* Codebase understanding
* File manipulation
* Shell execution
* Git integration
* Context management
* Terminal UX
* Provider abstraction
* Reliable and safe agent behavior

---

# Core Principles

## 1. Go First

Noryn is a Go project.

Prefer the Go standard library whenever it provides a reasonable solution.

Avoid introducing dependencies without a clear reason.

Use idiomatic Go:

* Small interfaces
* Explicit error handling
* Composition over inheritance
* Simple data structures
* Context propagation
* Clear package responsibilities
* Avoid unnecessary abstractions

Do not introduce patterns simply because they are common in other languages.

---

## 2. Keep the Architecture Simple

Noryn should follow clean architecture principles without overengineering.

The main goal is separation of responsibilities and dependency inversion.

Do not create abstractions before they are needed.

Avoid:

* Interfaces with only one trivial implementation when no substitution is required
* Excessive factories
* Unnecessary repositories
* Deep package hierarchies
* Generic abstractions with no concrete use case

Prefer simple, explicit code.

---

# Architecture

The current conceptual architecture is:

```text
User
 │
 ▼
Terminal UI
 │
 ▼
Agent
 │
 ├──────────────► LLM
 │
 ├──────────────► Tools
 │                  ├── Filesystem
 │                  ├── Search
 │                  ├── Shell
 │                  └── Git
 │
 └──────────────► Context
```

The agent coordinates the workflow.

The LLM provides reasoning and generates tool calls.

Tools perform operations in the real environment.

The agent feeds tool results back to the LLM.

---

# Project Structure

The intended project structure is:

```text
noryn/
├── cmd/
│   └── noryn/
│       └── main.go
│
├── internal/
│   ├── agent/
│   │   ├── agent.go
│   │   ├── loop.go
│   │   └── context.go
│   │
│   ├── llm/
│   │   ├── client.go
│   │   └── providers/
│   │
│   ├── tools/
│   │   ├── tool.go
│   │   ├── filesystem.go
│   │   ├── search.go
│   │   ├── shell.go
│   │   └── git.go
│   │
│   ├── project/
│   │   ├── project.go
│   │   └── scanner.go
│   │
│   ├── context/
│   │   └── ...
│   │
│   ├── tui/
│   │   └── ...
│   │
│   └── config/
│       └── ...
│
├── go.mod
├── go.sum
└── AGENTS.md
```

The structure can evolve as the project grows. Do not create empty packages simply to match this structure.

---

# Agent

The `agent` package contains the core agent loop.

The agent is responsible for coordinating:

```text
User Input
    ↓
Build Context
    ↓
Send Request to LLM
    ↓
Receive Response
    ↓
Does the LLM request a Tool?
    │
    ├── No → Return Response
    │
    └── Yes
          ↓
       Execute Tool
          ↓
       Get Result
          ↓
       Add Result to Context
          ↓
       Call LLM Again
```

The agent should not know how individual tools are implemented.

The agent should not directly manipulate the filesystem or execute shell commands.

The agent should communicate with those capabilities through the tool abstraction.

---

# LLM Layer

The `llm` package is responsible for communication with language models.

The agent should depend on an abstraction rather than directly depending on a specific provider.

Conceptually:

```go
type Client interface {
    Chat(ctx context.Context, request Request) (Response, error)
}
```

Possible implementations may include:

```text
LLM Client
├── OpenAI
├── Anthropic
├── Gemini
├── OpenAI-compatible APIs
└── Local models
```

The exact provider implementation should remain isolated from the agent.

Adding a new provider should not require changing the agent loop.

---

# Tools

Tools are capabilities that the LLM can request from Noryn.

Initial tools should include:

```text
read_file
write_file
edit_file
list_directory
search_code
run_command
git_status
git_diff
```

Every tool should have:

* A clear name
* A description
* A well-defined input schema
* A predictable output
* Explicit error handling

Tools must not silently perform destructive operations.

---

# Filesystem

Filesystem tools allow the agent to understand and modify the current project.

The agent should be able to:

* List directories
* Read files
* Create files
* Modify files
* Delete files when explicitly allowed
* Search project contents

The implementation must consider:

* Relative paths
* Absolute paths
* Path traversal
* Binary files
* Large files
* Missing files
* Permission errors

Never assume a file exists.

Never silently overwrite important files without the appropriate permission or confirmation mechanism.

---

# Shell Execution

Noryn should eventually allow the LLM to execute shell commands.

Example:

```text
go test ./...
git diff
go build ./...
```

Shell execution is a sensitive capability.

Commands should be visible to the user before or during execution.

Potentially destructive operations should require explicit user approval.

Do not implement unrestricted hidden shell execution.

---

# Git

Git integration should help the agent understand and safely modify repositories.

Initial capabilities:

```text
git status
git diff
git diff --cached
git log
```

Later capabilities may include:

```text
git add
git commit
```

Git operations should never discard user changes automatically.

Noryn must preserve changes that existed before the agent started working.

---

# Context Management

Context management is a core part of Noryn.

The agent should not blindly send the entire repository to the LLM.

Instead, it should progressively gather relevant information.

Example:

```text
User request
     ↓
Project structure
     ↓
Relevant files
     ↓
Search results
     ↓
Specific file contents
     ↓
Tool results
     ↓
LLM context
```

The system should eventually support:

* Conversation history
* Project context
* File context
* Tool results
* Context limits
* Context pruning
* Summarization
* Relevant-file selection

---

# Agent Behavior

The agent should be autonomous enough to complete development tasks but transparent enough for the developer to understand what it is doing.

For a request such as:

```text
Fix the authentication bug.
```

The agent should generally:

```text
1. Understand the request
2. Inspect the project
3. Locate relevant code
4. Read the necessary files
5. Determine a possible solution
6. Modify the code
7. Run relevant tests
8. Inspect failures
9. Iterate if necessary
10. Explain what changed
```

Do not make arbitrary changes unrelated to the user's request.

Prefer the smallest change that correctly solves the problem.

---

# Code Changes

Before modifying code:

1. Inspect the existing implementation.
2. Understand the surrounding code.
3. Follow existing project conventions.
4. Avoid unnecessary refactoring.
5. Make the smallest reasonable change.

After modifying code:

1. Format the code.
2. Run relevant tests.
3. Run static checks when available.
4. Inspect the resulting diff.
5. Report relevant failures clearly.

Never claim that tests passed if they were not actually executed.

---

# Error Handling

Errors should be explicit and useful.

Prefer:

```go
if err != nil {
    return fmt.Errorf("read project file: %w", err)
}
```

over silently ignoring errors.

Do not use panic for normal runtime errors.

Errors should retain their original cause where possible.

---

# Concurrency

Use concurrency only when it provides a real benefit.

Potential areas where concurrency may eventually be useful:

* Parallel file searches
* Tool execution where operations are independent
* Streaming
* Multiple background tasks

Do not introduce goroutines merely to make code "more concurrent".

Always consider cancellation through `context.Context`.

---

# Configuration

Configuration should be separated from application logic.

Sensitive information such as API keys must never be hardcoded.

Use environment variables or configuration files.

Never commit secrets.

---

# Testing

Tests should focus on behavior.

Important areas to test:

* Agent loop
* Tool execution
* Tool input validation
* LLM request/response handling
* Context management
* Filesystem operations
* Shell command handling
* Configuration
* Error cases

Use dependency injection where it makes testing substantially easier.

Do not create abstractions solely for the sake of achieving test coverage.

---

# Dependencies

Before adding a dependency:

1. Check whether the standard library is sufficient.
2. Check whether the dependency is actively maintained.
3. Check whether it solves a meaningful problem.
4. Keep the dependency surface small.

Prefer well-established Go libraries when a dependency is justified.

---

# Development Commands

The project should remain compatible with standard Go tooling.

Run:

```bash
go fmt ./...
```

Run tests:

```bash
go test ./...
```

Build:

```bash
go build ./...
```

Run the application during development:

```bash
go run ./cmd/noryn
```

---

# Development Roadmap

Noryn should be developed incrementally.

## Phase 1 — Basic CLI

* [ ] CLI entry point
* [ ] Interactive input
* [ ] Basic output
* [ ] Configuration

## Phase 2 — LLM Integration

* [ ] LLM client abstraction
* [ ] Provider implementation
* [ ] Chat requests
* [ ] Response handling
* [ ] Streaming

## Phase 3 — Agent Loop

* [ ] Agent state
* [ ] Message history
* [ ] LLM → tool decision
* [ ] Tool execution
* [ ] Tool results
* [ ] Iterative execution

## Phase 4 — Development Tools

* [ ] Read files
* [ ] List directories
* [ ] Search code
* [ ] Write files
* [ ] Edit files
* [ ] Shell execution
* [ ] Git integration

## Phase 5 — Context

* [ ] Project discovery
* [ ] Relevant file selection
* [ ] Context limits
* [ ] Context pruning
* [ ] Conversation management

## Phase 6 — Safety

* [ ] Tool permissions
* [ ] Command approval
* [ ] Destructive-operation confirmation
* [ ] Clear execution display

## Phase 7 — Advanced Agent Capabilities

* [ ] Multiple LLM providers
* [ ] Better codebase indexing
* [ ] Improved context selection
* [ ] Session persistence
* [ ] Agent configuration
* [ ] Tool extensibility

---

# Important Rules for the Coding Agent

When working on Noryn:

1. Read this file before making architectural changes.
2. Inspect existing code before modifying it.
3. Do not rewrite working code unnecessarily.
4. Do not introduce ML/DL functionality unless explicitly requested.
5. Do not add Python as a runtime dependency.
6. Keep the core implementation in Go.
7. Prefer the standard library when practical.
8. Keep packages focused.
9. Avoid premature abstractions.
10. Never hide shell commands from the user.
11. Never silently destroy user data or existing Git changes.
12. Run tests after meaningful changes.
13. Do not claim success without verification.
14. Keep changes focused on the requested task.
15. Update documentation when architectural behavior changes.

---

# Long-Term Vision

Noryn should become a capable, extensible AI coding environment that runs directly in the terminal.

The long-term architecture should allow:

```text
                    Noryn
                      │
              ┌───────┴───────┐
              │               │
            Agent             TUI
              │
       ┌──────┼──────┐
       │      │      │
      LLM   Tools  Context
       │      │      │
       │      │      ├── Project
       │      │      ├── Files
       │      │      └── History
       │      │
       │      ├── Filesystem
       │      ├── Search
       │      ├── Shell
       │      └── Git
       │
       └── Multiple Providers
```

Noryn should remain primarily a **Go software engineering project with AI capabilities**, not a machine learning research project.
