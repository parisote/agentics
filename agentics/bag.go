package agentics

import (
	"maps"
	"sync"
)

type Bag[T any] struct {
	m  map[string]T
	mu sync.RWMutex
}

func NewBag[T any]() *Bag[T]        { return &Bag[T]{m: make(map[string]T)} }
func (b *Bag[T]) Get(k string) T    { b.mu.RLock(); v := b.m[k]; b.mu.RUnlock(); return v }
func (b *Bag[T]) Set(k string, v T) { b.mu.Lock(); b.m[k] = v; b.mu.Unlock() }
func (b *Bag[T]) All() map[string]T { b.mu.RLock(); defer b.mu.RUnlock(); return maps.Clone(b.m) }
