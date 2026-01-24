package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_Allow(t *testing.T) {
	tb := NewTokenBucket(5, 10) // 5 tokens/sec, capacity 10

	// Initial burst should allow 10 requests
	for i := 0; i < 10; i++ {
		if !tb.Allow("test_user") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 11th request should be denied (no tokens left)
	if tb.Allow("test_user") {
		t.Error("Request 11 should be denied")
	}

	// Wait for token regeneration (200ms = 1 token at 5/sec)
	time.Sleep(200 * time.Millisecond)

	// Should allow 1 more request
	if !tb.Allow("test_user") {
		t.Error("Request after sleep should be allowed")
	}
}

func TestTokenBucket_AllowN(t *testing.T) {
	tb := NewTokenBucket(10, 20)

	// Allow 5 tokens
	if !tb.AllowN("test_user", 5) {
		t.Error("First AllowN(5) should succeed")
	}

	// Allow 10 more tokens
	if !tb.AllowN("test_user", 10) {
		t.Error("Second AllowN(10) should succeed")
	}

	// This should fail (only 5 tokens remaining)
	if tb.AllowN("test_user", 10) {
		t.Error("Third AllowN(10) should fail")
	}
}

func TestTokenBucket_MultipleKeys(t *testing.T) {
	tb := NewTokenBucket(5, 5)

	// Each user should have independent quota
	users := []string{"user1", "user2", "user3"}
	for _, user := range users {
		for i := 0; i < 5; i++ {
			if !tb.Allow(user) {
				t.Errorf("User %s request %d should be allowed", user, i+1)
			}
		}
		// 6th request should be denied
		if tb.Allow(user) {
			t.Errorf("User %s request 6 should be denied", user)
		}
	}
}

func TestTokenBucket_Reset(t *testing.T) {
	tb := NewTokenBucket(2, 5)

	// Use all tokens
	for i := 0; i < 5; i++ {
		tb.Allow("test_user")
	}

	// Should be rate limited
	if tb.Allow("test_user") {
		t.Error("Request should be denied before reset")
	}

	// Reset quota
	tb.Reset("test_user")

	// Should allow requests again
	if !tb.Allow("test_user") {
		t.Error("Request should be allowed after reset")
	}
}

func TestSlidingWindow_Allow(t *testing.T) {
	sw := NewSlidingWindow(5, 1*time.Second) // 5 requests per second

	// First 5 requests should succeed
	for i := 0; i < 5; i++ {
		if !sw.Allow("test_user") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	if sw.Allow("test_user") {
		t.Error("Request 6 should be denied")
	}

	// Wait for window to slide (1.1 seconds)
	time.Sleep(1100 * time.Millisecond)

	// Should allow new requests
	if !sw.Allow("test_user") {
		t.Error("Request after window slide should be allowed")
	}
}

func TestSlidingWindow_AllowN(t *testing.T) {
	sw := NewSlidingWindow(10, 2*time.Second)

	// Allow 5 requests
	if !sw.AllowN("test_user", 5) {
		t.Error("First AllowN(5) should succeed")
	}

	// Allow 3 more
	if !sw.AllowN("test_user", 3) {
		t.Error("Second AllowN(3) should succeed")
	}

	// This should fail (only 2 remaining)
	if sw.AllowN("test_user", 5) {
		t.Error("Third AllowN(5) should fail")
	}
}

func TestSlidingWindow_GetStats(t *testing.T) {
	sw := NewSlidingWindow(10, 1*time.Second)

	// Make some requests
	for i := 0; i < 7; i++ {
		sw.Allow("test_user")
	}

	stats := sw.GetStats("test_user")
	if stats["requests"] != 7 {
		t.Errorf("Expected 7 requests, got %v", stats["requests"])
	}
	if stats["remaining"] != 3 {
		t.Errorf("Expected 3 remaining, got %v", stats["remaining"])
	}
}

func TestFixedWindow_Allow(t *testing.T) {
	fw := NewFixedWindow(5, 1*time.Second)

	// First 5 requests should succeed
	for i := 0; i < 5; i++ {
		if !fw.Allow("test_user") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	if fw.Allow("test_user") {
		t.Error("Request 6 should be denied")
	}

	// Wait for window to reset
	time.Sleep(1100 * time.Millisecond)

	// Should allow new requests
	if !fw.Allow("test_user") {
		t.Error("Request after window reset should be allowed")
	}
}

func TestFixedWindow_AllowN(t *testing.T) {
	fw := NewFixedWindow(10, 2*time.Second)

	// Allow 6 requests
	if !fw.AllowN("test_user", 6) {
		t.Error("First AllowN(6) should succeed")
	}

	// Allow 4 more (exactly at limit)
	if !fw.AllowN("test_user", 4) {
		t.Error("Second AllowN(4) should succeed")
	}

	// This should fail (limit reached)
	if fw.AllowN("test_user", 1) {
		t.Error("Third AllowN(1) should fail")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tb := NewTokenBucket(100, 200)
	sw := NewSlidingWindow(100, 1*time.Second)
	fw := NewFixedWindow(100, 1*time.Second)

	limiters := []struct {
		name    string
		limiter Limiter
	}{
		{"TokenBucket", tb},
		{"SlidingWindow", sw},
		{"FixedWindow", fw},
	}

	for _, l := range limiters {
		t.Run(l.name, func(t *testing.T) {
			var wg sync.WaitGroup
			allowed := 0
			var mu sync.Mutex

			// Spawn 50 goroutines
			for i := 0; i < 50; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for j := 0; j < 10; j++ {
						if l.limiter.Allow("concurrent_user") {
							mu.Lock()
							allowed++
							mu.Unlock()
						}
					}
				}(i)
			}

			wg.Wait()

			// Should allow some requests (not all 500 due to limit)
			if allowed == 0 {
				t.Error("No requests were allowed")
			}
			if allowed > 200 {
				t.Errorf("Too many requests allowed: %d (expected <= 200)", allowed)
			}
		})
	}
}

func BenchmarkTokenBucket_Allow(b *testing.B) {
	tb := NewTokenBucket(1000, 2000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow("bench_user")
	}
}

func BenchmarkSlidingWindow_Allow(b *testing.B) {
	sw := NewSlidingWindow(1000, 1*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Allow("bench_user")
	}
}

func BenchmarkFixedWindow_Allow(b *testing.B) {
	fw := NewFixedWindow(1000, 1*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fw.Allow("bench_user")
	}
}

func BenchmarkConcurrentTokenBucket(b *testing.B) {
	tb := NewTokenBucket(10000, 20000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			tb.Allow("bench_user")
			i++
		}
	})
}
