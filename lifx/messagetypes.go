package lifx

const (
	MsgTypeUnknown              = 0
	MsgTypeSetLabel             = 24
	MsgTypeStateLabel           = 25
	MsgTypeGetVersion           = 32
	MsgTypeStateVersion         = 33
	MsgTypeSetColorMessage      = 102
	MsgTypeSetColorZonesMessage = 501
	MsgTypeSetAccessPoint       = 305

	ProtocolNumber     = 1024
	LabelLength        = 32
	WifiSSIDLength     = 32
	WifiPasswordLength = 64
)
