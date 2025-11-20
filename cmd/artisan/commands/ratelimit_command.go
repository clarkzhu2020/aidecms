package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/chenyusolar/aidecms/pkg/ratelimit"
)

// RateLimitCommand demonstrates rate limiting functionality
func RateLimitCommand(args []string) {
	if len(args) > 0 && args[0] == "demo" {
		runRateLimitDemo()
		return
	}

	fmt.Println("Rate Limit Commands:")
	fmt.Println("  artisan ratelimit demo - Run rate limit demonstration")
}

func runRateLimitDemo() {
	fmt.Println("=== Rate Limit Demo ===\n")

	// Demo 1: Token Bucket
	fmt.Println("1. Token Bucket (5 tokens/sec, burst 10)")
	fmt.Println("   - Good for smooth traffic flow")
	fmt.Println("   - Allows bursts within capacity\n")

	tb := ratelimit.NewTokenBucket(5, 10)

	// Try 15 requests rapidly
	allowed, denied := 0, 0
	for i := 1; i <= 15; i++ {
		if tb.Allow("user1") {
			allowed++
			fmt.Printf("   Request %2d: ✓ ALLOWED (tokens used: %d)\n", i, allowed)
		} else {
			denied++
			fmt.Printf("   Request %2d: ✗ DENIED  (no tokens)\n", i)
		}
	}
	fmt.Printf("   Result: %d allowed, %d denied\n\n", allowed, denied)

	// Demo 2: Sliding Window
	fmt.Println("2. Sliding Window (10 req/5sec)")
	fmt.Println("   - Precise time-based limiting")
	fmt.Println("   - Prevents request bursts\n")

	sw := ratelimit.NewSlidingWindow(10, 5*time.Second)

	// Try 15 requests
	allowed, denied = 0, 0
	for i := 1; i <= 15; i++ {
		if sw.Allow("user2") {
			allowed++
			fmt.Printf("   Request %2d: ✓ ALLOWED\n", i)
		} else {
			denied++
			fmt.Printf("   Request %2d: ✗ DENIED\n", i)
		}
	}
	fmt.Printf("   Result: %d allowed, %d denied\n", allowed, denied)

	// Show stats
	stats := sw.GetStats("user2")
	fmt.Printf("   Stats: %d/%d requests used\n\n", stats["requests"], stats["limit"])

	// Demo 3: Fixed Window
	fmt.Println("3. Fixed Window (8 req/10sec)")
	fmt.Println("   - Simple counter per window")
	fmt.Println("   - Resets at window boundary\n")

	fw := ratelimit.NewFixedWindow(8, 10*time.Second)

	// Try 12 requests
	allowed, denied = 0, 0
	for i := 1; i <= 12; i++ {
		if fw.Allow("user3") {
			allowed++
			fmt.Printf("   Request %2d: ✓ ALLOWED\n", i)
		} else {
			denied++
			fmt.Printf("   Request %2d: ✗ DENIED\n", i)
		}
	}
	fmt.Printf("   Result: %d allowed, %d denied\n\n", allowed, denied)

	// Demo 4: Concurrent requests
	fmt.Println("4. Concurrent Load Test")
	fmt.Println("   - 50 goroutines, 100 requests total")
	fmt.Println("   - Limit: 30 requests/second\n")

	tb2 := ratelimit.NewTokenBucket(30, 60)

	var wg sync.WaitGroup
	totalAllowed := 0
	totalDenied := 0
	var mu sync.Mutex

	start := time.Now()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 2; j++ {
				if tb2.Allow(fmt.Sprintf("user%d", id)) {
					mu.Lock()
					totalAllowed++
					mu.Unlock()
				} else {
					mu.Lock()
					totalDenied++
					mu.Unlock()
				}
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("   Duration: %v\n", elapsed)
	fmt.Printf("   Allowed:  %d requests\n", totalAllowed)
	fmt.Printf("   Denied:   %d requests\n", totalDenied)
	fmt.Printf("   Rate:     %.1f req/sec\n\n", float64(totalAllowed)/elapsed.Seconds())

	// Demo 5: Multiple keys with same limiter
	fmt.Println("5. Multiple Users with Shared Limiter")
	fmt.Println("   - Each user has independent quota")
	fmt.Println("   - Limit: 3 requests per user\n")

	tb3 := ratelimit.NewTokenBucket(1, 3)
	users := []string{"alice", "bob", "charlie"}

	for _, user := range users {
		allowed := 0
		denied := 0
		for i := 0; i < 5; i++ {
			if tb3.Allow(user) {
				allowed++
			} else {
				denied++
			}
		}
		fmt.Printf("   User %-7s: %d allowed, %d denied\n", user, allowed, denied)
	}
	fmt.Println()

	// Demo 6: Burst handling comparison
	fmt.Println("6. Burst Handling Comparison")
	fmt.Println("   - All: 5 req limit in 2 seconds")
	fmt.Println("   - Testing 8 immediate requests\n")

	tb4 := ratelimit.NewTokenBucket(3, 5) // rate=3/s, capacity=5
	sw2 := ratelimit.NewSlidingWindow(5, 2*time.Second)
	fw2 := ratelimit.NewFixedWindow(5, 2*time.Second)

	algorithms := []struct {
		name    string
		limiter ratelimit.Limiter
	}{
		{"Token Bucket", tb4},
		{"Sliding Window", sw2},
		{"Fixed Window", fw2},
	}

	for _, alg := range algorithms {
		allowed := 0
		for i := 0; i < 8; i++ {
			if alg.limiter.Allow("test_user") {
				allowed++
			}
		}
		fmt.Printf("   %-15s: %d/8 allowed\n", alg.name, allowed)
	}
	fmt.Println()

	// Summary
	fmt.Println("=== Algorithm Comparison ===")
	fmt.Println("┌─────────────────┬──────────────────┬─────────────────┬───────────────┐")
	fmt.Println("│ Algorithm       │ Pros             │ Cons            │ Use Case      │")
	fmt.Println("├─────────────────┼──────────────────┼─────────────────┼───────────────┤")
	fmt.Println("│ Token Bucket    │ Smooth flow      │ Slightly complex│ API limiting  │")
	fmt.Println("│                 │ Burst support    │                 │ General use   │")
	fmt.Println("├─────────────────┼──────────────────┼─────────────────┼───────────────┤")
	fmt.Println("│ Sliding Window  │ Precise limits   │ Memory overhead │ Strict limits │")
	fmt.Println("│                 │ No burst issues  │                 │ Premium APIs  │")
	fmt.Println("├─────────────────┼──────────────────┼─────────────────┼───────────────┤")
	fmt.Println("│ Fixed Window    │ Simple & fast    │ Burst at edge   │ Basic limits  │")
	fmt.Println("│                 │ Low memory       │                 │ Public APIs   │")
	fmt.Println("└─────────────────┴──────────────────┴─────────────────┴───────────────┘")

	fmt.Println("\n✓ Rate limit demo completed!")
}
