package main

import (
	"os"

	"github.com/aaronmurniadi/genicam-codegen/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
