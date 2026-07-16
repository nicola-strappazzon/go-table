package main

import (
	"math"

	"github.com/fatih/color"
)

// Format define cómo se presenta el valor de una columna al imprimir.
type Format int

const (
	Plain      Format = iota // el valor tal cual
	Percentage               // añade "%" al valor
	Duration                 // interpreta el valor como segundos y lo humaniza
	Bytes                    // interpreta el valor como bytes y lo humaniza (GiB, ...)
)

// ColorRule pinta el valor de una celda con Color cuando el valor de la
// columna cumple Condition (p. ej. ">= 50"). Las reglas se evalúan en orden
// y gana la primera que se cumple.
type ColorRule struct {
	Condition string
	Color     color.Attribute
}

type Column struct {
	Alignment Alignment
	Color     color.Attribute // color por defecto si no cumple ninguna regla
	Colors    []ColorRule     // reglas condicionales, evaluadas en orden
	Format    Format          // formato de presentación del valor
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
