package znet

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/zinx/utils"
	"github.com/zinx/ziface"
)

// IServer的接口实现
type Server struct {
	Name       string
	IPVersion  string
	IP         string
	Port       int
	MsgHandler ziface.IMessageHandle // 当前server注册的连接对应的处理处理业务的router
	ConnMgr    ziface.IConnManager   // 该server的连接管理器
	// 添加该server 创建连接之后自动调用 hook 函数
	OnConnStart func(ziface.IConnection)
	// 添加该server 创建连接之后自动调用 hook 函数
	OnConnStop func(ziface.IConnection)
	// 是否开启连接的心跳检测器，为true的话，此服务器的每个连接都会默认开启
	UseHeartBeat bool

	cId uint32 // 每来一个连接给分配一个cI使用原子方法进行自增
}

// 初始化 Server 模块
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:         utils.GlobalObject.Name,
		IPVersion:    "tcp4",
		IP:           utils.GlobalObject.Host,
		Port:         utils.GlobalObject.Port,
		MsgHandler:   NewMessageHandle(),
		ConnMgr:      NewConnManager(),
		OnConnStart:  func(ziface.IConnection) {},
		OnConnStop:   func(ziface.IConnection) {}, // 给所有的连接注册两个空的钩子函数，如果开发者不自己提供的话
		UseHeartBeat: true,
	}
	s.ConnMgr.SetServer(s) // 设置连接管理模块对应的server
	if s.UseHeartBeat {
		s.AddRouter(utils.MSGID_HEARTBEAT, &HeartbeatDefaultRouter{})
	}
	return s
}

// 开始服务器
func (s *Server) Start() {
	// 开启一个tcp 服务器
	// 开启工作池
	s.MsgHandler.StartWokerPool()
	go s.GetConnMgr().ConnManage() // 开启连接管理器 管理连接增加和删除的方法
	fmt.Println("[Server] started WokerPool")
	go func() { // 防止start 函数等待连接阻塞
		// 1 获取一个 TCP 的addr
		fmt.Printf("[start] Server %s Listenner at IP %s, Port %d is starting\n", s.Name, s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port)) // 得到一个tcp句柄
		if err != nil {
			log.Panic(err)
		}
		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("start Zinx server %s successfully\n", s.Name)
		// 3 阻塞，等待客户端连接，处理客户端连接业务，读写
		for {
			// 如果有客户端连接，此函数便有返回了，然后将conn传给我们自定的连接对象
			conn, err := listenner.AcceptTCP()
			if err != nil {
				log.Println("accept err:", err)
				continue
			}
			// 判断当前连接个数是否超过最大值，
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个错误数据，告诉它连不上了
				fmt.Println("********************* 连接数到达上限了")
				conn.Close()
				continue
			}
			// 客户端连接server 成功
			newcId := atomic.AddUint32(&s.cId, 1)
			dealConn := NewConnection(conn, newcId, s.MsgHandler, s.ConnMgr.GetConnMgrChan())
			if s.UseHeartBeat {
				s.bindHeartBeatChecker(dealConn)
			}
			go dealConn.Start() // 每个客户端应该异步开启连接，所以这里需要使用 goroutine
		}
	}()
}

// 结束服务器
func (s *Server) Stop() {
	fmt.Println("Zinx Stopepd !!!")
	s.ConnMgr.Clear()
	s.MsgHandler.ClearTaskQueue()
}

// 运行服务器
func (s *Server) Serve() {
	// 启动server 的服务功能
	s.Start()
	// TODO 完成额外业务
	// 用户调用Serve后，即开启服务器了，此时应阻塞住，
	select {}
}

// 路由功能：给当前的服务注册一个路由功能，供客户端的连接使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

// 得到连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 设置该server 创建连接之后自动调用的 hook 函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 调用该server 创建连接之后自动调用的 hook 函数，具体的调用是在连接管理模块中执行的
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	} else {
		fmt.Println("server 尚未注册 OnConnStart 方法 ")
	}
}

// 设置该server 断开连接之前自动调用的 hook 函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用该server 断开连接之前自动调用的 hook 函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	} else {
		fmt.Println("server 尚未注册 OnConnStop方法 ")
	}
}

// 给连接绑定心跳检测器
func (s *Server) bindHeartBeatChecker(conn ziface.IConnection) {
	// 设定最大值和最小值，具体连接的发送间隔去其中的随机数。因为设定唯一值会使所有连接同时发心跳包，当连接过多时会导致突发流量
	source := rand.NewSource(int64(conn.GetConnID())) // 根据流ID 生成随机数种子。
	randNumGenetor := rand.New(source)
	randomSendInterval := randNumGenetor.Intn(utils.GlobalObject.MaxSendInterval-utils.GlobalObject.MinSendInterval) + utils.GlobalObject.MinSendInterval
	conn.BindHeartBeatChecker(NewHeartbeatChecher(conn, time.Duration(randomSendInterval)*time.Second))
}
