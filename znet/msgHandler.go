package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandle struct {
	//存放每个msgID所对应的方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的wordker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! Need to register")
	}
	//2 根据MsgID 调度对应router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的API方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, "success!")
}

//启动一个worker工作池(开启工作池的动作只能发生一次)
func (mh *MsgHandle) StartWorkerPoll() {
	//根据PoolSize分别开启worker，每个worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//为当前的worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started ...")

	for {
		select {
		//如果有消息过来，出列一个客户端的request，执行业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给taskqueue，由worker处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不同的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), "request MsgID = ",
		request.GetMsgID(), "to WorkerID = ", workerID)
	mh.TaskQueue[workerID] <- request
	//2将消息发送给对应的worker的taskQueue即可
}
