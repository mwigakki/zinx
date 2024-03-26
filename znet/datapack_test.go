package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDatapack(t *testing.T) {
	// 模拟服务器goroutine，从客户端读取数据
	lisenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err :", err)
		return
	}
	go func() {
		conn, err := lisenner.Accept()
		if err != nil {
			fmt.Println("server listen err :", err)
			return
		}
		go func(conn net.Conn) {
			// 处理客户端请求
			dp := NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData) // 从 conn 中读 len(headData) 个字节，读满为止
				if err != nil {
					fmt.Println("server read err :", err)
					return
				}
				msgHead, err := dp.Unpack(headData, conn)
				if err != nil {
					fmt.Println("dp Unpack err :", err)
					return
				}
				if msgHead.GetMsgLen() > 0 {
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen()) // 这里MsgLen()即使非常长，超过1500字节一个包的长度了也无所谓，tcp流推给应用层的都是按顺序的，继续是多个包
					_, err = io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server read err :", err)
						return
					}
					fmt.Printf("----> 读取 msg 成功 ； msgId = %d, msgLen=%d, msg=%s\n", msg.Id, msg.DataLen, string(msg.Data))
				}
			}
		}(conn)
	}()

	// 模拟客户端goroutine
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("net dial err :", err)
		return
	}
	// 创建封装的包对象
	dp := NewDataPack()
	// 模拟两个数据包粘包过程，即封装两个message
	msg1 := &Message{1, 5, []byte("hello")}
	msg2 := &Message{2, 4, []byte("Zinx")}
	msgBytes1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("msg1 Pack err :", err)
		return
	}
	msgBytes2, _ := dp.Pack(msg2)
	// 把两个包强行粘在一起
	sendData := append(msgBytes1, msgBytes2...)
	_, err = conn.Write(sendData)
	if err != nil {
		fmt.Println("client write err :", err)
		return
	}
	// 客户端阻塞
	select {}
}
