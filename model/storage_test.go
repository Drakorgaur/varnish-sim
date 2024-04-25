//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

package model

import (
	"fmt"
	"testing"
)

func newCacheStorage(size int) (*CacheStorage[string, int], error) {
	return NewCacheStorage[string, int](size)
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestCreation(t *testing.T) {
	store, err := newCacheStorage(100)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if store.Size() != 100 {
		t.Fatalf("error: %v is not equal to 100", store.Size())
	}
}

func TestLRU1(t *testing.T) {
	store, err := newCacheStorage(100)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for i := 0; i < 10; i++ {
		store.Store(fmt.Sprintf("key%d", i), 10)
	}

	if store.Stored() != 100 {
		t.Fatalf("error: not stored 100, but %v", store.Stored())
	}

	store.Store("key10", 10)

	if store.Stored() != 100 {
		t.Fatalf("error: not stored 100, but %v", store.Stored())
	}

	if _, ok := store.Get("key0"); ok {
		t.Fatalf("error: key0 should be removed")
	}

	for i := 1; i < 11; i++ {
		if _, ok := store.Get(fmt.Sprintf("key%d", i)); !ok {
			t.Fatalf("error: key%d should be stored", i)
		}
	}
}

func TestLRU2(t *testing.T) {
	store, err := newCacheStorage(105)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	for i := 0; i < 10; i++ {
		store.Store(fmt.Sprintf("key%d", i), 10)
	}

	if store.Stored() != 100 {
		t.Fatalf("error: not stored 100, but %v", store.Stored())
	}

	store.Store("key10", 3)

	if store.Stored() != 103 {
		t.Fatalf("error: not stored 103, but %v", store.Stored())
	}

	if _, ok := store.Get("key0"); !ok {
		t.Fatalf("error: key0 should not be removed")
	}

	store.Store("key11", 36)

	if store.Stored() != 99 {
		t.Fatalf("error: not stored 99, but %v", store.Stored())
	}

	for i := 1; i < 5; i++ {
		if _, ok := store.Get(fmt.Sprintf("key%d", i)); ok {
			t.Fatalf("error: key%d should be removed", i)
		}
	}
}
