package queue

import (
	"container/list"
	"errors"
	"sync"
)

type Queue struct {
	queue         *list.List
	maxLen        int
	unfinishedCnt int

	lock      sync.Mutex
	queueCond *sync.Cond
}

func New() *Queue {
	q := &Queue{}
	q.Init()
	return q
}

func (q *Queue) Init() {
	q.queue = list.New()
	q.maxLen = -1
	q.queueCond = sync.NewCond(&q.lock)
}

func (q *Queue) SetMaxLen(maxLen int) {
	q.lock.Lock()
	q.maxLen = maxLen
	q.lock.Unlock()
}

func (q *Queue) Append(item interface{}) error {
	var err error

	q.queueCond.L.Lock()
	if q.maxLen != -1 && q.queue.Len() >= q.maxLen {
		err = errors.New("Queue is full")
	} else {
		q.queue.PushBack(item)
		q.unfinishedCnt++

		if q.unfinishedCnt == 1 {
			q.queueCond.Signal()
		}
		err = nil
	}
	q.queueCond.L.Unlock()

	return err
}

func (q *Queue) Remove() interface{} {
	q.queueCond.L.Lock()
	for q.queue.Len() == 0 && q.unfinishedCnt > 0 {
		q.queueCond.Wait()
	}

	if q.unfinishedCnt <= 0 {
		q.queueCond.L.Unlock()
		return nil
	}

	item := q.queue.Front()
	q.queue.Remove(item)
	q.queueCond.L.Unlock()
	return item.Value
}

func (q *Queue) Len() int {
	var len int
	q.lock.Lock()
	len = q.queue.Len()
	q.lock.Unlock()
	return len
}

func (q *Queue) Done() {
	q.queueCond.L.Lock()
	newCnt := q.unfinishedCnt - 1
	q.unfinishedCnt = newCnt
	q.queueCond.L.Unlock()

	if newCnt <= 0 {
		q.queueCond.Broadcast()
	}
}

func (q *Queue) Join() {
	q.queueCond.L.Lock()
	for q.unfinishedCnt > 0 {
		q.queueCond.Wait()
	}
	q.queueCond.L.Unlock()
}
