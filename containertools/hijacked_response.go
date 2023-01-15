package containertools

import (
	"bufio"
	"net"
)

// HijackedResponse represents the response object that is returned
// when attaching to an exec process within a container
type HijackedResponse interface {
	Conn() net.Conn
	Reader() *bufio.Reader
}

type HijackedResponseOption func(*hijackedResponse)

func WithConn(conn net.Conn) HijackedResponseOption {
	return func(hr *hijackedResponse) {
		hr.conn = conn
	}
}

func WithReader(reader *bufio.Reader) HijackedResponseOption {
	return func(hr *hijackedResponse) {
		hr.reader = reader
	}
}

func NewHijackedResponse(opts ...HijackedResponseOption) HijackedResponse {
	hr := &hijackedResponse{}

	for _, opt := range opts {
		opt(hr)
	}

	return hr
}

var _ HijackedResponse = &hijackedResponse{}

// hijackedResponse is a generic implementation of the HijackedResponse interface
type hijackedResponse struct {
	conn   net.Conn
	reader *bufio.Reader
}

func (hr *hijackedResponse) Conn() net.Conn {
	return hr.conn
}

func (hr *hijackedResponse) Reader() *bufio.Reader {
	return hr.reader
}
