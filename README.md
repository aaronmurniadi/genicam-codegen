# genicam-codegen

Generate type-safe Go bindings from GenICam XML camera descriptions.

## Install

```bash
go install github.com/aaronmurniadi/genicam-codegen@latest
```

This installs the `genicam-codegen` binary to `$GOPATH/bin` (or `$HOME/go/bin`).

## Usage

```bash
# genicam.xml in the current directory (default input)
genicam-codegen -o ./camera -pkg camera

# explicit XML path
genicam-codegen -i genicam.xml -o ./camera -pkg camera
```

Generates a single Go file: `{output}/{pkg}.go` (e.g. `./camera/camera.go`).

```go
import cam "your/module/camera"

// Auto-discover first GigE camera on each call (arv-tool style)
cam := cam.New()
port, err := cam.EthernetTransferCtl.GetTCPPort()

// Or target a specific camera
cam = cam.NewWithIP("en0", "192.168.1.108")
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-i` | Path to GenICam XML file | `genicam.xml` |
| `-o` | Output directory or `.go` file path | `./genicam` |
| `-pkg` | Go package name | `genicam` |
| `-runtime` | Runtime import path | module `pkg/runtime` |
| `-visibility` | Minimum feature visibility: `beginner`, `expert`, `guru` | `beginner` |
| `-v` | Verbose output | |


## Extract GenICam XML from a camera

Use [Aravis](https://github.com/AravisProject/aravis) `arv-tool-0.8` to dump the camera's GenICam description before generating bindings.

### Install Aravis

```bash
# macOS
brew install aravis

# Ubuntu / Debian
sudo apt-get install aravis-tools
```

### Discover the camera

Connect the camera to your network, then list devices:

```bash
arv-tool-0.8
```
### Dump the XML

By IP address:

```bash
arv-tool-0.8 -a 192.168.1.108 genicam > genicam.xml
```

Or if you only have one camera connected:

```bash
arv-tool-0.8 genicam > genicam.xml
```

### Generate bindings

```bash
genicam-codegen -i genicam.xml -o ./genicam
```

See [examples/README.md](examples/README.md) for a full generate-and-use walkthrough.

## Layout

| Path | Description |
|------|-------------|
| `main.go` | CLI entry point |
| `cmd/genicam-codegen` | Alternate CLI entry point (same binary) |
| `internal/cli` | CLI implementation |
| `pkg/parser` | GenICam XML parser |
| `pkg/generator` | Go code generator |
| `pkg/runtime` | `NodeMap` interface, mock and GigE implementations |
| `pkg/gige/control` | Pure-Go GigE Vision control transport (vendored) |
| `examples/` | Sample XML, generate/use demos |

## Third-party attribution

The GigE Vision control code in `pkg/gige/control/` is derived from [dougwatson/gige](https://github.com/dougwatson/gige) (commit `0eacca41f48a`, October 2022).

Copyright (c) 2022 Doug Watson. Licensed under the [MIT License](pkg/gige/LICENSE).

Original repository: https://github.com/dougwatson/gige

## Publish

Tagged releases are indexed on [pkg.go.dev](https://pkg.go.dev/github.com/aaronmurniadi/genicam-codegen).

```bash
git tag v0.0.8
git push origin v0.0.8
```

Install a specific version:

```bash
go install github.com/aaronmurniadi/genicam-codegen@v0.0.8
```

## Changelog

### v0.0.8

- Generated methods auto-discover GigE cameras per call (arv-tool style)
- Resolve `pValue` register addresses from XML into generated `featureAddresses` map
- Add `New()`, `NewWithIP()`, and `NewWithNodeMap()` constructors
- Add `GigeConfig.WithNodeMap` for ephemeral connect/disconnect in runtime

### v0.0.7

- Emit a single importable Go file (`{pkg}.go`) instead of multiple files
- Default `-i` to `genicam.xml` in the current directory
- Wire `Uncategorized` features onto `Device`
- Honor `-runtime` import path in generated code

### v0.0.6

- Restore root `main.go` so `go install github.com/aaronmurniadi/genicam-codegen@latest` works
- Extract CLI logic into `internal/cli` (shared by root and `cmd/genicam-codegen`)

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
