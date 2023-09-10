package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")
	//3s后发起测试请求，给服务器端开启服务的机会
	time.Sleep(2 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err,exit!", err)
		return
	}
	for {
		_, err := conn.Write([]byte("hello ZINX2"))
		if err != nil {
			fmt.Println("write error", err)
			break
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		fmt.Printf("server call back:%s,cnt=%d\n", buf[:cnt], cnt)
		time.Sleep(1 * time.Second)
	}
}
