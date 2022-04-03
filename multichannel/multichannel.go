// see https://github.com/SierraSoftworks/multicast
package multichannel

import (
	"sync"

	"github.com/YukiMiyatake/HereWeGo/generics/list"
)

type Listener[T any] struct {
	C chan T
	E *list.Element[*Listener[T]]
}

func (s *Listener[T]) Close() {
	close(s.C)
}

type Channel[T any] struct {
	C chan<- T
	c chan T
	l *list.List[*Listener[T]]
	m sync.Mutex
}

func New[T any]() *Channel[T] {
	c := make(chan T)
	l := list.New[*Listener[T]]()

	go func() {
		for v := range c {
			if l != nil {
				for e := l.Front(); e != nil; e = e.Next() {
					e.Value.C <- v
				}
			}
		}

		if l != nil {
			for e := l.Front(); e != nil; e = e.Next() {
				e.Value.Close()
			}
		}
	}()

	return &Channel[T]{
		C: c,
		c: c,
		l: l,
	}
}

func (c *Channel[T]) Add() *Listener[T] {
	c.m.Lock()
	defer c.m.Unlock()

	l := &Listener[T]{C: make(chan T, 0)}
	e := c.l.PushBack(l)
	l.E = e

	return l
}

func (c *Channel[T]) Remove(l *Listener[T]) {
	c.m.Lock()
	defer c.m.Unlock()

	c.l.Remove(l.E)
	l.Close()
}

func (c *Channel[T]) Close() {
	c.m.Lock()
	defer c.m.Unlock()

	close(c.c)
}
