// 解析BTSnoop文件格式
// 文件格式参考：https://fte.com/webhelpii/hsu/Content/Technical_Information/BT_Snoop_File_Format.htm
package main

import (
	"cmd/btsnooper.go/pkg/hci"
	"fmt"
	"io/ioutil"
	"os"
	"wangdalian/btsnooper/pkg/btsnoop"
)

// BTSnoop文件格式定义

const ()

var (
	FilePath = "./data/btsnoop_hci.log"
)

func main() {

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

	// 打印测试HCI包解析
	testAttWriteCmd(btsnooper)
	testLeExtendCreateConnection(btsnooper)
	testLeEnhancedConnectionComplete(btsnooper)
}

func testLeEnhancedConnectionComplete(btsnooper *btsnoop.FileParser) {
	const btsnoopPacketRecordIndex = 455
	btsnooper.Print(btsnoopPacketRecordIndex)
	hciPktBuf := btsnooper.PacketRecordList[btsnoopPacketRecordIndex]
	if len(hciPktBuf.Payload) <= 0 {
		fmt.Println("invalid hci packet")
		return
	}
	hciPktType := hciPktBuf.Payload[0]
	hciParserResult := hci.HciPktParse(hciPktType, hciPktBuf.Payload[1:])
	fmt.Printf("%#v", hciParserResult)
}

func testLeExtendCreateConnection(btsnooper *btsnoop.FileParser) {
	const btsnoopPacketRecordIndex = 453
	btsnooper.Print(btsnoopPacketRecordIndex)
	hciPktBuf := btsnooper.PacketRecordList[btsnoopPacketRecordIndex]
	if len(hciPktBuf.Payload) <= 0 {
		fmt.Println("invalid hci packet")
		return
	}
	hciPktType := hciPktBuf.Payload[0]
	hciParserResult := hci.HciPktParse(hciPktType, hciPktBuf.Payload[1:])
	fmt.Printf("%#v", hciParserResult)
}

func testAttWriteCmd(btsnooper *btsnoop.FileParser) {
	const btsnoopPacketRecordIndex = 1598
	btsnooper.Print(btsnoopPacketRecordIndex)
	hciPktBuf := btsnooper.PacketRecordList[btsnoopPacketRecordIndex]
	if len(hciPktBuf.Payload) <= 0 {
		fmt.Println("invalid hci packet")
		return
	}
	hciPktType := hciPktBuf.Payload[0]
	hciParserResult := hci.HciPktParse(hciPktType, hciPktBuf.Payload[1:])
	fmt.Printf("%#v", hciParserResult)
}
