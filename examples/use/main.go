// Command use demonstrates calling generated camera bindings with MockNodeMap.
//
// Generate bindings first:
//
//	go run ./examples/generate
//
// Then run this example:
//
//	go run ./examples/use
package main

import (
	"fmt"
	"log"

	camera "github.com/aaronmurniadi/genicam-codegen/examples/generated"
	"github.com/aaronmurniadi/genicam-codegen/pkg/runtime"
)

func main() {
	mock := runtime.NewMockNodeMap()
	mock.Seed("AcquisitionFrameCount", int64(42))

	cam := camera.New(mock)

	if err := cam.AcquisitionControl.AcquisitionStart(); err != nil {
		log.Fatalf("AcquisitionStart: %v", err)
	}

	count, err := cam.AcquisitionControl.GetAcquisitionFrameCount()
	if err != nil {
		log.Fatalf("GetAcquisitionFrameCount: %v", err)
	}

	fmt.Printf("AcquisitionFrameCount = %d\n", count)
	fmt.Printf("Commands executed: %v\n", mock.ExecutedCommands())
}
