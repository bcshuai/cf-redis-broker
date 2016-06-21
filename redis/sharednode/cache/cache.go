package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type CacheUpdater interface {
	Keys() []string
	Value(key string) interface{}
}

func NewCache(autoRefresh bool, refreshInterval int, updateHandler CacheUpdater) *Cache {
	cache := &Cache{
		autoRefresh:     autoRefresh,
		refreshInterval: time.Duration(refreshInterval) * time.Millisecond,
		updateHandler:   updateHandler,
	}

	cache.storage = make(map[string]interface{})

	return cache
}

type Cache struct {
	autoRefresh     bool                   //auto-refresh the cache, or not
	refreshInterval time.Duration          //a time interval to refresh the cache
	updateHandler   CacheUpdater           //a interface used to update the cache
	storage         map[string]interface{} //the backend storage for this cache
	lock            sync.RWMutex

	started  int32     //a value to indicated weather the auto-refreshing is started or not
	stopChan chan bool //a channel used to wait the auto refreshing goroutine to stop
}

func (cache *Cache) Start() {

	cache.updateCache() //update the cache when start

	if cache.autoRefresh {
		//the auto-refreshing feature is enabled
		cache.startAutoRefresh()
	}
}

func (cache *Cache) Get(key string) interface{} {
	cache.lock.RLock()
	val, ok := cache.storage[key]
	cache.lock.RUnlock()

	if ok {
		return val
	} else {
		return nil
	}
}

func (cache *Cache) Put(key string, val interface{}) {
	if val == nil {
		return
	}
	cache.lock.Lock()
	cache.storage[key] = val
	cache.lock.Unlock()
}

func (cache *Cache) Stop() {
	if cache.autoRefresh {
		cache.stopAutoRefresh()
	}
}

func (cache *Cache) IsAutoRefreshingStarted() bool {
	isStarted := atomic.LoadInt32(&cache.started)
	if isStarted != 0 {
		return true
	} else {
		return false
	}
}

func (cache *Cache) startAutoRefresh() {
	if cache.IsAutoRefreshingStarted() {
		return
	}

	cache.stopChan = make(chan bool, 1)
	atomic.StoreInt32(&cache.started, 1)

	go func() {
		for {
			if cache.IsAutoRefreshingStarted() {
				time.AfterFunc(cache.refreshInterval, func() {
					cache.updateCache()
				})
			} else {
				cache.stopChan <- true
			}
		}

	}()
}

func (cache *Cache) stopAutoRefresh() {
	if !cache.IsAutoRefreshingStarted() {
		return
	}
	atomic.StoreInt32(&cache.started, 0) //now the auto refreshing is disabled
	//now wait the auto refreshing goroutine to end
	if cache.stopChan != nil {
		<-cache.stopChan
	}
}

func (cache *Cache) updateCache() {
	if cache.updateHandler == nil {
		return
	}

	keys := cache.updateHandler.Keys()
	valueChan := make(chan interface{}, len(keys))

	for _, cacheKey := range keys {
		go func(ck string) {
			value := cache.updateHandler.Value(ck)
			valueChan <- value
		}(cacheKey)
	}

	for _, key := range keys {
		val := <-valueChan
		if val == nil {
			continue
		}
		cache.lock.Lock()
		cache.storage[key] = val
		cache.lock.Unlock()
	}

	close(valueChan)
}
