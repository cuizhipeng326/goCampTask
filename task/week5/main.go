// 作业题：参考 Hystrix 实现一个滑动窗口计数器。
package main

import "sync/atomic"

const (
	EventTypeSuccess = iota
	EventTypeFailure
	EventTypeTimeout
	EventTypeRejection
)

type Bucket struct {
	success, failure, timeout, rejection int64
	next                                 *Bucket
}

func (b *Bucket) AddEvent(eventType int) {
	switch eventType {
	case EventTypeSuccess:
		atomic.AddInt64(&b.success, 1)
	case EventTypeFailure:
		atomic.AddInt64(&b.failure, 1)
	case EventTypeTimeout:
		atomic.AddInt64(&b.timeout, 1)
	case EventTypeRejection:
		atomic.AddInt64(&b.rejection, 1)
	}
}

func (b *Bucket) GetEvent(eventType int) int64 {
	switch eventType {
	case EventTypeSuccess:
		return atomic.LoadInt64(&b.success)
	case EventTypeFailure:
		return atomic.LoadInt64(&b.failure)
	case EventTypeTimeout:
		return atomic.LoadInt64(&b.timeout)
	case EventTypeRejection:
		return atomic.LoadInt64(&b.rejection)
	}

	return 0
}

func (b *Bucket) Reset() {
	atomic.StoreInt64(&b.success, 0)
	atomic.StoreInt64(&b.failure, 0)
	atomic.StoreInt64(&b.timeout, 0)
	atomic.StoreInt64(&b.rejection, 0)
}

func NewBucket(tail *Bucket) *Bucket {
	return &Bucket{next: tail}
}

type Window struct {
	buckets []Bucket
	count   int // 窗口内桶个数
}

func (w *Window) ResetWindow() {
	for _, b := range w.buckets {
		b.Reset()
	}
}

func (w *Window) ResetBucket(offset int) {
	w.buckets[offset%w.count].Reset()
}

func (w *Window) SumByType(eventType int) int64 {
	sum := int64(0)
	for _, b := range w.buckets {
		sum += b.GetEvent(eventType)
	}

	return sum
}

func (w *Window) AddEvent(offset, eventType int) {
	w.buckets[offset%w.count].AddEvent(eventType)
}

func NewWindow(count int) *Window {
	w := &Window{make([]Bucket, count), count}
	for i := 0; i < count; i++ {
		if i == count-1 {
			w.buckets[i].next = &w.buckets[0]
		} else {
			w.buckets[i].next = &w.buckets[i+1]
		}
	}

	return w
}
