package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          //以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                       //启动Worker工作池
	SendMsgToTaskQueue(request IRequest)    //将消息交给TaskQueue 由Worker进行处理
}
