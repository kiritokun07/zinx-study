package ziface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动服务器方法
	Start()
	// Stop 停止服务器方法
	Stop()
	// Serve 开启业务服务方法
	Serve()
	// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端连接处理使用
	AddRouter(msgId uint32, router IRouter)
	// GetConnMgr 得到连接管理器
	GetConnMgr() IConnManager
	// SetOnConnStart 设置该 Server 连接创建时的 Hook 函数
	SetOnConnStart(func(conn IConnection))
	// SetOnConnStop 设置该 Server 连接断开时的 Hook 函数
	SetOnConnStop(func(conn IConnection))
	// CallOnConnStart 调用连接 OnConnStart Hook 函数
	CallOnConnStart(conn IConnection)
	// CallOnConnStop 调用连接 OnConnStop Hook 函数
	CallOnConnStop(conn IConnection)
}
