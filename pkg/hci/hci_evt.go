// hci cmd处理

package hci

import "encoding/binary"

// Evt列表
const (
	HCI_EVT_LE_META_EVENT = 0x3E
)

const (
	NO_SUB_EVENT = -1
)

// HCI_EVT_LE_META_EVENT子类型
const (
	LE_ENHANCED_CONNECTION_COMPLETE_EVENT = 0x0A
)

type HciEvtPktParseResult struct {
	Code         int
	EventCode    uint8
	SubEventCode uint8
	Ret          interface{}
}
type HciEvtPktParser func(EventCode uint8, SubEventCode int, hciEvtPktPayloadBuf []byte) HciEvtPktParseResult

// 二维parser map
// 没有二级的则使用-1做索引
var HciEvtPktParserMap map[uint8]map[int]HciEvtPktParser = map[uint8]map[int]HciEvtPktParser{
	HCI_EVT_LE_META_EVENT: {
		LE_ENHANCED_CONNECTION_COMPLETE_EVENT: LeEnhancedConnectionCompleteEventParser,
	},
}

func HciEvtPktParse(EventCode uint8, hciEvtPktPayloadBuf []byte) HciEvtPktParseResult {
	subEventCode := NO_SUB_EVENT
	if EventCode == HCI_EVT_LE_META_EVENT {
		subEventCode = int(hciEvtPktPayloadBuf[0])
	}
	parser, ok := HciEvtPktParserMap[EventCode][subEventCode]
	if !ok {
		parser = HciEvtPktDefaultParser
	}
	parsed := parser(EventCode, subEventCode, hciEvtPktPayloadBuf)
	parsed.EventCode = EventCode
	parsed.SubEventCode = uint8(subEventCode)
	return parsed
}

func HciEvtPktDefaultParser(EventCode uint8, SubEventCode int, hciEvtPktPayloadBuf []byte) HciEvtPktParseResult {
	return HciEvtPktParseResult{Code: HCI_PKT_RET_CODE_NOT_SUPPORT}
}

type LeEnhancedConnectionCompleteEvent struct {
	SubEventCode                  uint8
	Status                        uint8
	ConnectionHandle              uint16
	Role                          uint8
	PeerAddressType               uint8
	PeerAddress                   [6]uint8
	LocalResolvablePrivateAddress [6]uint8
	PeerResolvablePrivateAddress  [6]uint8
	ConnInterval                  uint16
	ConnLatency                   uint16
	SupervisionTimeout            uint16
	MasterClockAccuracy           uint8
}

// 包含了连接成功后的handle，和对端地址，记录下来用于相关操作信息
func LeEnhancedConnectionCompleteEventParser(EventCode uint8, SubEventCode int, hciEvtPktPayloadBuf []byte) HciEvtPktParseResult {
	pktIndex := 0
	pkt := LeEnhancedConnectionCompleteEvent{}
	pkt.SubEventCode = hciEvtPktPayloadBuf[pktIndex]
	pktIndex += binary.Size(pkt.SubEventCode)
	pkt.Status = hciEvtPktPayloadBuf[pktIndex]
	pktIndex += binary.Size(pkt.Status)
	pkt.ConnectionHandle = binary.LittleEndian.Uint16(hciEvtPktPayloadBuf[pktIndex:])
	pktIndex += binary.Size(pkt.ConnectionHandle)
	pkt.Role = hciEvtPktPayloadBuf[pktIndex]
	pktIndex += binary.Size(pkt.Role)
	pkt.PeerAddressType = hciEvtPktPayloadBuf[pktIndex]
	pktIndex += binary.Size(pkt.PeerAddressType)
	copy(pkt.PeerAddress[:], hciEvtPktPayloadBuf[pktIndex:])
	for index := 0; index < len(pkt.PeerAddress)/2; index++ {
		pkt.PeerAddress[index], pkt.PeerAddress[len(pkt.PeerAddress)-1-index] = pkt.PeerAddress[len(pkt.PeerAddress)-1-index], pkt.PeerAddress[index]
	}
	pktIndex += len(pkt.PeerAddress) * binary.Size(pkt.PeerAddress[0])
	copy(pkt.LocalResolvablePrivateAddress[:], hciEvtPktPayloadBuf[pktIndex:])
	for index := 0; index < len(pkt.LocalResolvablePrivateAddress)/2; index++ {
		pkt.LocalResolvablePrivateAddress[index], pkt.LocalResolvablePrivateAddress[len(pkt.LocalResolvablePrivateAddress)-1-index] = pkt.LocalResolvablePrivateAddress[len(pkt.LocalResolvablePrivateAddress)-1-index], pkt.LocalResolvablePrivateAddress[index]
	}
	pktIndex += len(pkt.LocalResolvablePrivateAddress) * binary.Size(pkt.LocalResolvablePrivateAddress[0])
	copy(pkt.PeerResolvablePrivateAddress[:], hciEvtPktPayloadBuf[pktIndex:])
	for index := 0; index < len(pkt.PeerResolvablePrivateAddress)/2; index++ {
		pkt.PeerResolvablePrivateAddress[index], pkt.PeerResolvablePrivateAddress[len(pkt.PeerResolvablePrivateAddress)-1-index] = pkt.PeerResolvablePrivateAddress[len(pkt.PeerResolvablePrivateAddress)-1-index], pkt.PeerResolvablePrivateAddress[index]
	}
	pktIndex += len(pkt.PeerResolvablePrivateAddress) * binary.Size(pkt.PeerResolvablePrivateAddress[0])
	pkt.ConnInterval = binary.LittleEndian.Uint16(hciEvtPktPayloadBuf[pktIndex:])
	pktIndex += binary.Size(pkt.ConnInterval)
	pkt.ConnLatency = binary.LittleEndian.Uint16(hciEvtPktPayloadBuf[pktIndex:])
	pktIndex += binary.Size(pkt.ConnLatency)
	pkt.SupervisionTimeout = binary.LittleEndian.Uint16(hciEvtPktPayloadBuf[pktIndex:])
	pktIndex += binary.Size(pkt.SupervisionTimeout)
	pkt.MasterClockAccuracy = hciEvtPktPayloadBuf[pktIndex]
	return HciEvtPktParseResult{Code: HCI_PKT_RET_CODE_OK, Ret: pkt}
}
