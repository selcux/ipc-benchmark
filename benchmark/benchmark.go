package benchmark

import "time"

type MeasureType string

const (
	Thr MeasureType = "Throughput"
	Lat MeasureType = "Latency"
)

type RunResult struct {
	Messages    int
	SuccessRate float32
	Size        int
}

type MeasuredResult struct {
	*RunResult
	Duration time.Duration
	Type     MeasureType
}

type Benchmark interface {
	Lat(lat Latency) (*MeasuredResult, error)
	Thr(thr Throughput) (*MeasuredResult, error)
}
