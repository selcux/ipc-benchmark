package tcp

import (
	"context"
	"log"

	"github.com/selcux/ipc-benchmark/benchmark"
	"golang.org/x/sync/errgroup"
)

const address = "127.0.0.1:9999"

type Tcp struct {
	size  int
	count int
}

func (t *Tcp) PingPong() (*benchmark.RunResult, error) {
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	readyCh := make(chan struct{})
	srvArgs := ServerArgs{
		maxDataSize:   t.size,
		rotationCount: t.count,
	}
	server := NewServer(srvArgs)
	server.SetHandler(serverHandler)

	g.Go(func() error {
		err := server.Listen(address)
		if err != nil {
			return err
		}

		readyCh <- struct{}{}
		return server.Serve()
	})

	clientCh := make(chan ResultArgs)
	clArgs := ClientArgs{
		maxDataSize:   t.size,
		rotationCount: t.count,
		dataCh:        clientCh,
	}
	client := NewClient(clArgs)
	client.SetHandler(clientHandler)

	g.Go(func() error {
		<-readyCh
		return client.Connect(address)
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	successMsg := 0
	totalBytes := 0
	for i := 0; i < t.count; i++ {
		recv := <-clientCh
		if recv.err == nil {
			successMsg++
			totalBytes += recv.messageLen
		}
	}

	log.Printf("received %d bytes", totalBytes)
	log.Println("End of connection")

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
