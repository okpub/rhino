package bytes

import (
	"fmt"
	"reflect"
)

//write
func (this *ByteArray) Wobj(v interface{}) {
	switch p := v.(type) {
	case WriteObj:
		p.Write(this)
	//case IWriter:
	//	p.WriteTo(this, 0) //this Will grow
	default:
		this.encode(reflect.ValueOf(v))
	}
}

//read
func (this *ByteArray) Robj(v interface{}) {
	switch p := v.(type) {
	case ReadObj:
		p.Read(this)
	//case IWriter:
	//	this.WriteTo(p, 0) //p Will grow
	default:
		this.decode(reflect.ValueOf(v))
	}
}

func (b *ByteArray) encode(fd reflect.Value) {
	switch fd.Kind() {
	case reflect.Ptr:
		b.encode(fd.Elem())
	case reflect.Interface:
		b.encode(fd.Elem())
	case reflect.Struct:
		for i := 0; i < fd.NumField(); i++ {
			if item := fd.Field(i); item.CanSet() {
				b.encode(item)
			}
		}
	case reflect.Array:
		for i := 0; i < fd.Len(); i++ {
			b.encode(fd.Index(i))
		}
	case reflect.Slice:
		for i := 0; i < fd.Len(); i++ {
			b.encode(fd.Index(i))
		}
	case reflect.Bool:
		b.Wbool(fd.Bool())
	case reflect.String:
		b.Wstr(fd.String())
	case reflect.Int8:
		b.Wint8(int8(fd.Int()))
	case reflect.Int16:
		b.Wint16(int16(fd.Int()))
	case reflect.Int32:
		b.Wint32(int32(fd.Int()))
	case reflect.Int64:
		b.Wint64(fd.Int())
	case reflect.Uint8:
		b.Wuint8(uint8(fd.Uint()))
	case reflect.Uint16:
		b.Wuint16(uint16(fd.Uint()))
	case reflect.Uint32:
		b.Wuint32(uint32(fd.Uint()))
	case reflect.Uint64:
		b.Wuint64(fd.Uint())
	case reflect.Int:
		b.Wint32(int32(fd.Int()))
	case reflect.Uint:
		b.Wuint32(uint32(fd.Uint()))
	default:
		panic(fmt.Errorf("fail write type %s", fd.Kind().String()))
	}
}

func (b *ByteArray) decode(fd reflect.Value) {
	switch fd.Kind() {
	case reflect.Ptr:
		b.decode(fd.Elem())
	case reflect.Interface:
		b.decode(fd.Elem())
	case reflect.Struct:
		for i := 0; i < fd.NumField(); i++ {
			if item := fd.Field(i); item.CanSet() {
				b.decode(item)
			}
		}
	case reflect.Array:
		for i := 0; i < fd.Len(); i++ {
			b.decode(fd.Index(i))
		}
	case reflect.Slice:
		for i := 0; i < fd.Len(); i++ {
			b.decode(fd.Index(i))
		}
	case reflect.Bool:
		fd.SetBool(b.Rbool())
	case reflect.String:
		fd.SetString(b.Rstr())

	case reflect.Int:
		fd.SetInt(int64(b.Rint32()))
	case reflect.Int8:
		fd.SetInt(int64(b.Rint8()))
	case reflect.Int16:
		fd.SetInt(int64(b.Rint16()))
	case reflect.Int32:
		fd.SetInt(int64(b.Rint32()))
	case reflect.Int64:
		fd.SetInt(b.Rint64())

	case reflect.Uint:
		fd.SetUint(uint64(b.Ruint32()))
	case reflect.Uint8:
		fd.SetUint(uint64(b.Ruint8()))
	case reflect.Uint16:
		fd.SetUint(uint64(b.Ruint16()))
	case reflect.Uint32:
		fd.SetUint(uint64(b.Ruint32()))
	case reflect.Uint64:
		fd.SetUint(b.Ruint64())

	default:
		panic(fmt.Errorf("fail read type %s", fd.Kind().String()))
	}
}

func makeBytes(n int) []byte {
	defer func() {
		if recover() != nil {
			panic(fmt.Errorf("bytes.Buffer: too large"))
		}
	}()
	return make([]byte, n)
}
