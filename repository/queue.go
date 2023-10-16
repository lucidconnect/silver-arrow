package repository

import (
	"errors"

	"github.com/OrlovEvgeny/go-mcache"
	"github.com/edwingeng/deque/v2"
	"github.com/lucidconnect/silver-arrow/repository/models"
)

type Queue struct {
	queue *deque.Deque[models.Subscription]
}

func NewDeque() *Queue {
	return &Queue{
		queue: deque.NewDeque[models.Subscription](),
	}
}

func (q *Queue) Read() (models.Subscription, error) {
	sub, ok := q.queue.TryPopFront()
	if !ok {
		return sub, errors.New("queue is empty")
	}

	return sub, nil
}

func (q *Queue) Write(sub models.Subscription) {
	q.queue.PushBack(sub)
}

type Cache struct {
	cache *mcache.CacheDriver
}

func NewMCache() *Cache {
	return &Cache{
		cache: mcache.New(),
	}
}

func (c *Cache) Set(key string, value interface{}) error {
	return c.cache.Set(key, value, mcache.TTL_FOREVER)
}

func (c *Cache) Get(key string) (interface{}, error) {
	v, ok := c.cache.Get(key)
	if !ok {
		return nil, errors.New("key does not exist")
	}
	return v, nil
}
