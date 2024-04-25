//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

package model

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

// Director is an interface for a director that based on an internal logic
// picks a backend to handle a request.
type Director interface {
	// AddBackend adds a backend to the director, which will be participating
	// in the internal logic of the director.
	AddBackend(WebInterface)

	// GetBackend returns a backend based on the internal logic of the director.
	GetBackend(string) WebInterface

	// Backends returns the list of backends that the director is managing.
	Backends() []WebInterface
}

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

// Sum64 returns the hash of data.
func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}

// ShardDirector is a director that uses consistent hashing to distribute
// requests to backends.
type ShardDirector struct {
	// registered backends
	backends []WebInterface

	// consistent hashing instance
	hashing *consistent.Consistent
}

// Backends returns the list of backends that the director is managing.
func (d *ShardDirector) Backends() []WebInterface {
	return d.backends
}

// NewShardDirector is a constructor for ShardDirector
func NewShardDirector() *ShardDirector {
	cfg := consistent.Config{
		Hasher: hasher{},
	}

	return &ShardDirector{
		backends: make([]WebInterface, 0),
		hashing:  consistent.New(nil, cfg),
	}
}

// AddBackend adds a backend to the director, which will be participating
// in the internal logic of the director.
// Note: hashing is a consistent hashing instance. On appending a new backend
// hashing is update a circle of backends.
func (d *ShardDirector) AddBackend(w WebInterface) {
	d.hashing.Add(w)
	d.backends = append(d.backends, w)
}

// GetBackend returns a backend based on the internal logic of the director.
func (d *ShardDirector) GetBackend(req string) WebInterface {
	// `LocateKey` returns a Member Interface that holds up a Hostname
	// of backend.
	// Note: member de facto is a WebInterface instance
	member := d.hashing.LocateKey([]byte(req))

	return member.(WebInterface)
}

// RoundRobinDirector is a director that uses round-robin to distribute
type RoundRobinDirector struct {
	backends []WebInterface

	index int
}

// Backends returns the list of backends that the director is managing.
func (d *RoundRobinDirector) Backends() []WebInterface {
	return d.backends
}

// NewRoundRobinDirector is a constructor for RoundRobinDirector
func NewRoundRobinDirector() *RoundRobinDirector {
	return &RoundRobinDirector{
		backends: make([]WebInterface, 0),
	}
}

// AddBackend adds a backend to the director, which will be participating
func (d *RoundRobinDirector) AddBackend(w WebInterface) {
	d.backends = append(d.backends, w)
}

// GetBackend returns a backend based on the internal logic of the director.
func (d *RoundRobinDirector) GetBackend(_ string) WebInterface {
	backend := d.backends[d.index]
	d.index = (d.index + 1) % len(d.backends)

	return backend
}
