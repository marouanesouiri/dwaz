/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD-3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package shardedmap

import (
	"hash/maphash"
	"sync"
)

type shard[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

type Map[K comparable, V any] struct {
	shards    []*shard[K, V]
	shardMask uint64
	seed      maphash.Seed
}

func New[K comparable, V any](shardCount int) Map[K, V] {
	shardCount = max(32, shardCount)
	n := 1
	for n < shardCount {
		n <<= 1
	}
	shards := make([]*shard[K, V], n)
	for i := range shards {
		shards[i] = &shard[K, V]{
			items: make(map[K]V),
		}
	}
	return Map[K, V]{
		shards:    shards,
		shardMask: uint64(n - 1),
		seed:      maphash.MakeSeed(),
	}
}

func (m Map[K, V]) getShard(key K) *shard[K, V] {
	idx := maphash.Comparable(m.seed, key) & m.shardMask
	return m.shards[idx]
}

func (m Map[K, V]) Get(key K) (V, bool) {
	s := m.getShard(key)
	s.RLock()
	val, ok := s.items[key]
	s.RUnlock()
	return val, ok
}

func (m Map[K, V]) Set(key K, value V) {
	s := m.getShard(key)
	s.Lock()
	s.items[key] = value
	s.Unlock()
}

func (m Map[K, V]) Delete(key K) (V, bool) {
	s := m.getShard(key)
	s.Lock()
	val, ok := s.items[key]
	if ok {
		delete(s.items, key)
	}
	s.Unlock()
	return val, ok
}

func (m Map[K, V]) Has(key K) bool {
	s := m.getShard(key)
	s.RLock()
	_, ok := s.items[key]
	s.RUnlock()
	return ok
}

func (m Map[K, V]) Clear() {
	for _, s := range m.shards {
		s.Lock()
		clear(s.items)
		s.Unlock()
	}
}

func (m Map[K, V]) Range(fn func(key K, val V) bool) {
	for _, s := range m.shards {
		s.RLock()
		for k, v := range s.items {
			if !fn(k, v) {
				s.RUnlock()
				return
			}
		}
		s.RUnlock()
	}
}

func (m Map[K, V]) Len() int {
	total := 0
	for _, s := range m.shards {
		s.RLock()
		total += len(s.items)
		s.RUnlock()
	}
	return total
}

func (m Map[K, V]) Compute(key K, fn func(val V, exists bool) V) V {
	s := m.getShard(key)
	s.Lock()
	defer s.Unlock()
	current, exists := s.items[key]
	newVal := fn(current, exists)
	s.items[key] = newVal
	return newVal
}
