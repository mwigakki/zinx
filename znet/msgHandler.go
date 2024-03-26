package znet

import (
	"fmt"

	"github.com/zinx/utils"
	"github.com/zinx/ziface"
)

// 每个server有一个MessageHandle属性，只是这个属性会同时传给所有connection
type MessageHandle struct {
	// 每一个消息ID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 业务工作 worker池的worker  数量
	WorkerPoolSize uint32
	// 负责 worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
}

func NewMessageHandle() *MessageHandle {
	return &MessageHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 在全局配置文件中读取应该，
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度，执行对应的router消息处理方法
func (m *MessageHandle) DoMsgHandler(req ziface.IRequest) {
	msgId := req.GetMsgId()
	handler, has := m.Apis[msgId]
	if !has {
		fmt.Println("[WARNING] api msg id [", msgId, "] is NOT FOUND! need register!")
		return
	}
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

// 给server添加具体的router 处理逻辑
func (m *MessageHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	m.Apis[msgID] = router
	fmt.Println("添加", router, " 成功")
}

// 启动一个worker 工作池（只能执行一次）
func (m *MessageHandle) StartWokerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 给当前的worker 对应的 channel 消息队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的worker，阻塞等待消息
		go m.StartWoker(i, m.TaskQueue[i])
	}
	// 也可以尝试 所有worker 共享一个队列
}

// 启动一个worker
func (m *MessageHandle) StartWoker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("[msgHandler] worker id = ", workerId, " is started ...")
	// 不断阻塞等待消息任务
	for {
		request := <-taskQueue // 从任务队列中获得请求
		m.DoMsgHandler(request)
	}
}

// 将消息交给 一个 taskQueue ，让一个worker 去处理
func (m *MessageHandle) SendMsgToTaskQueue(req ziface.IRequest) {
	// 1 将消息平均分配给不同的 worker
	// 根据每个客户端的连接 ID 来分配 worker
	workerId := req.GetConnection().GetConnID() % m.WorkerPoolSize
	// 2 将消息发送给对应的worker 的TaskQueue
	m.TaskQueue[workerId] <- req
}

// 清除消息队列的方法，当server 主动关闭时需要调用
func (m *MessageHandle) ClearTaskQueue() {
	for i := range m.TaskQueue {
		close(m.TaskQueue[i]) // channel类型的数据都需要关闭
	}
}
