package mysql

import (
	"fmt"
)

type AddrOption func(*Address)

//mysql addr
type Address struct {
	Host string
	Port int
	User string
	Pwd  string
	Name string
}

func (this *Address) Addr() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", this.User, this.Pwd, fmt.Sprint(this.Host, ":", this.Port), this.Name)
}

func (this Address) Copy(opts ...AddrOption) *Address {
	for _, o := range opts {
		o(&this)
	}
	return &this
}
