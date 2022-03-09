package ble

// OGF列表
const (
	HCI_CMD_OGF_LE_CONTROLLER_CMD = 0x08
)

// LE CONTROLLER COMMANDS OCF列表
const (
	HCI_LE_EXTENDED_CREATE_CONNECTION = 0x01
)

var BleParserMap map[int]map[int]string {

}

type BleParser interface {
	Parse(hciPayloadBuf []byte)
	Print()
}

type HCI_LE_EXTENDED_CREATE_CONNECTION struct {
	Handle       uint16
	PbFlag       uint8
	BcFlag       uint8
	DataTotalLen uint16
	Data         []byte
}