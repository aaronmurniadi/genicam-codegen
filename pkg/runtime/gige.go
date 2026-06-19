package runtime

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/aaronmurniadi/genicam-codegen/pkg/gige/control"
)

// GigeNodeMap implements NodeMap using the pure-Go GigE Vision control layer
// in pkg/gige/control.
//
// Feature names are resolved against the camera's RegisterBank (Width, Height,
// AcquisitionStart, …). For GenICam features not in that bank, supply hex
// register addresses via Addresses.
//
//	conn := control.GetControlConnection(hostIP)
//	cam := control.NewCamera(hostIP, cameraIP, control.STREAM_PORT, modelName)
//	nm := runtime.NewGigeNodeMap(cam, conn)
type GigeNodeMap struct {
	mu sync.Mutex
	cam control.Camera
	conn *net.UDPConn
	// Addresses maps GenICam feature names to 8-digit hex register addresses.
	Addresses map[string]string
}

// NewGigeNodeMap wraps an open gige control connection and camera descriptor.
func NewGigeNodeMap(cam control.Camera, conn *net.UDPConn) *GigeNodeMap {
	return &GigeNodeMap{
		cam:  cam,
		conn: conn,
	}
}

// OpenGigeNodeMap discovers a camera and returns a ready NodeMap.
// device is the local NIC (e.g. "en0"); pass "" to auto-detect.
// cameraIP is the camera address or "255.255.255.255" to broadcast-discover.
func OpenGigeNodeMap(device, cameraIP string) (*GigeNodeMap, error) {
	if device == "" {
		device = control.DetectInterface(cameraIP)
		if device == "" {
			return nil, fmt.Errorf("gige: unable to detect network interface")
		}
	}
	hostIP := control.GetIP(device)
	conn := control.GetControlConnection(hostIP)
	cameraMap, err := control.GetCameraMap(conn, &net.IPAddr{IP: net.ParseIP(cameraIP)})
	if err != nil {
		return nil, fmt.Errorf("gige: discover cameras: %w", err)
	}
	if len(cameraMap) == 0 {
		return nil, fmt.Errorf("gige: no cameras found at %s", cameraIP)
	}
	for _, cam := range cameraMap {
		return NewGigeNodeMap(cam, conn), nil
	}
	return nil, fmt.Errorf("gige: no cameras found")
}

func (g *GigeNodeMap) registerName(feature string) (string, error) {
	if addr, ok := g.Addresses[feature]; ok {
		return fmt.Sprintf("R[%s]", addr), nil
	}
	if g.cam.RegisterBank.GetAddress(feature) != "" {
		return feature, nil
	}
	return "", fmt.Errorf("gige: unknown feature %q", feature)
}

func (g *GigeNodeMap) readUint32(feature string) (uint64, error) {
	name, err := g.registerName(feature)
	if err != nil {
		return 0, err
	}
	hex, err := g.cam.GetCameraRegister(name, g.conn)
	if err != nil {
		return 0, fmt.Errorf("gige: read %s: %w", feature, err)
	}
	if hex == "9999" {
		return 0, fmt.Errorf("gige: no response reading %s", feature)
	}
	return control.Hex2Num(hex), nil
}

func (g *GigeNodeMap) writeRegister(feature string, value string) error {
	name, err := g.registerName(feature)
	if err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := g.cam.SendRegisterValues(g.conn, g.cam.CameraIp, name+"="+value); err != nil {
		return fmt.Errorf("gige: write %s: %w", feature, err)
	}
	return nil
}

func (g *GigeNodeMap) ExecuteCommand(feature string) error {
	return g.writeRegister(feature, "1")
}

func (g *GigeNodeMap) GetInteger(feature string) (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, err := g.readUint32(feature)
	return int64(v), err
}

func (g *GigeNodeMap) SetInteger(feature string, value int64) error {
	return g.writeRegister(feature, strconv.FormatInt(value, 10))
}

func (g *GigeNodeMap) GetFloat(feature string) (float64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, err := g.readUint32(feature)
	if err != nil {
		return 0, err
	}
	return float64(math.Float32frombits(uint32(v))), nil
}

func (g *GigeNodeMap) SetFloat(feature string, value float64) error {
	bits := math.Float32bits(float32(value))
	return g.writeRegister(feature, strconv.FormatUint(uint64(bits), 10))
}

func (g *GigeNodeMap) GetBoolean(feature string) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, err := g.readUint32(feature)
	return v != 0, err
}

func (g *GigeNodeMap) SetBoolean(feature string, value bool) error {
	if value {
		return g.SetInteger(feature, 1)
	}
	return g.SetInteger(feature, 0)
}

func (g *GigeNodeMap) GetEnumeration(feature string) (int64, error) {
	return g.GetInteger(feature)
}

func (g *GigeNodeMap) SetEnumeration(feature string, value int64) error {
	return g.SetInteger(feature, value)
}

func (g *GigeNodeMap) GetString(feature string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	v, err := g.readUint32(feature)
	if err != nil {
		return "", err
	}
	b := []byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
	return strings.TrimRight(string(b), "\x00"), nil
}

func (g *GigeNodeMap) SetString(feature string, value string) error {
	if net.ParseIP(value) != nil {
		return g.writeRegister(feature, value)
	}
	return fmt.Errorf("gige: SetString only supports IPv4 values for %q", feature)
}
