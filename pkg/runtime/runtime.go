package runtime

// Package runtime defines the NodeMap interface that generated code relies on,
// and ships GigE Vision transport with arv-tool-style auto-discovery:
//
//   - GigeConfig.WithNodeMap – discover, connect, call, disconnect per operation
//   - GigeNodeMap          – persistent GigE Vision connection
//   - MockNodeMap          – in-memory implementation for unit tests
//
// Generated code calls WithNodeMap automatically when no device IP is configured.

import (
	"sync"
)

// ──────────────────────────────────────────────────────────────────────────────
// NodeMap interface – the only surface the generated code touches
// ──────────────────────────────────────────────────────────────────────────────

// NodeMap abstracts read/write access to a camera's GenICam feature tree.
// Any transport (GigE Vision, Baumer, Basler Pylon, …) can be wrapped behind this
// interface so the generated code stays portable.
type NodeMap interface {
	// Commands
	ExecuteCommand(feature string) error

	// Integer features
	GetInteger(feature string) (int64, error)
	SetInteger(feature string, value int64) error

	// Float features
	GetFloat(feature string) (float64, error)
	SetFloat(feature string, value float64) error

	// Boolean features
	GetBoolean(feature string) (bool, error)
	SetBoolean(feature string, value bool) error

	// Enumeration features (stored as int64 ordinal)
	GetEnumeration(feature string) (int64, error)
	SetEnumeration(feature string, value int64) error

	// String features
	GetString(feature string) (string, error)
	SetString(feature string, value string) error
}

// ──────────────────────────────────────────────────────────────────────────────
// MockNodeMap – deterministic in-memory implementation for testing/simulation
// ──────────────────────────────────────────────────────────────────────────────

// MockNodeMap stores feature values in memory.  It is safe for concurrent use.
type MockNodeMap struct {
	mu       sync.RWMutex
	integers map[string]int64
	floats   map[string]float64
	booleans map[string]bool
	enums    map[string]int64
	strings  map[string]string
	commands []string // log of executed commands
	// Hooks – set these to intercept feature access in tests
	OnExecuteCommand func(feature string) error
	OnSetInteger     func(feature string, value int64) error
	OnSetFloat       func(feature string, value float64) error
}

// NewMockNodeMap creates a ready-to-use MockNodeMap.
func NewMockNodeMap() *MockNodeMap {
	return &MockNodeMap{
		integers: make(map[string]int64),
		floats:   make(map[string]float64),
		booleans: make(map[string]bool),
		enums:    make(map[string]int64),
		strings:  make(map[string]string),
	}
}

func (m *MockNodeMap) ExecuteCommand(feature string) error {
	if m.OnExecuteCommand != nil {
		return m.OnExecuteCommand(feature)
	}
	m.mu.Lock()
	m.commands = append(m.commands, feature)
	m.mu.Unlock()
	return nil
}

func (m *MockNodeMap) GetInteger(feature string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.integers[feature], nil
}

func (m *MockNodeMap) SetInteger(feature string, value int64) error {
	if m.OnSetInteger != nil {
		return m.OnSetInteger(feature, value)
	}
	m.mu.Lock()
	m.integers[feature] = value
	m.mu.Unlock()
	return nil
}

func (m *MockNodeMap) GetFloat(feature string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.floats[feature], nil
}

func (m *MockNodeMap) SetFloat(feature string, value float64) error {
	if m.OnSetFloat != nil {
		return m.OnSetFloat(feature, value)
	}
	m.mu.Lock()
	m.floats[feature] = value
	m.mu.Unlock()
	return nil
}

func (m *MockNodeMap) GetBoolean(feature string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.booleans[feature], nil
}

func (m *MockNodeMap) SetBoolean(feature string, value bool) error {
	m.mu.Lock()
	m.booleans[feature] = value
	m.mu.Unlock()
	return nil
}

func (m *MockNodeMap) GetEnumeration(feature string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enums[feature], nil
}

func (m *MockNodeMap) SetEnumeration(feature string, value int64) error {
	m.mu.Lock()
	m.enums[feature] = value
	m.mu.Unlock()
	return nil
}

func (m *MockNodeMap) GetString(feature string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.strings[feature], nil
}

func (m *MockNodeMap) SetString(feature string, value string) error {
	m.mu.Lock()
	m.strings[feature] = value
	m.mu.Unlock()
	return nil
}

// ExecutedCommands returns the list of commands that have been executed.
func (m *MockNodeMap) ExecutedCommands() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, len(m.commands))
	copy(out, m.commands)
	return out
}

// Seed pre-populates mock values for testing.
func (m *MockNodeMap) Seed(feature string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch v := value.(type) {
	case int64:
		m.integers[feature] = v
	case int:
		m.integers[feature] = int64(v)
	case float64:
		m.floats[feature] = v
	case bool:
		m.booleans[feature] = v
	case string:
		m.strings[feature] = v
	}
}
