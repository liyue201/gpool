package main

import (
	"fmt"
	"github.com/liyue201/gpool"
	"time"
)

type Req struct {
	a int
	b int
}

type Reply struct {
	req *Req
	c   int
	err error
}

func add(request interface{}, replyCh interface{}) {
	r := request.(*Req)
	rpCh := replyCh.(chan *Reply)
	rp := &Reply{
		req: r,
		c:   r.a + r.b,
	}
	rpCh <- rp
	time.Sleep(time.Second / 10)
}

func main() {
	p := gpool.NewPool(3, 1000)

	repChan := make(chan *Reply, 100)
	for i := 0; i < 10; i++ {
		req := &Req{i, 1}
		job := gpool.JobParam{Request: req, Reply: repChan, Func: add}
		p.AddJob(&job)
	}
	for i := 0; i < 10; i++ {
		ch := <-repChan
		fmt.Printf("%v + %v = %v\n", ch.req.a, ch.req.b, ch.c)
	}
	p.Stop()
}
