# genicam-codegen

Generate type-safe Go bindings from GenICam XML camera descriptions.

## Install

```bash
go install github.com/aaronmurniadi/genicam-codegen@latest
```

This installs the `genicam-codegen` binary to your `$GOPATH/bin` (or `$HOME/go/bin`).

## Usage

```bash
genicam-codegen -i genicam.xml -o ./genicam
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-i` | Path to GenICam XML file (required) | |
| `-o` | Output directory | `./genicam` |
| `-pkg` | Go package name | `genicam` |
| `-runtime` | Runtime import path | module `pkg/runtime` |
| `-visibility` | Minimum feature visibility: `beginner`, `expert`, `guru` | `beginner` |
| `-v` | Verbose output | |

See [examples/README.md](examples/README.md) for a full generate-and-use walkthrough.

## Packages

| Path | Description |
|------|-------------|
| `main.go` | CLI entry point |
| `pkg/parser` | GenICam XML parser |
| `pkg/generator` | Go code generator |
| `pkg/runtime` | `NodeMap` interface, mock and GigE implementations |
| `pkg/gige/control` | Pure-Go GigE Vision control transport (vendored) |

## Third-party attribution

The GigE Vision control code in `pkg/gige/control/` is derived from [dougwatson/gige](https://github.com/dougwatson/gige) (commit `0eacca41f48a`, October 2022).

Copyright (c) 2022 Doug Watson. Licensed under the [MIT License](pkg/gige/LICENSE).

Original repository: https://github.com/dougwatson/gige
