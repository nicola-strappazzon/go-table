package main

import "github.com/fatih/color"

func main() {
	tbl := New()
	tbl.Title("Cluster nodes")
	tbl.Add("com-prd-clickhouse-node01", "10.10.1.1", 16, 137438953472, 10, 64.2, 1486080)
	tbl.Add("com-prd-clickhouse-node02", "10.10.1.2", 16, 137438953472, 91, 64.2, 3600)
	tbl.Add("com-prd-clickhouse-node03", "10.10.1.3", 16, 137438953472, 72, 64.2, 1563840)
	tbl.Add("com-prd-clickhouse-node04", "10.10.1.4", 16, 137438953472, 60, 64.2, 1563840)
	tbl.Column(0, Column{Name: "NAME"})
	tbl.Column(1, Column{Name: "IP"})
	tbl.Column(2, Column{Name: "CORES"})
	tbl.Column(3, Column{Name: "RAM", Format: Bytes, Width: 7})
	tbl.Column(4, Column{
		Name:      "CPU",
		Format:    Percentage,
		Alignment: Right,
		Width:     5,
		Color:     color.FgGreen,
		Colors: []ColorRule{
			{">= 80", color.FgRed},
			{">= 70", color.FgYellow},
		},
	})
	tbl.Column(5, Column{Name: "DISK"})
	tbl.Column(6, Column{
		Name:   "UpTime",
		Format: Duration,
		Color:  color.FgGreen,
		Colors: []ColorRule{
			{"<= 3600", color.FgRed},
			{"<= 86400", color.FgYellow},
		},
	})
	// tbl.Column(2, Column{ZeroFill: true, Precision: 3, Scale: 1})
	tbl.Margin(Margin{Left: 2})
	// tbl.FilterBy(1, ">= 2000")
	tbl.SortBy(4)
	tbl.Print()
	// fmt.Printf("%s\nRows filtered: %d/%d\n", strings.Repeat("─", tbl.Width()), tbl.Filtered(), tbl.Count())
	// fmt.Printf("Rate = Sum: %f, Min: %f, Max: %f\n", tbl.Sum(2), tbl.Min(2), tbl.Max(2))
	// fmt.Printf("Table size: %dx%d\n", terminal.Height(), terminal.Width())
}
