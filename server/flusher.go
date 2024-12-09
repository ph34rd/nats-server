package server

import (
	"sync"
)

type flusher struct {
	in chan *client
	wg sync.WaitGroup
}

func startFlusher(cnt int) *flusher {
	var f flusher
	f.in = make(chan *client, cnt)
	for i := 0; i < cnt; i++ {
		f.wg.Add(1)
		go f.worker()
	}
	return &f
}

func (f *flusher) worker() {
	defer f.wg.Done()

	for cp := range f.in {
		cp.out.sg.Signal()
	}
}

func (f *flusher) Q(cp *client) {
	f.in <- cp
}

func (f *flusher) Stop() {
	close(f.in)
	f.wg.Wait()
}
