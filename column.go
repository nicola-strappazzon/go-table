package table

import (
	"math"

	"github.com/fatih/color"
)

// Format defines how a column's value is presented when printed.
type Format int

const (
	Plain      Format = iota // the value as is
	Percentage               // appends "%" to the value
	Duration                 // reads the value as seconds and humanizes it
	Bytes                    // reads the value as bytes and humanizes it (GiB, ...)
)

// ColorRule paints a cell value with Color when the column value satisfies
// Condition (e.g. ">= 50"). Rules are evaluated in order and the first match
// wins.
type ColorRule struct {
	Condition string
	Color     color.Attribute
}

type Column struct {
	Alignment Alignment
	Color     color.Attribute // default color when no rule matches
	Colors    []ColorRule     // conditional rules, evaluated in order
	Format    Format          // value presentation format
	Name      string
	Precision int
	Scale     int
	Truncate  int
	Width     int
	ZeroFill  bool
}

func (c *Column) SetName(s string) {
	(*c).Name = s
}

func (c Column) GetWidth() int {
	return int(math.Max(float64(len(c.Name)), float64(c.Width)))
}

// toField builds a Field with value and the column's presentation, so the
// table can reuse Field's format and color logic when rendering.
func (c Column) toField(value any) Field {
	return Field{
		Value:     value,
		Name:      c.Name,
		Format:    c.Format,
		Color:     c.Color,
		Colors:    c.Colors,
		Precision: c.Precision,
		Scale:     c.Scale,
		ZeroFill:  c.ZeroFill,
		Alignment: c.Alignment,
	}
}
