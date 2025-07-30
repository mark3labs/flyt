# Flyt Codebase Guidelines for AI Agents

## Build/Test Commands
- Run all tests: `go test -race ./...`
- Run single test: `go test -race -v -run TestName`
- Format code: `go fmt ./...`
- Lint code: `go vet ./...`
- Build: `go build`

## Code Style
- Package: Always start with package comment describing purpose
- Imports: Group stdlib, then external (none used), then internal
- Naming: Use Go conventions (camelCase for private, PascalCase for public)
- Interfaces: Define minimal interfaces, embed for composition
- Error handling: Return explicit errors, no panic in library code
- Comments: Document all exported types/functions with godoc format
- Testing: Use table-driven tests, test both success and error cases
- Concurrency: Use sync primitives directly, no abstraction layers
- Dependencies: ZERO external dependencies - stdlib only

## Project Structure
- Single package design (package flyt)
- Core types: Node, Flow, Shared, Params, Action
- Helper functions for batch processing patterns
- Comprehensive test coverage with TestNode implementations