package gpool

import (
	"testing"
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

func TestPool(t *testing.T) {
	p := NewPool(3, 10000)

	rep := make(chan *Reply, 1000)
	for i := 0; i < 100; i++ {
		req := &Req{i, 1}
		job := JobParam{Request: req, Reply: rep, Func: add}
		p.AddJob(&job)
	}

	count := 0
	for ch := range rep {
		t.Logf("%v + %v = %v", ch.req.a, ch.req.b, ch.c)
		count++
		if count == 100 {
			break
		}
	}
	p.Stop()
}

func TestStop(t *testing.T) {
	p := NewPool(3, 1000)

	rep := make(chan *Reply, 1000)
	for i := 0; i < 100; i++ {
		req := &Req{i, 1}
		job := JobParam{Request: req, Reply: rep, Func: add}
		p.AddJob(&job)
	}

	go func() {
		count := 0
		for ch := range rep {
			t.Logf("%v + %v = %v", ch.req.a, ch.req.b, ch.c)
			count++
		}
		t.Logf("count = %v", count)
	}()
	time.Sleep(time.Second)
	p.Stop()
	time.Sleep(time.Second)
	close(rep)
	time.Sleep(time.Second)
}

func TestStopGracefully(t *testing.T) {
	p := NewPool(3, 1000)

	rep := make(chan *Reply, 1000)
	for i := 0; i < 60; i++ {
		req := &Req{i, 1}
		job := JobParam{Request: req, Reply: rep, Func: add}
		p.AddJob(&job)
	}

	go func() {
		count := 0
		for ch := range rep {
			t.Logf("%v + %v = %v", ch.req.a, ch.req.b, ch.c)
			count++
		}
		t.Logf("count = %v", count)
	}()
	time.Sleep(time.Second)
	p.StopGracefully()
	time.Sleep(time.Second * 4)
	close(rep)
	time.Sleep(time.Second)
}
