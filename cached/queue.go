package cached

import (
	"container/list"
	"sync"
)

type FIFOQueue struct {
	lock sync.RWMutex
	ll   *list.List
}

func NewFIFOQueue() *FIFOQueue {
	fifoQueue := &FIFOQueue{}
	fifoQueue.ll = list.New()
	return fifoQueue
}

func (f *FIFOQueue) Enqueue(val any) *list.Element {
	f.lock.Lock()
	defer f.lock.Unlock()
	e := f.ll.PushFront(val)
	return e
}

func (f *FIFOQueue) Dequeue() any {
	f.lock.Lock()
	defer f.lock.Unlock()
	e := f.ll.Back()
	f.ll.Remove(e)
	return e.Value
}

func (f *FIFOQueue) Size() int {
	return f.ll.Len()
}
