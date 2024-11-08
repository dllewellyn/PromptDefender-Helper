package cache

import (
	"context"
	"sync"
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

func TestHashKey(t *testing.T) {
	key := "testKey"
	expectedHash := "15291f67d99ea7bc578c3544dadfbb991e66fa69cb36ff70fe30e798e111ff5f"
	hashedKey := hashKey(key)

	if hashedKey != expectedHash {
		t.Fatalf("expected %v, got %v", expectedHash, hashedKey)
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewInMemoryCache()
	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := cache.Set(ctx, key, value)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cache.Get(ctx, key)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		}()
	}

	wg.Wait()
}
