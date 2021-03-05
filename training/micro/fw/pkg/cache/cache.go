package cache

import (
	"github.com/micro/micro/v3/service/logger"
	microStore "github.com/micro/micro/v3/service/store"
	"time"
)

// cacheDefaultExpiration defaults time for a value in the cache to expire
const cacheDefaultExpiration = 2 * time.Hour

// cacheMaxListLimit defines maximum number of items to pull from the cache when doing a list of keys
const cacheMaxListLimit = 200

type CacheOptions func(*Cache) error

func WithDatabase(database string) CacheOptions {
	return func(c *Cache) error {
		c.database = database
		return nil
	}
}

func WithPrefix(prefix string) CacheOptions {
	return func(c *Cache) error {
		c.prefix = prefix
		return nil
	}
}

func WithExpiryDuration(expiryDuration time.Duration) CacheOptions {
	return func(c *Cache) error {
		if expiryDuration == 0 {
			expiryDuration = cacheDefaultExpiration
		}
		c.expiryDuration = expiryDuration
		return nil
	}
}

type Cache struct {
	store          microStore.Store
	database       string
	prefix         string
	expiryDuration time.Duration
}

func NewCache(store microStore.Store, opts ...CacheOptions) (*Cache, error) {
	if store == nil {
		store = microStore.DefaultStore
	}
	c := &Cache{store: store}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Cache) Close() error {
	return c.store.Close()
}

func (c *Cache) Set(key, value string) error {
	rec := microStore.Record{
		Key:    key,
		Value:  []byte(value),
		Expiry: c.expiryDuration,
	}
	prefixOptions := microStore.WriteTo(c.database, c.prefix)
	err := c.store.Write(&rec, prefixOptions)
	if err != nil {
		logger.Errorf("Set key %v failed %v", key, err)
		return err
	}
	logger.Infof("Set cache success %v:%v", key, value)
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	prefixOptions := microStore.ReadFrom(c.database, c.prefix)
	recs, err := c.store.Read(key, prefixOptions)
	if err != nil {
		logger.Errorf("Get key %v failed: %v", key, err)
		return "", err
	}
	if len(recs) > 0 {
		return string(recs[0].Value), nil
	}
	return "", nil
}

func (c *Cache) Delete(key string) error {
	prefixOptions := microStore.DeleteFrom(c.database, c.prefix)
	err := c.store.Delete(key, prefixOptions)
	if err != nil {
		logger.Errorf("Get key %v failed: %v", key, err)
		return err
	}
	return nil
}

func (c *Cache) List(prefix string, limit uint, offset uint) ([]string, error) {
	if limit > cacheMaxListLimit {
		logger.Errorf("CacheTooManyValuesToList: %v", cacheMaxListLimit)
		limit = cacheMaxListLimit
	}
	prefixOptions := microStore.ListFrom(c.database, c.prefix)
	pfOptions := microStore.ListPrefix(prefix)
	limitOptions := microStore.ListLimit(limit)
	offsetOptions := microStore.ListOffset(offset)
	recs, err := c.store.List(prefixOptions, pfOptions, limitOptions, offsetOptions)
	if err != nil {
		logger.Errorf("List prefix %v failed: %v", prefix, err)
		return nil, err
	}
	return recs, nil
}
