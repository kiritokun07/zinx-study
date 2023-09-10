package main

import (
	"fmt"
	"myzinx/ziface"
	"myzinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	//一定要先定义基础路由 BaseRouter
	znet.BaseRouter
}

func (r *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Println("PreHandle ping error")
	}
}

func (r *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
	if err != nil {
		fmt.Println("Handle ping error")
	}
}

func (r *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
	if err != nil {
		fmt.Println("PostHandle ping error")
	}
}

// Server模块的测试函数
func main() {
	//1.创建一个Server句柄s
	s := znet.NewServer("[zinx v0.1]")
	s.AddRouter(&PingRouter{})
	//2.开启服务
	s.Serve()
}
