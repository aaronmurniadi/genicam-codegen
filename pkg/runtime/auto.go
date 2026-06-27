package runtime

// GigeConfig holds optional GigE connection settings for ephemeral camera access.
// Zero value: auto-discover the first camera on the network for each call.
type GigeConfig struct {
	// Device is the local NIC name (e.g. "en0"). Empty auto-detects.
	Device string
	// CameraIP is the camera address. Empty broadcasts to discover any camera.
	CameraIP string
}

// WithNodeMap opens a GigE connection, runs fn, then disconnects.
// When CameraIP is empty the first discovered camera is used (arv-tool style).
func (c GigeConfig) WithNodeMap(addrs map[string]string, fn func(NodeMap) error) error {
	cameraIP := c.CameraIP
	if cameraIP == "" {
		cameraIP = "255.255.255.255"
	}
	nm, err := OpenGigeNodeMap(c.Device, cameraIP)
	if err != nil {
		return err
	}
	if len(addrs) > 0 {
		if nm.Addresses == nil {
			nm.Addresses = make(map[string]string, len(addrs))
		}
		for k, v := range addrs {
			nm.Addresses[k] = v
		}
	}
	defer nm.Close()
	return fn(nm)
}
