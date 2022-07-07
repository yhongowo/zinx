package ziface

import "net"

//连接模块的抽象层
type Iconnection interface {
	//启动连接 让当前的连接准备开始工作
	Start()
	//停止链接 结束当前连接的工作
	Stop()
	//获取当前连接的绑定socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前连接模块的ID
	GetConnID() uint32
	//获取远程客户端的TCP状态 IP PORT
	RemoteAddr() net.Addr
	//发送数据 将数据发送给远程客户端
	SendMsg(msgId uint32, data []byte) error
	//设置连接属性
	SetProperty(key string, value interface{})
	//获取连接属性
	GetProperty(key string) (interface{}, error)
	//移除连接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
