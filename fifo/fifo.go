package fifo

import (
	"bufio"
	"github.com/selcux/ipc-benchmark/benchmark"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/selcux/ipc-benchmark/util"
)

const chunkSize int = 1024

type Fifo struct {
	namedPipe string
	size      int
	count     int
}

func (f *Fifo) Produce() error {
	file, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file.Close()

	data, err := util.GenRandomBytes(f.size)
	if err != nil {
		return errors.Wrap(err, "unable to create random byte array")
	}

	for i := 0; i < f.count; i++ {
		_, err := f.send(file, data)
		if err != nil {
			return errors.Wrap(err, "could not send data")
		}
	}

	return nil
}

func (f *Fifo) Consume() (*benchmark.RunResult, error) {
	// Open named pipe for reading
	file, err := os.OpenFile(f.namedPipe, os.O_RDONLY, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file.Close()

	successMsg := 0
	totalBytes := 0
	for i := 0; i < f.count; i++ {
		data, err := f.receive(file)
		/*
			if err != nil {
				return nil, errors.Wrap(err, "unable to receive data")
			}
		*/
		if err == nil {
			successMsg++
			totalBytes += len(data)
		}
	}

	return &benchmark.RunResult{
		Messages:    successMsg,
		SuccessRate: float32(successMsg) / float32(f.count),
		Size:        totalBytes,
	}, nil
}

func (f *Fifo) PingPong() (*benchmark.RunResult, error) {
	file1, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file1.Close()

	file2, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file2.Close()

	successMsg := 0
	totalBytes := 0
	for i := 0; i < f.count; i++ {
		transferredBytes, err := f.pingPong(file1, file2)
		/*
			if err != nil {
				return nil, errors.Wrap(err, "communication error")
			}
		*/
		if err == nil {
			successMsg++
			totalBytes += transferredBytes
		}
	}

	return &benchmark.RunResult{
		Messages:    successMsg,
		SuccessRate: float32(successMsg) / float32(f.count),
		Size:        totalBytes,
	}, nil
}

func (f *Fifo) pingPong(file1, file2 *os.File) (int, error) {
	errCh := make(chan error)
	defer close(errCh)

	data, err := util.GenRandomBytes(f.size)
	if err != nil {
		return 0, errors.Wrap(err, "unable to create random byte array")
	}

	go func(e chan error, data []byte) {
		_, err := f.send(file1, data)
		e <- err
	}(errCh, data)

	recv1, err := f.receive(file2)
	if err != nil {
		return 0, errors.Wrap(err, "could not received message - pong")
	}
	err = <-errCh
	if err != nil {
		return 0, errors.Wrap(err, "could not send message - ping")
	}

	go func(e chan error, data []byte) {
		_, err := f.send(file2, data)
		e <- err
	}(errCh, recv1)

	recv2, err := f.receive(file1)
	if err != nil {
		return 0, errors.Wrap(err, "could not received message - ping")
	}
	err = <-errCh
	if err != nil {
		return 0, errors.Wrap(err, "could not send message - ping")
	}

	return len(recv2), nil
}

func (f *Fifo) send(file *os.File, data []byte) (int, error) {
	n, err := file.Write(data)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to write to file %s", file.Name())
	}

	return n, nil
}

func (f *Fifo) receive(file *os.File) ([]byte, error) {
	reader := bufio.NewReader(file)
	data := make([]byte, f.size)
	chunk := make([]byte, chunkSize)

	count := 0

	for count < f.size {
		n, err := reader.Read(chunk)
		if n == 0 {
			if err == nil {
				continue
			}

			if err == io.EOF {
				break
			}

			return nil, errors.Wrapf(err, "error reading %s", file.Name())
		}

		copy(data[count:], chunk)
		count += n
	}

	return data, nil
}

func NewFifo(namedPipe string, size, count int) *Fifo {
	return &Fifo{
		namedPipe: namedPipe,
		size:      size,
		count:     count,
	}
}
