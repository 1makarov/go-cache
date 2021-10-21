package cache

import (
	"sync"
	"time"
)

type handler struct {
	ch chan string

	expire map[string]time.Time

	m sync.Mutex
}

func initHandler(c chan string) *handler {
	we := make(map[string]time.Time)

	h := &handler{ch: c, expire: we}

	go h.run()

	return h
}

func (h *handler) run() {
	for {
		t := time.Now()

		for k, v := range h.expire {
			if v.Before(t) {
				h.send(k)
				h.delete(k)
			}
		}

	}
}

func (h *handler) add(k string, t time.Time) {
	h.m.Lock()
	h.expire[k] = t
	h.m.Unlock()
}

func (h *handler) delete(k string) {
	h.m.Lock()
	delete(h.expire, k)
	h.m.Unlock()
}

func (h *handler) send(k string) {
	h.ch <- k
}
