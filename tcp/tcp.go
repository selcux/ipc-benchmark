package tcp

import (
	"context"

	"github.com/selcux/ipc-benchmark/benchmark"
	"golang.org/x/sync/errgroup"
)

const address = "127.0.0.1:9999"

type Tcp struct {
	size   int
	count  int
	server *Server
	client *Client
}

func (t *Tcp) PrepareConnections() {
	errCh := make(chan error, 1)
	srvArgs := ServerArgs{
		maxDataSize:   t.size,
		rotationCount: t.count,
		errCh:         errCh,
	}
	server := NewServer(srvArgs)
	t.server = server

	clientCh := make(chan ResultArgs)
	clArgs := ClientArgs{
		maxDataSize:   t.size,
		rotationCount: t.count,
		dataCh:        clientCh,
	}
	client := NewClient(clArgs)
	t.client = client
}

func (t *Tcp) init() error {
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	readyCh := make(chan struct{})

	g.Go(func() error {
		err := t.server.Listen(address)
		if err != nil {
			return err
		}

		readyCh <- struct{}{}
		return t.server.Serve()
	})

	g.Go(func() error {
		<-readyCh
		return t.client.Connect(address)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if err := <-t.server.args.errCh; err != nil {
		return err
	}

	return nil
}

func (t *Tcp) PingPong() (*benchmark.RunResult, error) {
	t.server.SetHandler(serverLatencyHandler)
	t.client.SetHandler(clientLatencyHandler)

	defer func() { t.server.Close() }()

	err := t.init()
	if err != nil {
		return nil, err
	}

	clientCh := t.client.DataCh()
	successMsg := 0
	totalBytes := 0
	for i := 0; i < t.count; i++ {
		recv := <-clientCh
		if recv.err == nil {
			successMsg++
			totalBytes += recv.messageLen
		}
	}

	return &benchmark.RunResult{
		Messages:    successMsg,
		SuccessRate: float32(successMsg) / float32(t.count),
		Size:        totalBytes,
	}, nil
}

func (t *Tcp) Throughput() (*benchmark.RunResult, error) {
	t.server.SetHandler(serverThroughputHandler)
	t.client.SetHandler(clientThroughputHandler)

	defer func() { t.server.Close() }()

	err := t.init()
	if err != nil {
		return nil, err
	}

	dataCh := t.client.DataCh()
	successMsg := 0
	totalBytes := 0
	for i := 0; i < t.count; i++ {
		recv := <-dataCh
		if recv.err == nil {
			successMsg++
			totalBytes += recv.messageLen
		}
	}

	return &benchmark.RunResult{
		Messages:    successMsg,
		SuccessRate: float32(successMsg) / float32(t.count),
		Size:        totalBytes,
	}, nil
}

func NewTcp(size, count int) *Tcp {
	return &Tcp{
		size:  size,
		count: count,
	}
}
