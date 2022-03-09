package btsnoop

// https://fte.com/webhelpii/hsu/Content/Technical_Information/BT_Snoop_File_Format.htm
// https://datatracker.ietf.org/doc/rfc1761/

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// BtsnooperFileHeader DataType 类型定义
const (
	DATATYPE_RESERVED_MIN = 0
	DATATYPE_RESERVED_MAX = 1000
	DATATYPE_HCI_UNEN     = 1001
	DATATYPE_HCI_UART     = 1002
	DATATYPE_HCI_BSCP     = 1003
	DATATYPE_HCI_SERIAL   = 1004
	DATATYPE_UNSIGNED_MIN = 1005
	DATATYPE_UNSIGNED_MAX = 4294967295
)

// BtsnooperFileHeader DataType 对应字符串
var DataTypeStrMap = map[int]string{
	DATATYPE_HCI_UNEN:   "Un-encapsulated HCI (H1)",
	DATATYPE_HCI_UART:   "HCI UART (H4)",
	DATATYPE_HCI_BSCP:   "HCI BSCP",
	DATATYPE_HCI_SERIAL: "HCI Serial (H5)",
}

var (
	// 文件头部固定标识
	BTSOP_IDENTY = [8]byte{0x62, 0x74, 0x73, 0x6E, 0x6F, 0x6F, 0x70, 0x00}

	// 文件头部固定版本号
	BTSOP_VERNUM uint32 = 0x01
)

// 文件头定义
type FileHeader struct {
	Identy   [8]byte // 文件标识，固定为：62 74 73 6E 6F 6F 70 00
	VerNum   uint32  // 版本号，固定为1
	DataType uint32  // 数据包中的数据链路报头类型
}

// 文件包项
type PacketRecord struct {
	OriginLen   uint32
	IncludedLen uint32
	PacketFlags uint32
	CumuDrops   uint32
	TimestampMs uint64
	Payload     []byte
}

// 解析后的内容
type FileParser struct {
	FileHeader       FileHeader     // 文件头
	PacketRecordList []PacketRecord // Packet Record List
}

func NewFileParser() *FileParser {
	return &FileParser{}
}

// 判断是否支持的文件
func (fp *FileParser) IsSupport() bool {
	return bytes.Equal(fp.FileHeader.Identy[:], BTSOP_IDENTY[:]) && fp.FileHeader.VerNum == BTSOP_VERNUM
}

// 文件解析
func (fp *FileParser) Parse(buf []byte) error {
	index := 0
	copy(fp.FileHeader.Identy[:], buf[index:index+len(fp.FileHeader.Identy)])
	index += len(fp.FileHeader.Identy)
	fp.FileHeader.VerNum = binary.BigEndian.Uint32(buf[index:])
	index += binary.Size(fp.FileHeader.VerNum)
	fp.FileHeader.DataType = binary.BigEndian.Uint32(buf[index:])
	index += binary.Size(fp.FileHeader.DataType)
	if !fp.IsSupport() {
		return fmt.Errorf("not support file type")
	}

	for {
		if index >= len(buf) {
			break
		}
		pkt := PacketRecord{}
		pkt.OriginLen = binary.BigEndian.Uint32(buf[index:])
		index += binary.Size(pkt.OriginLen)
		pkt.IncludedLen = binary.BigEndian.Uint32(buf[index:])
		index += binary.Size(pkt.IncludedLen)
		pkt.PacketFlags = binary.BigEndian.Uint32(buf[index:])
		index += binary.Size(pkt.PacketFlags)
		pkt.CumuDrops = binary.BigEndian.Uint32(buf[index:])
		index += binary.Size(pkt.CumuDrops)
		pkt.TimestampMs = binary.BigEndian.Uint64(buf[index:])
		index += binary.Size(pkt.TimestampMs)
		pkt.Payload = make([]byte, pkt.IncludedLen)
		copy(pkt.Payload, buf[index:uint32(index)+pkt.IncludedLen])
		index += int(pkt.IncludedLen)
		fp.PacketRecordList = append(fp.PacketRecordList, pkt)
	}

	return nil
}

// 自定义输出
func (fp *FileParser) Print(index int) {
	fmt.Println("\n----------------------------")
	fmt.Println("Identification Pattern:", string(fp.FileHeader.Identy[:]))
	fmt.Printf("Version Number: %#x\n", fp.FileHeader.VerNum)
	fmt.Printf("Datalink Type: %#x %s\n", fp.FileHeader.DataType, DataTypeStrMap[int(fp.FileHeader.DataType)])
	if index == -1 {
		for index := 0; index < len(fp.PacketRecordList); index++ {
			fmt.Println("----------------------------")
			fmt.Println("BTSnoop Packet Record Index:", index)
			fmt.Println("----------------------------")
			fmt.Println(" Original Length:", fp.PacketRecordList[index].OriginLen)
			fmt.Println(" Included Length:", fp.PacketRecordList[index].IncludedLen)
			fmt.Printf(" Packet Flags: %#x\n", fp.PacketRecordList[index].PacketFlags)
			fmt.Println(" Cumulative Drops:", fp.PacketRecordList[index].CumuDrops)
			fmt.Println(" Timestamp Microseconds:", fp.PacketRecordList[index].TimestampMs)
			fmt.Println(" Packet Data: ", hex.EncodeToString(fp.PacketRecordList[index].Payload))
		}
	} else {
		fmt.Println("----------------------------")
		fmt.Println("BTSnoop Packet Record Index:", index)
		fmt.Println("----------------------------")
		fmt.Println(" Original Length:", fp.PacketRecordList[index].OriginLen)
		fmt.Println(" Included Length:", fp.PacketRecordList[index].IncludedLen)
		fmt.Printf(" Packet Flags: %#x\n", fp.PacketRecordList[index].PacketFlags)
		fmt.Println(" Cumulative Drops:", fp.PacketRecordList[index].CumuDrops)
		fmt.Println(" Timestamp Microseconds:", fp.PacketRecordList[index].TimestampMs)
		fmt.Println(" Packet Data: ", hex.EncodeToString(fp.PacketRecordList[index].Payload))
	}
}
