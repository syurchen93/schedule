package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/maypok86/otter"
)

type CacheInterface interface {
	Get(key string) (string, bool)
	Set(key string, value string) bool
	Delete(key string)
}

var cachePool CacheInterface

func InitCache(ttl time.Duration, size int) {
	var err error
	cachePool, err = otter.MustBuilder[string, string](size).
		CollectStats().
		Cost(func(key string, value string) uint32 {
			return 1
		}).
		WithTTL(ttl).
		Build()

	if err != nil {
		panic(err)
	}
}

func GetCacheItem(key string, v interface{}) error {
	value, found := cachePool.Get(key)
	if !found {
		return fmt.Errorf("key not found in cache")
	}

	err := json.Unmarshal([]byte(value), v)
	if err != nil {
		return fmt.Errorf("error unmarshalling value: %w", err)
	}

	return nil
}

func GetCacheString(key string) (string, bool) {
	return cachePool.Get(key)
}

func SetCacheItem(key string, value interface{}) {
	var strValue string

	if reflect.TypeOf(value).Kind() == reflect.String {
		strValue = value.(string)
	} else {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			fmt.Println("Error marshalling value:", err)
			return
		}
		strValue = string(jsonValue)
	}

	cachePool.Set(key, strValue)
}

func DeleteCacheItem(key string) {
	cachePool.Delete(key)
}
