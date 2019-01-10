package lifx

func booltouint16(b bool) uint16 {
	if b {
		return 0xffff
	}
	return 0
}
