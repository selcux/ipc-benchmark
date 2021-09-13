package fifo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/selcux/ipc-benchmark/util"
)

type Fifo struct {
	namedPipe string
	size      int
	count     int
}

func (f *Fifo) Produce() (int64, error) {
	fmt.Println("Opening named pipe for writing")
	stdout, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return 0, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer stdout.Close()
	fmt.Println("Writing")

	var count int64 = 0
	for i := 0; i < f.count; i++ {
		data, err := util.GenRandomBytes(f.size)
		if err != nil {
			return 0, errors.Wrap(err, "unable to create random byte array")
		}

		stdout.Write(data)
		count += int64(len(data))
	}

	return count, nil
}

func (f *Fifo) Consume() (int64, error) {
	// Open named pipe for reading
	fmt.Println("Opening named pipe for reading")
	stdout, err := os.OpenFile(f.namedPipe, os.O_RDONLY, 0600)
	if err != nil {
		return 0, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer stdout.Close()
	fmt.Println("Reading")
	fmt.Println("Waiting for someone to write something")

	reader := bufio.NewReader(stdout)
	chunk := make([]byte, f.size)

	var count int64 = 0
	rcount := 0
	for i := 0; i < f.count; i++ {
		if rcount, err = reader.Read(chunk); err != nil {
			break
		}
		count += int64(rcount)
	}

	if err != nil && err != io.EOF {
		log.Fatalln("Error reading ", stdout, ": ", err)
		return 0, errors.Wrapf(err, "error reading %s", stdout)
	}

	return count, nil
}

func NewFifo(namedPipe string, size, count int) *Fifo {
	return &Fifo{
		namedPipe: namedPipe,
		size:      size,
		count:     count,
	}
}
