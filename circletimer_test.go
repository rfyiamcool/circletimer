package circletimer

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircleTimer(t *testing.T) {
	var (
		wg      = sync.WaitGroup{}
		counter = 0
	)

	wg.Add(100)

	fn := func() {
		counter++
		if counter <= 100 {
			wg.Done()
		}
		t.Log("trigger -> ", counter)
	}

	ct := NewCircleTimer(50*time.Millisecond, fn)
	err := ct.Start()
	assert.Equal(t, err, nil)

	err = ct.Start()
	assert.Error(t, err)

	wg.Wait()
	ct.Stop()

	time.Sleep(200 * time.Millisecond) // 等待收尾

	cur := counter // 是否还有运行
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, cur, counter)
	t.Log("end")
}

func TestCircleTimerReset(t *testing.T) {
	var (
		wg      = sync.WaitGroup{}
		counter = 0
		last    time.Time
	)

	wg.Add(100)

	fn := func() {
		counter++
		if counter <= 100 {
			wg.Done()
		}
		if last.IsZero() {
			t.Log("trigger -> ", counter)
			last = time.Now()
			return
		}

		since := time.Since(last)
		t.Log("trigger -> ", counter, since.String())
		last = time.Now()
	}

	ct := NewCircleTimer(50*time.Millisecond, fn)
	ct.Start()

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			n := i*10 + 10
			duration := time.Duration(n) * time.Millisecond
			time.Sleep(1 * time.Second)
			time.AfterFunc(duration, func() {
				t.Log("reset", duration.String())
				ct.Reset(duration)
			})
		}()
	}

	wg.Wait()
	ct.Stop()
	time.Sleep(1 * time.Second) // 等待完全stop

	cur := counter
	time.Sleep(1 * time.Second)
	assert.Equal(t, cur, counter)
	t.Log("end")
}

func TestCircleTimerStop(t *testing.T) {
	var (
		counter = 0
	)

	fn := func() {
		counter++
		t.Log("trigger -> ", counter)
	}

	ct := NewCircleTimer(100*time.Millisecond, fn)
	ct.Start()

	time.Sleep(500 * time.Millisecond) // 约为 4个
	ct.Stop()

	time.Sleep(1 * time.Second) // wait
	assert.LessOrEqual(t, counter, 6)
}

func TestCircleTimerStopReset(t *testing.T) {
	var (
		counter = 0
	)

	fn := func() {
		counter++
		t.Log("trigger -> ", counter)
	}

	ct := NewCircleTimer(100*time.Millisecond, fn)
	ct.Start()

	time.Sleep(500 * time.Millisecond) // 约为 4个
	ct.Stop()
	ct.Stop()
	ct.Stop()

	err := ct.Reset(1 * time.Millisecond)
	assert.NotNil(t, err)

	time.Sleep(1 * time.Second) // wait
	assert.LessOrEqual(t, counter, 6)

	ct.Start()
	time.Sleep(500 * time.Millisecond) // 约为 4个
	ct.Stop()
	assert.LessOrEqual(t, counter, 12)
}
