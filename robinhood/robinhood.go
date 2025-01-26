package fastmap

import (
	"fmt"
	"hash/maphash"
)

type RobinHoodMap[K comparable, V any] struct {
	entries    []entry[K, V]
	size       int
	mask       uint64
	loadFactor float64
	hasher     maphash.Hash
}

type entry[K comparable, V any] struct {
	key      K
	value    V
	distance uint8
	occupied bool
}

func NewRobinHoodMap[K comparable, V any]() *RobinHoodMap[K, V] {
	initialSize := 8
	return &RobinHoodMap[K, V]{
		entries:    make([]entry[K, V], initialSize),
		mask:       uint64(initialSize - 1),
		loadFactor: 0.75,
		hasher:     maphash.Hash{},
	}
}

func (m *RobinHoodMap[K, V]) hash(key K) uint64 {
	m.hasher.Write([]byte(fmt.Sprintf("%v", key)))
	defer m.hasher.Reset()
	return m.hasher.Sum64()
}

func (m *RobinHoodMap[K, V]) Put(key K, value V) {
	if float64(m.size+1)/float64(len(m.entries)) > m.loadFactor {
		m.resize()
	}

	hash := m.hash(key)
	index := hash & m.mask
	dist := uint8(0)

	for {
		entry := &m.entries[index]

		if !entry.occupied {
			entry.key = key
			entry.value = value
			entry.distance = dist
			entry.occupied = true
			m.size++
			return
		}

		if entry.key == key {
			entry.value = value
			return
		}

		// Robin Hood: rich (current entry) vs poor (new entry)
		if dist > entry.distance {
			// Swap entries
			key, entry.key = entry.key, key
			value, entry.value = entry.value, value
			dist, entry.distance = entry.distance, dist
		}

		dist++
		index = (index + 1) & m.mask
	}
}

func (m *RobinHoodMap[K, V]) Get(key K) (V, bool) {
	hash := m.hash(key)
	index := hash & m.mask
	dist := uint8(0)

	for {
		entry := &m.entries[index]
		if !entry.occupied || dist > entry.distance {
			var zero V
			return zero, false
		}
		if entry.key == key {
			return entry.value, true
		}
		dist++
		index = (index + 1) & m.mask
	}
}

func (m *RobinHoodMap[K, V]) Remove(key K) bool {
	hash := m.hash(key)
	index := hash & m.mask
	dist := uint8(0)

	for {
		entry := &m.entries[index]
		if !entry.occupied || dist > entry.distance {
			return false
		}
		if entry.key == key {
			// Found the entry to remove
			m.size--

			// Backward shift deletion
			nextIndex := (index + 1) & m.mask
			for {
				nextEntry := &m.entries[nextIndex]
				if !nextEntry.occupied || nextEntry.distance == 0 {
					entry.occupied = false
					return true
				}
				*entry = *nextEntry
				entry.distance--
				index = nextIndex
				nextIndex = (nextIndex + 1) & m.mask
				entry = &m.entries[index]
			}
		}
		dist++
		index = (index + 1) & m.mask
	}
}

func (m *RobinHoodMap[K, V]) resize() {
	oldEntries := m.entries
	newSize := len(m.entries) * 2
	m.entries = make([]entry[K, V], newSize)
	m.mask = uint64(newSize - 1)
	m.size = 0

	for i := range oldEntries {
		if oldEntries[i].occupied {
			m.Put(oldEntries[i].key, oldEntries[i].value)
		}
	}
}

func (m *RobinHoodMap[K, V]) Size() int {
	return m.size
}

func (m *RobinHoodMap[K, V]) Clear() {
	m.entries = make([]entry[K, V], 8)
	m.mask = 7
	m.size = 0
}
