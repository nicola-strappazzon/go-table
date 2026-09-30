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

func TestPaddingWidthNeverNegative(t *testing.T) {
	tbl := New().(*table)
	tbl.Padding(2).Margin(Margin{Left: 3})

	if got := tbl.paddingWidth(); got != 0 {
		t.Errorf("paddingWidth() = %d, want 0", got)
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
