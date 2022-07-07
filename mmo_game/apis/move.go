package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

//玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	//resolve protobuf from cli
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move : Position Unmarshal Error: ", err)
		return
	}
	//get cli's playerID
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error: ", err)
		return
	}

	fmt.Printf("playerID = %d,move(%f,%f,%f,%f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
	player := WorldMgrObj.GetPlayerByPid(pid.(int32))
	//broadcast and update player's position
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
