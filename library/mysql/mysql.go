package mysql

import (
	"database/sql"
	"fmt"
	"time"
	//init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//default values
var (
	defaultAddr = Address{
		Host: "localhost",
		Port: 3306,
		User: "root",
		Pwd:  "123456",
		Name: "test",
	}

	defaultConfig = Options{
		Addr:            defaultAddr.Addr(),
		MaxIdle:         10,
		MaxActive:       100,
		IdleTimeout:     0,
		MaxConnLifetime: 0,
	}
)

func NewMysql(opts ...Option) *Options {
	return defaultConfig.Copy(opts...)
}

//class mysql Config
type Option func(*Options)

type Options struct {
	Addr            string
	MaxIdle         int
	MaxActive       int
	MaxConnLifetime time.Duration
	IdleTimeout     time.Duration //not used
}

func (this Options) Copy(opts ...Option) *Options {
	for _, o := range opts {
		o(&this)
	}
	return &this
}

func (this *Options) Open() *sql.DB {
	db, err := sql.Open("mysql", this.Addr)
	if err != nil {
		panic(fmt.Errorf("mysql underlying engine err:" + err.Error()))
	}
	return this.Filler(db)
}

func (this *Options) Filler(db *sql.DB) *sql.DB {
	db.SetMaxIdleConns(this.MaxIdle)
	db.SetMaxOpenConns(this.MaxActive)
	db.SetConnMaxLifetime(this.MaxConnLifetime)
	return db
}

//选项 更改mysql地址
func OptionAddr(addr Address) Option {
	return func(p *Options) {
		p.Addr = addr.Addr()
	}
}
