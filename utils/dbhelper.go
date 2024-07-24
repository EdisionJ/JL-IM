package utils

import (
	"IM/globle"
	"IM/service/enum"
	"context"
	"encoding/json"
	"errors"
	"github.com/barkimedes/go-deepcopy"
	"github.com/redis/go-redis/v9"
)

func Get(cacheKey string, value any, queryFunc func() (any, error)) error {
	err := GetFromCache(cacheKey, value)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err := GetAndSet(cacheKey, value, queryFunc)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}
func GetAndSet(cacheKey string, value any, queryFunc func() (any, error)) error {
	data, err := queryFunc()
	if err != nil {
		return err
	}
	value, err = deepcopy.Anything(data)
	if err != nil {
		return err
	}
	err = SetToCache(cacheKey, value)
	if err != nil {
		return err
	}
	return nil
}
func SetToCache(key string, value any) error {
	ctx := context.Background()
	byteData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = globle.Rdb.Set(ctx, key, byteData, enum.CacheTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetFromCache(key string, value any) error {
	ctx := context.Background()
	data, err := globle.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), value)
	if err != nil {
		return err
	}
	return nil
}

func RemoveCacheData(key string) error {
	ctx := context.Background()
	err := globle.Rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
