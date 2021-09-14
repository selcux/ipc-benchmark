package fifo

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
)

type FifoBench struct {
	size  int
	count int
}

func (fb *FifoBench) Latency() (int, error) {
	tmpDir, err := ioutil.TempDir("", "named-pipes")
	if err != nil {
		return 0, errors.Wrap(err, "could not create temp dir 'named-pipes'")
	}
	// Create named pipe
	namedPipe := filepath.Join(tmpDir, "ping_pong")
	syscall.Mkfifo(namedPipe, 0600)

	fifo := NewFifo(namedPipe, fb.size, fb.count)
	err = fifo.PingPong()
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (fb *FifoBench) Throughput() (int, error) {
	tmpDir, err := ioutil.TempDir("", "named-pipes")
	if err != nil {
		return 0, errors.Wrap(err, "could not create temp dir 'named-pipes'")
	}
	// Create named pipe
	namedPipe := filepath.Join(tmpDir, "stdout")
	syscall.Mkfifo(namedPipe, 0600)

	fifo := NewFifo(namedPipe, fb.size, fb.count)
	countCh := make(chan int64)
	errCh := make(chan error)

	go func(ff *Fifo, count chan int64, e chan error) {
		c, err := ff.Produce()
		if err != nil {
			e <- err
			close(count)
			return
		}

		count <- c
		close(e)
	}(fifo, countCh, errCh)

	cCount, cErr := fifo.Consume()
	pCount := <-countCh
	pErr := <-errCh

	if pErr != nil {
		return 0, pErr
	}

	if cErr != nil {
		return 0, cErr
	}

	fmt.Printf("Total produced bytes: %d\n", pCount)
	fmt.Printf("Total consumed bytes: %d\n", cCount)

	return 0, nil
}

func NewFifoBench(size, count int) *FifoBench {
	return &FifoBench{
		size:  size,
		count: count,
	}
}
