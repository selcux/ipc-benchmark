package benchmark

import "time"

type BenchType string
type IpcName string

const (
	Thr BenchType = "Throughput"
	Lat BenchType = "Latency"
)

const (
	Fifo IpcName = "Named Pipe"
	Tcp  IpcName = "TCP"
	Ucp  IpcName = "UDP"
)

type RunResult struct {
	Messages    int
	SuccessRate float32
	Size        int
}

type MeasuredResult struct {
	*RunResult
	Duration time.Duration
	Type     BenchType
	Ipc      IpcName
}

type Benchmark interface {
	Latency() (*MeasuredResult, error)
	Throughput() (*MeasuredResult, error)
}
