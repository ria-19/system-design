package cache

// represents a single item in the cache
type Node struct {
	key   string
	value interface{} // or string, or any type you choose
	prev  *Node
	next  *Node
}

// represents manager; it holds the two DS (map and list) and cache's settings
type LRUCache struct {
	capacity int
	cache    map[string]*Node // HashMap: key -> pointer to Node; Important!!
	head     *Node            // Dummy head (LRU side)
	tail     *Node            // Dummy tail (MRU side)
}

func NewLRUCache(capacity int) *LRUCache {
	//Initialize with dummy head and tail
	dummy_head := &Node{} // default is {key: "", value: nil, prev: nil, next: nil}
	dummy_tail := &Node{}

	dummy_head.next = dummy_tail // head.next should point to tail initially
	dummy_tail.prev = dummy_head // tail.prev should point to head initially

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
		head:     dummy_head,
		tail:     dummy_tail,
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	if node, ok := c.cache[key]; ok {
		c.moveToTail(node) // Mark as recently used
		return node.value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value interface{}) {
	// check if node exists
	node, ok := c.cache[key]

	// if exists -> update and move to tail
	if ok {
		node.value = value
		c.moveToTail(node)
		return
	}

	// if doesn't exist -> create and insert new node

	// create new node
	node = &Node{key: key, value: value}

	// insert new node

	// if capacity has reached -> evict c.head.next {remove from ll and cache}
	if len(c.cache) == c.capacity {
		lru := c.head.next        // Capture reference BEFORE mutation
		c.removeNode(c.head.next) // Mutates c.head.next
		delete(c.cache, lru.key)
	}

	// add new node at tail and in cache for tracking
	c.addToTail(node)
	c.cache[key] = node
}

// Helper methods

// Remove node from linked list
func (c *LRUCache) removeNode(node *Node) {
	// remove link from prev and next nodes
	node.prev.next = node.next
	node.next.prev = node.prev

	// remove link from node itself
	node.prev = nil
	node.next = nil
}

// Add node before tail dummy node
func (c *LRUCache) addToTail(node *Node) {
	node.prev = c.tail.prev
	node.next = c.tail

	node.prev.next = node
	c.tail.prev = node
}

// Remove then add to tail
func (c *LRUCache) moveToTail(node *Node) {
	c.removeNode(node)
	c.addToTail(node)
}
