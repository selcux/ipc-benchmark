package benchmark

type Latency interface {
	PingPong() (*RunResult, error)
}
