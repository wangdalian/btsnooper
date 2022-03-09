// hci acl处理

package hci

import (
	"encoding/binary"
)

const (
	ATT_WRITE_REQUEST = 0x12
)

type HciAclPktParseResult struct {
	Code int
	Ret  interface{}
}

type AttPktParser func(OpCode uint8, attPayloadBuf []byte) HciAclPktParseResult

var AttPktParserMap map[int]AttPktParser = map[int]AttPktParser{
	ATT_WRITE_REQUEST: AttPktWriteRequestParser,
}

func AttPktDefaultParser(OpCode uint8, attPayloadBuf []byte) HciAclPktParseResult {
	return HciAclPktParseResult{Code: HCI_PKT_RET_CODE_NOT_SUPPORT}
}

type AttWriteRequest struct {
	OpCode uint8
	Handle uint16
	Value  []byte
}

func AttPktWriteRequestParser(OpCode uint8, attPayloadBuf []byte) HciAclPktParseResult {
	pkt := AttWriteRequest{}
	pkt.OpCode = OpCode

	pktIndex := 0
	pkt.Handle = binary.LittleEndian.Uint16(attPayloadBuf)
	pktIndex += binary.Size(pkt.Handle)
	pkt.Value = make([]byte, len(attPayloadBuf[pktIndex:]))
	copy(pkt.Value, attPayloadBuf[pktIndex:])
	return HciAclPktParseResult{Code: HCI_PKT_RET_CODE_OK, Ret: pkt}
}

func ConnectionOrientedChannelsInBasicFrame(hciAclPktPayloadBuf []byte) HciAclPktParseResult {
	pktIndex := 0
	Length := binary.LittleEndian.Uint16(hciAclPktPayloadBuf[pktIndex:])
	pktIndex += binary.Size(Length)
	ChannelId := binary.LittleEndian.Uint16(hciAclPktPayloadBuf[pktIndex:])
	pktIndex += binary.Size(ChannelId)
	payloadBuf := make([]byte, Length)
	copy(payloadBuf, hciAclPktPayloadBuf[pktIndex:])

	// ATT OpCode
	// BLUETOOTH SPECIFICATION Version 4.2 [Vol 3, Part F] 3.3.1 Attribute PDU Format
	attPktIndex := 0
	var OpCode uint8 = payloadBuf[attPktIndex] & 0x3f
	attPktIndex += binary.Size(OpCode)
	// AuthenticationSignature := payloadBuf[0] >> 0x07
	attPayloadBuf := payloadBuf[attPktIndex:]
	parser, ok := AttPktParserMap[int(OpCode)]
	if !ok {
		parser = AttPktDefaultParser
	}
	return parser(OpCode, attPayloadBuf)
}

// BLUETOOTH SPECIFICATION Version 4.2 [Vol 3, Part A] 3 DATA PACKET FORMAT
func HciAclPktParse(hciAclPktPayloadBuf []byte) HciAclPktParseResult {
	// TODO: 其他报文类型处理
	return ConnectionOrientedChannelsInBasicFrame(hciAclPktPayloadBuf)
}
