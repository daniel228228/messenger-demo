package service

import (
	"context"
	"sync"

	"messenger.api/go/api"
)

var cachedUsers *cachedUserList

func initCachedUsers(s *service) {
	cachedUsers = &cachedUserList{
		s:    s,
		list: make(map[string]*api.User),
	}
}

type cachedUserList struct {
	s    *service
	list map[string]*api.User
	mtx  sync.RWMutex
}

func (c *cachedUserList) get(id string) *api.User {
	c.mtx.RLock()
	v, ok := c.list[id]
	c.mtx.RUnlock()

	if !ok {
		resp, err := c.s.users.GetUser(context.Background(), &api.GetUserRequest{
			UserId: id,
		})
		if err != nil {
			return nil
		}

		v = resp.User

		c.mtx.Lock()
		c.list[id] = v
		c.mtx.Unlock()
	}

	return v
}

func (c *cachedUserList) clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for k := range c.list {
		delete(c.list, k)
	}
}
