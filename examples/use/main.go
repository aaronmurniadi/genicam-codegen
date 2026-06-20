// Command use demonstrates calling generated camera bindings over GigE Vision.
//
// Generate bindings first:
//
//	go run ./examples/generate
//
// Then run against a connected camera:
//
//	go run ./examples/use -a 192.168.1.108
//	go run ./examples/use -d en0 -a 192.168.1.108
package main

import (
	"flag"
	"fmt"
	"log"

	camera "github.com/aaronmurniadi/genicam-codegen/examples/generated"
	"github.com/aaronmurniadi/genicam-codegen/pkg/runtime"
)

func main() {
	device := flag.String("d", "", "local network interface (e.g. en0); auto-detect if empty")
	cameraIP := flag.String("a", "192.168.1.108", "camera IP address")
	flag.Parse()

	nm, err := runtime.OpenGigeNodeMap(*device, *cameraIP)
	if err != nil {
		log.Fatalf("open camera: %v", err)
	}

	cam := camera.New(nm)

	width, err := cam.MonoImageFormatControl.GetSensorWidth()
	if err != nil {
		log.Fatalf("GetSensorWidth: %v", err)
	}
	height, err := cam.MonoImageFormatControl.GetSensorHeight()
	if err != nil {
		log.Fatalf("GetSensorHeight: %v", err)
	}

	fmt.Printf("SensorWidth  = %d\n", width)
	fmt.Printf("SensorHeight = %d\n", height)
}
