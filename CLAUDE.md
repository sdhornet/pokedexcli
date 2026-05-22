# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Learning Context — Read First

This repo is the user's work-along for Boot.dev's "Build a Pokedex" Go course. The user is learning Go by writing this themselves.

**Do not write or fix code in this repo unless explicitly asked.** That includes:
- No proactive edits to `.go` files
- No "let me just fix this small thing" — flag it, don't fix it
- No completing half-written functions, even when the next step is obvious
- No refactoring "while you're in there"

What is welcome:
- Explanations of Go language features, stdlib behavior, idioms
- Guidance on approach when asked ("how would you structure X?")
- Pointing out a bug or design issue in their code — describe it, don't patch it
- Answering conceptual questions about what the existing code does

When the user asks for help, default to explanation and Socratic guidance over showing code. If they explicitly ask for a code sample, keep it minimal and illustrative — not a drop-in replacement for what they're writing.

## Commands

- Build: `go build`
- Run REPL: `go run .`
- Test: `go test ./...`
- Single test: `go test -run TestCleanInput`

## Architecture

Single `package main` in the repo root. The REPL loop in `main.go` reads stdin via `bufio.Scanner`, normalizes each line with `cleanInput` (in `repl.go`), and dispatches the first token through a `map[string]cliCommand` returned by `getCommands()`.

Adding a command means: define a `func() error` callback, register it in `getCommands()` with a `cliCommand{name, description, callback}` entry. `commandHelp` auto-lists everything in the map, so new commands surface in `help` for free.

Go version pinned to 1.26.1 in `go.mod`. No external dependencies yet.
