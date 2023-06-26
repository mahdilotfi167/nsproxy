package cache

import (
	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	gocache "github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"log"
	"net/url"
	"nsproxy/config"
	"strconv"
	"strings"
	"time"
)

func NewCacheManager(config *config.CacheConfig) *cache.Cache[string] {
	var cacheManager *cache.Cache[string]

	if config.CacheURL == "" {
		client := gocache.New(gocache.NoExpiration, 10*time.Minute)

		store := gocache_store.NewGoCache(client)
		cacheManager = cache.New[string](store)
	} else {
		parse, err := url.Parse(config.CacheURL)
		if err != nil {
			log.Fatalf("Unable to parse cache-url")
		}
		db, err := strconv.Atoi(strings.TrimPrefix(parse.Path, "/"))
		if err != nil {
			log.Fatalf("Bad cache db name")
		}

		store := redis_store.NewRedis(redis.NewClient(&redis.Options{
			Addr: parse.Host,
			DB:   db,
		}))
		cacheManager = cache.New[string](store)
	}

	return cacheManager
}
