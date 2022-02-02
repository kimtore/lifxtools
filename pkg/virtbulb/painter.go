package virtbulb

import (
	"fmt"

	"github.com/dorkowscy/lyslix/lifx"
	log "github.com/sirupsen/logrus"
)

type ErrUnsupportedPacket error

func PaintPacket(cv Canvas, packet *lifx.Packet) error {
	switch payload := packet.Payload.(type) {

	case *lifx.SetColorMessage:
		cv.Draw(-1, -1, payload.Color)

	case *lifx.SetColorZonesMessage:
		log.Debugf("Draw(%d, %d, %v)", int(payload.StartIndex), int(payload.EndIndex), payload.Color)
		cv.Draw(int(payload.StartIndex), int(payload.EndIndex), payload.Color)

	default:
		return ErrUnsupportedPacket(fmt.Errorf("unsupported packet type %T: %v", payload, payload))
	}

	return nil
}
