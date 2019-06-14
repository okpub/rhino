package bit

type Buffer struct {
	p   int
	buf []byte
}

func (this *Buffer) SetBuffer(b []byte) {
	this.p = 0
	this.buf = b
}

func (this *Buffer) Pos() int {
	return this.p
}

func (this *Buffer) Seek(p int) {
	this.p = p
}

func (this *Buffer) SeekBegin() {
	this.p = 0
}

func (this *Buffer) SeekEnd() int {
	this.p = this.Len()
	return this.p
}

func (this *Buffer) Move(n int) {
	this.p += n
}

func (this *Buffer) Len() int {
	return len(this.buf)
}

func (this *Buffer) LenSet(n int) {
	this.SeekBegin()
	if m := len(this.buf); n > m {
		this.grow(n - m)
	} else {
		this.buf = this.buf[:n]
	}
}

func (this *Buffer) Available() int {
	return len(this.buf) - this.p
}

func (this *Buffer) Bit(n int) byte {
	return this.buf[n]
}

func (this *Buffer) BitSet(n int, b byte) {
	this.buf[n] = b
}

func (this *Buffer) Write(b []byte) (n int) {
	this.grow(len(b))
	n = copy(this.buf[this.p:], b)
	this.Move(n)
	return
}

func (this *Buffer) Read(b []byte) (n int) {
	if n = len(b); n > 0 {
		n = copy(b, this.buf[this.p:])
		this.Move(n)
	}
	return
}

func (this *Buffer) grow(size int) int {
	if n := size - this.Available(); n > 0 {
		buf := make([]byte, len(this.buf)+n)
		copy(buf, this.buf)
		this.buf = buf
		return n
	}
	return 0
}

func (this *Buffer) Payload() []byte {
	return this.buf[this.p:]
}

func (this *Buffer) Bytes() []byte {
	return this.buf
}

func (this *Buffer) String() string {
	return string(this.buf)
}
