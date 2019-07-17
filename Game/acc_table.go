package main

import (
	"sync"
)

// AccountsMap ... List of users who plays on GS
type AccountsMap struct {
	mx sync.RWMutex
	m  map[uint64]uint32
}

// NewAccountsMap ... Creates AccountsMap object
func NewAccountsMap() *AccountsMap {
	return &AccountsMap{
		m: make(map[uint64]uint32),
	}
}

// Set ... adds new connection to list
func (c *AccountsMap) Set(accID uint64, connID uint32) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[accID] = connID
}

// Get ... retrieves value from AccountsMap
func (c *AccountsMap) Get(accID uint64) (uint32, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[accID]
	return val, ok
}

// Remove ... deletes acc from map
func (c *AccountsMap) Remove(accID uint64) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.m, accID)
}

// AccountsMap ... List of users who plays on GS
type SessionMap struct {
	mx sync.RWMutex
	m  map[uint32]*Session
}

// NewAccountsMap ... Creates AccountsMap object
func NewSessionMap() *SessionMap {
	return &SessionMap{
		m: make(map[uint32]*Session),
	}
}

// Set ... adds new connection to list
func (c *SessionMap) Set(sessID uint32, sess *Session) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[sessID] = sess
}

// Get ... retrieves value from AccountsMap
func (c *SessionMap) Get(sessID uint32) (*Session, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	val, ok := c.m[sessID]
	return val, ok
}

// Remove ... deletes acc from map
func (c *SessionMap) Remove(sessID uint32) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.m, sessID)
}
