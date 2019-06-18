package mysql

import (
	"context"
	"fmt"
	"time"
)

type UID_T int

type TestObj struct {
	UID  UID_T     `mysql:"uid"`
	Bat  int       `mysql:"begin_at"`
	Eat  int       `mysql:"end_at"`
	Num  int8      `mysql:"num"`
	Name string    `mysql:"name"`
	Date time.Time `mysql:"time"`
}

func init() {
	db := NewMysql().Open()

	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)

	items, err := Table(db.QueryContext(ctx, "select * from test"))
	for _, item := range items {
		obj := TestObj{}
		item.Obj(&obj)
		fmt.Printf("%+v \n", obj)
	}
	fmt.Println(err)
	fmt.Printf("%p \n", &defaultAddr)
	fmt.Printf("%p \n", &defaultAddr)
	fmt.Println("借宿")
	time.Sleep(time.Millisecond)
}
