package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	adr, timeout := handleParams()

	client := NewTelnetClient(adr, *timeout, os.Stdin, os.Stdout)
	if e := client.Connect(); e != nil {
		log.Panic(fmt.Errorf("failed to connect: %w", e))
	} else {
		fmt.Fprintf(os.Stderr, "...Connected to %s\n", adr)
	}
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go sending(client, stop)
	go receiving(client, stop)

	<-ctx.Done()
}

func sending(c TelnetClient, stop context.CancelFunc) {
	defer stop()
	if e := c.Send(); e == nil {
		fmt.Fprintln(os.Stderr, "...EOF")
	} else {
		fmt.Fprintf(os.Stderr, "sending error: %s\n", e.Error())
	}
}

func receiving(c TelnetClient, stop context.CancelFunc) {
	defer stop()

	if e := c.Receive(); e == nil {
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	} else {
		fmt.Fprintf(os.Stderr, "receiving error: %s", e.Error())
	}
}

func handleParams() (address string, timeout *time.Duration) {
	timeout = flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Panic("Usage expected: go-telnet [--timeout 10s] <HOST> <PORT>")
	}

	if _, e := strconv.ParseInt(flag.Arg(1), 10, 16); e != nil {
		log.Panic("Connection port isn't numeric")
	}

	return net.JoinHostPort(flag.Arg(0), flag.Arg(1)), timeout
}
