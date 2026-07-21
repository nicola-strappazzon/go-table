package table

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

const DOWN_ARROW = "\u2193"

const (
	Left Alignment = iota
	Center
	Right
)

type Value any
type Alignment int

type Stats struct {
	Min float64
	Max float64
	Sum float64
	Len int
}

type Margin struct {
	Left int
}

type Table interface {
	Add(...any) Table
	Column(int, Column) Table
	Count() int
	FilterBy(int, string) Table
	Filtered() int
	Margin(Margin) Table
	Max(int) float64
	Min(int) float64
	Padding(uint) Table
	Print()
	SetWidth(int) Table
	SortBy(int) Table
	Sum(int) float64
	Title(string) Table
	Total(int) Table
	Width() int
	Rows() Rows
}

type table struct {
	columns    map[int]Column
	count      int
	filtered   int
	filters    map[int]string
	padding    uint
	margin     Margin
	rows       Rows
	sort       []int
	sortColumn int
	stats      map[int]Stats
	title      string
	totals     map[int]bool
	width      int
	fixedWidth int
}

func New() Table {
	t := table{}
	t.filters = map[int]string{}
	t.columns = map[int]Column{}
	t.stats = map[int]Stats{}
	t.totals = map[int]bool{}
	t.padding = uint(2)
	t.margin = Margin{}

	return &t
}

func (t *table) Title(title string) Table {
	t.title = title
	return t
}

func (t *table) Add(vals ...any) Table {
	row := Row{}
	for _, val := range vals {
		row.Fields = append(row.Fields, Field{Value: val})
	}
	t.rows = append(t.rows, row)
	t.count++
	return t
}

func (t *table) Padding(p uint) Table {
	t.padding = p
	return t
}

func (t *table) Column(index int, column Column) Table {
	t.columns[index] = column
	return t
}

func (t *table) FilterBy(index int, condition string) Table {
	t.filters[index] = condition

	return t
}

func (t *table) Margin(margin Margin) Table {
	t.margin = margin
	return t
}

func (t *table) SortBy(index int) Table {
	t.sortColumn = index
	t.SetColumnName(
		index,
		fmt.Sprintf("%s%s", t.columns[index].Name, DOWN_ARROW),
	)

	return t
}

func (t *table) Total(index int) Table {
	t.totals[index] = true
	return t
}

func (t *table) SetColumnName(index int, name string) {
	tmp := t.columns[index]
	tmp.SetName(name)
	(*t).columns[index] = tmp
}

func (t *table) Print() {
	t.resetVariables()
	t.clearColumnsValues()
	t.filterRows()
	t.calculateColumnStats()
	t.calculateTableWidth()
	t.printTitle()
	t.printTitleSpacing()
	t.printHeader()
	t.sortRows()
	t.printRows()
	t.printFooter()
}

func (t *table) Filtered() int {
	return t.filtered
}

func (t *table) Count() int {
	return t.count
}

func (t *table) Min(index int) float64 {
	return t.stats[index].Min
}

func (t *table) Max(index int) float64 {
	return t.stats[index].Max
}

func (t *table) Sum(index int) float64 {
	return t.stats[index].Sum
}

func (t *table) Width() int {
	return t.width
}

// SetWidth forces the table width (used for the title header and the separator
// line). Pass 0 to restore the automatic width computed from the columns.
// Note: it does not truncate columns; if the forced width is smaller than the
// content, rows will extend past the decorative lines.
func (t *table) SetWidth(w int) Table {
	t.fixedWidth = w
	return t
}

func (t *table) Rows() Rows {
	return t.rows
}

func (t *table) printTitle() {
	printTitleLine(t.title, t.width)
	printTitleSeparator(t.width)
}

func (t *table) printTitleSpacing() {
	if t.title != "" {
		// fmt.Println()
	}
}

// printTitleLine prints the "─── title ───" header fitted to width. The number
// of dashes is never negative, even when width is smaller than the title.
func printTitleLine(title string, width int) {
	dashes := width - utf8.RuneCountInString(title) - 4
	if dashes < 0 {
		dashes = 0
	}

	c := color.New(color.FgWhite, color.Bold)
	c.Println("━━", title, strings.Repeat("━", dashes))
}

func printTitleSeparator(width int) {
	if width < 2 {
		width = 2
	}
	color.New(color.FgWhite, color.Bold).Println(strings.Repeat(" ", 2) + strings.Repeat("─", width-2))
}

func (t *table) printHeader() {
	c := color.New(color.FgWhite, color.Bold)
	c.Println(t.buildHeader())
	c.DisableColor()
	c.Println(strings.Repeat(" ", t.margin.Left) + strings.Repeat("┄", t.width-t.margin.Left))
}

func (t *table) printRows() {
	for index, row := range t.rows {
		fmt.Println(t.buildRow(index, row))
	}
}

func (t *table) printFooter() {
	if !t.hasTotals() {
		return
	}

	c := color.New(color.FgWhite, color.Bold)
	c.DisableColor()
	c.Println(strings.Repeat(" ", t.margin.Left) + strings.Repeat("┄", t.width-t.margin.Left))
	c.Println(t.buildFooter())
}

func (t *table) sortRows() {
	t.rows.SortBy(t.sortColumn)
}

func (t *table) buildHeader() (p string) {
	for i := 0; i < len(t.columns); i++ {
		if t.columns[i].Alignment == Right {
			p = p + t.lenOffset(t.columns[i].Name, t.stats[i].Len) + t.columns[i].Name + t.printPadding()
		} else {
			p = p + strings.Repeat(" ", t.margin.Left) + t.columns[i].Name + t.lenOffset(t.columns[i].Name, t.stats[i].Len) + t.printPadding()
		}
	}
	return p
}

func (t *table) buildRow(rowIndex int, row Row) (p string) {
	for columnIndex, columnValue := range row.Fields {
		p = p + t.buildColumn(rowIndex, columnIndex, columnValue)
	}
	return strings.TrimRight(p, " ")
}

func (t *table) buildFooter() (p string) {
	for i := 0; i < len(t.columns); i++ {
		p = p + t.buildFooterColumn(i)
	}
	return strings.TrimRight(p, " ")
}

func (t *table) buildColumn(rowIndex, columnIndex int, columnValue Field) string {
	pf := t.columns[columnIndex].toField(columnValue.Value)

	field := pf.ToString()
	if rowIndex >= 0 {
		field = pf.Render()
	}

	colored := pf.Colorize(field)

	if pf.Alignment == Right {
		return t.lenOffset(field, t.stats[columnIndex].Len) + colored + t.printPadding()
	}

	return strings.Repeat(" ", t.margin.Left) + colored + t.lenOffset(field, t.stats[columnIndex].Len) + t.printPadding()
}

func (t *table) buildFooterColumn(columnIndex int) string {
	field := ""
	colored := ""

	if t.totals[columnIndex] {
		pf := t.columns[columnIndex].toField(t.stats[columnIndex].Sum)
		field = pf.Render()
		colored = color.New(color.FgWhite, color.Bold).Sprint(field)
	}

	if t.columns[columnIndex].Alignment == Right {
		return t.lenOffset(field, t.stats[columnIndex].Len) + colored + t.printPadding()
	}

	return strings.Repeat(" ", t.margin.Left) + colored + t.lenOffset(field, t.stats[columnIndex].Len) + t.printPadding()
}

func (t *table) hasTotals() bool {
	for _, enabled := range t.totals {
		if enabled {
			return true
		}
	}
	return false
}

func (t *table) resetVariables() {
	t.filtered = 0
	t.width = 0
}

func (t *table) clearColumnsValues() {
	for rowIndex, rowValue := range t.rows {
		for columnIndex := range rowValue.Fields {
			(*t).rows[rowIndex].Fields[columnIndex].Clear()
			(*t).rows[rowIndex].Fields[columnIndex].Truncate(t.columns[columnIndex].Truncate)
		}
	}
}

func (t *table) filterRows() {
	for index := len(t.rows) - 1; index >= 0; index-- {
		for columnIndex, columnValue := range t.rows[index].Fields {
			if condition, ok := t.filters[columnIndex]; ok {
				if !columnValue.EvalCondition(condition) {
					t.rows.Remove(index)
					t.filtered++
				}
			}
		}
	}
}

func (t *table) calculateColumnStats() {
	for index, value := range t.columns {
		var v float64

		if len(t.rows) > 0 {
			v = t.rows[0].Fields[index].ToFloat64()
		}

		t.stats[index] = Stats{
			Len: t.calculateColumnLen(index, t.calculateColumnBaseLen(index, value)),
			Min: t.calculateColumnMin(index, v),
			Max: t.calculateColumnMax(index, v),
			Sum: t.calculateColumnSum(index),
		}
	}
}

func (t *table) calculateColumnBaseLen(index int, column Column) int {
	width := column.GetWidth()
	if t.totals[index] {
		total := column.toField(t.calculateColumnSum(index)).Render()
		width = int(math.Max(float64(width), float64(utf8.RuneCountInString(total))))
	}
	return width
}

func (t *table) calculateColumnMin(index int, value float64) float64 {
	for x := 1; x < len(t.rows); x++ {
		value = math.Min(value, t.rows[x].Fields[index].ToFloat64())
	}

	return value
}

func (t *table) calculateColumnMax(index int, value float64) float64 {
	for x := 1; x < len(t.rows); x++ {
		value = math.Max(value, t.rows[x].Fields[index].ToFloat64())
	}

	return value
}

func (t *table) calculateColumnSum(index int) (sum float64) {
	for x := 0; x < len(t.rows); x++ {
		sum = sum + t.rows[x].Fields[index].ToFloat64()
	}

	return sum
}

func (t *table) calculateColumnLen(index int, value int) int {
	for x := 0; x < len(t.rows); x++ {
		value = int(math.Max(float64(value), float64(t.rows[x].Fields[index].Len())))
	}

	return value
}

func (t *table) calculateTableWidth() {
	if t.fixedWidth > 0 {
		t.width = t.fixedWidth
		return
	}

	for _, stats := range t.stats {
		t.width = t.width + stats.Len + int(t.padding)
	}
}

func (t *table) lenOffset(s string, w int) string {
	l := w - utf8.RuneCountInString(s)
	if l < 0 {
		return ""
	}

	return strings.Repeat(" ", l)
}

func (t *table) printPadding() string {
	return strings.Repeat(" ", int(t.padding)-t.margin.Left)
}
