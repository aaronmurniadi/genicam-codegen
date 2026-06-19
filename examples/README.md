# Examples

Generate type-safe Go bindings from the sample GenICam XML in this directory (`genicam.xml`).

## Generate bindings

From the repository root:

```bash
go run ./examples/generate
```

Or after installing the CLI:

```bash
genicam-codegen -i examples/genicam.xml -o examples/generated -pkg camera
```

This writes three files into `examples/generated/`:

- `doc.go` — package documentation
- `enums.go` — enumeration types and constants
- `genicam.go` — category structs and feature methods

### Options

```bash
go run ./examples/generate -visibility Expert -v
```

Or regenerate via `go generate`:

```bash
go generate ./examples/generate
```

## Use generated bindings

After generating, run the mock-camera demo:

```bash
go run ./examples/use
```

Expected output:

```
AcquisitionFrameCount = 42
Commands executed: [AcquisitionStart]
```

## Connect to a real camera

Swap `MockNodeMap` for `GigeNodeMap` when talking to a GigE Vision camera:

```go
nm, err := runtime.OpenGigeNodeMap("en0", "192.168.1.108")
if err != nil {
    log.Fatal(err)
}
cam := camera.New(nm)
```

See `pkg/gige/control` for GigE transport details. Upstream source: [dougwatson/gige](https://github.com/dougwatson/gige).
