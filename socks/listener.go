package socks

import "net"

type Listener struct {
	tcpListener *net.TCPListener
	localAddr   *net.TCPAddr

	acceptChan chan *Conn
}

func ListenSocks(addr string) (*Listener, error) {
	l := new(Listener)

	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil { return nil, err }
	l.localAddr = laddr

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil { return nil, err }
	l.tcpListener = listener

	l.acceptChan = make(chan *Conn)

	go l.acceptLoop()

	return l, nil
}

func (l *Listener) acceptLoop() {
	for {
		conn, err := l.tcpListener.AcceptTCP()
		if err != nil { return }
		c, err := NewConn(l, conn)
		if err != nil { continue }
		l.acceptChan <- c
	}
}
func (l *Listener) Accept() (*Conn, error) {
	return <- l.acceptChan, nil
}

func (l *Listener) Close() {
	l.tcpListener.Close()
}