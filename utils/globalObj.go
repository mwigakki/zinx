package utils

import (
	"encoding/json"
	"os"

	"github.com/zinx/ziface"
)

/*
	存储一切zinx 框架使用的全局参数，供其他模块使用
	一些参数应该通过 zinx.json由用户去配置
*/

// 消息ID 定义
const (
	MSGID_HEARTBEAT   = 0
	MSGID_GENERAL_MSG = 1
	MSGID_FILE        = 2
)

type GlobalObj struct {
	// server 的配置
	TcpServer ziface.IServer // 当前zinx 全局的server对象
	Host      string         // 当前服务器监听的IP
	Port      int            // 当前服务器监听的tcp 端口
	Name      string         // 当前服务器名称
	// zinx 的配置
	Version          string // 当前 zinx 版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作worker 池的goroutine 数量
	MaxWorkerTaskLen uint32 // 每个worker 对应的消息任务队列最多可以有多少个任务
	// 心跳检测器配置,定义全局的心跳包发送间隔
	// （设定最大值和最小值，具体连接的发送间隔去其中的随机数。因为设定唯一值会使所有连接同时发心跳包，当连接过多时会导致突发流量）
	MinSendInterval int
	MaxSendInterval int // 以秒为单位
}

// 定义全局对外的globalObj 对象
var GlobalObject *GlobalObj

// 提供init方法 初始化对象
func init() {
	GlobalObject = &GlobalObj{ // 现在配置一些默认值
		Name:             "Zinx Server Appp",
		Host:             "127.0.0.1",
		Port:             8999,
		Version:          "V0.9",
		MaxConn:          1000,
		MaxPackageSize:   1024,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MinSendInterval:  3,
		MaxSendInterval:  4,
	}
	// GlobalObject.Reload("")
}

func (g *GlobalObj) Reload(filePath string) {
	if filePath == "" {
		filePath = "conf/zinx.json"
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject) // 会自动将data 接触成的数据装入 对象中
	if err != nil {
		panic(err)
	}
}
