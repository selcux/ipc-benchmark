package render

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/selcux/ipc-benchmark/benchmark"
)

func Table(results ...*benchmark.MeasuredResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	//table.SetRowLine(true)
	table.SetHeader([]string{"IPC", "Type", "Transferred Messages", "Duration", "Success Rate"})

	for _, r := range results {
		table.Append([]string{
			fmt.Sprint(r.Ipc),
			fmt.Sprint(r.Type),
			fmt.Sprintf("%d", r.Messages),
			fmt.Sprintf("%s", r.Duration),
			fmt.Sprintf("%.2f", r.SuccessRate*100.0),
		})
	}

	table.Render()
}
