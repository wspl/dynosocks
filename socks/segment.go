package socks

import (
	"net"
	"bytes"
	"github.com/pkg/errors"
	"encoding/binary"
)

var RequestSegment = []byte{0x05, 0x01, 0x00}
var ResponseSegment = []byte{0x05, 0x00}

func ParseAddrSegment(buf []byte) (*net.TCPAddr, error) {
	if !bytes.HasPrefix(buf, RequestSegment) {
		return nil, errors.New("invalid addr segment")
	}
	switch buf[3] {
	case 0x01:
		return &net.TCPAddr{
			IP: net.IP(buf[4:8]),
			Port: int(binary.BigEndian.Uint16(buf[8:10])),
		}, nil
	case 0x03:
		ip, err := net.LookupIP(string(buf[5:5+buf[4]]))
		if err != nil { return nil, err }
		return &net.TCPAddr{
			IP: ip[0],
			Port: int(binary.BigEndian.Uint16(buf[5+buf[4]:7+buf[4]])),
		}, nil
	}
	return nil, errors.New("unsupport addr type")
}

func CreateAddrSegment(addr *net.TCPAddr) []byte {
	buf := new(bytes.Buffer)
	buf.Write(ResponseSegment)
	buf.Write([]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00})
	portBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(portBuf, uint16(addr.Port))
	buf.Write(portBuf)
	return buf.Bytes()
}