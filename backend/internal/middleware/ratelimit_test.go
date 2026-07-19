package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_CleanupEvictsIdle(t *testing.T) {
	rl := &RateLimiter{visitors: map[string]*visitor{}, rate: 1, burst: 1}

	rl.getLimiter("1.2.3.4")
	rl.visitors["1.2.3.4"].lastSeen = time.Now().Add(-2 * visitorTTL)
	rl.getLimiter("5.6.7.8")

	rl.cleanup(time.Now())

	if _, ok := rl.visitors["1.2.3.4"]; ok {
		t.Fatal("idle visitor should be evicted")
	}
	if _, ok := rl.visitors["5.6.7.8"]; !ok {
		t.Fatal("recent visitor should remain")
	}
}
