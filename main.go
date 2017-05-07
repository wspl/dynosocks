package main

import (
	"dynosocks/socks"
	"net"

	"github.com/armon/go-socks5"
)

func main() {
	conf := &socks5.Config{}
	ss, _ := socks5.New(conf)
	ss.ListenAndServe("unix", "/tmp/dynosocks_socks5.sock")
}

func mfasfdsaain() {
	go func() {
		laddr, _ := net.ResolveTCPAddr("tcp", ":0")
		server, _ := socks.ListenSocks(":1087")
		//print(err.Error())
		go func() {
			for {
				conn, _ := server.Accept()
				go func() {
					conn.Ready()
					println("accept connection to:", conn.Target.String())
					tcpConn, err := net.DialTCP("tcp", laddr, conn.Target)
					if err != nil {
						conn.Close()
						return
					}

					go func() {
						buf := make([]byte, 4096)
						for {
							length, err := conn.Read(buf)
							if length == 0 || err != nil {
								tcpConn.Close()
								conn.Close()
								return
							}
							tcpConn.Write(buf[:length])
						}
					}()

					go func() {
						buf := make([]byte, 4096)
						for {
							length, err := tcpConn.Read(buf)
							if length == 0 || err != nil {
								tcpConn.Close()
								conn.Close()
								return
							}
							conn.Write(buf[:length])
						}
					}()
				}()
			}
		}()
	}()

	<- make(chan bool)
}