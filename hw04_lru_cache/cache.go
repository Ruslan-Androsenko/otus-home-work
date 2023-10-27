package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    *sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type Node struct {
	Key   Key
	Value interface{}
}

// NewCache Создать новый кэш.
func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    &sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set Добавить элемент в кэш по указанному ключу.
func (elem *lruCache) Set(key Key, value interface{}) bool {
	elem.mutex.Lock()
	defer elem.mutex.Unlock()

	if item, ok := elem.items[key]; !ok {
		elem.items[key] = elem.queue.PushFront(&Node{
			Key:   key,
			Value: value,
		})

		if elem.queue.Len() > elem.capacity {
			lastItem := elem.queue.Back()

			if nodeItem, nodeOk := lastItem.Value.(*Node); nodeOk {
				delete(elem.items, nodeItem.Key)
				elem.queue.Remove(lastItem)
			}
		}
	} else if nodeItem, nodeOk := item.Value.(*Node); nodeOk {
		nodeItem.Value = value
		elem.queue.MoveToFront(item)

		return true
	}

	return false
}

// Get Получить элемент из кэша по указанному ключу.
func (elem *lruCache) Get(key Key) (interface{}, bool) {
	elem.mutex.Lock()
	defer elem.mutex.Unlock()

	if item, ok := elem.items[key]; ok {
		elem.queue.MoveToFront(item)

		if nodeItem, nodeOk := item.Value.(*Node); nodeOk {
			return nodeItem.Value, ok
		}
	}

	return nil, false
}

// Clear Очистить кэш.
func (elem *lruCache) Clear() {
	elem.mutex.Lock()
	defer elem.mutex.Unlock()

	elem.queue = NewList()
	elem.items = make(map[Key]*ListItem, elem.capacity)
}
