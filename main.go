package main

import (
	"github.com/armon/go-socks5"
	//kcp "github.com/xtaci/kcp-go"
	"net"
	"time"
	"github.com/xtaci/kcp-go"
	"io"
)

func main() {
	go server()
	time.Sleep(1000)
	go cli()

	<- make(chan bool)
}

func cli() {
	laddr, _ := net.ResolveTCPAddr("tcp", ":1087")
	tcpListener, _ := net.ListenTCP("tcp", laddr)

	for {
		tcpConn, _ := tcpListener.Accept()
		conn, _ := kcp.DialWithOptions("127.0.0.1:9980", nil, 10, 3)
		go pipe(conn, tcpConn)
	}
}

func server() {
	go func() {
		conf := &socks5.Config{}
		ss, _ := socks5.New(conf)
		ss.ListenAndServe("tcp", "127.0.0.1:9527")
	}()

	listener, _ := kcp.ListenWithOptions(":9980", nil, 10, 3)
	//listener, _ := net.ListenTCP("tcp", zzbl_a)
	for {
		conn, _ := listener.AcceptKCP()
		saddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9527")
		laddr, _ := net.ResolveTCPAddr("tcp", ":0")
		socksCli, _ := net.DialTCP("tcp", laddr, saddr)

		println("remote accept", conn.RemoteAddr().String())
		go pipe(conn, socksCli)
	}
}

func pipe(s1, s2 io.ReadWriteCloser) {
	defer s1.Close()
	defer s2.Close()

	s1Close := make(chan struct{})
	go func() { io.Copy(s1, s2); close(s1Close) }()

	s2Close := make(chan struct{})
	go func() { io.Copy(s2, s1); close(s2Close) }()

	select {
	case <-s1Close:
	case <-s2Close:
	}
}