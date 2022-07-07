package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (b *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Router Handle...")
	//读取客户端数据，再回写ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), " data= ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("Hello Welcome to Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func (b *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//读取客户端数据，再回写ping..ping..ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), " data= ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn ziface.Iconnection) {
	fmt.Println("===>DoConnectionBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Set conn property                             ...")
	conn.SetProperty("name", "harry")
	conn.SetProperty("age", 18)

}

func DoConnectionLost(conn ziface.Iconnection) {
	fmt.Println("===>DoConnectionLost is called...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is lost..")

	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("name:", name)
	}
	if age, err := conn.GetProperty("age"); err == nil {
		fmt.Println("age:", age)
	}
}

func main() {
	//1.create server
	s := znet.NewServer()
	//2.register hook function
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//3.add customized router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//4.run server
	s.Serve()
}
