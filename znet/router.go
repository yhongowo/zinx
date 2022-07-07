package znet

import (
	"zinx/ziface"
)

//实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

//有的业务不需要PreHandle或PostHandle业务
func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
