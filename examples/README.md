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

After generating, run the GigE camera example (camera must be reachable on the network):

```bash
go run ./examples/use -a 192.168.1.108
go run ./examples/use -d en0 -a 192.168.1.108
```

Expected output (values depend on your camera):

```
Width  = 2448
Height = 2048
```

Features not in the built-in GigE register bank need hex addresses on `GigeNodeMap.Addresses`. See `pkg/gige/control` for transport details.
