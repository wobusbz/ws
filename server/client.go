package server

import (
	"sync"
)

type Client struct {
	id     string
	client map[string]*SocketClient
	me     sync.Mutex
}

func newClinet() *Client {
	return &Client{
		client: make(map[string]*SocketClient, 1000),
	}
}

func (c *Client) addClient(sc *SocketClient) {
	c.me.Lock()
	if sc.ID != "" {
		c.client[sc.ID] = sc
	}
	c.me.Unlock()
}

func (c *Client) delClient(sc *SocketClient) {
	c.me.Lock()
	if sc.ID != "" {
		if _, ok := c.client[sc.ID]; ok {
			c.client[sc.ID] = nil
			delete(c.client, sc.ID)
		}
	}
	c.me.Unlock()
}

func (c *Client) getClient(id string) (sc *SocketClient) {
	c.me.Lock()
	defer c.me.Unlock()
	if v, ok := c.client[id]; ok {
		sc = v
		return
	}
	return nil
}
