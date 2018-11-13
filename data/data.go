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
	Read(key string) (*Record, error)
	Save(r *Record) error
	List() ([]*Record, error)
	Delete(key string) error
}

// Record represents a data record
type Record struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}

type Option func(o *Options)
