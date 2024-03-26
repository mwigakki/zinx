package ziface

type IDataPack interface {
	GetHeadLen() uint32 // 得到应用层包头的长度
	Pack(IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
