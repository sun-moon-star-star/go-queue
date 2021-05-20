package go_queue

import (
	"container/list"
	"errors"
	"sync"
)

type ConcurrentQueue struct {
	queue  *list.List
	maxLen int

	lock      sync.Mutex
	queueCond *sync.Cond
}

func NewConcurrentQueue() *ConcurrentQueue {
	q := &ConcurrentQueue{}
	q.init()
	return q
}

func (q *ConcurrentQueue) init() {
	q.queue = list.New()
	q.maxLen = -1
	q.queueCond = sync.NewCond(&q.lock)
}

func (q *ConcurrentQueue) SetMaxLen(maxLen int) {
	q.lock.Lock()
	q.maxLen = maxLen
	q.lock.Unlock()
}

func (q *ConcurrentQueue) Append(item interface{}) error {
	var err error

	q.queueCond.L.Lock()
	if q.maxLen != -1 && q.queue.Len() >= q.maxLen {
		err = errors.New("Queue is full")
	} else {
		q.queue.PushBack(item)
		err = nil
	}
	q.queueCond.L.Unlock()

	return err
}

func (q *ConcurrentQueue) Remove() interface{} {
	q.queueCond.L.Lock()
	for q.queue.Len() == 0 {
		q.queueCond.Wait()
	}

	item := q.queue.Front()
	q.queue.Remove(item)
	q.queueCond.L.Unlock()
	return item.Value
}

func (q *ConcurrentQueue) Len() int {
	var len int
	q.lock.Lock()
	len = q.queue.Len()
	q.lock.Unlock()
	return len
}
