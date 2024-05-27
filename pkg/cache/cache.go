package cache

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNotExists = errors.New("value not exists")

type Options struct {
	Addr string
}

type Cache struct {
	cl *redis.Client
}

func NewCache(opt Options) (*Cache, error) {
	red := redis.NewClient(&redis.Options{
		Addr: opt.Addr,
	})

	c := Cache{cl: red}

	return &c, nil
}

func (c *Cache) SetString(ctx context.Context, key, val string) error {
	cmd := c.cl.Set(ctx, key, val, redis.KeepTTL)

	return cmd.Err()
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	cmd := c.cl.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return "", ErrNotExists
		}

		return "", err
	}

	str, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotExists
		}

		return "", err
	}

	return str, nil
}

func (c *Cache) SetHMap(ctx context.Context, key, field, val string) error {
	cmd := c.cl.HSet(ctx, key, []string{field, val})
	return cmd.Err()
}

func (c *Cache) GetValHMap(ctx context.Context, key, field string) (string, error) {
	cmd := c.cl.HGet(ctx, key, field)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotExists
		}

		return "", err
	}

	return cmd.String(), nil
}

func (c *Cache) AddToList(ctx context.Context, key string, val ...any) error {
	cmd := c.cl.SAdd(ctx, key, val...)
	err := cmd.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetList(ctx context.Context, key string) ([]string, error) {
	cmd := c.cl.SMembers(ctx, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotExists
		}

		return nil, err
	}

	list, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotExists
		}

		return nil, err
	}

	return list, nil
}

func (c *Cache) IsExists(ctx context.Context, key string) (bool, error) {
	cmd := c.cl.Exists(ctx, key)
	err := cmd.Err()
	if err != nil {
		if err == redis.Nil {
			return false, ErrNotExists
		}

		return false, err
	}

	i, err := cmd.Uint64()
	if err != nil {
		if err == redis.Nil {
			return false, ErrNotExists
		}

		return false, err
	}

	return i != 0, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	cmd := c.cl.Del(ctx, key)
	err := cmd.Err()
	if err != nil {
		return err
	}

	_, err = cmd.Uint64()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Ping(ctx context.Context) error {
	cmd := c.cl.Ping(ctx)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Close() error {
	return c.cl.Close()
}
