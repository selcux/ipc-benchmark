package fifo

import (
	"io/ioutil"
	"path/filepath"
	"syscall"
	"time"

	"github.com/selcux/ipc-benchmark/benchmark"

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
		Type:      benchmark.Lat,
		Ipc:       benchmark.Fifo,
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
	errCh := make(chan error)

	start := time.Now()
	go func(ff *Fifo, e chan error) {
		err := ff.Produce()
		if err != nil {
			e <- err
			return
		}

		close(e)
	}(fifo, errCh)

	cResult, cErr := fifo.Consume()
	pErr := <-errCh

	if pErr != nil {
		return nil, pErr
	}

	if cErr != nil {
		return nil, cErr
	}

	duration := time.Since(start)

	return &benchmark.MeasuredResult{
		RunResult: &benchmark.RunResult{
			Messages:    int(float64(cResult.Messages) / duration.Seconds()),
			SuccessRate: cResult.SuccessRate,
			Size:        cResult.Size,
		},
		Duration: time.Second,
		Type:     benchmark.Thr,
		Ipc:      benchmark.Fifo,
	}, nil
}

func NewFifoBench(size, count int) *Bench {
	return &Bench{
		size:  size,
		count: count,
	}
}
