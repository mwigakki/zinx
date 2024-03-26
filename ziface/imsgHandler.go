package ziface

// 此接口要放在在server 中
type IMessageHandle interface {
	// 调度，执行对应的router消息处理方法
	DoMsgHandler(IRequest)
	// 给server添加具体的router 处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// 启动一个worker 工作池（只能执行一次）
	StartWokerPool()
	// 将消息交给 一个 taskQueue ，让一个worker 去处理
	SendMsgToTaskQueue(request IRequest)
	// 清除消息队列的方法，当server 主动关闭时需要调用，channel类型的数据都需关闭
	ClearTaskQueue()
}
