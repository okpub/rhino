package network

import (
	"fmt"

	"github.com/okpub/rhino/bytes"
)

/*
 * 基础网络包
 * msgID: 0表示异步, >0表示同步消息，(必须回传, 不回传会导致客户端阻塞)
 */
type SocketPacket struct {
	bytes.ByteArray
	Topic int
	MsgID int64 //默认异步(0)
}

//write
func (this *SocketPacket) WriteBegin() *SocketPacket {
	this.SeekBegin()
	this.Wint(0)
	this.Wobj(this) //默认拷贝
	return this
}

func (this *SocketPacket) With(args ...interface{}) *SocketPacket {
	for _, v := range args {
		this.Wobj(v)
	}
	return this
}

func (this *SocketPacket) Flush() []byte {
	this.SeekBegin()
	this.Wint(this.Len() - NET_Paylen) //write paylen
	//this.Wobj(this)                    //write header
	this.SeekEnd()
	return this.Bytes()
}

func (this *SocketPacket) Sync() bool {
	return this.MsgID > 0
}

func (this *SocketPacket) Cmd() int {
	return this.Topic
}

//read
func (this *SocketPacket) ReadBegin() {
	this.Seek(NET_Paylen)
	this.Robj(this)
}

func (this *SocketPacket) String() string {
	return fmt.Sprintf("[cmd=%#X len=%d id=%d]", this.Cmd(), this.Len(), this.MsgID)
}

/*
 * 读包
 */
func ReadBegin(b []byte) *SocketPacket {
	this := &SocketPacket{}
	this.SetBuffer(b)
	this.ReadBegin()
	return this
}

/*
 * 写包(这里可以通过规定增长提高效率)
 */
func WriteBegin(cmd int, args ...interface{}) *SocketPacket {
	this := &SocketPacket{Topic: cmd}
	return this.WriteBegin().With(args...)
}

//同步回复
func WriteSync(id int64, cmd int, args ...interface{}) *SocketPacket {
	this := &SocketPacket{MsgID: id, Topic: cmd}
	return this.WriteBegin().With(args...)
}

//同样返回
func WriteCopy(pack *SocketPacket, args ...interface{}) *SocketPacket {
	this := &SocketPacket{MsgID: pack.MsgID, Topic: pack.Cmd()}
	return this.WriteBegin().With(args...)
}
