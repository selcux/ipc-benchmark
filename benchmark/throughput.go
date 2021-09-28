package benchmark

type Throughput interface {
	Produce() error
	Consume() (*RunResult, error)
}
