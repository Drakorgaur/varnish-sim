package model

import lru "github.com/hashicorp/golang-lru/v2"

type Storage[C comparable, V Numeric] interface { // TODO: any other name?
	Size() V
	Store(C, V) bool
	Get(C) (V, bool)
}

// CacheStorage is a storage that uses LRU cache
type CacheStorage[K comparable, V Numeric] struct {
	cache *lru.Cache[K, V]

	size   V
	stored V
}

func (s *CacheStorage[K, V]) Size() V {
	return s.size
}

func (s *CacheStorage[K, V]) Stored() V {
	return s.stored
}

// Store stores a value in the cache
// if the value is bigger than the cache size, it returns false
// if the value can be stored, it returns if object was nuked.
func (s *CacheStorage[K, V]) Store(k K, v V) bool {
	// check if size of cache allows to store a new artifact
	size := s.cache.Len()

	if v > s.size {
		// if the object is bigger than the cache size, we cannot store it
		return false
	}

	nuked := false
	if s.stored+v <= s.size {
		// if we can store, resize cache to store more keys.
		s.cache.Resize(size + 1)
	} else {
		// we have to remove some objects from cache to store the new one
		nuked = true

		// we remove the oldest object from cache
		// and keep the size of cache the same
		for {
			// get the oldest object, to see if it will be removed
			// third argument is a boolean - `is the list is empty`,
			// we do not count that our list will be empty when limit is reached,
			// so it is omitted.
			oldKey, oldValue, _ := s.cache.GetOldest()

			// remove the old object from cache
			// we can also use method `RemoveOldest`, but as we have no Lock on this,
			// there is possibility that the object will be removed by another goroutine;
			// project is not prepared and not tested for multithreading,
			// so we remove the key that is corresponding to the value we got.
			s.cache.Remove(oldKey)

			// decrease the size of stored objects
			s.stored -= oldValue

			// repeat removing objects until we can store the new one
			if s.stored+v <= s.size {
				break
			}
		}
		// on this point we have free space in both storages. LRU and Abstraction wrapper(this/s/self object)
		// we now allowed to store the new one.
	}

	s.stored += v
	s.cache.Add(k, v)

	return nuked
}

func (s *CacheStorage[K, V]) Get(k K) (V, bool) {
	return s.cache.Get(k)
}

func NewCacheStorage[K comparable, V Numeric](size V) (*CacheStorage[K, V], error) {
	// set size of lru to 1 key. LRU cache will be resized based on size of CacheStorage.
	// As we need to watch size of stored objects, not the count.
	cache, err := lru.New[K, V](1)
	if err != nil {
		return nil, err
	}

	return &CacheStorage[K, V]{cache, size, 0}, nil
}
