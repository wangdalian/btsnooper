// hci evt处理

package hci

import (
	"encoding/binary"
)

const (
	HCI_CMD_OGF_LE_CONTROLLER_CMD = 0x08
)

const (
	HCI_LE_EXTENDED_CREATE_CONNECTION = 0x0043
	// ...
)

// BLUETOOTH SPECIFICATION Version 4.2 [Vol 2, Part E] 7.8.12 LE Create Connection Command
// HCI_LE_Create_Connection

// 连接参数，针对不同的连接模式配置
type ConnectionInitialting struct {
	ScanInterval      uint16
	ScanWindow        uint16
	ConnIntervalMin   uint16
	ConnIntervalMax   uint16
	ConnLatency       uint16
	SupervisionTimout uint16
	MinimumCeLength   uint16
	MaximumCeLength   uint16
}

// BLUETOOTH SPECIFICATION Version 5.0 | Vol 2, Part E 7.8.66 LE Extended Create Connection Command
// HCI_LE_Extended_Create_Connection
type HciLeExtendedCreateConnection struct {
	InitiatingFilterPolicy uint8
	OwnAddressType         uint8
	PeerAddressType        uint8
	PeerAddress            [6]byte
	InitialtingPhys        uint8

	// 根据InitialtingPhys配置，InitialtingPhys bit置1的话，对应的ConnectionInitialtingList index设置
	ConnectionInitialtingList [3]ConnectionInitialting
}

type HciCmdPktParseResult struct {
	Code int
	Ret  interface{}
}
type HciCmdPktParser func(OpCodeOgf uint8, OpCodeOcf uint16, hciCmdPktPayloadBuf []byte) HciCmdPktParseResult

// 二维parser map
var HciCmdPktParserMap map[uint8]map[uint16]HciCmdPktParser = map[uint8]map[uint16]HciCmdPktParser{
	HCI_CMD_OGF_LE_CONTROLLER_CMD: {
		HCI_LE_EXTENDED_CREATE_CONNECTION: HciLeExtendedCreateConnectionParser,
	},
}

func HciCmdPktParse(OpCodeOgf uint8, OpCodeOcf uint16, hciCmdPktPayloadBuf []byte) HciCmdPktParseResult {
	parser, ok := HciCmdPktParserMap[OpCodeOgf][OpCodeOcf]
	if !ok {
		parser = HciCmdPktDefaultParser
	}
	return parser(OpCodeOgf, OpCodeOcf, hciCmdPktPayloadBuf)
}

func HciCmdPktDefaultParser(OpCodeOgf uint8, OpCodeOcf uint16, hciCmdPktPayloadBuf []byte) HciCmdPktParseResult {
	return HciCmdPktParseResult{Code: HCI_PKT_RET_CODE_NOT_SUPPORT}
}

func HciLeExtendedCreateConnectionParser(OpCodeOgf uint8, OpCodeOcf uint16, hciCmdPktPayloadBuf []byte) HciCmdPktParseResult {
	pkt := HciLeExtendedCreateConnection{}
	bufIndex := 0
	pkt.InitiatingFilterPolicy = hciCmdPktPayloadBuf[bufIndex]
	bufIndex += binary.Size(pkt.InitiatingFilterPolicy)
	pkt.OwnAddressType = hciCmdPktPayloadBuf[bufIndex]
	bufIndex += binary.Size(pkt.OwnAddressType)
	pkt.PeerAddressType = hciCmdPktPayloadBuf[bufIndex]
	bufIndex += binary.Size(pkt.PeerAddressType)
	copy(pkt.PeerAddress[:], hciCmdPktPayloadBuf[bufIndex:bufIndex+len(pkt.PeerAddress)])
	for index := 0; index < len(pkt.PeerAddress)/2; index++ {
		pkt.PeerAddress[index], pkt.PeerAddress[len(pkt.PeerAddress)-1-index] = pkt.PeerAddress[len(pkt.PeerAddress)-1-index], pkt.PeerAddress[index]
	}
	bufIndex += binary.Size(pkt.PeerAddress)
	pkt.InitialtingPhys = hciCmdPktPayloadBuf[bufIndex]
	bufIndex += binary.Size(pkt.InitialtingPhys)
	for index := 0; index < binary.Size(pkt.InitialtingPhys); index++ {
		if pkt.InitialtingPhys&(0x01<<uint8(index)) != 0 {
			conn := ConnectionInitialting{}
			conn.ScanInterval = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.ScanInterval)
			conn.ScanWindow = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.ScanWindow)
			conn.ConnIntervalMin = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.ConnIntervalMin)
			conn.ConnIntervalMax = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.ConnIntervalMax)
			conn.ConnLatency = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.ConnLatency)
			conn.SupervisionTimout = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.SupervisionTimout)
			conn.MinimumCeLength = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.MinimumCeLength)
			conn.MaximumCeLength = binary.LittleEndian.Uint16(hciCmdPktPayloadBuf[bufIndex:])
			bufIndex += binary.Size(conn.MaximumCeLength)
			pkt.ConnectionInitialtingList[index] = conn
		}
	}
	return HciCmdPktParseResult{Code: HCI_PKT_RET_CODE_OK, Ret: pkt}
}
