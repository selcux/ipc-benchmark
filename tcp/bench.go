package tcp

import (
	"time"

	"github.com/selcux/ipc-benchmark/benchmark"
)

type Bench struct {
	size  int
	count int
}

func (tb *Bench) Latency() (*benchmark.MeasuredResult, error) {
	tcp := NewTcp(tb.size, tb.count)
	tcp.PrepareConnections()

	start := time.Now()
	runResult, err := tcp.PingPong()
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	return &benchmark.MeasuredResult{
		RunResult: runResult,
		Duration:  duration,
		Type:      benchmark.Lat,
		Ipc:       benchmark.Tcp,
	}, nil
}

func (tb *Bench) Throughput() (*benchmark.MeasuredResult, error) {
	tcp := NewTcp(tb.size, tb.count)
	tcp.PrepareConnections()

	start := time.Now()
	result, err := tcp.Throughput()
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	return &benchmark.MeasuredResult{
		RunResult: &benchmark.RunResult{
			Messages:    int(float64(result.Messages) / duration.Seconds()),
			SuccessRate: result.SuccessRate,
			Size:        result.Size,
		},
		Duration: time.Second,
		Type:     benchmark.Thr,
		Ipc:      benchmark.Tcp,
	}, nil
}

func NewTcpBench(size, count int) *Bench {
	return &Bench{
		size:  size,
		count: count,
	}
}
