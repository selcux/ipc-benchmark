package tcp

import (
	"log"
	"net"

	"github.com/pkg/errors"
)

type ResultArgs struct {
	messageLen int
	err        error
}

type ClientArgs struct {
	conn          net.Conn
	maxDataSize   int
	rotationCount int
	dataCh        chan ResultArgs
}

type Client struct {
	handler func(args ClientArgs)
	args    ClientArgs
}

func (c *Client) SetHandler(f func(args ClientArgs)) {
	c.handler = f
}

func (c *Client) Connect(address string) error {
	log.Printf("Dialing... %s", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "could not connect to %s", address)
	}

	if c.handler == nil {
		return nil
	}

	c.args.conn = conn
	c.handler(c.args)

	return nil
}

func NewClient(args ClientArgs) *Client {
	return &Client{args: args}
}
