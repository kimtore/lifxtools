package lifx

// From https://lan.developer.lifx.com/docs/changing-a-device

import (
	"io"
	"strings"
)

// SetLabel - Packet 24
//
// This packet lets you set the label on the device. The label is a string you assign to the device and will be displayed as the name of the device in the LIFX mobile apps.
//
// Will return one StateLabel (25) message.
type SetLabelMessage struct {
	Label string
}

func DecodeSetLabelMessage(r io.Reader) (*SetLabelMessage, error) {
	buf := make([]byte, LabelLength)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	return &SetLabelMessage{
		Label: strings.Split(string(buf), "\x00")[0],
	}, nil
}

func (m *SetLabelMessage) Write(w io.Writer) error {
	_, err := WriteString(w, m.Label, LabelLength)
	return err
}

func (m *SetLabelMessage) Len() int {
	return 32
}

func (m *SetLabelMessage) Type() uint16 {
	return MsgTypeSetLabel
}
