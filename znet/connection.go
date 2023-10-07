package znet

import (
	"errors"
	"fmt"
	"io"
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
	//该连接的处理方法 router
	//Router ziface.IRouter
	MsgHandler ziface.IMsgHandle
	//告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲chan 用于读写两个协程之间的消息通信
	msgChan chan []byte
}

// StartReader 处理conn读数据的协程
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		//创建拆包解包对象
		dp := NewDataPack()

		//读取客户端的msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆包 得到msgId和dataLen后放到msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			c.ExitBuffChan <- true
			continue
		}
		//根据dataLen 读取data 放到msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg dta error:", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		//我:获取完整msg的过程可以再次封装
		msg.SetData(data)

		//得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// StartWriter 写消息协程，将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[writer goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error:", err, " conn writer exit")
				return
			}
		case <-c.ExitBuffChan:
			//conn已经关闭
			return
		}
	}
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) ziface.IConnection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		//Router:       router,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//1.开启用户从客户端读取数据流程的goroutine
	go c.StartReader()
	//2.开启用于写回客户端数据流程的goroutine
	go c.StartWriter()

	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// Stop 停止连接，结束当前连接状态
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	//将data封包,并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msgId=", msgId)
		return errors.New("pack error msg")
	}
	//写回客户端
	//将之前直接回写给 conn.Write 的方法改为发送给chan
	c.msgChan <- msg
	return nil
}
