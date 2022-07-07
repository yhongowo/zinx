package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只负责测试Datapack拆包封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	//创建一个go 负责从客户端处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err: ", err)
			}
			go func(conn net.Conn) {
				//handle request
				//------unpack process------
				dp := NewDataPack()
				for {
					//1.read from conn, read head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err: ", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg有数据，进行二次读取
						//2.read from conn, read data by datalen
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						//read completed
						fmt.Println("-->recv msgID: ", msg.ID, "datalen: ", msg.Datalen, "data:", string(msg.Data))
					}
				}
			}(conn)

		}
	}()
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Client err: ", err)
		return
	}
	//创建一个封包对象
	dp := NewDataPack()
	//模拟黏包过程，封装两个msg一起发送
	//first msg package
	msg1 := &Message{
		ID:      0,
		Datalen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	//second msg package
	msg2 := &Message{
		ID:      1,
		Datalen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err: ", err)
		return
	}
	//put them together, send to server
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	//阻塞
	select {}
}
