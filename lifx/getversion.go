package lifx

type GetVersionMessage struct {
	emptyMessage
}

func (m *GetVersionMessage) Type() uint16 {
	return MsgTypeGetVersion
}

type StateVersionMessage struct {
	emptyMessage
	Vendor    uint32
	Product   uint32
	Reserved6 [4]byte
}

func (m *StateVersionMessage) Type() uint16 {
	return MsgTypeStateVersion
}
