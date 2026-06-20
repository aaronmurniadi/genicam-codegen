package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aaronmurniadi/genicam-codegen/pkg/parser"
)

func TestParseMinimal(t *testing.T) {
	path := filepath.Join("testdata", "minimal.xml")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rd, err := parser.Parse(f)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if rd.ModelName != "TestCam" {
		t.Fatalf("ModelName = %q, want TestCam", rd.ModelName)
	}
	if got := len(rd.Nodes); got != 4 { // 3 features + 1 category
		t.Fatalf("len(Nodes) = %d, want 4", got)
	}
	if rd.Nodes["Width"] == nil {
		t.Fatal("missing Width node")
	}
	if rd.Nodes["Width"].Visibility != parser.VisiBeginner {
		t.Fatalf("Width visibility = %q, want Beginner", rd.Nodes["Width"].Visibility)
	}
}
