package handler

import (
	"context"
	"examples/cache/api"
	"fmt"
	"time"

	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/logger"
)

// Cache ...
type Cache struct {
	cache cache.Cache
}

// NewCache ...
func NewCache(opts ...cache.Option) *Cache {
	c := cache.NewCache(opts...)
	return &Cache{c}
}

// Get ...
func (c *Cache) Get(ctx context.Context, in *api.GetRequest, out *api.GetResponse) error {
	logger.Info("Received Cache.Get request: %v", in)

	v, e, err := c.cache.Context(ctx).Get(in.Key)
	if err != nil {
		return err
	}

	out.Value = fmt.Sprintf("%v", v)
	out.Expiration = e.String()

	return nil
}

// Put ...
func (c *Cache) Put(ctx context.Context, req *api.PutRequest, rsp *api.PutResponse) error {
	logger.Infof("Received Cache.Put request: %v", req)

	d, err := time.ParseDuration(req.Duration)
	if err != nil {
		return err
	}

	if err := c.cache.Context(ctx).Put(req.Key, req.Value, d); err != nil {
		return err
	}

	return nil
}

// Delete ...
func (c *Cache) Delete(ctx context.Context, req *api.DeleteRequest, rsp *api.DeleteResponse) error {
	logger.Infof("Received Cache.Delete request: %v", req)

	if err := c.cache.Context(ctx).Delete(req.Key); err != nil {
		return err
	}

	return nil
}
