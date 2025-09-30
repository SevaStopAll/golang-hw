package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type MyTelnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (m *MyTelnetClient) Connect() error {
	conn, e := net.DialTimeout("tcp", m.address, m.timeout)
	if e != nil {
		return fmt.Errorf("failed to dial: %w", e)
	}

	m.conn = conn
	return nil
}

func (m *MyTelnetClient) Close() error {
	return m.conn.Close()
}

func (m *MyTelnetClient) Send() error {
	_, e := io.Copy(m.conn, m.in)
	return e
}

func (m *MyTelnetClient) Receive() error {
	_, e := io.Copy(m.out, m.conn)
	return e
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &MyTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
