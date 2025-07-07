package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lruCache *lruCache) Set(key Key, value interface{}) bool {
	got, ok := lruCache.items[key]
	switch ok {
	case true:
		lruCache.queue.MoveToFront(got)
		lruCache.items[key].Value = value
		return ok
	case false:
		if lruCache.queue.Len() < lruCache.capacity {
			pushed := lruCache.queue.PushFront(value)
			lruCache.items[key] = pushed
		} else {
			oldestCachedValue := lruCache.queue.Back().Value
			for key, val := range lruCache.items {
				if val.Value == oldestCachedValue {
					delete(lruCache.items, key)
				}
			}
			lruCache.queue.Remove(lruCache.queue.Back())
			pushed := lruCache.queue.PushFront(value)
			lruCache.items[key] = pushed
		}
		return ok
	}
	return false
}

func (lruCache *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := lruCache.items[key]; ok {
		lruCache.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (lruCache *lruCache) Clear() {
}
