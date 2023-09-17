package stringset

import "sync"

type Stringset struct {
	sync.RWMutex
	kv map[string]struct{}
}

func New() *Stringset {
	return &Stringset{
		kv: make(map[string]struct{}),
	}
}

func FromElements(el ...string) *Stringset {
	ss := New()
	for _, e := range el {
		ss.Lock()
		defer ss.Unlock()
		ss.kv[e] = struct{}{}
	}

	return ss
}

func (ss *Stringset) Add(s string) bool {
	ss.RLock()
	defer ss.RUnlock()
	if _, ok := ss.kv[s]; ok {
		return false
	}

	ss.Lock()
	ss.kv[s] = struct{}{}
	ss.Unlock()

	return true
}

func (ss *Stringset) AsSlice() []string {
	ss.RLock()
	defer ss.RUnlock()
	keys := make([]string, 0, len(ss.kv))
	for k := range ss.kv {
		keys = append(keys, k)
	}
	return keys
}

func (ss *Stringset) Contains(s string) bool {
	ss.RLock()
	defer ss.RUnlock()
	_, ok := ss.kv[s]
	return ok
}
