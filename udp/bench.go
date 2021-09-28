package udp

import (
	"github.com/selcux/ipc-benchmark/benchmark"
	"time"
)

type Bench struct {
	size  int
	count int
}

func (b *Bench) Latency() (*benchmark.MeasuredResult, error) {
	udp := NewUdp(b.size, b.count)

	start := time.Now()
	result, err := udp.PingPong()
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	return &benchmark.MeasuredResult{
		RunResult: result,
		Duration:  duration,
		Type:      benchmark.Lat,
		Ipc:       benchmark.Ucp,
	}, nil
}

func (b *Bench) Throughput() (*benchmark.MeasuredResult, error) {
	udp := NewUdp(b.size, b.count)

	start := time.Now()
	result, err := udp.Throughput()
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
		Ipc:      benchmark.Ucp,
	}, nil
}

func NewUdpBench(size, count int) *Bench {
	return &Bench{
		size:  size,
		count: count,
	}
}
