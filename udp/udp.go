package udp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/selcux/ipc-benchmark/benchmark"
	"github.com/selcux/ipc-benchmark/util"
	"net"
)

const address = "127.0.0.1:"

type Udp struct {
	size  int
	count int
}

func (u *Udp) PingPong() (*benchmark.RunResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverAddr, err := u.echoServer(ctx, address)
	if err != nil {
		return nil, errors.Wrap(err, "could not create server")
	}

	client, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, errors.Wrap(err, "could not create client")
	}
	defer func() { _ = client.Close() }()

	data, err := util.GenRandomBytes(u.size)
	if err != nil {
		return nil, err
	}

	successMsg := 0
	totalBytes := 0
	for i := 0; i < u.count; i++ {
		_, err = client.WriteTo(data, serverAddr)
		if err != nil {
			return nil, errors.Wrap(err, "could not send message to server")
		}

		data, _, err := receive(client, serverAddr, u.size)
		if err != nil {
			continue
		}

		successMsg++
		totalBytes += len(data)
	}

	return &benchmark.RunResult{
		Messages:    successMsg,
		SuccessRate: float32(successMsg) / float32(u.count),
		Size:        totalBytes,
	}, nil
}

func (u *Udp) Throughput() (*benchmark.RunResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, errors.Wrap(err, "could not create client")
	}
	defer func() { _ = client.Close() }()

	serverAddr, err := u.thrServer(ctx, address, client.LocalAddr())
	if err != nil {
		return nil, errors.Wrap(err, "could not create server")
	}

	successMsg := 0
	totalBytes := 0
	for i := 0; i < u.count; i++ {
		data, _, err := receive(client, serverAddr, u.size)
		if err != nil {
			continue
		}

		successMsg++
		totalBytes += len(data)
	}

	return &benchmark.RunResult{
		Messages:    successMsg,
		SuccessRate: float32(successMsg) / float32(u.count),
		Size:        totalBytes,
	}, nil
}

func (u *Udp) echoServer(ctx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "binding to udp %s", addr)
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		for i := 0; i < u.count; i++ {
			data, clientAddr, err := receive(s, nil, u.size)
			if err != nil {
				panic(err)
			}

			_, err = s.WriteTo(data, clientAddr) // server to client
			if err != nil {
				panic(err)
			}
		}
	}()

	return s.LocalAddr(), nil
}

func (u *Udp) thrServer(ctx context.Context, addr string, clientAddr net.Addr) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "binding to udp %s", addr)
	}

	data, err := util.GenRandomBytes(u.size)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		for i := 0; i < u.count; i++ {
			_, err = s.WriteTo(data, clientAddr) // server to client
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}

func receive(conn net.PacketConn, sourceAddr net.Addr, maxSize int) ([]byte, net.Addr, error) {
	buf := make([]byte, maxSize)

	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to receive data")
	}

	if sourceAddr != nil && addr.String() != sourceAddr.String() {
		return nil, addr, errors.Errorf("received reply from %q instead of %q", addr, sourceAddr)
	}

	return buf[:n], addr, nil
}

func NewUdp(size, count int) *Udp {
	return &Udp{size: size, count: count}
}
