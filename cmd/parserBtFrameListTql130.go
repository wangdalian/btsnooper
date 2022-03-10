// 从抓包中过滤出写入操作
// 1. 过滤上报的连接状态信息
// 2. 过滤出对应的写入信息

package main

import (
	"cmd/btsnooper.go/pkg/hci"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"wangdalian/btsnooper/pkg/btsnoop"
)

func main() {
	FilePath := "./data/btsnoop_hci.log"

	// 打开文件
	file, err := os.Open(FilePath)
	if err != nil {
		fmt.Printf("open file error: %s %v", FilePath, err)
		return
	}
	defer file.Close()

	// 读取文件
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("read file error: %s %v", FilePath, err)
		return
	}

	btsnooper := btsnoop.NewFileParser()
	if btsnooper == nil {
		fmt.Printf("create btsnooper failed")
		return
	}

	err = btsnooper.Parse(content)
	if err != nil {
		fmt.Printf("parse error: %s %s", FilePath, err)
		return
	}

	// 解析所有的包
	// 1. 处理PKT_TYPE_HCI_EVT LE_ENHANCED_CONNECTION_COMPLETE_EVENT，获取连接信息(主要是对端地址，Connection Handle)
	// 2. 处理PKT_TYPE_HCI_ACL ATT_WRITE_REQUEST，获取写Handle操作(主要通过Connection Handle匹配)
	var leEnhancedConnectionCompleteEventList []hci.LeEnhancedConnectionCompleteEvent
	var attWriteRequestList []hci.AttWriteRequest
	var attWriteRequestConnHandleList []uint16
	for _, hciPktBuf := range btsnooper.PacketRecordList {
		hciPktType := hciPktBuf.Payload[0]
		hciParserResult := hci.HciPktParse(hciPktType, hciPktBuf.Payload[1:])
		if hciParserResult.Code == hci.HCI_PKT_RET_CODE_OK { // 解析成功后
			if hciParserResult.HciPktType == hci.PKT_TYPE_HCI_EVT {
				parsed, _ := hciParserResult.Ret.(hci.HciEvt)
				hciEvtPktParseResult, _ := parsed.PayloadParsedResult.(hci.HciEvtPktParseResult)
				if hciEvtPktParseResult.Code == hci.HCI_PKT_RET_CODE_OK {
					if hciEvtPktParseResult.EventCode == hci.HCI_EVT_LE_META_EVENT && hciEvtPktParseResult.SubEventCode == hci.LE_ENHANCED_CONNECTION_COMPLETE_EVENT {
						parsed, _ := hciEvtPktParseResult.Ret.(hci.LeEnhancedConnectionCompleteEvent)
						leEnhancedConnectionCompleteEventList = append(leEnhancedConnectionCompleteEventList, parsed)
					}
				}
			} else if hciParserResult.HciPktType == hci.PKT_TYPE_HCI_ACL {
				parsed, _ := hciParserResult.Ret.(hci.HciAcl)
				hciAclPktParseResult, _ := parsed.PayloadParsedResult.(hci.HciAclPktParseResult)
				if hciAclPktParseResult.Code == hci.HCI_PKT_RET_CODE_OK {
					if hciAclPktParseResult.OpCode == hci.ATT_WRITE_REQUEST {
						attWriteRequestConnHandleList = append(attWriteRequestConnHandleList, parsed.Handle)
						parsed, _ := hciAclPktParseResult.Ret.(hci.AttWriteRequest)
						attWriteRequestList = append(attWriteRequestList, parsed)
					}
				}
			}
		}
	}

	for _, item := range leEnhancedConnectionCompleteEventList {
		fmt.Printf("LE_ENHANCED_CONNECTION_COMPLETE_EVENT: %d %s\n", item.ConnectionHandle, hex.EncodeToString(item.PeerAddress[:]))
	}

	for index, item := range attWriteRequestList {
		fmt.Printf("ATT_WRITE_REQUEST: %d %d %s\n", attWriteRequestConnHandleList[index], item.Handle, hex.EncodeToString(item.Value))
	}
}
