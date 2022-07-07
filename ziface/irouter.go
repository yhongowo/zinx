package ziface

/*
路由的抽象接口
路由里的数据都是IRequest
*/

type IRouter interface {
	//处理conn业务前的hook方法
	PreHandle(request IRequest)
	//处理conn业务的主方法
	Handle(request IRequest)
	//处理conn业务之后的hook方法
	PostHandle(request IRequest)
}
