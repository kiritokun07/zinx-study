package ziface

// IRequest 把客户端请求的连接信息和请求的数据包装到Request里
type IRequest interface {
	GetConnection() IConnection //获取请求连接信息
	GetMsgId() uint32           //获取请求消息ID
	GetData() []byte            //获取请求消息的数据
}
