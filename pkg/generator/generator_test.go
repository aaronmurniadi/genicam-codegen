package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aaronmurniadi/genicam-codegen/pkg/generator"
	"github.com/aaronmurniadi/genicam-codegen/pkg/parser"
)

func TestNormalizeVisibility(t *testing.T) {
	tests := []struct {
		in   string
		want string
		err  bool
	}{
		{"", "Beginner", false},
		{"beginner", "Beginner", false},
		{"Expert", "Expert", false},
		{"guru", "Guru", false},
		{"invalid", "", true},
	}
	for _, tc := range tests {
		got, err := generator.NormalizeVisibility(tc.in)
		if tc.err {
			if err == nil {
				t.Fatalf("NormalizeVisibility(%q) expected error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("NormalizeVisibility(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("NormalizeVisibility(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestGenerateVisibilityFilter(t *testing.T) {
	path := filepath.Join("..", "parser", "testdata", "minimal.xml")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rd, err := parser.Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	beginner, err := generator.Generate(rd, generator.Options{
		PackageName: "cam",
		Visibility:  "beginner",
	})
	if err != nil {
		t.Fatal(err)
	}
	src := string(beginner["genicam.go"])
	if !strings.Contains(src, "GetWidth") {
		t.Fatal("beginner output missing GetWidth")
	}
	if strings.Contains(src, "GetGuruFeature") {
		t.Fatal("beginner output should not include guru feature")
	}

	guru, err := generator.Generate(rd, generator.Options{
		PackageName: "cam",
		Visibility:  "guru",
	})
	if err != nil {
		t.Fatal(err)
	}
	src = string(guru["genicam.go"])
	if !strings.Contains(src, "GetGuruFeature") {
		t.Fatal("guru output missing GetGuruFeature")
	}
}

func TestGenerateFormatsOutput(t *testing.T) {
	path := filepath.Join("..", "parser", "testdata", "minimal.xml")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rd, err := parser.Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	files, err := generator.Generate(rd, generator.Options{PackageName: "cam"})
	if err != nil {
		t.Fatal(err)
	}
	for name, src := range files {
		if !strings.HasPrefix(string(src), "// Code generated") {
			t.Fatalf("%s missing generated header", name)
		}
		if !strings.Contains(string(src), "package cam") {
			t.Fatalf("%s missing package declaration", name)
		}
	}
}
