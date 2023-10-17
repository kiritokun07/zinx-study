package ziface

import "net"

// 定义连接接口
type IConnection interface {
	// Start 启动连接，让当前连接开始工作
	Start()
	// Stop 停止连接，结束当前连接状态
	Stop()
	// GetTCPConnection 从当前连接获取原始的 socket TCPConn
	GetTCPConnection() *net.TCPConn
	// GetConnId 获取当前连接ID
	GetConnId() uint32
	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// SendMsg 直接将 Message 数据发送给远程的TCP客户端 (无缓冲 我:可以把基本的interface封装起来)
	SendMsg(msgId uint32, data []byte) error
	//直接将 Message 数据发送给远程的TCP客户端 (有缓冲）
	SendBuffMsg(msgId uint32, data []byte) error
}

// HandFunc 定义一个统一处理连接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
