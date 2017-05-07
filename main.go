package main

import (
	"github.com/armon/go-socks5"
	"net"
	"io"
	"time"
	"github.com/xtaci/kcp-go"
)

func main() {
	go server()
	time.Sleep(1000)
	go cli()

	<- make(chan bool)
}

func setKCP(conn *kcp.UDPSession) {
	conn.SetStreamMode(true)
	conn.SetWriteDelay(true)
	conn.SetNoDelay(1, 10, 2, 1)
	conn.SetWindowSize(128, 512)
	conn.SetMtu(1350)
	conn.SetACKNoDelay(true)
}

func getBlockCrypt() kcp.BlockCrypt {
	bc, _ := kcp.NewXTEABlockCrypt([]byte{1,2,3,4,5,6,8})
	return bc
}

func cli() {
	laddr, _ := net.ResolveTCPAddr("tcp", ":1087")
	tcpListener, _ := net.ListenTCP("tcp", laddr)

	for i := 0; i< 8; i++ {
		go func() {
			for {
				tcpConn, _ := tcpListener.Accept()
				conn, _ := kcp.DialWithOptions("batman.vecsight.com:9980", getBlockCrypt(), 10, 3)

				setKCP(conn)

				go xPipe(conn, tcpConn)
			}
		}()
	}
}

func server() {
	go func() {
		conf := &socks5.Config{}
		ss, _ := socks5.New(conf)
		ss.ListenAndServe("tcp", "127.0.0.1:9527")
	}()

	listener, _ := kcp.ListenWithOptions(":9980", getBlockCrypt(), 10, 3)
	//listener, _ := net.ListenTCP("tcp", zzbl_a)

	for i := 0; i< 8; i++ {
		go func() {
			for {
				conn, _ := listener.AcceptKCP()
				setKCP(conn)
				saddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9527")
				laddr, _ := net.ResolveTCPAddr("tcp", ":0")
				socksCli, _ := net.DialTCP("tcp", laddr, saddr)

				println("remote accept", conn.RemoteAddr().String())
				go xPipe(conn, socksCli)
			}
		}()
	}
}

func xPipe(s1, s2 io.ReadWriteCloser) {
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