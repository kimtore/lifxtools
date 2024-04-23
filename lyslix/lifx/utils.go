package lifx

import (
	"bytes"
	"io"
	"strings"
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

func MACAdressToFrameAddress(mac []byte) uint64 {
	return 0 |
		(uint64(mac[0]) << 0) |
		(uint64(mac[1]) << 8) |
		(uint64(mac[2]) << 16) |
		(uint64(mac[3]) << 24) |
		(uint64(mac[4]) << 32) |
		(uint64(mac[5]) << 40)
}

func WriteString(w io.Writer, s string, length int) (int64, error) {
	r := strings.NewReader(s)
	buf := &bytes.Buffer{}
	_, err := io.CopyN(buf, r, int64(length)-1)
	if err != nil && err != io.EOF {
		return int64(buf.Len()), err
	}
	for buf.Len() < length {
		buf.WriteByte(0)
	}
	return io.Copy(w, buf)
}

func ReadString(r io.Reader, maxLength int) (string, error) {
	buf := make([]byte, maxLength)
	_, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	return strings.Split(string(buf), "\x00")[0], nil
}
