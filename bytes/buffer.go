package bytes

func NewBuffer(b []byte) Buffer {
	return Buffer{buf: b, n: len(b)}
}

type Buffer struct {
	p   int
	n   int
	buf []byte
}

func (this *Buffer) SetBuffer(b []byte) {
	this.p = 0
	this.buf = b
	this.n = len(b)
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

func (this *Buffer) Reset() {
	this.p = 0
	this.n = 0
}

func (this *Buffer) Next(n int) {
	this.p += n
}

func (this *Buffer) Len() int {
	return this.n
}

func (this *Buffer) CapLen() int {
	return len(this.buf)
}

func (this *Buffer) LenSet(n int) {
	this.n = n
	if m := this.CapLen(); n > m {
		this.grow(n - m)
	} else {
		this.buf = this.buf[:n]
	}
	if this.p > n {
		this.p = n
	}
}

func (this *Buffer) Bit(n int) byte {
	return this.buf[n]
}

func (this *Buffer) BitSet(n int, b byte) {
	this.buf[n] = b
}

func (this *Buffer) Available() int {
	return this.Len() - this.p
}

func (this *Buffer) Write(b []byte) (n int, err error) {
	this.grow(len(b))
	n = copy(this.payload(), b)
	this.Next(n)
	return
}

func (this *Buffer) Read(b []byte) (n int, err error) {
	if n = len(b); n > 0 {
		n = copy(b, this.payload())
		this.Next(n)
	}
	return
}

func (this *Buffer) Bytes() []byte {
	return this.buf[:this.n]
}

func (this *Buffer) String() string {
	return string(this.Bytes())
}

//private
func (this *Buffer) grow(size int) int {
	if n := size - this.Available(); n > 0 {
		newBuf := makeBytes(this.CapLen() + n)
		copy(newBuf, this.buf)
		this.buf = newBuf
		this.n = len(newBuf)
		return n
	}
	return 0
}

func (this *Buffer) payload() []byte {
	return this.buf[this.p:this.n]
}
