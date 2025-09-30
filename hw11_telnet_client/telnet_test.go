package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type params struct {
	messages2Send    []string
	messages2Receive []string
}

type pipes struct {
	toServerR   *io.PipeReader
	fromClientW *io.PipeWriter
	toClientR   *io.PipeReader
	fromServerW *io.PipeWriter
}

type telnetTestSuite struct {
	suite.Suite
	params   params
	pipes    pipes
	listener net.Listener
	client   TelnetClient
}

func (s *telnetTestSuite) SetupTest() {
	s.pipes = pipes{}

	s.pipes.toServerR, s.pipes.fromClientW = io.Pipe()
	s.pipes.toClientR, s.pipes.fromServerW = io.Pipe()

	var err error
	s.listener, err = net.Listen("tcp", "127.0.0.1:")
	s.NoError(err)

	s.client = NewTelnetClient(s.listener.Addr().String(), 10*time.Second, s.pipes.toServerR, s.pipes.fromServerW)
	s.NoError(s.client.Connect())
}

func (s *telnetTestSuite) TearDownTest() {
	s.NoError(s.listener.Close())
	s.NoError(s.client.Close())
}

func (s *telnetTestSuite) TestBasicCase() {
	s.params = params{
		messages2Send:    []string{"hello"},
		messages2Receive: []string{"world"},
	}

	s.RunTest()
}

func (s *telnetTestSuite) TestSeveralMessagesCase() {
	s.params = params{
		messages2Send:    []string{"hello", "cruel", "world"},
		messages2Receive: []string{"bye-bye", "calm", "heaven"},
	}

	s.RunTest()
}

func (s *telnetTestSuite) RunTest() {
	conn, err := s.listener.Accept()
	s.NoError(err)
	s.NotNil(conn)
	defer func() { s.NoError(conn.Close()) }()

	var wg sync.WaitGroup
	wg.Add(4)

	go func() { // send messages from client
		defer wg.Done()
		for _, m := range s.params.messages2Send {
			s.pipes.fromClientW.Write([]byte(m + "\n"))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() { // read messages from client
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for i := 0; i < len(s.params.messages2Send); i++ {
			scanner.Scan()
			s.Equal(s.params.messages2Send[i], scanner.Text())
		}
		s.NoError(scanner.Err())
	}()

	go func() { // send messages from server
		defer wg.Done()
		for _, m := range s.params.messages2Receive {
			conn.Write([]byte(m + "\n"))
			time.Sleep(time.Second)
		}
	}()

	go func() { // read messages from server
		defer wg.Done()
		scanner := bufio.NewScanner(s.pipes.toClientR)
		for i := 0; i < len(s.params.messages2Receive); i++ {
			scanner.Scan()
			s.Equal(s.params.messages2Receive[i], scanner.Text())
		}
		s.NoError(scanner.Err())
	}()

	go func() {
		s.client.Send()
	}()

	go func() {
		s.client.Receive()
	}()

	wg.Wait()
}

func TestTelnetClient(t *testing.T) {
	suite.Run(t, new(telnetTestSuite))
}

func TestTelnetClientTimeout(t *testing.T) {
	client := NewTelnetClient("1.1.1.1:1234", 3*time.Second, os.Stdin, os.Stdout)
	var wg sync.WaitGroup
	var completed int32
	var err error
	wg.Add(1)
	go func(c *int32, e *error) {
		defer wg.Done()
		*e = client.Connect()
		atomic.AddInt32(c, 1)
	}(&completed, &err)

	require.Eventuallyf(t, func() bool {
		return !atomic.CompareAndSwapInt32(&completed, 1, 0)
	}, 3*time.Second, 2*time.Second, "")

	require.Eventuallyf(t, func() bool {
		return atomic.CompareAndSwapInt32(&completed, 1, 1)
	}, 3*time.Second, 2*time.Second, "")

	wg.Wait()

	var result net.Error
	require.Truef(t, errors.As(err, &result), "Expected net.Error, but got %v", err)
	require.Truef(t, result.Timeout(), "Expected timeout error but got %s", result.Error())
}
