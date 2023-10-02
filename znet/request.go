package znet

import "myzinx/ziface"

type Request struct {
	conn ziface.IConnection //已经和客户端建立好连接
	msg  ziface.IMessage    //客户端请求的数据
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
