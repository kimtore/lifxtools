package virtbulb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/dorkowscy/lyslix/lifx"
)

type PacketStreamer struct {
	C      chan *lifx.Packet
	Errors chan error
	conn   net.PacketConn
	ctx    context.Context
}

const channelSize = 16
const maxPacketSize = 256

var ErrInvalidLifxPacket = errors.New("unable to decode LIFX packet")

func NewPacketStreamer(ctx context.Context, conn net.PacketConn) *PacketStreamer {
	p := &PacketStreamer{
		ctx:    ctx,
		conn:   conn,
		C:      make(chan *lifx.Packet, channelSize),
		Errors: make(chan error, channelSize),
	}
	go p.stream()
	return p
}

func (p *PacketStreamer) stream() {
	packetBuffer := make([]byte, maxPacketSize)

	defer close(p.C)
	defer close(p.Errors)

	for p.ctx.Err() == nil {
		bytesRead, _, err := p.conn.ReadFrom(packetBuffer)
		if err != nil {
			p.Errors <- fmt.Errorf("udp socket received fatal error: %w", err)
			return
		}

		r := bytes.NewReader(packetBuffer[:bytesRead])
		packet, err := lifx.DecodePacket(r)
		if err != nil {
			p.Errors <- ErrInvalidLifxPacket
			continue
		}

		p.C <- packet
	}
}
