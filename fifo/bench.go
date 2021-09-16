package fifo

import (
	"fmt"
	"github.com/selcux/ipc-benchmark/benchmark"
	"io/ioutil"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

type Bench struct {
	size  int
	count int
}

func (fb *Bench) Latency() (*benchmark.MeasuredResult, error) {
	tmpDir, err := ioutil.TempDir("", "named-pipes")
	if err != nil {
		return nil, errors.Wrap(err, "could not create temp dir 'named-pipes'")
	}
	// Create named pipe
	namedPipe := filepath.Join(tmpDir, "ping_pong")
	err = syscall.Mkfifo(namedPipe, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "could not make fifo of %s", namedPipe)
	}

	fifo := NewFifo(namedPipe, fb.size, fb.count)
	start := time.Now()
	runResult, err := fifo.PingPong()
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	return &benchmark.MeasuredResult{
		RunResult: runResult,
		Duration:  duration,
		Type: benchmark.Lat,
	}, nil
}

func (fb *Bench) Throughput() (*benchmark.MeasuredResult, error) {
	tmpDir, err := ioutil.TempDir("", "named-pipes")
	if err != nil {
		return nil, errors.Wrap(err, "could not create temp dir 'named-pipes'")
	}
	// Create named pipe
	namedPipe := filepath.Join(tmpDir, "stdout")
	err = syscall.Mkfifo(namedPipe, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "could not make fifo of %s", namedPipe)
	}

	fifo := NewFifo(namedPipe, fb.size, fb.count)
	resultCh := make(chan *benchmark.RunResult)
	errCh := make(chan error)

	start := time.Now()
	go func(ff *Fifo, count chan *benchmark.RunResult, e chan error) {
		c, err := ff.Produce()
		if err != nil {
			e <- err
			close(count)
			return
		}

		count <- c
		close(e)
	}(fifo, resultCh, errCh)

	cResult, cErr := fifo.Consume()
	pResult := <-resultCh
	pErr := <-errCh

	if pErr != nil {
		return nil, pErr
	}

	if cErr != nil {
		return nil, cErr
	}

	fmt.Printf("Total produced bytes: %+v\n", pResult)
	fmt.Printf("Total consumed bytes: %+v\n", cResult)

	duration := time.Since(start)

	return &benchmark.MeasuredResult{
		RunResult: cResult,
		Duration:  duration,
		Type: benchmark.Thr,
	}, nil
}

func NewFifoBench(size, count int) *Bench {
	return &Bench{
		size:  size,
		count: count,
	}
}
