package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Serve 运行服务器
	Serve()
	// Stop 停止服务器
	Stop()
	//路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(msgfID uint32, router IRouter)
	//获取当前server的连接管理器
	GetConnMgr() IConnManager
	//register
	SetOnConnStart(func(conn Iconnection))

	SetOnConnStop(func(conn Iconnection))
	//call
	CallOnConnStart(conn Iconnection)

	CallOnConnStop(conn Iconnection)
}
