package znet

// 定义应用层的消息结构体来解决粘包的问题
// 将请求的消息封装在message中
type Message struct {
	DataLen uint32 // 不包含头（1字节）的数据长度
	Id      uint32
	Data    []byte
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetDate() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetMsgLen(dataLen uint32) {
	m.DataLen = dataLen
}
func (m *Message) SetDate(data []byte) {
	m.Data = data
}
