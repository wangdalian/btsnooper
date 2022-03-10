// HCI协议相关定义

package hci

import (
	"encoding/binary"
)

// HCI 数据类型
// BLUETOOTH SPECIFICATION Version 4.2 [Vol 4, Part A] 2 PROTOCOL
const (
	PKT_TYPE_HCI_CMD  = 0x01
	PKT_TYPE_HCI_ACL  = 0x02
	PKT_TYPE_HCI_SYNC = 0x03
	PKT_TYPE_HCI_EVT  = 0x04
)

// HCI ACL DATA数据格式
// BLUETOOTH SPECIFICATION Version 4.2 [Vol 2, Part E] 5.4.2 HCI ACL Data Packets
type HciAcl struct {
	Handle              uint16
	PbFlag              uint8
	BcFlag              uint8
	DataTotalLen        uint16
	Data                []byte
	PayloadParsedResult interface{} // Data解析后的结果
}

// BLUETOOTH SPECIFICATION Version 4.2 [Vol 2, Part E]
// 5.4.1 HCI Command Packet
type HciCmd struct {
	OpCode              uint16
	OpCodeOcf           uint16
	OpCodeOgf           uint8
	ParamTotalLen       uint8
	Data                []byte
	PayloadParsedResult interface{} // Data解析后的结果
}

type HciSync struct {
	// TODO
}

type HciEvt struct {
	EventCode            uint8
	ParameterTotalLength uint8
	EventParameterList   []byte
	PayloadParsedResult  interface{} // Data解析后的结果
}

const (
	HCI_PKT_RET_CODE_OK          = 0
	HCI_PKT_RET_CODE_NOT_SUPPORT = 1001 // 不支持
)

type HciPktParseResult struct {
	Code       int
	HciPktType uint8
	Ret        interface{}
}
type HciPktParser func(hciPktType byte, hciPayloadBuf []byte) HciPktParseResult

var HciPktParserMap map[int]HciPktParser = map[int]HciPktParser{
	PKT_TYPE_HCI_ACL: HciPktAclParser,
	PKT_TYPE_HCI_CMD: HciPktCmdParser,
	PKT_TYPE_HCI_EVT: HciPktEvtParser,
}

func HciPktParse(hciPktType byte, hciPayloadBuf []byte) HciPktParseResult {
	parser, ok := HciPktParserMap[int(hciPktType)]
	if !ok {
		parser = HciDefaultParser
	}
	parsed := parser(hciPktType, hciPayloadBuf)
	parsed.HciPktType = hciPktType
	return parsed
}

func HciDefaultParser(hciPktType byte, hciPayloadBuf []byte) HciPktParseResult {
	return HciPktParseResult{Code: HCI_PKT_RET_CODE_NOT_SUPPORT}
}

func HciPktEvtParser(hciPktType byte, hciPayloadBuf []byte) HciPktParseResult {
	pkt := HciEvt{}
	pkt.EventCode = hciPayloadBuf[0]
	pkt.ParameterTotalLength = hciPayloadBuf[1]
	pkt.EventParameterList = make([]byte, pkt.ParameterTotalLength)
	copy(pkt.EventParameterList, hciPayloadBuf[2:])
	pkt.PayloadParsedResult = HciEvtPktParse(pkt.EventCode, pkt.EventParameterList)
	return HciPktParseResult{Code: HCI_PKT_RET_CODE_OK, Ret: pkt}
}

func HciPktCmdParser(hciPktType byte, hciPayloadBuf []byte) HciPktParseResult {
	pkt := HciCmd{}
	pkt.OpCode = binary.LittleEndian.Uint16(hciPayloadBuf)
	pkt.OpCodeOcf = pkt.OpCode & 0x03ff
	pkt.OpCodeOgf = uint8(pkt.OpCode >> 10 & 0xfc)
	pkt.ParamTotalLen = hciPayloadBuf[2]
	pkt.Data = make([]byte, pkt.ParamTotalLen)
	copy(pkt.Data, hciPayloadBuf[3:])
	pkt.PayloadParsedResult = HciCmdPktParse(pkt.OpCodeOgf, pkt.OpCodeOcf, pkt.Data)
	return HciPktParseResult{Code: HCI_PKT_RET_CODE_OK, Ret: pkt}
}

func HciPktAclParser(hciPktType byte, hciPayloadBuf []byte) HciPktParseResult {
	pkt := HciAcl{}
	pkt.Handle = binary.LittleEndian.Uint16(hciPayloadBuf) & 0x0fff
	pkt.PbFlag = (hciPayloadBuf[0] & 0xf0) >> 0x04
	pkt.BcFlag = (hciPayloadBuf[0] & 0xf0) >> 0x06
	pkt.DataTotalLen = binary.LittleEndian.Uint16(hciPayloadBuf[2:])
	pkt.Data = make([]byte, pkt.DataTotalLen)
	copy(pkt.Data, hciPayloadBuf[4:])
	pkt.PayloadParsedResult = HciAclPktParse(pkt.Data)
	return HciPktParseResult{Code: HCI_PKT_RET_CODE_OK, Ret: pkt}
}
