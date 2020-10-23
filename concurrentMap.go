/*
Copyright 2020 The Board of Trustees of The Leland Stanford Junior University

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file also includes code from 
https://github.com/orcaman/concurrent-map/blob/master/concurrent_map.go
which is licensed under the MIT license (Copyright (c) 2014 streamrail)
*/


package lzr

import (
	"sync"
	//"fmt"
	//"os"
)

var SHARD_COUNT = 4096

// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type pState []*pStateShared



// A "thread" safe string to anything map.
type pStateShared struct {
	items        map[string]*packet_state
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Creates a new concurrent map.
func NewpState() pState {
	m := make(pState, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &pStateShared{items: make(map[string]*packet_state)}
	}
	return m
}

// GetShard returns shard under given key
func (m pState) GetShard(key string) *pStateShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (m pState) Insert(key string, p * packet_state) {
	shard := m.GetShard(key)
	shard.Lock()

	shard.items[key] = p
	shard.Unlock()
}

// Get retrieves an element from map under given key.
func (m pState) Get(key string) (*packet_state, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// Get item from shard.
	pstate, ok := shard.items[key]
	shard.RUnlock()
	return pstate, ok
}

// Count returns the number of elements within the map.
func (m pState) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		/*if len(shard.items) != 0{
			fmt.Fprintln(os.Stderr,shard.items)
		}*/
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
	delete(shard.items, key)
	shard.Unlock()
}


/* FOR PACKET_METADATA */
//is Processing for goPackets
func (m pState) IsStartProcessing( p * packet_metadata ) ( bool,bool ) {
    // Get shard
	pKey := constructKey(p)
    shard := m.GetShard(pKey)
    shard.Lock()
    // Get item from shard.
    p_out, ok := shard.items[pKey]
    if !ok {
		shard.Unlock()
        return false,false
    }
	if !p_out.Packet.Processing {
		p_out.Packet.startProcessing()
		shard.Unlock()
		return true,true
	}
    shard.Unlock()
    return true, false

}

func (m pState) StartProcessing( p * packet_metadata ) bool {

    // Get shard
	pKey := constructKey(p)
    shard := m.GetShard(pKey)
    shard.RLock()
    // See if element is within shard.
    p_out, ok := shard.items[pKey]
    if !ok {
		shard.RUnlock()
        return false
    }
    p_out.Packet.startProcessing()
    shard.RUnlock()
    return ok

}

func (m pState) FinishProcessing( p * packet_metadata ) bool {

    // Get shard
	pKey := constructKey(p)
    shard := m.GetShard(pKey)
    shard.Lock()
    // See if element is within shard.
    p_out, ok := shard.items[pKey]
    if !ok {
		shard.Unlock()
        return false
    }
    p_out.Packet.finishedProcessing()
	shard.Unlock()
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
