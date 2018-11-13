// Package data is an interface for key-value storage.
package data

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

// Data is a distributed key-value store interface
type Data interface {
	Get(key string) (*Item, error)
	Del(key string) error
	Put(item *Item) error
	List() ([]*Item, error)
}

type Item struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}

type Option func(o *Options)
