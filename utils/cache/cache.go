/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package cache

import (
	"context"
	"sync"
)

type ServiceCache struct {
	mu      sync.RWMutex
	control CacheControl
	key     string
}

func (c *ServiceCache) SetCacheControl(control CacheControl) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.control = control
	return nil
}

func (c *ServiceCache) GetCacheControl() CacheControl {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.control
}

func (c *ServiceCache) SetKey(key string) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.control != Apply {
		return
	}
	c.key = key
	return
}

func (c *ServiceCache) GetKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.control != Apply {
		return ""
	}
	return c.key
}

func (c *ServiceCache) Invalidate(context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.key = ""
	return nil
}

// NewServiceCache creates a service cache.
func NewServiceCache() IServerCache {
	return &ServiceCache{
		mu:      sync.RWMutex{},
		control: Apply,
		key:     "",
	}
}
