// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

// Package cache implements caching of RDAP Service Registry files.
//
// There are two separate implementations: MemoryCache and DiskCache.
package cache

import "time"

type FileState int
const (
	// File is not in the cache.
	Absent FileState = iota

	// File is in the cache. The latest version has already accessed (Load or Saved()).
	Good

	// File is in the cache. A newer version of is available to be Load()'ed.
	//
	// This is used by DiskCache, which uses a shared cache directory.
	ShouldReload

	// File is in the cache, but has expired. It still can be Load()'ed.
	Expired
)

// A RegistryCache implements a cache of Service Registry files.
type RegistryCache interface {
	Load(filename string) ([]byte, error)
	Save(filename string, data []byte) error

	State(filename string) FileState

	SetTimeout(timeout time.Duration)
}
