package bit

import (
	"encoding/binary"
)

func New() IBytes {
	return &ByteArray{}
}

func With(b []byte) IBytes {
	this := &ByteArray{}
	this.SetBuffer(b)
	return this
}

func WithSize(n int) IBytes {
	return With(make([]byte, n))
}

//bytes
type ByteArray struct {
	Buffer
	endian binary.ByteOrder
}

func (this *ByteArray) EndianSet(v binary.ByteOrder) {
	this.endian = v
}

func (this *ByteArray) Endian() binary.ByteOrder {
	if this.endian == nil {
		return binary.BigEndian //默认大端
	}
	return this.endian
}

//特殊的io(int=int64 int=int32)
func (this *ByteArray) Rint() int {
	return int(this.Rint32())
}

func (this *ByteArray) Wint(v int) {
	this.Wint32(int32(v))
}

//read
func (this *ByteArray) Rbool() bool {
	return this.Ruint8() != 0
}

func (this *ByteArray) Rint8() int8 {
	return int8(this.Ruint8())
}

func (this *ByteArray) Rint16() int16 {
	return int16(this.Ruint16())
}

func (this *ByteArray) Rint32() int32 {
	return int32(this.Ruint32())
}

func (this *ByteArray) Rint64() int64 {
	return int64(this.Ruint64())
}

func (this *ByteArray) Ruint8() (v uint8) {
	v = this.Bit(this.p)
	this.Move(1)
	return
}

func (this *ByteArray) Ruint16() (v uint16) {
	v = this.Endian().Uint16(this.Payload())
	this.Move(2)
	return
}

func (this *ByteArray) Ruint32() (v uint32) {
	v = this.Endian().Uint32(this.Payload())
	this.Move(4)
	return
}

func (this *ByteArray) Ruint64() (v uint64) {
	v = this.Endian().Uint64(this.Payload())
	this.Move(8)
	return
}

func (this *ByteArray) Rstr() (str string) {
	if n := this.Ruint32(); n > 0 {
		b := make([]byte, n)
		this.Read(b)
		str = string(b)
	}
	return
}

//write int
func (this *ByteArray) Wbool(ok bool) {
	if ok {
		this.Wuint8(1)
	} else {
		this.Wuint8(0)
	}
}

func (this *ByteArray) Wint8(v int8) {
	this.Wuint8(uint8(v))
}

func (this *ByteArray) Wint16(v int16) {
	this.Wuint16(uint16(v))
}

func (this *ByteArray) Wint32(v int32) {
	this.Wuint32(uint32(v))
}

func (this *ByteArray) Wint64(v int64) {
	this.Wuint64(uint64(v))
}

//write uint
func (this *ByteArray) Wuint8(v uint8) {
	this.grow(1)
	this.BitSet(this.p, v)
	this.Move(1)
}

func (this *ByteArray) Wuint16(v uint16) {
	this.grow(2)
	this.Endian().PutUint16(this.Payload(), v)
	this.Move(2)
}

func (this *ByteArray) Wuint32(v uint32) {
	this.grow(4)
	this.Endian().PutUint32(this.Payload(), v)
	this.Move(4)
}

func (this *ByteArray) Wuint64(v uint64) {
	this.grow(8)
	this.Endian().PutUint64(this.Payload(), v)
	this.Move(8)
}

func (this *ByteArray) Wstr(str string) {
	b := []byte(str)
	if n := len(b); n > 0 {
		this.Wuint32(uint32(n))
		this.Write(b)
	} else {
		this.Wuint32(0)
	}
}

/*
 * 读取一部分写入到对象
 */
func (this *ByteArray) ReadAny(b IBuffer, n int) int {
	if n < 1 { //read all
		n = b.Write(this.buf[this.p:])
	} else {
		n = b.Write(this.buf[this.p : this.p+n])
	}
	this.Move(n)
	return n
}
