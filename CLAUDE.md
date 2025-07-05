# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Testing
- Run all tests: `go test`
- Update golden files: `go test -update`
- Run specific test: `go test -run TestName`

### Building and Development
- Build: `go build`
- Format code: `go fmt ./...`
- Vet code: `go vet ./...`

## Architecture

Gilt is a Go library for golden file testing with a flexible, interface-based architecture centered around generic types.

### Core Components

**Golden[TActual, TExpected]**: The main struct that orchestrates golden file testing. It uses composition of interfaces to handle different aspects:
- `IsUpdater`: Determines when to update golden files (controlled by `-update` flag)
- `Saver[TActual]`: Converts actual values to bytes for storage
- `Loader[TExpected]`: Reads and converts golden file data to expected values
- `ReadWriter`: Handles file I/O operations

**Specialized Constructors**:
- `NewStringGolden()`: For string-based golden tests
- `NewBytesGolden()`: For byte slice golden tests  
- `NewJSONGolden[T,U]()`: For JSON serializable types
- `New[T,U]()`: Generic constructor with custom options

### File Organization

Golden files are stored in `testdata/{TestName}/golden/` with `.golden` extension by default. The path structure can be customized using `GoldenFileWithPathHandler` option.

### Key Patterns

The library uses a functional options pattern (`Option[TActual, TExpected]`) for configuration and interface composition for extensibility. Tests call `golden.Assert()` with actual values, expected test names, and custom assertion functions.