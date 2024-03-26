package znet

import (
	"fmt"

	"github.com/zinx/ziface"
)

// 定义实现 IRouter 的基类，但这个基类不具体实现任何方法
// 用户实现自己的router时，先嵌入这个基类，然后对这个基类的方法进行重写就好了
// 有点像适配层，对接口进行隔离
// 当然用户也可以不继承，直接实现 IRouter 也行。
// 只是实现接口的话就必须三个方法都实现，继承基类的话就可以按需求实现想实现的方法即可
type BaseRouter struct{}

// 在处理 conn 业务之前的钩子方法
func (br *BaseRouter) PreHandle(reqeust ziface.IRequest) {}

// 处理 conn 业务的主方法
func (br *BaseRouter) Handle(reqeust ziface.IRequest) {}

// 在处理 conn 业务之后的钩子方法
func (br *BaseRouter) PostHandle(reqeust ziface.IRequest) {}

// 默认的 收到心跳包 回包时的路由处理
type HeartbeatDefaultRouter struct {
	BaseRouter
}

// 重写 BaseRouter 的Handle 方法
func (br *HeartbeatDefaultRouter) Handle(reqeust ziface.IRequest) {
	fmt.Printf("收到连接（id=%d）对端的心跳回包, msg= %s \n", reqeust.GetConnection().GetConnID(), string(reqeust.GetData()))
}
