# genicam-codegen

Generate type-safe Go bindings from GenICam XML camera descriptions.

## Usage

```bash
go run . -xml camera.xml -out ./genicam -pkg genicam
```

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
