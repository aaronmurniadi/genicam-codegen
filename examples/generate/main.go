// Command generate writes type-safe Go camera bindings from the example GenICam XML.
//
// Run from the repository root:
//
//	go run ./examples/generate
//
// Or from this directory:
//
//	go run .
//
//go:generate go run . -o ../generated
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aaronmurniadi/genicam-codegen/pkg/generator"
	"github.com/aaronmurniadi/genicam-codegen/pkg/parser"
)

func defaultXMLPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "genicam.xml"
	}
	// examples/generate/main.go → examples/genicam.xml
	return filepath.Join(filepath.Dir(file), "..", "genicam.xml")
}

func defaultOutDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "examples/generated"
	}
	return filepath.Join(filepath.Dir(file), "..", "generated")
}

func main() {
	xmlPath := flag.String("i", defaultXMLPath(), "path to GenICam XML")
	outDir := flag.String("o", defaultOutDir(), "output directory for generated Go files")
	pkg := flag.String("pkg", "camera", "Go package name")
	runtimeImport := flag.String("runtime", "github.com/aaronmurniadi/genicam-codegen/pkg/runtime", "runtime import path")
	visibility := flag.String("visibility", "Beginner", "minimum visibility: Beginner|Expert|Guru|All")
	verbose := flag.Bool("v", false, "verbose output")
	flag.Parse()

	f, err := os.Open(*xmlPath)
	if err != nil {
		log.Fatalf("open xml: %v", err)
	}
	defer f.Close()

	if *verbose {
		log.Printf("parsing %s …", *xmlPath)
	}

	rd, err := parser.Parse(f)
	if err != nil {
		log.Fatalf("parse xml: %v", err)
	}

	opts := generator.Options{
		PackageName:   *pkg,
		RuntimeImport: *runtimeImport,
		Visibility:    *visibility,
	}

	files, err := generator.Generate(rd, opts)
	if err != nil {
		log.Fatalf("generate: %v", err)
	}

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
