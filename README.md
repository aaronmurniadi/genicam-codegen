# genicam-codegen

Generate type-safe Go bindings from GenICam XML camera descriptions.

## Install

```bash
go install github.com/aaronmurniadi/genicam-codegen@latest
```

This installs the `genicam-codegen` binary to `$GOPATH/bin` (or `$HOME/go/bin`).

## Usage

```bash
genicam-codegen -i genicam.xml -o ./genicam
genicam-codegen -i genicam.xml -o ./genicam -visibility guru
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

## Development

```bash
go test ./...
go run . -i examples/genicam.xml -o examples/generated -pkg camera
```

See [examples/README.md](examples/README.md) for a full generate-and-use walkthrough.

## Layout

| Path | Description |
|------|-------------|
| `cmd/genicam-codegen` | CLI entry point |
| `main.go` | CLI entry point (also `cmd/genicam-codegen`)
| `pkg/parser` | GenICam XML parser |
| `pkg/generator` | Go code generator |
| `pkg/runtime` | `NodeMap` interface, mock and GigE implementations |
| `pkg/gige/control` | Pure-Go GigE Vision control transport (vendored) |
| `examples/` | Sample XML, generate/use demos |

## Third-party attribution

The GigE Vision control code in `pkg/gige/control/` is derived from [dougwatson/gige](https://github.com/dougwatson/gige) (commit `0eacca41f48a`, October 2022).

Copyright (c) 2022 Doug Watson. Licensed under the [MIT License](pkg/gige/LICENSE).

Original repository: https://github.com/dougwatson/gige

## Changelog

### v0.0.5

- Move CLI to `cmd/genicam-codegen` (standard Go project layout)
- Add parser and generator tests with `testdata`
- Implement visibility filtering in the generator (`beginner`, `expert`, `guru`)
- Remove debug stdout from `MockNodeMap`
- Fix godoc comments on exported packages and APIs
- Expand `.gitignore` and remove committed binary artifact
- Update install path: `go install github.com/aaronmurniadi/genicam-codegen@latest`

### v0.0.4

- Add `-visibility` flag to filter generated features by GenICam visibility level
- Switch `examples/use` to a real GigE camera via `GigeNodeMap`
- Normalize visibility flag values to lowercase

### v0.0.3

- Set `go` version to 1.22 for pkg.go.dev indexing compatibility

### v0.0.2

- Add MIT license

### v0.0.1

- Rename module to `github.com/aaronmurniadi/genicam-codegen`
- Change CLI flags to `-i` (input XML) and `-o` (output directory)
- Vendor GigE Vision control code from [dougwatson/gige](https://github.com/dougwatson/gige) into `pkg/gige/control`
- Add `GigeNodeMap` runtime and code generation examples
