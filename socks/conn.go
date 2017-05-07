package socks

import (
	"net"
	"bytes"
)

type SetupConn func(*net.TCPAddr) (interface{}, error)

func NewConn(listener *Listener, conn *net.TCPConn) (*Conn, error) {
	c := new(Conn)
	c.listener = listener
	c.tcpConn = conn

	c.allReady = make(chan bool)

	go c.localReadyRead()

	return c, nil
}

type Conn struct {
	listener *Listener
	tcpConn  *net.TCPConn
	Target   *net.TCPAddr

	allReady       chan bool
}

func (c *Conn) localReadyRead() {
	counter := 0
	for counter < 2{
		buf := make([]byte, 1024)
		bufSize, err := c.tcpConn.Read(buf)
		if err != nil { continue }
		buf = buf[:bufSize]
		c.handleLocalRead(buf)
		counter++
	}
}

func (c *Conn) Ready() {
	<- c.allReady
}

func (c *Conn) handleLocalRead(buf []byte) {
	if bytes.Equal(buf, RequestSegment) {
		c.tcpConn.Write(ResponseSegment)
	} else if c.Target == nil {
		target, err := ParseAddrSegment(buf)
		if err != nil { return }

		c.tcpConn.Write(CreateAddrSegment(c.listener.localAddr))
		c.Target = target
		c.allReady <- true
	}
}

func (c *Conn) Read(buf []byte) (int, error) {
	return c.tcpConn.Read(buf)
}

func (c *Conn) Write(buf []byte) (int, error) {
	return c.tcpConn.Write(buf)
}

func (c *Conn) Close() {
	c.tcpConn.Close()
}