package circletimer

import (
	"errors"
	"sync"
	"time"
)

type CircleTimer struct {
	interval time.Duration
	running  bool
	callback func()
	timer    *time.Timer

	sync.Mutex
}

func NewCircleTimer(interval time.Duration, callback func()) *CircleTimer {
	if interval < time.Millisecond {
		interval = time.Millisecond
	}

	ct := &CircleTimer{
		interval: interval,
		callback: callback,
		running:  false,
	}
	return ct
}

func (ct *CircleTimer) Start() error {
	ct.Mutex.Lock()
	defer ct.Mutex.Unlock()

	var (
		handle func()
		again  func()
	)

	if ct.running {
		return errors.New("already running")
	}

	ct.running = true

	handle = func() {
		if !ct.running {
			return
		}

		ct.callback()
		again()
	}
	again = func() {
		ct.timer = time.AfterFunc(ct.interval, handle)
	}

	again()
	return nil
}

func (ct *CircleTimer) Reset(interval time.Duration) error {
	ct.Mutex.Lock()
	defer ct.Mutex.Unlock()

	if !ct.running {
		return errors.New("already stopped")
	}

	ct.interval = interval
	ct.timer.Reset(interval)
	return nil
}

func (ct *CircleTimer) Stop() {
	ct.Mutex.Lock()
	defer ct.Mutex.Unlock()

	ct.running = false
	ct.timer.Stop()
}
