package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler *MsgHandle
	//connection manager
	ConnMgr ziface.IConnManager
	//创建连接后自动调用的hook函数
	OnConnStart func(conn ziface.Iconnection)
	//销毁连接前自动调用的hook函数
	OnConnStop func(conn ziface.Iconnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listen at IP: %s, Port:%d is starting",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		//0.开启消息队列以及Worker工作池
		s.MsgHandler.StartWorkerPoll()
		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error: ", err)
			return
		}

		//2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("Start Zinx server success,", s.Name, "Listening...")
		var cid uint32
		cid = 0

		//3.阻塞的等待客户端连接，处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}
			//设置最大连接个数的判断，如果超过最大链接，则关闭此链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO response error to Client
				fmt.Println("========>>Too many connections = ", utils.GlobalObject.MaxConn, "<<========")
				conn.Close()
				continue
			}
			//将处理新连接的业务方法和conn绑定，得到连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前的连接业务
			go dealConn.Start()
		}
	}()

}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//阻塞
	select {}
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name: ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

//register OnConnStart
func (s *Server) SetOnConnStart(f func(conn ziface.Iconnection)) {
	s.OnConnStart = f
}

//register OnConnStop
func (s *Server) SetOnConnStop(f func(conn ziface.Iconnection)) {
	s.OnConnStop = f
}

//do OnConnStart
func (s *Server) CallOnConnStart(conn ziface.Iconnection) {
	if s.OnConnStart != nil {
		fmt.Println("---->Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//do OnConnStop
func (s *Server) CallOnConnStop(conn ziface.Iconnection) {
	if s.OnConnStop != nil {
		fmt.Println("---->Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
