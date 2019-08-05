package bytes

import (
	"io"
)

type IWriter interface {
	io.Writer
	Wbool(bool)
	Wstr(string)
	//int
	Wint(int)
	Wint8(int8)
	Wint16(int16)
	Wint32(int32)
	Wint64(int64)
	//uint
	//Wuint(uint)
	Wuint8(uint8)
	Wuint16(uint16)
	Wuint32(uint32)
	Wuint64(uint64)
	//other
	Wobj(interface{})
	WriteTo(io.Writer, int) int
}

type IReader interface {
	io.Reader
	Rbool() bool
	Rstr() string
	//int
	Rint() int
	Rint8() int8
	Rint16() int16
	Rint32() int32
	Rint64() int64
	//uint
	//Ruint() uint
	Ruint8() uint8
	Ruint16() uint16
	Ruint32() uint32
	Ruint64() uint64
	//other
	Robj(interface{})
}

type IBuffer interface {
	IReader
	IWriter
	//pos
	Pos() int
	Seek(int)
	SeekBegin()
	SeekEnd() int
	//len
	Reset()
	Len() int
	LenSet(int)
	CapLen() int
	Available() int
	//bytes
	Bytes() []byte
}

//i/o object
type ReadObj interface {
	Read(IReader)
}

type WriteObj interface {
	Write(IWriter)
}
