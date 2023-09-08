package znet

import (
	"fmt"
	"myzinx/ziface"
	"net"
)

type Connection struct {
	//当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	//当前连接的ID 即SessionID 全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法API
	handleAPI ziface.HandFunc
	//告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
}

// StartReader 处理conn读数据的协程
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		//将最大的数据读到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf err", err)
			c.ExitBuffChan <- true
			continue
		}
		//调用当前连接业务（这里执行的是当前conn绑定的handle方法
		if err = c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID", c.ConnID, " handle is error")
			c.ExitBuffChan <- true
			return
		}
	}
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandFunc) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callbackApi,
		ExitBuffChan: make(chan bool, 1),
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//开启处理该连接读取客户端数据之后的请求业务
	go c.StartReader()
	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	//1.如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，则在此刻应该显示调用
	//关闭socket连接
	_ = c.Conn.Close()
	//通知从缓冲队列读数据的业务，该连接已经关闭
	c.ExitBuffChan <- true
	//关闭该连接的全部管道
	close(c.ExitBuffChan)
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnId 获取当前连接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
