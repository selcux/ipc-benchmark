package benchmark

type Throughput interface {
	Produce() (*RunResult, error)
	Consume() (*RunResult, error)
}
