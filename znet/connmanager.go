package znet

import (
	"errors"
	"fmt"
	"myzinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息
	connLock    sync.RWMutex                  //读写连接的读写锁
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源 map加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn连接添加到 ConnManager 中
	c.connections[conn.GetConnId()] = conn

	fmt.Println("connection add to ConnManager successfully:conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源 map加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除连接信息
	delete(c.connections, conn.GetConnId())

	fmt.Println("connection remove connId = ", conn.GetConnId(), " successfully:conn num = ", c.Len())
}

func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	//保护共享资源 map加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	conn, ok := c.connections[connId]
	if ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	//保护共享资源 map加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connId, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connId)
	}

	fmt.Println("clear all connections successfully:conn num = ", c.Len())
}

func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
