package tcp

import (
	"net"

	"github.com/pkg/errors"
)

type ServerArgs struct {
	conn          net.Conn
	maxDataSize   int
	rotationCount int
	errCh         chan error
}

type Server struct {
	handler  func(args ServerArgs)
	args     ServerArgs
	listener net.Listener
}

func (s *Server) SetHandler(f func(args ServerArgs)) {
	s.handler = f
}

func (s *Server) Listen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "unable to bind address %s", address)
	}

	s.listener = listener

	return nil
}

func (s *Server) Serve() error {
	conn, err := s.listener.Accept()
	if err != nil {
		return errors.Wrap(err, "could not accept the connection")
	}
	defer func() { conn.Close() }()

	if s.handler == nil {
		return nil
	}

	s.args.conn = conn
	s.handler(s.args)

	return nil
}

func (s *Server) ErrCh() chan error {
	return s.args.errCh
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func NewServer(args ServerArgs) *Server {
	return &Server{args: args}
}
