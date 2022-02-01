package virtbulb

import (
	"fmt"

	"github.com/dorkowscy/lyslix/lifx"
)

type ErrUnsupportedPacket error

func PaintPacket(cv Canvas, packet *lifx.Packet) error {
	switch payload := packet.Payload.(type) {

	case *lifx.SetColorMessage:
		cv.Draw(-1, -1, payload.Color)

	case *lifx.SetColorZonesMessage:
		cv.Draw(int(payload.StartIndex)-1, int(payload.EndIndex)-1, payload.Color)

	default:
		return ErrUnsupportedPacket(fmt.Errorf("unsupported packet type %T: %v", payload, payload))
	}

	return nil
}
