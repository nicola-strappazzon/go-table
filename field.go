package table

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Field is a cell value together with its presentation (label, format and
// color). The presentation is optional: the Table fills it from its Column when
// rendering, and a standalone Row sets it directly on each Field.
type Field struct {
	Value     any
	Name      string          // label (used when printing a vertical Row)
	Format    Format          // value presentation format
	Color     color.Attribute // default color when no rule matches
	Colors    []ColorRule     // conditional rules, evaluated in order
	Precision int
	Scale     int
	ZeroFill  bool
	Alignment Alignment
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

func (f Field) ZeroFilled(precision, scale int) string {
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

// ToDuration reads the cell value as seconds and returns it humanized with the
// largest applicable unit: "45s", "3m", "1h", "17.2d".
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

// ToBytes reads the cell value as bytes and returns it humanized with binary
// units (base 1024): "512B", "64KiB", "128GiB".
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
	operand := f.ToString()
	if f.IsString() {
		operand = strconv.Quote(operand)
	}

	fs := token.NewFileSet()
	tv, err := types.Eval(
		fs,
		nil,
		token.NoPos,
		fmt.Sprintf("%s %s", operand, condition))

	if err != nil || tv.Value == nil || tv.Value.Kind() != constant.Bool {
		return false
	}

	return constant.BoolVal(tv.Value)
}

// Render returns the value formatted according to Format and ZeroFill, without
// coloring or aligning.
func (f Field) Render() string {
	field := f.ToString()

	switch f.Format {
	case Percentage:
		field = fmt.Sprintf("%s%%", f.ToString())
	case Duration:
		field = f.ToDuration()
	case Bytes:
		field = f.ToBytes()
	}

	if f.ZeroFill {
		field = f.ZeroFilled(f.Precision, f.Scale)
	}

	return field
}

// Colorize paints text with the first ColorRule whose condition the Field value
// satisfies; if none match it uses Color. It returns text unchanged when no
// color is configured (width is measured separately on the uncolored text so
// the ANSI codes do not break alignment).
func (f Field) Colorize(text string) string {
	attr := f.Color
	for _, rule := range f.Colors {
		if f.EvalCondition(rule.Condition) {
			attr = rule.Color
			break
		}
	}

	if attr == color.Reset {
		return text
	}

	return color.New(attrs256(attr)...).Sprint(text)
}

func attrs256(attr color.Attribute) []color.Attribute {
	switch attr {
	case color.FgGreen:
		return []color.Attribute{38, 5, 34}
	case color.FgRed:
		return []color.Attribute{38, 5, 203}
	case color.FgYellow:
		return []color.Attribute{38, 5, 208}
	default:
		return []color.Attribute{attr}
	}
}
