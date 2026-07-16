package main

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type Field struct {
	Value any
}

func (f Field) IsString() bool {
	return reflect.ValueOf(f.Value).Kind() == reflect.String
}

func (f Field) ToString() string {
	return fmt.Sprintf("%v", f.Value)
}

func (f Field) ToFloat64() float64 {
	return InterfaceToFloat64(f.Value)
}

func (f Field) Len() int {
	if f.IsString() {
		return len(f.ToString())
	}
	return 0
}

func (f *Field) Truncate(t int) {
	if f.IsString() && t < f.Len() && t > 0 {
		(*f).Value = f.ToString()[:t]
	}
}

func (f *Field) Clear() {
	if f.IsString() {
		t := f.ToString()
		t = strings.TrimSpace(t)
		t = strings.ReplaceAll(t, "\n", " ")
		t = strings.ReplaceAll(t, "\r", " ")
		t = strings.ReplaceAll(t, "  ", " ")

		(*f).Value = t
	}
}

func (f Field) ZeroFill(precision, scale int) string {
	var m float64
	var n float64

	if v, err := strconv.ParseFloat(f.ToString(), 64); err == nil {
		n = v
		if v < 1e-6 {
			m = 1e6
		} else if v < 1e-3 {
			m = 1e3
		} else {
			m = 1
		}
	}
	return fmt.Sprintf(
		"%[1]*.[2]*[3]f",
		precision,
		scale, n*m)
}

// ToDuration interpreta el valor de la celda como segundos y lo devuelve
// humanizado con la unidad mayor que aplique: "45s", "3m", "1h", "17.2d".
func (f Field) ToDuration() string {
	s := f.ToFloat64()

	round1 := func(x float64) float64 { return math.Round(x*10) / 10 }

	switch {
	case s < 60:
		return fmt.Sprintf("%gs", round1(s))
	case s < 3600:
		return fmt.Sprintf("%gm", round1(s/60))
	case s < 86400:
		return fmt.Sprintf("%gh", round1(s/3600))
	default:
		return fmt.Sprintf("%gd", round1(s/86400))
	}
}

// ToBytes interpreta el valor de la celda como bytes y lo devuelve humanizado
// con unidades binarias (base 1024): "512B", "64KiB", "128GiB".
func (f Field) ToBytes() string {
	b := f.ToFloat64()
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	i := 0
	for b >= 1024 && i < len(units)-1 {
		b /= 1024
		i++
	}

	return fmt.Sprintf("%g %s", math.Round(b*10)/10, units[i])
}

func (f Field) EvalCondition(condition string) bool {
	fs := token.NewFileSet()
	tv, err := types.Eval(
		fs,
		nil,
		token.NoPos,
		fmt.Sprintf("%s %s", f.ToString(), condition))

	if err != nil || tv.Value == nil || tv.Value.Kind() != constant.Bool {
		return false
	}

	return constant.BoolVal(tv.Value)
}
