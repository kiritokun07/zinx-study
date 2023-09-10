package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"myzinx/utils"
	"myzinx/ziface"
)

// DataPack 封包拆包类示例，暂时不需要成员
type DataPack struct {
}

// NewDataPack 封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//Id uint32（4字节）+DataLen uint32（4字节）
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放[]byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//写data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	//创建一个输入二进制数据的 ioReader
	dataBuff := bytes.NewReader(data)
	//只解压head的信息，得到dataLen和msgId
	msg := &Message{}
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen的长度是否超出允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data receive")
	}
	//这里只需要吧head的数据拆包出来就可以了，然后通过head的长度，再从conn读取一次数据
	return msg, nil
}
