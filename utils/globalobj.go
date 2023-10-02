package utils

import (
	"encoding/json"
	"fmt"
	"myzinx/ziface"
	"os"
)

// GlobalObj 存储一切有关zinx框架的全局参数，供其他模块使用
// 一些参数也可以通过用户根据zinx.json来配置
// 相当于 serviceContext
type GlobalObj struct {
	TcpServer ziface.IServer //当前zinx的全局 Server 对象
	Host      string         `json:"host"`    //当前服务器主机IP
	TcpPort   int            `json:"tcpPort"` //当前服务器主机监听端口号
	Name      string         `json:"name"`    //当前服务器名称
	Version   string         `json:"version"` //当前zinx版本号

	MaxPacketSize uint32 `json:"maxPacketSize"` //读取数据包的最大值
	MaxConn       int    `json:"maxConn"`       //当前服务器主机允许的最大连接个数
}

// GlobalObject 定义一个全局的对象
var GlobalObject *GlobalObj

// Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	fmt.Printf("reload json:%s\n", data)
	//注意原生的json.Unmarshal 在反序列化long时会丢失经度，要使用go-zero的jsonx.Unmarshal
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供init方法，默认加载
func init() {
	//初始化GlobalObject变量
	GlobalObject = &GlobalObj{
		Name:          "ZinxServerApp",
		Version:       "v0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
		//TcpServer:     nil,
	}
	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
