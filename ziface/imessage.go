package ziface

// 将请求的消息封装在message中
type IMessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetDate() []byte

	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetDate([]byte)
}
