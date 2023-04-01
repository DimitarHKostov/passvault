package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	cacheManager *CacheManager
)

type CacheManager struct {
	cache *cache.Cache
}

func Get() *CacheManager {
	if cacheManager == nil {
		cache := cache.New(5*time.Minute, 10*time.Minute)
		cacheManager = &CacheManager{cache: cache}
	}

	return cacheManager
}

func (cm *CacheManager) Set(key, value string) {
	cm.cache.Set(key, value, cache.DefaultExpiration)
}

func (cm *CacheManager) Get(key string) (string, bool) {
	cacheValue, found := cm.cache.Get(key)

	return cacheValue.(string), found
}
