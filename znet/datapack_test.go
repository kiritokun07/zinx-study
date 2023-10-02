package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"unsafe"
)

// 测试封包拆包 Server
func TestDataPack_Server(t *testing.T) {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	//创建服务器goroutine，负责从客户端goroutine读取黏包的数据，进行解析
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}
		//处理客户端请求
		go func(conn net.Conn) {
			//创建封包拆包对象dp
			dp := NewDataPack()
			for {
				//1.先读出流中的head部分
				headData := make([]byte, dp.GetHeadLen())
				//ReadFull会把msg填充满为止
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error")
					break
				}
				//2.将headData字节流拆包到msg中
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack err:", err)
					return
				}
				//3.根据dataLen从io中读取字节流
				if msgHead.GetDataLen() > 0 {
					//msg有data数据,需要再次读取data数据
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}
					fmt.Printf("===>Receive msg:id=%d,len=%d,data=%s\n", msg.Id, msg.DataLen, string(msg.Data))
				}
			}
		}(conn)
	}
}

func B2s(b []byte) string {
	return unsafe.String(&b[0], unsafe.IntegerType(len(b)))
}

func S2b(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// 测试封包拆包 Client
func TestDataPack_Client(t *testing.T) {
	//客户端goroutine 负责模拟黏包的数据,然后进行发送
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	//1.创建一个封包对象dp
	dp := NewDataPack()
	//2.封装一个msg1包
	msg1 := &Message{
		Id:      0,
		DataLen: 5,
		//Data:    []byte{'h', 'e', 'l', 'l', 'o'},
		Data: []byte("hello"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}
	//3.封装一个msg2包
	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte("world!!"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}
	//4.将sendData1和sendData2拼接到一起，组成黏包
	sendData1 = append(sendData1, sendData2...)
	//5.向服务器端写数据
	conn.Write(sendData1)
	//客户端阻塞
	//select {}
	time.Sleep(5 * time.Second)
}
