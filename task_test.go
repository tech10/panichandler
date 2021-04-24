package panic_handler

import (
	"sync"
	"testing"
)

type tt struct {
	sync.Mutex
	wg    sync.WaitGroup
	value string
}

func (t *tt) DoPanicTask(i *Info) {
	t.Lock()
	defer t.Unlock()
	t.value = i.String()
	t.wg.Done()
}

func Test_HandleTask(t *testing.T) {
	ts := &tt{}
	ts.wg.Add(1)
	go func() {
		defer HandleTask(ts)
		panic("Test task interface.")
	}()
	ts.wg.Wait()
	ts.Lock()
	defer ts.Unlock()
	if ts.value == "" {
		t.Fatal("Value is blank and shouldn't have been. A panic should have been caught.")
	}
	t.Log("Panic caught. Data received:\n" + ts.value)
}
