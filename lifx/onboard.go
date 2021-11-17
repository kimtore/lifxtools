package lifx

import (
	"encoding/binary"
	"io"
)

//goland:noinspection ALL
const (
	Security_Open           uint8 = 0x01
	Security_WEP_PSK        uint8 = 0x02
	Security_WPA_TKIP_PSK   uint8 = 0x03
	Security_WPA_AES_PSK    uint8 = 0x04
	Security_WPA2_AES_PSK   uint8 = 0x05
	Security_WPA2_TKIP_PSK  uint8 = 0x06
	Security_WPA2_MIXED_PSK uint8 = 0x07
)

// SetAccessPoint - Packet 305
//
// Sets WiFi configuration. Undocumented at LIFX. Figures.
//
// Reverse engineered using https://github.com/tserong/lifx-hacks/blob/master/onboard.py
type SetAccessPointMessage struct {
	Reserved1 byte // Really don't know what this is. Set to 0x02 from the example.
	SSID      string
	PSK       string
	Security  byte
}

func DecodeSetAccessPointMessage(r io.Reader) (*SetAccessPointMessage, error) {
	var err error
	msg := &SetAccessPointMessage{}

	err = binary.Read(r, binary.LittleEndian, &msg.Reserved1)
	if err != nil {
		return nil, err
	}

	msg.SSID, err = ReadString(r, WifiSSIDLength)
	if err != nil {
		return nil, err
	}

	msg.PSK, err = ReadString(r, WifiPasswordLength)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.LittleEndian, &msg.Security)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *SetAccessPointMessage) Write(w io.Writer) error {
	var err error

	err = binary.Write(w, binary.LittleEndian, m.Reserved1)
	if err != nil {
		return err
	}

	_, err = WriteString(w, m.SSID, WifiSSIDLength)
	if err != nil {
		return err
	}

	_, err = WriteString(w, m.PSK, WifiPasswordLength)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, m.Security)
}

func (m *SetAccessPointMessage) Len() int {
	return WifiSSIDLength + WifiPasswordLength + 2
}

func (m *SetAccessPointMessage) Type() uint16 {
	return MsgTypeSetAccessPoint
}
