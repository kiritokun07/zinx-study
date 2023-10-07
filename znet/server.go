package znet

import (
	"errors"
	"fmt"
	"myzinx/utils"
	"myzinx/ziface"
	"net"
)

// Server IServer 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tpc4 or other
	IPVersion string
	//服务绑定的IP地址
	Ip string
	//服务绑定的端口
	Port int
	//当前 Server 由用户绑定回调 router ，也就是 Server 注册的连接对应的处理业务
	//Router ziface.IRouter
	//当前 Server 的消息管理模块，用来绑定 MsgId 和对应的处理方法
	msgHandler ziface.IMsgHandle
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	obj := utils.GlobalObject
	return &Server{
		Name:       obj.Name, //从全局参数获取
		IPVersion:  "tcp4",
		Ip:         obj.Host,    //从全局参数获取
		Port:       obj.TcpPort, //从全局参数获取
		msgHandler: NewMsgHandle(),
		//Router:    nil,
	}
}

// CallBackToClient 定义当前客户端连接的handle API
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle]CallBackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START]Server name: %s,listener at %s:%d is starting\n", s.Name, s.Ip, s.Port)
	fmt.Printf("[Zinx]Version:%s,MaxConn:%d,MaxPacketSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	//fmt.Printf("[START]Server listener at IP:%s,Port %d,is starting\n", s.Ip, s.Port)
	//开启一个go去做服务器端Listener业务
	go func() {
		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		//2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err:", err)
			return
		}
		//已经监听成功
		fmt.Println("start zinx server [", s.Name, "] success,now listening...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		//3.启动server网络连接业务
		for {
			//3.1阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//3.2 TODO 设置服务器最大连接控制，如果超过最大连接，则关闭此连接
			//3.3 处理该新连接请求的业务方法，此时handler和conn应该是绑定的
			dealConn := NewConnection(conn, cid, s.msgHandler)
			cid++
			//3.4 启动当前连接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP]zinx server,name:", s.Name)
	//TODO Server.Stop() 将需要清理的连接信息或其他信息一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 如果在启动服务时还要处理其他事情，则可以在这里添加
	//阻塞，否则主协程退出，listener的协程将会退出
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	//s.Router = router
	fmt.Println("Add Router success!")
}
