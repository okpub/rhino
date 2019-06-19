package cache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

//default values
var (
	defaultConfig = Options{
		Wait:            false, //如果超过最大并发，那么不会等待，直接返回错误conn
		MaxIdle:         10,    //空闲连接数目
		MaxActive:       100,   //最大连接数目
		IdleTimeout:     0,
		MaxConnLifetime: 0,
		DialTimeout:     time.Second * 3,
		Addr:            "192.168.0.100:6379",
	}
)

type Option func(*Options)

//class redis Options
type Options struct {
	Addr            string
	Wait            bool
	MaxIdle         int
	MaxActive       int
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
	DialTimeout     time.Duration
}

/*
* 修改个别配置，获得新的配置表来启动redis
 */
func (this Options) Copy(opts ...Option) *Options {
	for _, o := range opts {
		o(&this)
	}
	return &this
}

/*
* 通过配置表来填充redis池的配置选项
 */
func (this *Options) Filler(p *redis.Pool) *redis.Pool {
	var (
		addr    = this.Addr
		timeout = this.DialTimeout
	)
	p.Wait = this.Wait
	p.MaxIdle = this.MaxIdle
	p.MaxActive = this.MaxActive
	p.IdleTimeout = this.IdleTimeout
	p.MaxConnLifetime = this.MaxConnLifetime
	p.Dial = func() (redis.Conn, error) { return redis.Dial("tcp", addr, redis.DialConnectTimeout(timeout)) }
	return p
}

func (this *Options) String() string {
	return fmt.Sprintf("[addr=%s min=%d max=%d]", this.Addr, this.MaxIdle, this.MaxActive)
}

/*
* 通过配置新建一个连接池
 */
func (this *Options) Open() *redis.Pool {
	return this.Filler(&redis.Pool{})
}

/*
*通过默认配置修改一些选项来获得新配置
 */
func NewRedis(opts ...Option) *Options {
	return defaultConfig.Copy(opts...)
}

//选项(使用新配置)
func OptionAddr(addr string) Option {
	return func(p *Options) {
		p.Addr = addr
	}
}

func OptionHost(host string) Option {
	return func(p *Options) {
		p.Addr = host + ":6379"
	}
}
