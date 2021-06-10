package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/rfyiamcool/circletimer"
)

func main() {
	var (
		incr   int32 = 0
		notify       = make(chan struct{}, 0)
	)

	handle := func() {
		num := atomic.AddInt32(&incr, 1)
		fmt.Printf("time: %v ; incr: %v  \n", time.Now().Format("2006-01-02 15:04:05.000"), num)

		if incr >= 10 {
			notify <- struct{}{}
		}
	}

	timer := circletimer.NewCircleTimer(1*time.Second, handle)

	time.AfterFunc(5*time.Second, func() {
		timer.Reset(100 * time.Millisecond) // timer reset
	})

	timer.Start()
	<-notify
	timer.Stop()

	// check if timer is still running ?
	time.Sleep(3 * time.Second)
	fmt.Println(atomic.LoadInt32(&incr))
}
