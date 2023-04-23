package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func RenderTable(header []string, rows [][]string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColumnSeparator(" ")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(rows)

	return table
}
