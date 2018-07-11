package main

import (
	"sync"
)

// cache the results of sqs stats calls locally.
// better would be to use redis or firebase

type StatsCache struct {
	sync.RWMutex
	c map[string]QueueStats
}

var Cache StatsCache

func (sc *StatsCache) Init() {
	sc.c = make(map[string]QueueStats)
}

func CacheGet() map[string]QueueStats {
	Cache.RLock()
	defer Cache.RUnlock()
	return Cache.c
}

func CacheGetSingle(name string) QueueStats {
	Cache.RLock()
	defer Cache.RUnlock()
	return Cache.c[name]
}

func CacheSet(q QueueStats) {
	Cache.Lock()
	Cache.Set(q)
	Cache.Unlock()
}

func (sc *StatsCache) Set(q QueueStats) {
	name := q.Name
	if myQ, ok := sc.c[name]; ok {
		myQ.Copy(q)
		sc.c[name] = myQ
		return
	}
	sc.c[name] = q
}
