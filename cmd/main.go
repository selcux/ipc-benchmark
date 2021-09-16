package main

import (
	"fmt"
	"github.com/selcux/ipc-benchmark/report/render"
	"log"

	"github.com/selcux/ipc-benchmark/fifo"
	"github.com/selcux/ipc-benchmark/util"
)

func main() {
	args := util.NewArgs()

	fifoBench := fifo.NewFifoBench(args.Size, args.Count)

	fmt.Println("----- FIFO -----")
	fmt.Println(">> Throughput <<")
	tResult, err := fifoBench.Throughput()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(">> Latency <<")
	lResult, err := fifoBench.Latency()
	if err != nil {
		log.Fatalln(err)
	}

	render.Table(tResult, lResult)
}
