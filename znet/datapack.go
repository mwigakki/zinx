package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/zinx/utils"
	"github.com/zinx/ziface"
)

/**
封包，拆包 模块，用于解决tcp 粘包问题
TLV 格式， [TYPE | HEADER LENGTH | BODY LENGTH | HEADER | BODY ]
*/

// 我感觉这里DataPack 有点多余，它的方法完全可以交给 Message 去完成
type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
func (dp *DataPack) GetHeadLen() uint32 { // 得到应用层包头的长度
	// DataLen uint32(4字节) + Id uint32(4字节)
	return 8
}

// 相当于结构体的序列化
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{}) // 创建一个空的缓冲
	// 把msgLen，msgId，data 按顺序写入缓冲
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetDate()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 相当于将字节切片反序列化为结构体
func (dp *DataPack) Unpack(headData []byte, conn net.Conn) (ziface.IMessage, error) {
	// 先读 head（len和id）的信息
	buf := bytes.NewReader(headData)
	msg := &Message{}
	if err := binary.Read(buf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断发来的包是否大于了我们设定的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		fmt.Println(utils.GlobalObject.MaxPackageSize, msg.DataLen)
		return nil, errors.New("too large msg data received ! ")
	}
	// 继续从 conn 中读取data
	msg.Data = make([]byte, msg.GetMsgLen())
	_, err := io.ReadFull(conn, msg.Data)
	if err != nil {
		fmt.Println("datapack read msg.Data err :", err)
		return nil, err
	}
	return msg, nil
}
