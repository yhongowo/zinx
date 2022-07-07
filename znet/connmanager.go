package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
 Connection Manager
*/

type ConnManager struct {
	connections map[uint32]ziface.Iconnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.Iconnection),
	}
}

func (c *ConnManager) Add(conn ziface.Iconnection) {
	//protect map, Write Lock
	c.connLock.Lock()
	defer c.connLock.Unlock()
	//add conn to map
	c.connections[conn.GetConnID()] = conn
	fmt.Println("conn ID = ", conn.GetConnID(), "Connections add to connManager success: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn ziface.Iconnection) {
	//protect map, lock
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//delete
	delete(c.connections, conn.GetConnID())
	fmt.Println("conn ID = ", conn.GetConnID(), "Remove from connManager success: conn num = ", c.Len())

}

func (c *ConnManager) Get(ConnID uint32) (ziface.Iconnection, error) {
	//read lock
	c.connLock.RLock()
	defer c.connLock.Unlock()
	if conn, ok := c.connections[ConnID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connections not found")
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	fmt.Println("Clear all connections success! conn num: ", c.Len())
}
