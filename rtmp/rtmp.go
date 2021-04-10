package rtmp

import (
	"fmt"
	"net"
)

type RTMP struct {
	conn      rtmpConn
	message   map[uint32]Message //message stream id
	chunkSize uint32
}

func NewRTMP(conn net.Conn) (rtmp *RTMP) {
	return &RTMP{
		conn: rtmpConn{
			conn: conn,
		},
	}
}

func (rtmp *RTMP) Handler() {
	err := NewServer().Handshake(rtmp.conn)
	if err != nil {
		fmt.Println("handshake error:", err)
		return
	}

	for {
		chunk, err := ParseChunk(rtmp.conn)
		if err != nil {
			fmt.Println("NewChunk error:", err)
			continue
		}

		var message Message
		var ok bool
		if message, ok = rtmp.message[chunk.MessageStreamID]; !ok {
			message, err = ParseMessage(rtmp, chunk)
			if err != nil {
				fmt.Println("NewMessage error:", err)
				continue
			}
			fmt.Printf("message:%+v\n", message)
			//TODO save MessageStreamID to rtmp.message
		} else {
			message.Combine(chunk)
		}
		message.Do(rtmp.conn)
		break
	}
}
