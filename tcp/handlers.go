package tcp

import (
	"io"
	"net"

	"github.com/pkg/errors"
	"github.com/selcux/ipc-benchmark/util"
)

const bufSize = 1024

func serverLatencyHandler(args ServerArgs) {
	for i := 0; i < args.rotationCount; i++ {
		data, err := receive(args.conn, args.maxDataSize)
		if err != nil {
			go func() {
				args.errCh <- errors.Wrap(err, "could not receive data")
			}()
			return
		}

		_, err = args.conn.Write(data)
		if err != nil {
			go func() {
				args.errCh <- errors.Wrap(err, "could not send data")
			}()
			return
		}
	}

	close(args.errCh)
}

func clientLatencyHandler(args ClientArgs) {
	for i := 0; i < args.rotationCount; i++ {
		data, err := util.GenRandomBytes(args.maxDataSize)
		if err != nil {
			go func() {
				args.dataCh <- ResultArgs{err: err}
			}()
			return
		}

		_, err = args.conn.Write(data)
		if err != nil {
			go func() {
				args.dataCh <- ResultArgs{
					err: errors.Wrap(err, "could not send data"),
				}
			}()
			return
		}

		recvData, err := receive(args.conn, args.maxDataSize)
		if err != nil {
			go func() {
				args.dataCh <- ResultArgs{
					err: errors.Wrap(err, "could not receive data"),
				}
			}()
			return
		}

		go func() {
			args.dataCh <- ResultArgs{
				messageLen: len(recvData),
				err:        nil,
			}
		}()
	}
}

func serverThroughputHandler(args ServerArgs) {
	data, err := util.GenRandomBytes(args.maxDataSize)
	if err != nil {
		go func() {
			args.errCh <- errors.Wrap(err, "unable to create random byte array")
		}()
		return
	}

	for i := 0; i < args.rotationCount; i++ {
		_, err = args.conn.Write(data)
		if err != nil {
			go func() {
				args.errCh <- errors.Wrap(err, "could not send data")
			}()
			return
		}
	}

	close(args.errCh)
}

func clientThroughputHandler(args ClientArgs) {
	for i := 0; i < args.rotationCount; i++ {
		recvData, err := receive(args.conn, args.maxDataSize)
		if err != nil {
			go func() {
				args.dataCh <- ResultArgs{
					err: errors.Wrap(err, "could not receive data"),
				}
			}()
			return
		}

		go func() {
			args.dataCh <- ResultArgs{
				messageLen: len(recvData),
				err:        nil,
			}
		}()
	}
}

func receive(conn net.Conn, maxSize int) ([]byte, error) {
	buf := make([]byte, bufSize)
	data := make([]byte, 0)
	totalSize := 0

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, errors.Wrap(err, "unable to read data")
			}

			break
		}

		totalSize += n
		data = append(data, buf[:n]...)
		if totalSize >= maxSize {
			break
		}
	}

	return data, nil
}
