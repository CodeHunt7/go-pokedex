package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// NewCache создает новый экземпляр кэша с заданным интервалом очистки
func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}

	go c.reapLoop(interval)
	
	return c
}

// Add добавляет новый элемент в кэш
func (c *Cache) Add(key string, val []byte) {
	// Блокируем мьютекс на время записи
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

// Get получает элемент из кэша по ключу
func (c *Cache) Get(key string) ([]byte, bool) {
	// Блокируем мьютекс на время чтения
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

// reapLoop периодически очищает устаревшие элементы из кэша
func (c *Cache) reapLoop(interval time.Duration) {
	// Делаем таймер
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Цикл ждет тика
	for t:= range ticker.C {
		c.mu.Lock()

		for key, entry := range c.cache {
			// Если элемент старше интервала, удаляем
			if t.Sub(entry.createdAt) > interval {
				delete(c.cache, key)
			}
		}

		c.mu.Unlock()
	}	
}