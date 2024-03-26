package main

import (
	"fmt"

	"github.com/zinx/utils"
	"github.com/zinx/ziface"
	"github.com/zinx/znet"
)

type myRouter struct {
	znet.BaseRouter
}

// 开发者自己实现的处理 conn 业务的主方法
func (mr *myRouter) Handle(reqeust ziface.IRequest) {
	conn := reqeust.GetConnection()
	data := reqeust.GetData() // 得到的只是数据，不包含message 的头
	fmt.Printf("来自连接id= %d, data = %s \n", conn.GetConnID(), string(data))
	// 数据回复

	err := conn.SendMsg(reqeust.GetMsgId()+1, data)
	if err != nil {
		fmt.Println("router handle err :", err)
		return
	}
}

// 钩子函数的调用具体是在连接管理模块中执行的
func myHookStart(conn ziface.IConnection) {
	fmt.Println("[自定义的 myHookStart] conn id = ", conn.GetConnID())
	// conn.SendMsg(1, []byte("连接服务器成功！"))
	// 给当前的连接设置一些属性
	conn.SetProperty("name", "连接的用户的名字 tom")
}
func myHookStop(conn ziface.IConnection) {
	name, _ := conn.GetProperty("name")
	fmt.Printf("[自定义的 myHookStop] goodbye! conn id = %d, name = %s \n", conn.GetConnID(), name)
}

/*
基于 zinx开发的服务器端应用程序
*/
func main() {
	// 1 创建一个server 句柄，使用 zinx 的api
	s := znet.NewServer("[zinx V1]")
	// 2 注册连接建立和断开的hook函数
	s.SetOnConnStart(myHookStart)
	s.SetOnConnStop(myHookStop)
	// 3 添加一个router
	s.AddRouter(utils.MSGID_GENERAL_MSG, &myRouter{})
	// 4 启动server
	s.Serve()
}
