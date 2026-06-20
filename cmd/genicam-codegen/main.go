// Command genicam-codegen generates type-safe Go code from a GenICam XML file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aaronmurniadi/genicam-codegen/pkg/generator"
	"github.com/aaronmurniadi/genicam-codegen/pkg/parser"
)

func main() {
	os.Exit(run())
}

func run() int {
	log.SetFlags(0)

	fs := flag.NewFlagSet("genicam-codegen", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: genicam-codegen -i <genicam.xml> -o <output_dir>\n\n")
		fs.PrintDefaults()
	}

	xmlPath := fs.String("i", "", "path to GenICam XML file (required)")
	outDir := fs.String("o", "./genicam", "output directory")
	pkg := fs.String("pkg", "genicam", "Go package name for generated code")
	runtimePath := fs.String("runtime", "github.com/aaronmurniadi/genicam-codegen/pkg/runtime", "import path of the runtime package")
	visibility := fs.String("visibility", "beginner", "minimum feature visibility: beginner, expert, guru")
	verbose := fs.Bool("v", false, "verbose output")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return 2
	}

	if *xmlPath == "" {
		fs.Usage()
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
		log.Printf("categories: %d", len(rd.Categories))
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

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	for name, src := range files {
		path := filepath.Join(*outDir, name)
		if err := os.WriteFile(path, src, 0o644); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}
		if *verbose {
			log.Printf("wrote %s (%d bytes)", path, len(src))
		}
	}

	fmt.Printf("Generated %d file(s) in %s\n", len(files), *outDir)
	fmt.Printf("  Package    : %s\n", *pkg)
	fmt.Printf("  Model      : %s %s\n", rd.VendorName, rd.ModelName)
	fmt.Printf("  Visibility : %s\n", vis)
	return 0
}
