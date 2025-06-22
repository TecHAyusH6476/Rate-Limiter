package ratelimit

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/coocood/freecache"
)

type tokenBucket struct {
	Tokens         int   `json:"tokens"`
	LastRefillUnix int64 `json:"last_refill_unix"`
}

type RateLimiter struct {
	cache     *freecache.Cache
	rules     []RateLimitRule
	ruleIndex map[string]RateLimitRule
	mu        sync.Mutex
}

func NewRateLimiter(configPath string) (*RateLimiter, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	rl := &RateLimiter{
		cache:     freecache.NewCache(10 * 1024 * 1024), // 10MB
		rules:     cfg.RateLimits,
		ruleIndex: make(map[string]RateLimitRule),
	}
	for _, rule := range cfg.RateLimits {
		key := ruleKey(rule.Domain, rule.Descriptors)
		rl.ruleIndex[key] = rule
	}
	return rl, nil
}

func ruleKey(domain string, descriptors map[string]string) string {
	h := fnv.New64a()
	h.Write([]byte(domain))
	for k, v := range descriptors {
		h.Write([]byte(fmt.Sprintf("|%s=%s", k, v)))
	}
	return fmt.Sprintf("%x", h.Sum64())
}

func (rl *RateLimiter) Allow(domain string, descriptors map[string]string) (bool, error) {
	key := ruleKey(domain, descriptors)
	rule, ok := rl.ruleIndex[key]
	if !ok {
		// No rule, allow by default
		return true, nil
	}

	bucketKey := []byte("bucket:" + key)
	now := time.Now().Unix()
	unitSeconds := unitToSeconds(rule.RateLimit.Unit)
	maxTokens := rule.RateLimit.RequestsPerUnit

	rl.mu.Lock()
	defer rl.mu.Unlock()

	var bucket tokenBucket
	entry, err := rl.cache.Get(bucketKey)
	if err == nil {
		_ = json.Unmarshal(entry, &bucket)
	} else {
		bucket = tokenBucket{Tokens: maxTokens, LastRefillUnix: now}
	}

	// Refill tokens
	elapsed := now - bucket.LastRefillUnix
	if elapsed >= unitSeconds {
		bucket.Tokens = maxTokens
		bucket.LastRefillUnix = now
	}

	allowed := false
	if bucket.Tokens > 0 {
		bucket.Tokens--
		allowed = true
	}

	b, _ := json.Marshal(bucket)
	_ = rl.cache.Set(bucketKey, b, int(unitSeconds))

	return allowed, nil
}

func unitToSeconds(unit string) int64 {
	switch unit {
	case "second":
		return 1
	case "minute":
		return 60
	case "hour":
		return 3600
	case "day":
		return 86400
	default:
		return 60
	}
}
