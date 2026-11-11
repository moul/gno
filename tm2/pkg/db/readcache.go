package db

import (
	"container/list"
	"sync"
)

// ReadCacheDB is a bounded, write-through LRU cache of Get results in front
// of a DB. Every write made through it (Set, Delete, batches) updates the
// cache after the write lands, so a Get through it returns exactly what the
// wrapped DB would return. Writes made to the wrapped DB through any other
// handle are not seen: wrap only a DB this handle owns.
//
// The cache sits below the gas-charging layer (cache.Store charges on its
// own miss, from the length of the value its parent returns), so it changes
// how fast a value is read, never which value or how much gas.
type ReadCacheDB struct {
	DB

	mtx      sync.Mutex
	gen      uint64 // bumped by every write; a read inserts only if unchanged
	entries  map[string]*list.Element
	lru      *list.List // front is most recent
	size     int
	maxBytes int
}

type readCacheEntry struct {
	key   string
	value []byte // nil: the key does not exist
}

// readCacheEntryOverhead approximates the map, list and header bytes per entry.
const readCacheEntryOverhead = 96

var _ DB = (*ReadCacheDB)(nil)

// NewReadCacheDB wraps db with a cache holding at most maxBytes of keys and
// values. maxBytes <= 0 returns db unchanged.
func NewReadCacheDB(db DB, maxBytes int) DB {
	if maxBytes <= 0 {
		return db
	}
	return &ReadCacheDB{
		DB:       db,
		entries:  make(map[string]*list.Element),
		lru:      list.New(),
		maxBytes: maxBytes,
	}
}

// Get implements DB. The returned slice is a copy the caller may keep.
func (c *ReadCacheDB) Get(key []byte) ([]byte, error) {
	c.mtx.Lock()
	if el, ok := c.entries[string(key)]; ok {
		c.lru.MoveToFront(el)
		v := el.Value.(*readCacheEntry).value
		c.mtx.Unlock()
		return cloneBytes(v), nil
	}
	gen := c.gen
	c.mtx.Unlock()

	v, err := c.DB.Get(key)
	if err != nil {
		return nil, err
	}
	c.mtx.Lock()
	if c.gen == gen {
		c.putLocked(string(key), cloneBytes(v))
	}
	c.mtx.Unlock()
	return v, nil
}

// Has implements DB.
func (c *ReadCacheDB) Has(key []byte) (bool, error) {
	c.mtx.Lock()
	if el, ok := c.entries[string(key)]; ok {
		c.lru.MoveToFront(el)
		has := el.Value.(*readCacheEntry).value != nil
		c.mtx.Unlock()
		return has, nil
	}
	c.mtx.Unlock()
	return c.DB.Has(key)
}

// Set implements DB.
func (c *ReadCacheDB) Set(key, value []byte) error {
	return c.write([]readCacheOp{{key: key, value: value}}, func() error { return c.DB.Set(key, value) })
}

// SetSync implements DB.
func (c *ReadCacheDB) SetSync(key, value []byte) error {
	return c.write([]readCacheOp{{key: key, value: value}}, func() error { return c.DB.SetSync(key, value) })
}

// Delete implements DB.
func (c *ReadCacheDB) Delete(key []byte) error {
	return c.write([]readCacheOp{{key: key, del: true}}, func() error { return c.DB.Delete(key) })
}

// DeleteSync implements DB.
func (c *ReadCacheDB) DeleteSync(key []byte) error {
	return c.write([]readCacheOp{{key: key, del: true}}, func() error { return c.DB.DeleteSync(key) })
}

// NewBatch implements DB.
func (c *ReadCacheDB) NewBatch() Batch {
	return &readCacheBatch{Batch: c.DB.NewBatch(), c: c}
}

// NewBatchWithSize implements DB.
func (c *ReadCacheDB) NewBatchWithSize(size int) Batch {
	return &readCacheBatch{Batch: c.DB.NewBatchWithSize(size), c: c}
}

type readCacheOp struct {
	key, value []byte
	del        bool
}

// write drops the keys, runs the write, then caches the written values.
// Bumping gen on both sides keeps a concurrent Get that read the old value
// from inserting it.
func (c *ReadCacheDB) write(ops []readCacheOp, do func() error) error {
	c.mtx.Lock()
	c.gen++
	for _, op := range ops {
		c.removeLocked(string(op.key))
	}
	c.mtx.Unlock()

	if err := do(); err != nil {
		return err
	}

	c.mtx.Lock()
	c.gen++
	for _, op := range ops {
		if op.del {
			c.putLocked(string(op.key), nil)
		} else {
			c.putLocked(string(op.key), cloneBytes(op.value))
		}
	}
	c.mtx.Unlock()
	return nil
}

func (c *ReadCacheDB) putLocked(key string, value []byte) {
	c.removeLocked(key)
	n := len(key) + len(value) + readCacheEntryOverhead
	if n > c.maxBytes {
		return
	}
	c.entries[key] = c.lru.PushFront(&readCacheEntry{key: key, value: value})
	c.size += n
	for c.size > c.maxBytes {
		c.removeLocked(c.lru.Back().Value.(*readCacheEntry).key)
	}
}

func (c *ReadCacheDB) removeLocked(key string) {
	el, ok := c.entries[key]
	if !ok {
		return
	}
	e := c.lru.Remove(el).(*readCacheEntry)
	delete(c.entries, key)
	c.size -= len(e.key) + len(e.value) + readCacheEntryOverhead
}

// readCacheBatch records its ops so the cache can follow once they land.
type readCacheBatch struct {
	Batch
	c   *ReadCacheDB
	ops []readCacheOp
}

func (b *readCacheBatch) Set(key, value []byte) error {
	if err := b.Batch.Set(key, value); err != nil {
		return err
	}
	b.ops = append(b.ops, readCacheOp{key: cloneBytes(key), value: cloneBytes(value)})
	return nil
}

func (b *readCacheBatch) Delete(key []byte) error {
	if err := b.Batch.Delete(key); err != nil {
		return err
	}
	b.ops = append(b.ops, readCacheOp{key: cloneBytes(key), del: true})
	return nil
}

func (b *readCacheBatch) Write() error {
	ops := b.ops
	b.ops = nil
	return b.c.write(ops, b.Batch.Write)
}

func (b *readCacheBatch) WriteSync() error {
	ops := b.ops
	b.ops = nil
	return b.c.write(ops, b.Batch.WriteSync)
}

func (b *readCacheBatch) Close() error {
	b.ops = nil
	return b.Batch.Close()
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	return append(make([]byte, 0, len(b)), b...)
}
