package main

import (
	"fmt"
	"io"
	"myzinx/znet"
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
		//发封包 message 消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMessage(0, []byte("zinx v0.5 client test message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err", err)
			return
		}

		//1.先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		//ReadFull会把msg填充满为止
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//2.将headData字节流拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack err:", err)
			return
		}
		//3.根据dataLen从io中读取字节流
		if msgHead.GetDataLen() > 0 {
			//msg有data数据,需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client unpack data err:", err)
				return
			}
			fmt.Printf("===>Receive msg:id=%d,len=%d,data=%s\n", msg.Id, msg.DataLen, string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
