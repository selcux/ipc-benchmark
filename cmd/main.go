package main

import (
	"github.com/selcux/ipc-benchmark/udp"
	"log"

	"github.com/selcux/ipc-benchmark/report/render"
	"github.com/selcux/ipc-benchmark/tcp"

	"github.com/selcux/ipc-benchmark/fifo"
	"github.com/selcux/ipc-benchmark/util"
)

/*
func init() {
	log.SetOutput(ioutil.Discard)
}
*/
func main() {
	args := util.NewArgs()

	fifoBench := fifo.NewFifoBench(args.Size, args.Count)
	fifoResultT, err := fifoBench.Throughput()
	if err != nil {
		log.Fatalln(err)
	}

	fifoResultL, err := fifoBench.Latency()
	if err != nil {
		log.Fatalln(err)
	}

	tcpBench := tcp.NewTcpBench(args.Size, args.Count)
	tcpResultT, err := tcpBench.Throughput()
	if err != nil {
		log.Fatalln(err)
	}

	tcpResultL, err := tcpBench.Latency()
	if err != nil {
		log.Fatalln(err)
	}

	udpBench := udp.NewUdpBench(args.Size, args.Count)
	udpResultT, err := udpBench.Throughput()
	if err != nil {
		log.Fatalln(err)
	}

	udpResultL, err := udpBench.Latency()
	if err != nil {
		log.Fatalln(err)
	}

	render.Table(fifoResultT, fifoResultL, tcpResultT, tcpResultL, udpResultT, udpResultL)
}
