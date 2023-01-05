package cached

import "container/list"

type LRUCache struct {
	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key   string
	value any
}

func NewLRUCache() *LRUCache {
	c := make(map[string]*list.Element)
	l := list.New()
	return &LRUCache{l, c}
}

// Set 函数添加一个缓存项到Cache对象中
func (c *LRUCache) Set(key string, value any) {
	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.ll = list.New()
	}
	// 如果缓存已经存在于Cache中，那么将该缓存项移到双向链表的最前端
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	
	// 将新添加的缓存项放入双向链表的最前端
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
}

func (c *LRUCache) GetEx(key string) (value any, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Get 方法获取具有指定键的缓存项
func (c *LRUCache) Get(key string) any {
	if c.cache == nil {
		return nil
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value
	}
	return nil
}

// Remove 方法移除具有指定键的缓存
func (c *LRUCache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest 移除双向链表中访问时间最远的那一项
func (c *LRUCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

// Len 方法获取Cache中包含的缓存项个数
func (c *LRUCache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear 清除整个Cache对象
func (c *LRUCache) Clear() {
	c.ll = nil
	c.cache = nil
}

func (c *LRUCache) GetKeys() []string {
	if c.cache == nil {
		return nil
	}
	keys := make([]string, 0)
	for k, _ := range c.cache {
		keys = append(keys, k)
	}
	return keys
}

func (c *LRUCache) GetAll() map[string]any {
	data := make(map[string]any)
	for k, v := range c.cache {
		data[k] = v.Value.(*entry).value
	}
	return data
}
