package benchmark

type Benchmark interface {
    Latency() (int, error)
	Throughput() (int, error)
}

type Producer interface {
	Produce() (int64, error)
}

type Consumer interface {
	Consume() (int64, error)
}
