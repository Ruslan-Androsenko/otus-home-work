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

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Telnet struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect: %w", err)
	}

	t.conn = conn
	return nil
}

func (t *Telnet) Close() error {
	err := t.in.Close()
	if err != nil {
		return fmt.Errorf("cannot closing input: %w", err)
	}

	if t.conn != nil {
		err = t.conn.Close()
		if err != nil {
			return fmt.Errorf("cannot closing connect: %w", err)
		}
	}

	return nil
}

func (t *Telnet) Send() error {
	return readAndWrite(t.in, t.conn)
}

func (t *Telnet) Receive() error {
	return readAndWrite(t.conn, t.out)
}

// Прочитать сообщение из входного источника, и записать его в выходной.
func readAndWrite(in io.ReadCloser, out io.Writer) error {
	buffer := make([]byte, 1024)
	n, errRead := in.Read(buffer)
	if hasBeenError(errRead) {
		return fmt.Errorf("cannot read from input: %w", errRead)
	}

	_, errWrite := out.Write(buffer[:n])
	if errWrite != nil {
		return fmt.Errorf("cannot write to output: %w", errWrite)
	}

	return errRead
}
