package main

import (
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *Client) Connect() (err error) {
	log.Printf("Connect to %s\n", c.address)
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	log.Printf("...Connected to %s\n", c.address)
	return nil
}

func (c *Client) Close() error {
	log.Printf("Close\n")
	err := c.conn.Close()
	if err != nil {
		log.Printf("Close with error %s\n", err.Error())
		return err
	}
	log.Printf("...Connection was closed by peer\n")
	return nil
}

func (c *Client) Send() (err error) {
	_, err = io.Copy(c.conn, c.in)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Receive() (err error) {
	_, err = io.Copy(c.out, c.conn)
	if err != nil {
		return err
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
