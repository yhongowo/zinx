package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	//Server
	TcpServer ziface.IServer //全局Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口
	Name      string         //服务器名称

	//Zinx
	Version          string //Zinx版本号
	MaxConn          int    //允许的最大连接数
	MaxPacketSize    uint32 //数据包最大值
	WorkerPoolSize   uint32 //当前业务工作worker池的goroutine数量
	MaxWorkerTaskLen uint32 //Zinx框架允许用户最多开辟多少个worker（限定条件）
}

//定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//初始化当前的GlobalObject
func init() {
	//如果配置文件未加载，默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.5",
		TcpPort:          9000,
		MaxConn:          1000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024, //每个worker的消息队列处理的最大任务数量
	}
	//尝试从conf/中加载
	GlobalObject.Reload()
}
