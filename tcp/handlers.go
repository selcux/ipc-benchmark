package tcp

import (
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/selcux/ipc-benchmark/util"
)

const bufSize = 1024

func serverHandler(args ServerArgs) {
	for i := 0; i < args.rotationCount; i++ {
		log.Println("Server receiving data")
		data, err := receive(args.conn, args.maxDataSize)
		if err != nil {
			panic(errors.Wrap(err, "could not receive data"))
		}

		log.Println("Server sending back data")
		_, err = args.conn.Write(data)
		if err != nil {
			panic(errors.Wrap(err, "could not send data"))
		}
	}
}

func clientHandler(args ClientArgs) {
	for i := 0; i < args.rotationCount; i++ {
		log.Println("Generating random data...")
		data, err := util.GenRandomBytes(args.maxDataSize)
		if err != nil {
			args.dataCh <- ResultArgs{err: err}
			return
		}

		log.Println("Client sending data")
		_, err = args.conn.Write(data)
		if err != nil {
			args.dataCh <- ResultArgs{
				err: errors.Wrap(err, "could not send data"),
			}
			return
		}

		log.Println("Client receiving data")
		recvData, err := receive(args.conn, args.maxDataSize)
		if err != nil {
			args.dataCh <- ResultArgs{
				err: errors.Wrap(err, "could not receive data"),
			}
			return
		}

		log.Println("Client data received")
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
