// genicam-codegen generates type-safe Go code from a GenICam XML file.
//
// Usage:
//
//	genicam-codegen -xml camera.xml -out ./genicam -pkg genicam
//
// Flags:
//
//	-xml        Path to the GenICam XML file (required)
//	-out        Output directory for generated files (default: "./genicam")
//	-pkg        Go package name (default: "genicam")
//	-runtime    Import path of the runtime package
//	            (default: "github.com/genicam-codegen/pkg/runtime")
//	-visibility Minimum visibility to emit: Beginner|Expert|Guru|All
//	            (default: Beginner)
//	-v          Verbose output
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/genicam-codegen/pkg/generator"
	"github.com/genicam-codegen/pkg/parser"
)

func main() {
	var (
		xmlPath    = flag.String("xml", "", "Path to GenICam XML file (required)")
		outDir     = flag.String("out", "./genicam", "Output directory")
		pkg        = flag.String("pkg", "genicam", "Go package name")
		runtime    = flag.String("runtime", "github.com/genicam-codegen/pkg/runtime", "Runtime import path")
		visibility = flag.String("visibility", "Beginner", "Minimum visibility: Beginner|Expert|Guru|All")
		verbose    = flag.Bool("v", false, "Verbose output")
	)
	flag.Parse()

	if *xmlPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// ── Parse ──────────────────────────────────────────────────────────────
	f, err := os.Open(*xmlPath)
	if err != nil {
		log.Fatalf("open xml: %v", err)
	}
	defer f.Close()

	if *verbose {
		log.Printf("Parsing %s …", *xmlPath)
	}

	rd, err := parser.Parse(f)
	if err != nil {
		log.Fatalf("parse xml: %v", err)
	}

	if *verbose {
		log.Printf("Model: %s %s", rd.VendorName, rd.ModelName)
		log.Printf("Nodes: %d", len(rd.Nodes))
		log.Printf("Categories: %d", len(rd.Categories))
	}

	// ── Generate ───────────────────────────────────────────────────────────
	opts := generator.Options{
		PackageName:   *pkg,
		RuntimeImport: *runtime,
		Visibility:    *visibility,
	}

	files, err := generator.Generate(rd, opts)
	if err != nil {
		log.Fatalf("generate: %v", err)
	}

	// ── Write ──────────────────────────────────────────────────────────────
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	for name, src := range files {
		path := filepath.Join(*outDir, name)
		if err := os.WriteFile(path, src, 0o644); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}
		if *verbose {
			log.Printf("  wrote %s (%d bytes)", path, len(src))
		}
	}

	fmt.Printf("Generated %d file(s) in %s\n", len(files), *outDir)
	fmt.Printf("  Package : %s\n", *pkg)
	fmt.Printf("  Model   : %s %s\n", rd.VendorName, rd.ModelName)
	fmt.Printf("  Features: %d\n", len(rd.Nodes))
}
