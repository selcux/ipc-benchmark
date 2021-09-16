package render

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/selcux/ipc-benchmark/benchmark"
)

func Table(results ...*benchmark.MeasuredResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Transferred Messages", "Success Rate (%)", "Duration (sec)"})

	for _, r := range results {
		table.Append([]string{
			fmt.Sprint(r.Type),
			fmt.Sprintf("%d", r.Messages),
			fmt.Sprintf("%.2f", r.SuccessRate*100.0),
			fmt.Sprintf("%s", r.Duration),
		})
	}

	table.Render()
}
