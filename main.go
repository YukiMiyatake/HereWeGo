package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/YukiMiyatake/HereWeGo/multichannel"
)

var ll *multichannel.Listener[interface{}]

func bus_func(i int, wg *sync.WaitGroup, ch *multichannel.Channel[interface{}]) {

	l := ch.Listen()
	wg.Done()
	defer wg.Done()

	if i == 1 {
		ll = l
	}

	for {
		select {
		case v := <-l.C:
			if v == nil {
				return
			}
			fmt.Println(i, v)
		}
	}

}

func main() {
	c := multichannel.New[interface{}]()
	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go bus_func(i, &wg, c)
	}

	wg.Wait()
	wg.Add(3) //結果受け取り用

	c.C <- "1 Hello World!"
	time.Sleep(time.Second)

	c.C <- "2 Hello World!"
	c.Remove(ll)
	time.Sleep(time.Second)

	c.C <- "3 Hello World!"
	time.Sleep(time.Second)
	c.C <- "4 Hello World!"

	c.Close()
	wg.Wait()

}
