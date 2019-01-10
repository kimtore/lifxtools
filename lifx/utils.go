package lifx

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

func MACAdressToFrameAddress(mac []byte) uint64 {
	return 0 |
		(uint64(mac[0]) << 0) |
		(uint64(mac[1]) << 8) |
		(uint64(mac[2]) << 16) |
		(uint64(mac[3]) << 24) |
		(uint64(mac[4]) << 32) |
		(uint64(mac[5]) << 40)
}
