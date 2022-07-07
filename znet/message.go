package znet

type Message struct {
	ID      uint32 //消息ID
	Datalen uint32 //消息长度
	Data    []byte //消息
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		Datalen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.Datalen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.SetMsgId(id)
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.Datalen = len
}
