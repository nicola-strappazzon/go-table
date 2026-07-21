package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nicola-strappazzon/go-table"
)

func main() {
	health := table.NewRow()
	health.SetWidth(80)
	health.Title("Cluster Health")
	health.Add(table.Field{Name: "Name", Value: "com-prd-clickhouse"})
	health.Add(table.Field{
		Name:  "Status",
		Value: "green",
		Colors: []table.ColorRule{
			{Condition: `== "green"`, Color: color.FgGreen},
		},
	})
	health.Add(table.Field{Name: "number_of_data_nodes", Value: 9})
	health.Add(table.Field{Name: "active_primary_shards", Value: 9})
	health.Add(table.Field{Name: "active_shards", Value: 312})
	health.Add(table.Field{Name: "relocating_shards", Value: 633})
	health.Add(table.Field{Name: "initializing_shards", Value: 0})
	health.Add(table.Field{Name: "unassigned_shards", Value: 0})
	health.Add(table.Field{Name: "number_of_pending_tasks", Value: 0})
	health.Add(table.Field{Name: "active_shards_percent_as_number", Value: 100, Format: table.Percentage})
	health.Print()
	fmt.Println()

	tblNodes := table.New()
	tblNodes.SetWidth(80)
	tblNodes.Title("Cluster nodes")
	tblNodes.Add("com-prd-clickhouse-node01", "10.10.1.1", 16, 137438953472, 10, 64.2, 1486080)
	tblNodes.Add("com-prd-clickhouse-node02", "10.10.1.2", 16, 137438953472, 91, 64.2, 3600)
	tblNodes.Add("com-prd-clickhouse-node03", "10.10.1.3", 16, 137438953472, 72, 64.2, 1563840)
	tblNodes.Add("com-prd-clickhouse-node04", "10.10.1.4", 16, 137438953472, 60, 64.2, 1563840)
	tblNodes.Column(0, table.Column{Name: "NAME"})
	tblNodes.Column(1, table.Column{Name: "IP"})
	tblNodes.Column(2, table.Column{Name: "CORES"})
	tblNodes.Column(3, table.Column{Name: "RAM", Format: table.Bytes, Width: 7})
	tblNodes.Column(4, table.Column{
		Name:      "CPU",
		Format:    table.Percentage,
		Alignment: table.Right,
		Width:     5,
		Color:     color.FgGreen,
		Colors: []table.ColorRule{
			{Condition: ">= 80", Color: color.FgRed},
			{Condition: ">= 70", Color: color.FgYellow},
		},
	})
	tblNodes.Column(5, table.Column{Name: "DISK"})
	tblNodes.Column(6, table.Column{
		Name:   "UpTime",
		Format: table.Duration,
		Color:  color.FgGreen,
		Colors: []table.ColorRule{
			{Condition: "<= 3600", Color: color.FgRed},
			{Condition: "<= 86400", Color: color.FgYellow},
		},
	})
	tblNodes.Margin(table.Margin{Left: 2})
	tblNodes.SortBy(4)
	tblNodes.Total(2)
	tblNodes.Total(3)
	tblNodes.Print()
}
