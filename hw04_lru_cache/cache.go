package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type pair struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	_, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(lc.items[key])
		lc.items[key].Value.(*pair).Value = value
		return true
	} else {
		if lc.queue.Len() == lc.capacity {
			delete(lc.items, lc.queue.Back().Value.(*pair).Key)
			lc.queue.Remove(lc.queue.Back())
		}
		p := pair{Key: key, Value: value}
		lc.items[key] = lc.queue.PushFront(&p)
		return false
	}
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := lc.items[key]
	if ok {
		return val.Value.(*pair).Value, ok
	} else {
		return nil, ok
	}

}

func (lc *lruCache) Clear() {
	lc.items = make(map[Key]*ListItem, lc.capacity)
	for lc.queue.Len() > 0 {
		lc.queue.Remove(lc.queue.Front())
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
