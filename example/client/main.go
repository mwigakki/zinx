package main

import (
	"fmt"
	"io"
	"net"

	"github.com/zinx/znet"
)

func main() {
	fmt.Println("client start ...")
	// 1 连接客户端，得到conn对象
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,err = ", err)
	}
	var i uint32
	for {
		i++
		// 2 使用 conn 链接进行读写
		// 客户端也需要写成 message 的形式
		msg := &znet.Message{
			DataLen: 7,
			Id:      0, // 测试心跳包
			Data:    []byte("hello!!")}
		dp := znet.NewDataPack()
		buf, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("client Pack err,err = ", err)
		}
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println("client write err,err = ", err)
			return
		}

		// 接收server的回复
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("client read err :", err)
			return
		}
		msgReceived, err := dp.Unpack(headData, conn)
		if err != nil {
			fmt.Println("client dp Unpack err :", err)
			return
		}
		fmt.Println("Server 回复: ", string(msgReceived.GetDate()))
		// time.Sleep(10 * time.Second)
	}
}
