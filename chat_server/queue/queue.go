package queue

import (
	"errors"
	"sync"
)

type Queue struct {
	front *item
	back  *item
	mutex *sync.Mutex
}

type item struct {
	next  *item
	value interface{}
}

func New() *Queue {
	return &Queue{nil, nil, &sync.Mutex{}}
}

func (q Queue) Empty() bool {
	defer q.mutex.Unlock()
	q.mutex.Lock()
	return q.front == nil
}

func (q *Queue) Push(value interface{}) {
	defer q.mutex.Unlock()
	q.mutex.Lock()

	if q.back == nil {
		q.back = &item{nil, value}
		q.front = q.back
	} else {
		q.back.next = &item{nil, value}
		q.back = q.back.next
	}
}

func (q *Queue) Pop() (value interface{}, err error) {
	defer q.mutex.Unlock()
	q.mutex.Lock()

	if q.front == nil {
		return nil, errors.New("Queue empty")
	}

	err = nil
	value = q.front.value

	q.front = q.front.next
	return
}
