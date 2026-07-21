package table

import (
	"fmt"
	"unicode/utf8"
)

// Row is a list of Fields with an optional title. It works both as a Table's
// internal row and as a standalone record that can be printed vertically
// ("label: value") with NewRow().Add(...).Print().
type Row struct {
	title      string
	fixedWidth int
	Fields     []Field
}

func NewRow() *Row {
	return &Row{}
}

func (r *Row) Title(title string) *Row {
	r.title = title
	return r
}

// SetWidth forces the Row width (the title header). Pass 0 to restore the
// automatic width computed from the content.
func (r *Row) SetWidth(width int) *Row {
	r.fixedWidth = width
	return r
}

func (r *Row) Add(field Field) *Row {
	r.Fields = append(r.Fields, field)
	return r
}

// Print prints the Row vertically: one line per Field as "label: value", using
// the title as header and applying each Field's format and color.
func (r *Row) Print() {
	values := r.renderValues()
	labelWidth := r.calculateLabelWidth()
	width := r.calculateWidth(labelWidth, values)

	printTitleLine(r.title, width)
	r.printFields(labelWidth, values)
}

// printFields prints one "label: value" line per Field, with the label right
// aligned to labelWidth and the value already formatted and colored.
func (r *Row) printFields(labelWidth int, values []string) {
	for i, field := range r.Fields {
		fmt.Printf("%*s: %s\n", labelWidth, field.Name, field.Colorize(values[i]))
	}
}

// renderValues returns the formatted value (Field.Render) of each Field, in the
// same order as r.Fields.
func (r *Row) renderValues() []string {
	values := make([]string, len(r.Fields))
	for i, field := range r.Fields {
		values[i] = field.Render()
	}
	return values
}

// calculateLabelWidth returns the width of the longest label (Field name), used
// to right align the labels column.
func (r *Row) calculateLabelWidth() (width int) {
	for _, field := range r.Fields {
		if l := utf8.RuneCountInString(field.Name); l > width {
			width = l
		}
	}
	return width
}

// calculateWidth returns the title header width. If forced with SetWidth it is
// honored; otherwise it uses the longest line ("label: value") with a minimum
// based on the title itself.
func (r *Row) calculateWidth(labelWidth int, values []string) int {
	if r.fixedWidth > 0 {
		return r.fixedWidth
	}

	width := utf8.RuneCountInString(r.title) + 6
	for _, value := range values {
		if w := labelWidth + 2 + utf8.RuneCountInString(value); w > width {
			width = w
		}
	}
	return width
}
