package znet

import "C"
import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

//连接模块
type Connection struct {
	//从属于哪个server
	TcpServer ziface.IServer
	//当前连接的socket
	Conn *net.TCPConn
	//连接ID
	ConnID uint32
	//连接状态
	isClosed bool
	//告知当前连接退出
	ExitChan chan bool
	//无缓冲管道，用于读写goroutine之间的消息通信
	msgChan chan []byte
	//消息的管理msgID和对应的处理业务api关系
	MsgHandler ziface.IMsgHandle
	//连接属性集合
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID: ", c.ConnID, "[Reader is exit!], remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流 8byte
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err: ", err)
			break
		}
		//拆包，得到MsgID, msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据dataLen 再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err: ", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的Request
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启工作池，将消息发送给worker工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

//写消息goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[Conn Writer exit!]")
	//不断阻塞的等待channel消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//have data to send to client
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.ExitChan:
			//reader exit, now exit writer
			return
		}
	}
}

//提供一个SendMsg方法 将要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	//将data封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("package err, msg id: ", msgId)
		return errors.New("pack error msg")
	}
	//send to writer
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID: ", c.ConnID)
	//启动当前连接的读数据业务
	go c.StartReader()
	//启动当前连接的写数据业务
	go c.StartWriter()

	//do hook function
	c.TcpServer.CallOnConnStart(c)

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID: ", c, c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true
	//do hook function, before conn is close
	c.TcpServer.CallOnConnStop(c)
	//close socket
	c.Conn.Close()
	//tell writer to close
	c.ExitChan <- true
	//close channel

	//从ConnManager中删除
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加连接属性
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//delete
	delete(c.property, key)
}
