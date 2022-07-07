package ziface

/***
封包，拆包
直接面向TCP链接中的数据流，处理TCP黏包问题
*/

type IDataPack interface {
	GetHeadLen() int32

	Pack(msg IMessage) ([]byte, error)

	Unpack([]byte) (IMessage, error)
}
