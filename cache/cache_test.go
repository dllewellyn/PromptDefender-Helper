package cache

import (
	"context"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	cache := NewInMemoryCache()
	ctx := context.Background()
	key := "testKey"
	value := "testValue"

	err := cache.Set(ctx, key, value)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != value {
		t.Fatalf("expected %v, got %v", value, got)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	cache := NewInMemoryCache()
	ctx := context.Background()
	key := "nonExistentKey"

	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != "" {
		t.Fatalf("expected empty string, got %v", got)
	}
}

func TestExists(t *testing.T) {
	cache := NewInMemoryCache()
	ctx := context.Background()
	key := "testKey"
	value := "testValue"

	exists, err := cache.Exists(ctx, key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exists {
		t.Fatalf("expected key to not exist")
	}

	err = cache.Set(ctx, key, value)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	exists, err = cache.Exists(ctx, key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Fatalf("expected key to exist")
	}
}
