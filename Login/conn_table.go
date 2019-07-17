package main

import "sync"

// ConnectionMap ... List of users who plays on GS
type ConnectionMap struct {
	mx sync.RWMutex
	m  map[uint32]*Session
}

// NewConnectionMap ... Creates ConnectionMap object
func NewConnectionMap() *ConnectionMap {
	return &ConnectionMap{
		m: make(map[uint32]*Session),
	}
}

// Set ... adds new connection to list
func (c *ConnectionMap) Set(connID uint32, sess *Session) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[connID] = sess
}

// Get ... retrieves value from ConnectionMap
func (c *ConnectionMap) Get(connID uint32) (*Session, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[connID]
	return val, ok
}

// Remove ... deletes connection from map
func (c *ConnectionMap) Remove(connID uint32) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.m, connID)
}
