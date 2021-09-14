package fifo

import (
	"bufio"
	"fmt"
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

func (f *Fifo) Produce() (int64, error) {
	fmt.Println("Opening named pipe for writing")
	file, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return 0, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file.Close()
	fmt.Println("Writing")

	var count int64 = 0
	for i := 0; i < f.count; i++ {
		scount, err := f.send(file)
		if err != nil {
			return 0, errors.Wrap(err, "could not receive data")
		}

		count += scount
	}

	return count, nil
}

func (f *Fifo) Consume() (int64, error) {
	// Open named pipe for reading
	fmt.Println("Opening named pipe for reading")
	file, err := os.OpenFile(f.namedPipe, os.O_RDONLY, 0600)
	if err != nil {
		return 0, errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file.Close()
	fmt.Println("Reading")
	fmt.Println("Waiting for someone to write something")

	var count int64 = 0
	for i := 0; i < f.count; i++ {
		data, err := f.receive(file)
		if err != nil {
			return 0, errors.Wrap(err, "unable to receive data")
		}

		count += int64(len(data))
	}

	return count, nil
}

func (f *Fifo) PingPong() error {
	/*
		errCh := make(chan error)

		go func(e chan error) {
			e <- f.ping()
		}(errCh)

		errPong := f.pong()
		if errPong != nil {
			return errors.Wrap(errPong, "could not pong")
		}

		errPing := <-errCh
		if errPing != nil {
			return errors.Wrap(errPing, "could not ping")
		}
	*/

	file1, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file1.Close()

	file2, err := os.OpenFile(f.namedPipe, os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s", f.namedPipe)
	}
	defer file2.Close()

	for i := 0; i < f.count; i++ {
		err = f.pingPong(file1, file2)
		if err != nil {
			return errors.Wrap(err, "communication error")
		}
	}

	return nil
}

func (f *Fifo) pingPong(file1, file2 *os.File) error {
	errCh := make(chan error)
	defer close(errCh)

	go func(e chan error) {
		_, err := f.send(file1)
		e <- err
	}(errCh)

	_, err := f.receive(file2)
	if err != nil {
		return errors.Wrap(err, "could not received message - pong")
	}
	err = <-errCh
	if err != nil {
		return errors.Wrap(err, "could not send message - ping")
	}

	go func(e chan error) {
		_, err := f.send(file2)
		e <- err
	}(errCh)

	_, err = f.receive(file1)
	if err != nil {
		return errors.Wrap(err, "could not received message - ping")
	}
	err = <-errCh
	if err != nil {
		return errors.Wrap(err, "could not send message - ping")
	}

	return nil
}

func (f *Fifo) send(file *os.File) (int64, error) {
	data, err := util.GenRandomBytes(f.size)
	if err != nil {
		return 0, errors.Wrap(err, "unable to create random byte array")
	}

	file.Write(data)

	return int64(len(data)), nil
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

			return nil, errors.Wrapf(err, "error reading %s", file)
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
