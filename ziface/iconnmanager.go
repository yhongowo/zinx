package ziface

type IConnManager interface {
	//add conn
	Add(conn Iconnection)
	//del conn
	Remove(conn Iconnection)
	//get conn by connID
	Get(ConnID uint32) (Iconnection, error)
	//get 连接总数
	Len() int
	//clear all conn
	ClearConn()
}
