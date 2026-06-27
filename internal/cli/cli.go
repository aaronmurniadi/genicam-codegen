// Package cli implements the genicam-codegen command-line interface.
package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aaronmurniadi/genicam-codegen/pkg/generator"
	"github.com/aaronmurniadi/genicam-codegen/pkg/parser"
)

// Run executes the genicam-codegen CLI and returns an exit code.
func Run(args []string) int {
	log.SetFlags(0)

	fs := flag.NewFlagSet("genicam-codegen", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: genicam-codegen -i genicam.xml [-o output]\n\n")
		fs.PrintDefaults()
	}

	xmlPath := fs.String("i", "genicam.xml", "path to GenICam XML file")
	out := fs.String("o", "./genicam", "output directory or .go file path")
	pkg := fs.String("pkg", "genicam", "Go package name for generated code")
	runtimePath := fs.String("runtime", "github.com/aaronmurniadi/genicam-codegen/pkg/runtime", "import path of the runtime package")
	visibility := fs.String("visibility", "beginner", "minimum feature visibility: beginner, expert, guru")
	verbose := fs.Bool("v", false, "verbose output")

	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}

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

	if *verbose {
		log.Printf("model: %s %s", rd.VendorName, rd.ModelName)
		log.Printf("nodes: %d", len(rd.Nodes))
	}

	vis, err := generator.NormalizeVisibility(*visibility)
	if err != nil {
		log.Fatalf("visibility: %v", err)
	}

	opts := generator.Options{
		PackageName:   *pkg,
		RuntimeImport: *runtimePath,
		Visibility:    vis,
	}

	files, err := generator.Generate(rd, opts)
	if err != nil {
		log.Fatalf("generate: %v", err)
	}

	dir, outPath := resolveOutput(*out, *pkg)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	for _, src := range files {
		if err := os.WriteFile(outPath, src, 0o644); err != nil {
			log.Fatalf("write %s: %v", outPath, err)
		}
		if *verbose {
			log.Printf("wrote %s (%d bytes)", outPath, len(src))
		}
	}

	fmt.Printf("Generated %s\n", outPath)
	fmt.Printf("  Package    : %s\n", *pkg)
	fmt.Printf("  Model      : %s %s\n", rd.VendorName, rd.ModelName)
	fmt.Printf("  Visibility : %s\n", vis)
	return 0
}
