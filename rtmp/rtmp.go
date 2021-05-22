package rtmp

import (
	"fmt"
	"io"
	"net"

	"github.com/SmartBrave/utils/easyio"
)

type RTMP struct {
	conn             easyio.EasyReadWriter
	lastChunk        map[uint32]*Chunk //csid
	peerMaxChunkSize int
	ownMaxChunkSize  int
}

func NewRTMP(conn net.Conn) (rtmp *RTMP) {
	return &RTMP{
		conn: rtmpConn{
			Conn: conn,
		},
		lastChunk:        make(map[uint32]*Chunk),
		peerMaxChunkSize: 128,
		ownMaxChunkSize:  128,
	}
}

func (rtmp *RTMP) Handler() {
	err := NewServer().Handshake(rtmp)
	if err != nil {
		fmt.Println("handshake error:", err)
		return
	}

	for {
		fmt.Println("-----------------------------------")
		err = ParseMessage(rtmp)
		if err == io.EOF {
			fmt.Println("disconnect")
			break
		}
		if err != nil {
			fmt.Println("ParseMessage error:", err)
			continue
		}
	}
}
