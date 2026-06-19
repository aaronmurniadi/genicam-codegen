// Package parser parses GenICam XML (SFNC / GenApi schema) into an intermediate
// representation (IR) that the generator can walk to emit Go source code.
package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ──────────────────────────────────────────────────────────────────────────────
// Raw XML structures – mirrors the GenApi XML schema
// ──────────────────────────────────────────────────────────────────────────────

// RegisterDescription is the root element of a GenICam XML file.
type rawRegisterDescription struct {
	XMLName     xml.Name   `xml:"RegisterDescription"`
	ModelName   string     `xml:"ModelName,attr"`
	VendorName  string     `xml:"VendorName,attr"`
	ToolTip     string     `xml:"ToolTip,attr"`
	SchemaMajor int        `xml:"SchemaMajorVersion,attr"`
	SchemaMinor int        `xml:"SchemaMinorVersion,attr"`
	SchemaPatch int        `xml:"SchemaPatchVersion,attr"`
	Group       []rawGroup `xml:"Group"`
	// Top-level nodes (without a Group wrapper)
	Categories   []rawCategory    `xml:"Category"`
	Integers     []rawInteger     `xml:"Integer"`
	IntRegs      []rawIntReg      `xml:"IntReg"`
	Floats       []rawFloat       `xml:"Float"`
	FloatRegs    []rawFloatReg    `xml:"FloatReg"`
	Booleans     []rawBoolean     `xml:"Boolean"`
	Commands     []rawCommand     `xml:"Command"`
	Enumerations []rawEnumeration `xml:"Enumeration"`
	Strings      []rawString      `xml:"String"`
	StringRegs   []rawStringReg   `xml:"StringReg"`
}

type rawGroup struct {
	Comment      string           `xml:"Comment,attr"`
	Categories   []rawCategory    `xml:"Category"`
	Integers     []rawInteger     `xml:"Integer"`
	IntRegs      []rawIntReg      `xml:"IntReg"`
	Floats       []rawFloat       `xml:"Float"`
	FloatRegs    []rawFloatReg    `xml:"FloatReg"`
	Booleans     []rawBoolean     `xml:"Boolean"`
	Commands     []rawCommand     `xml:"Command"`
	Enumerations []rawEnumeration `xml:"Enumeration"`
	Strings      []rawString      `xml:"String"`
	StringRegs   []rawStringReg   `xml:"StringReg"`
}

type rawCategory struct {
	Name        string   `xml:"Name,attr"`
	NameSpace   string   `xml:"NameSpace,attr"`
	ToolTip     string   `xml:"ToolTip,attr"`
	Description string   `xml:"Description"`
	PFeatures   []string `xml:"pFeature"`
}

type rawInteger struct {
	Name        string `xml:"Name,attr"`
	NameSpace   string `xml:"NameSpace,attr"`
	ToolTip     string `xml:"ToolTip,attr"`
	Description string `xml:"Description"`
	Min         string `xml:"Min"`
	Max         string `xml:"Max"`
	Inc         string `xml:"Inc"`
	Unit        string `xml:"Unit"`
	Value       string `xml:"Value"`
	PValue      string `xml:"pValue"`
	AccessMode  string `xml:"AccessMode"`
	Visibility  string `xml:"Visibility"`
}

type rawIntReg struct {
	rawInteger
}

type rawFloat struct {
	Name        string `xml:"Name,attr"`
	NameSpace   string `xml:"NameSpace,attr"`
	ToolTip     string `xml:"ToolTip,attr"`
	Description string `xml:"Description"`
	Min         string `xml:"Min"`
	Max         string `xml:"Max"`
	Unit        string `xml:"Unit"`
	Value       string `xml:"Value"`
	PValue      string `xml:"pValue"`
	AccessMode  string `xml:"AccessMode"`
	Visibility  string `xml:"Visibility"`
}

type rawFloatReg struct {
	rawFloat
}

type rawBoolean struct {
	Name        string `xml:"Name,attr"`
	NameSpace   string `xml:"NameSpace,attr"`
	ToolTip     string `xml:"ToolTip,attr"`
	Description string `xml:"Description"`
	Value       string `xml:"Value"`
	PValue      string `xml:"pValue"`
	AccessMode  string `xml:"AccessMode"`
	Visibility  string `xml:"Visibility"`
}

type rawCommand struct {
	Name         string `xml:"Name,attr"`
	NameSpace    string `xml:"NameSpace,attr"`
	ToolTip      string `xml:"ToolTip,attr"`
	Description  string `xml:"Description"`
	PValue       string `xml:"pValue"`
	CommandValue string `xml:"CommandValue"`
	AccessMode   string `xml:"AccessMode"`
	Visibility   string `xml:"Visibility"`
}

type rawEnumEntry struct {
	Name         string `xml:"Name,attr"`
	NameSpace    string `xml:"NameSpace,attr"`
	ToolTip      string `xml:"ToolTip,attr"`
	Value        string `xml:"Value"`
	NumericValue string `xml:"NumericValue"`
}

type rawEnumeration struct {
	Name        string         `xml:"Name,attr"`
	NameSpace   string         `xml:"NameSpace,attr"`
	ToolTip     string         `xml:"ToolTip,attr"`
	Description string         `xml:"Description"`
	EnumEntries []rawEnumEntry `xml:"EnumEntry"`
	PValue      string         `xml:"pValue"`
	AccessMode  string         `xml:"AccessMode"`
	Visibility  string         `xml:"Visibility"`
}

type rawString struct {
	Name        string `xml:"Name,attr"`
	NameSpace   string `xml:"NameSpace,attr"`
	ToolTip     string `xml:"ToolTip,attr"`
	Description string `xml:"Description"`
	Value       string `xml:"Value"`
	PValue      string `xml:"pValue"`
	AccessMode  string `xml:"AccessMode"`
	Visibility  string `xml:"Visibility"`
}

type rawStringReg struct {
	rawString
}

// ──────────────────────────────────────────────────────────────────────────────
// Intermediate Representation
// ──────────────────────────────────────────────────────────────────────────────

// NodeKind classifies a GenICam feature node.
type NodeKind string

const (
	KindCategory    NodeKind = "Category"
	KindInteger     NodeKind = "Integer"
	KindFloat       NodeKind = "Float"
	KindBoolean     NodeKind = "Boolean"
	KindCommand     NodeKind = "Command"
	KindEnumeration NodeKind = "Enumeration"
	KindString      NodeKind = "String"
)

// AccessMode mirrors GenICam access modes.
type AccessMode string

const (
	AccessRO AccessMode = "RO"
	AccessWO AccessMode = "WO"
	AccessRW AccessMode = "RW"
	AccessNA AccessMode = "NA"
)

// Visibility mirrors GenICam visibility levels.
type Visibility string

const (
	VisiBeginner  Visibility = "Beginner"
	VisiExpert    Visibility = "Expert"
	VisiGuru      Visibility = "Guru"
	VisiInvisible Visibility = "Invisible"
)

// EnumEntry represents a single enum value.
type EnumEntry struct {
	Name         string
	GoName       string // sanitised Go identifier
	NumericValue int64
	ToolTip      string
}

// Node is the unified IR for any GenICam feature.
type Node struct {
	Name        string
	GoName      string // sanitised Go identifier
	Kind        NodeKind
	ToolTip     string
	Description string
	AccessMode  AccessMode
	Visibility  Visibility

	// Category children
	Children []string

	// Integer / Float range
	Min  string
	Max  string
	Inc  string
	Unit string

	// Enumeration entries
	EnumEntries []EnumEntry

	// For grouping: which category this belongs to (filled in post-process)
	Category string
}

// RegisterDescription is the top-level IR produced by the parser.
type RegisterDescription struct {
	ModelName  string
	VendorName string
	Nodes      map[string]*Node // keyed by GenICam name
	Categories []*Node          // ordered top-level categories
}

// ──────────────────────────────────────────────────────────────────────────────
// Public API
// ──────────────────────────────────────────────────────────────────────────────

// Parse reads a GenICam XML document from r and returns the IR.
func Parse(r io.Reader) (*RegisterDescription, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var raw rawRegisterDescription
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("xml unmarshal: %w", err)
	}

	rd := &RegisterDescription{
		ModelName:  raw.ModelName,
		VendorName: raw.VendorName,
		Nodes:      make(map[string]*Node),
	}

	// Collect nodes from all groups + top-level
	collectors := []func(*RegisterDescription, rawGroup){collectGroup}
	topGroup := rawGroup{
		Categories:   raw.Categories,
		Integers:     raw.Integers,
		IntRegs:      raw.IntRegs,
		Floats:       raw.Floats,
		FloatRegs:    raw.FloatRegs,
		Booleans:     raw.Booleans,
		Commands:     raw.Commands,
		Enumerations: raw.Enumerations,
		Strings:      raw.Strings,
		StringRegs:   raw.StringRegs,
	}
	_ = collectors
	collectGroup(rd, topGroup)
	for _, g := range raw.Group {
		collectGroup(rd, g)
	}

	// Identify top-level categories
	for _, n := range rd.Nodes {
		if n.Kind == KindCategory {
			rd.Categories = append(rd.Categories, n)
		}
	}

	// Mark which category each node belongs to (shallow – first category that references it)
	for _, cat := range rd.Categories {
		markCategory(rd, cat, cat.GoName)
	}

	return rd, nil
}

func markCategory(rd *RegisterDescription, cat *Node, catGoName string) {
	for _, childName := range cat.Children {
		child, ok := rd.Nodes[childName]
		if !ok {
			continue
		}
		if child.Category == "" {
			child.Category = catGoName
		}
		if child.Kind == KindCategory {
			markCategory(rd, child, catGoName)
		}
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Collection helpers
// ──────────────────────────────────────────────────────────────────────────────

func collectGroup(rd *RegisterDescription, g rawGroup) {
	for _, v := range g.Categories {
		n := &Node{
			Name:    v.Name,
			GoName:  toGoIdent(v.Name),
			Kind:    KindCategory,
			ToolTip: v.ToolTip,
		}
		for _, pf := range v.PFeatures {
			n.Children = append(n.Children, strings.TrimSpace(pf))
		}
		rd.Nodes[v.Name] = n
	}
	for _, v := range g.Integers {
		rd.Nodes[v.Name] = integerNode(v)
	}
	for _, v := range g.IntRegs {
		rd.Nodes[v.Name] = integerNode(v.rawInteger)
	}
	for _, v := range g.Floats {
		rd.Nodes[v.Name] = floatNode(v)
	}
	for _, v := range g.FloatRegs {
		rd.Nodes[v.Name] = floatNode(v.rawFloat)
	}
	for _, v := range g.Booleans {
		rd.Nodes[v.Name] = boolNode(v)
	}
	for _, v := range g.Commands {
		rd.Nodes[v.Name] = commandNode(v)
	}
	for _, v := range g.Enumerations {
		rd.Nodes[v.Name] = enumNode(v)
	}
	for _, v := range g.Strings {
		rd.Nodes[v.Name] = stringNode(v)
	}
	for _, v := range g.StringRegs {
		rd.Nodes[v.Name] = stringNode(v.rawString)
	}
}

func integerNode(v rawInteger) *Node {
	return &Node{
		Name:        v.Name,
		GoName:      toGoIdent(v.Name),
		Kind:        KindInteger,
		ToolTip:     v.ToolTip,
		Description: v.Description,
		AccessMode:  parseAccess(v.AccessMode),
		Visibility:  parseVisibility(v.Visibility),
		Min:         v.Min,
		Max:         v.Max,
		Inc:         v.Inc,
		Unit:        v.Unit,
	}
}

func floatNode(v rawFloat) *Node {
	return &Node{
		Name:        v.Name,
		GoName:      toGoIdent(v.Name),
		Kind:        KindFloat,
		ToolTip:     v.ToolTip,
		Description: v.Description,
		AccessMode:  parseAccess(v.AccessMode),
		Visibility:  parseVisibility(v.Visibility),
		Min:         v.Min,
		Max:         v.Max,
		Unit:        v.Unit,
	}
}

func boolNode(v rawBoolean) *Node {
	return &Node{
		Name:        v.Name,
		GoName:      toGoIdent(v.Name),
		Kind:        KindBoolean,
		ToolTip:     v.ToolTip,
		Description: v.Description,
		AccessMode:  parseAccess(v.AccessMode),
		Visibility:  parseVisibility(v.Visibility),
	}
}

func commandNode(v rawCommand) *Node {
	return &Node{
		Name:        v.Name,
		GoName:      toGoIdent(v.Name),
		Kind:        KindCommand,
		ToolTip:     v.ToolTip,
		Description: v.Description,
		AccessMode:  AccessWO,
		Visibility:  parseVisibility(v.Visibility),
	}
}

func enumNode(v rawEnumeration) *Node {
	n := &Node{
		Name:        v.Name,
		GoName:      toGoIdent(v.Name),
		Kind:        KindEnumeration,
		ToolTip:     v.ToolTip,
		Description: v.Description,
		AccessMode:  parseAccess(v.AccessMode),
		Visibility:  parseVisibility(v.Visibility),
	}
	for i, e := range v.EnumEntries {
		val := int64(i)
		if e.Value != "" {
			fmt.Sscanf(e.Value, "%d", &val)
		} else if e.NumericValue != "" {
			fmt.Sscanf(e.NumericValue, "%d", &val)
		}
		n.EnumEntries = append(n.EnumEntries, EnumEntry{
			Name:         e.Name,
			GoName:       toGoIdent(e.Name),
			NumericValue: val,
			ToolTip:      e.ToolTip,
		})
	}
	return n
}

func stringNode(v rawString) *Node {
	return &Node{
		Name:        v.Name,
		GoName:      toGoIdent(v.Name),
		Kind:        KindString,
		ToolTip:     v.ToolTip,
		Description: v.Description,
		AccessMode:  parseAccess(v.AccessMode),
		Visibility:  parseVisibility(v.Visibility),
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Utility
// ──────────────────────────────────────────────────────────────────────────────

func parseAccess(s string) AccessMode {
	switch strings.TrimSpace(strings.ToUpper(s)) {
	case "RO":
		return AccessRO
	case "WO":
		return AccessWO
	case "RW", "":
		return AccessRW
	default:
		return AccessNA
	}
}

func parseVisibility(s string) Visibility {
	switch strings.TrimSpace(s) {
	case "Expert":
		return VisiExpert
	case "Guru":
		return VisiGuru
	case "Invisible":
		return VisiInvisible
	default:
		return VisiBeginner
	}
}

// toGoIdent converts a GenICam name to a valid exported Go identifier.
func toGoIdent(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "_"
	}
	var b strings.Builder
	upperNext := true
	for i, ch := range name {
		switch {
		case ch >= 'A' && ch <= 'Z':
			if i == 0 {
				upperNext = false
			}
			b.WriteRune(ch)
			upperNext = false
		case ch >= 'a' && ch <= 'z':
			if upperNext {
				b.WriteRune(ch - 32)
				upperNext = false
			} else {
				b.WriteRune(ch)
			}
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
			upperNext = false
		default:
			upperNext = true
		}
	}
	id := b.String()
	if id == "" {
		return "_"
	}
	// Ensure starts with a letter
	if id[0] >= '0' && id[0] <= '9' {
		id = "N" + id
	}
	return id
}
