package main

import (
	"log"

	"github.com/selcux/ipc-benchmark/fifo"
	"github.com/selcux/ipc-benchmark/util"
)

func main() {
	args := util.NewArgs()

	fifoBench := fifo.NewFifoBench(args.Size, args.Count)

	_, err := fifoBench.Throughput()
	if err != nil {
		log.Fatalln(err)
	}
}
