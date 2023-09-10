package ziface

type IMessage interface {
	GetMsgId() uint32   //获取消息ID
	GetData() []byte    //获取消息内容
	GetDataLen() uint32 //获取消息数据段长度

	SetMsgId(msgId uint32) //设置消息ID
	SetData(data []byte)   //设置消息内容
	SetDataLen(uint32)     //设置消息数据段长度
}
