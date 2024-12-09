package server

import (
	"sync"
	"unsafe"
)

type flusher struct {
	in []chan delivery
	wg sync.WaitGroup
}

type delivery struct {
	client                  *client
	prodIsMQTT              bool
	sub                     *subscription
	acc                     *Account
	subject, reply, mh, msg []byte
	gwrply                  bool
}

func startFlusher(cnt int) *flusher {
	var f flusher
	f.in = make([]chan delivery, cnt)
	for i := range f.in {
		f.wg.Add(1)
		f.in[i] = make(chan delivery)
		go f.worker(f.in[i])
	}
	return &f
}

func (f *flusher) worker(in <-chan delivery) {
	defer f.wg.Done()

	for d := range in {
		d.client.deliverMsg(d.prodIsMQTT, d.sub, d.acc, d.subject, d.reply, d.mh, d.msg, d.gwrply)
	}
}

func (f *flusher) Q(d delivery) {
	f.in[uintptr(unsafe.Pointer(d.client))%uintptr(cap(f.in))] <- d
}

func (f *flusher) Stop() {
	for i := range f.in {
		close(f.in[i])
	}
	f.wg.Wait()
}
