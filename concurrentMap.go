//https://github.com/orcaman/concurrent-map/blob/master/concurrent_map.go
package main

import (
	"sync"
	//"fmt"
)



var SHARD_COUNT = 4096

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type pState []*pStateShared

// A "thread" safe string to anything map.
type pStateShared struct {
	items        map[string]*packet_metadata
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func NewpState() pState {
	m := make(pState, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &pStateShared{items: make(map[string]*packet_metadata)}
	}
	return m
}

// GetShard returns shard under given key
func (m pState) GetShard(key string) *pStateShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m pState) Insert(key string, value * packet_metadata) (res * packet_metadata) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
	return res
}

// Get retrieves an element from map under given key.
func (m pState) Get(key string) (*packet_metadata, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Count returns the number of elements within the map.
func (m pState) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// IsEmpty checks if map is empty.
func (m pState) IsEmpty() bool {
	return m.Count() == 0
}

// Looks up an item under specified key
func (m pState) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove removes an element from the map.
func (m pState) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

/* FOR PACKET_METADATA */
//is Processing for goPackets
func (m pState) isStartProcessing( p * packet_metadata ) ( bool,bool ) {
    // Get shard
    shard := m.GetShard(p.Saddr)
    shard.Lock()
    defer shard.Unlock()
    // Get item from shard.
    p_out, ok := shard.items[p.Saddr]
    if !ok {
        return false,false
    }
	if !p_out.Processing {
		p_out.startProcessing()
		return true,true
	}
    return true, false

}

func (m pState) startProcessing( p * packet_metadata ) bool {

    // Get shard
    shard := m.GetShard(p.Saddr)
    shard.RLock()
    defer shard.RUnlock()
    // See if element is within shard.
    p_out, ok := shard.items[p.Saddr]
    if !ok {
        return false
    }
    p_out.startProcessing()
    return ok

}

func (m pState) finishProcessing( p * packet_metadata ) bool {

    // Get shard
    shard := m.GetShard(p.Saddr)
    shard.Lock()
    defer shard.Unlock()
    // See if element is within shard.
    p_out, ok := shard.items[p.Saddr]
    if !ok {
        return false
    }
    p_out.finishedProcessing()
    return ok

}

/* Meta functions */
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
