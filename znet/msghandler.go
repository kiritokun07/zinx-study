package znet

import (
	"fmt"
	"myzinx/utils"
	"myzinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter //存放每个MsgId所对应的处理方法的map属性
	WorkerPoolSize uint32                    //业务工作池 Worker 的数量
	TaskQueue      []chan ziface.IRequest    //Worker 负责取任务的消息队列
}

func (m *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker Id = ", workerId, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//如果有消息，则取出队列的 Request 并执行绑定的业务方法
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandle) StartWorkerPool() {
	//遍历需要启动 Worker 的数量，依次启动
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个 Worker 被启动
		//给当前 Worker 对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前 Worker 阻塞等待对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据connId 来分配当前的连接应该由哪个 Worker 负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerId
	workerId := request.GetConnection().GetConnId() % m.WorkerPoolSize
	fmt.Println("Add connId = ", request.GetConnection().GetConnId(), " request msgId = ", request.GetMsgId(), " to workerId = ", workerId)
	//将请求消息发送给任务队列
	m.TaskQueue[workerId] <- request
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not found!")
		return
	}
	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1.判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		fmt.Println("repeated api,msgId = ", msgId)
		return
	}
	//2.添加msg与api的绑定关系
	m.Apis[msgId] = router
	fmt.Println("add api msgId = ", msgId)
}

func NewMsgHandle() ziface.IMsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
