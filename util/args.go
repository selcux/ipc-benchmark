package util

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Size  int
	Count int
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s -s <message-size> -c <roundtrip-count>\n", os.Args[0])
}

func NewArgs() *Args {
	var size, count int

	flag.IntVar(&size, "s", 0, "Message size")
	flag.IntVar(&count, "c", 0, "Roundtrip count")

	flag.Usage = printUsage

	flag.Parse()

	if size < 1 || count < 1 {
		printUsage()
		os.Exit(1)
	}

	return &Args{Size: size, Count: count}
}
