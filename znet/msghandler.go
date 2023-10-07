package znet

import (
	"fmt"
	"myzinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter //存放每个MsgId所对应的处理方法的map属性
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
		Apis: make(map[uint32]ziface.IRouter),
	}
}
