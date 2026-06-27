package parser_test

import (
	"os"
	"path/filepath"
	"strings"
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

func TestParsePValueRegisterAddr(t *testing.T) {
	const xmlDoc = `<?xml version="1.0" encoding="utf-8"?>
<RegisterDescription ModelName="T" VendorName="T" xmlns="http://www.genicam.org/GenApi/Version_1_1">
	<Integer Name="TCPPort"><pValue>TCPPortReg</pValue></Integer>
	<IntReg Name="TCPPortReg"><Address>0x4E05C740</Address></IntReg>
</RegisterDescription>`
	rd, err := parser.Parse(strings.NewReader(xmlDoc))
	if err != nil {
		t.Fatal(err)
	}
	n := rd.Nodes["TCPPort"]
	if n == nil {
		t.Fatal("missing TCPPort")
	}
	if n.RegisterAddr != "4e05c740" {
		t.Fatalf("RegisterAddr = %q, want 4e05c740", n.RegisterAddr)
	}
}
