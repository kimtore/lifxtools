package lifx

import (
	"net"
	"time"
)

type client struct {
	addr   string
	source uint32
	target uint64 //todo
	conn   net.Conn
}

type Client interface {
	Close() error
	SetColor(color HBSK, fadeTime time.Duration) error
	SetColorZones(color HBSK, start, end uint8, fadeTime time.Duration) error
}

func NewClient(addr string) (Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return &client{
		addr:   addr,
		conn:   conn,
		source: 0xbeef,
	}, nil
}

func (c *client) send(payload Payload) error {
	packet := NewPacket(payload)
	return packet.Write(c.conn)
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) SetColor(color HBSK, fadeTime time.Duration) error {
	return c.send(&SetColorMessage{
		Color:    color,
		Duration: uint32(fadeTime * time.Millisecond),
	})
}

func (c *client) SetColorZones(color HBSK, start, end uint8, fadeTime time.Duration) error {
	return c.send(&SetColorZonesMessage{
		StartIndex: start,
		EndIndex:   end,
		Color:      color,
		Duration:   uint32(fadeTime * time.Millisecond),
		Apply:      MultiZoneApply,
	})
}
