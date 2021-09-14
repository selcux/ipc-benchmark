package main

import (
	"fmt"
	"log"

	"github.com/selcux/ipc-benchmark/fifo"
	"github.com/selcux/ipc-benchmark/util"
)

func main() {
	args := util.NewArgs()

	fifoBench := fifo.NewFifoBench(args.Size, args.Count)

	fmt.Println("----- FIFO -----")
	fmt.Println(">> Throughput <<")
	_, err := fifoBench.Throughput()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(">> Latency <<")
	_, err = fifoBench.Latency()
	if err != nil {
		log.Fatalln(err)
	}
}
