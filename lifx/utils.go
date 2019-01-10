package lifx

const (
	payloadPadding = " "
)

func booltouint8(b bool) uint8 {
	if b {
		return 0xff
	}
	return 0
}

func booltouint16(b bool) uint16 {
	if b {
		return 0xffff
	}
	return 0
}
