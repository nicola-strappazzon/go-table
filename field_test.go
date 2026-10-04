package table

import "testing"

func TestFieldRenderNumericPresentation(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  string
	}{
		{
			name:  "positive signed percentage",
			field: Field{Value: 0.424, Precision: 2, Signed: true, Suffix: "%"},
			want:  "+0.42%",
		},
		{
			name:  "negative signed percentage",
			field: Field{Value: -0.125, Precision: 2, Signed: true, Suffix: "%"},
			want:  "-0.12%",
		},
		{
			name:  "zero has no plus sign",
			field: Field{Value: 0.0, Precision: 2, Signed: true},
			want:  "0.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.Render(); got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTitleSeparatorRespectsLeftMargin(t *testing.T) {
	if got, want := titleSeparator(5, 0), "─────"; got != want {
		t.Errorf("no-margin separator = %q, want %q", got, want)
	}
	if got, want := titleSeparator(5, 2), "  ───"; got != want {
		t.Errorf("indented separator = %q, want %q", got, want)
	}
}

func TestTitleLineUsesHeavyJunctions(t *testing.T) {
	if got, want := titleLine("ETF", 12), "━━┫ ETF ┣━━━"; got != want {
		t.Errorf("titleLine() = %q, want %q", got, want)
	}
}

func TestPaddingWidthIsIndependentFromMargin(t *testing.T) {
	tbl := New().(*table)
	tbl.Padding(2).Margin(Margin{Left: 3})

	if got := tbl.paddingWidth(); got != 2 {
		t.Errorf("paddingWidth() = %d, want 2", got)
	}
}

func TestTitleSeparatorDefaultsToEnabled(t *testing.T) {
	tbl := New().(*table)
	if !tbl.titleSep {
		t.Error("title separator should be enabled by default")
	}
	tbl.TitleSeparator(false)
	if tbl.titleSep {
		t.Error("TitleSeparator(false) did not disable the separator")
	}
}

func TestTitleBoldDefaultsToEnabled(t *testing.T) {
	tbl := New().(*table)
	if !tbl.titleBold {
		t.Error("title bold should be enabled by default")
	}
	tbl.TitleBold(false)
	if tbl.titleBold {
		t.Error("TitleBold(false) did not disable bold")
	}
}

func TestCalculateColumnLenUsesRenderedValue(t *testing.T) {
	tbl := New().(*table)
	tbl.Add(int64(13_656_110_304), int64(355*1024*1024*1024))
	tbl.Column(0, Column{Alignment: Right})
	tbl.Column(1, Column{Format: Bytes, Alignment: Right})

	tbl.calculateColumnStats()

	if got, want := tbl.stats[0].Len, len("13656110304"); got != want {
		t.Errorf("plain numeric width = %d, want %d", got, want)
	}
	if got, want := tbl.stats[1].Len, len("355 GiB"); got != want {
		t.Errorf("formatted width = %d, want %d", got, want)
	}
}

func TestBuildRowKeepsPaddingWithLeftMargin(t *testing.T) {
	tbl := New().(*table)
	tbl.Add(int64(13_656_110_304), int64(508*1024*1024*1024))
	tbl.Column(0, Column{Alignment: Right})
	tbl.Column(1, Column{Format: Bytes, Alignment: Right})
	tbl.Margin(Margin{Left: 2}).Padding(2)
	tbl.calculateColumnStats()

	if got, want := tbl.buildRow(0, tbl.rows[0]), "  13656110304  508 GiB"; got != want {
		t.Errorf("buildRow() = %q, want %q", got, want)
	}
}

func TestFitWidthTruncatesConfiguredColumns(t *testing.T) {
	tbl := New().(*table)
	tbl.Add("a very long table name", "a very long index definition")
	tbl.Column(0, Column{Name: "TABLE", MaxWidth: 12})
	tbl.Column(1, Column{Name: "COLUMNS", MaxWidth: 14})
	tbl.Margin(Margin{Left: 2}).Padding(2).FitWidth(30)
	tbl.calculateColumnStats()
	tbl.fitColumns()

	if got, want := tbl.buildRow(0, tbl.rows[0]), "  a very lo...  a very lo..."; got != want {
		t.Errorf("buildRow() = %q, want %q", got, want)
	}
	if got, want := tbl.naturalWidth(), 30; got != want {
		t.Errorf("naturalWidth() = %d, want %d", got, want)
	}
}
