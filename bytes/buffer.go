package bytes

func NewBuffer(b []byte) Buffer {
	return Buffer{buf: b}
}

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

func (this *Buffer) Next(n int) {
	this.p += n
}

func (this *Buffer) Len() int {
	return len(this.buf)
}

func (this *Buffer) LenSet(n int) {
	if m := len(this.buf); n > m {
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
	return len(this.buf) - this.p
}

func (this *Buffer) Write(b []byte) (n int, err error) {
	this.grow(len(b))
	n = copy(this.buf[this.p:], b)
	this.Next(n)
	return
}

func (this *Buffer) Read(b []byte) (n int, err error) {
	if n = len(b); n > 0 {
		n = copy(b, this.buf[this.p:])
		this.Next(n)
	}
	return
}

func (this *Buffer) Bytes() []byte {
	return this.buf
}

func (this *Buffer) String() string {
	return string(this.buf)
}

//private
func (this *Buffer) grow(size int) int {
	if n := size - this.Available(); n > 0 {
		newBuf := makeBytes(len(this.buf) + n)
		copy(newBuf, this.buf)
		this.buf = newBuf
		return n
	}
	return 0
}

func (this *Buffer) payload() []byte {
	return this.buf[this.p:]
}
