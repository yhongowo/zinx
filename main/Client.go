package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

//模拟客户端
func main() {
	fmt.Println("Client start...")

	time.Sleep(1 * time.Second)
	//1.直接连接远程服务器 得到conn连接
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Client start err,exit!")
		return
	}

	for {
		//发送封包的msessage消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx v0.8 client test message")))
		if err != nil {
			fmt.Println("pack error: ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error")
			return
		}
		//服务器应该给我们回复一个message， msgID:1 ping..ping..ping

		//读取流中的head 得到ID, datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}
		//将二进制的head拆包到msg结构体
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			//msg has data, so read data
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data err: ", err)
				return
			}

			fmt.Println("----> Recv Server Msg: ID= ", msg.ID, "len= ", msg.Datalen, "data= ", string(msg.GetData()))
		}
		time.Sleep(1 * time.Second)
	}

	//2.连接调用Write 写数据
}
