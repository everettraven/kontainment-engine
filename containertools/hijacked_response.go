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

var _ HijackedResponse = &hijackedResponse{}

// hijackedResponse is a generic implementation of the HijackedResponse interface
type hijackedResponse struct {
	conn   net.Conn
	reader *bufio.Reader
}
