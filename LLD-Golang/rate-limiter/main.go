package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

/***************
 * Strategy API
 ***************/
type RateLimiter interface {
	AllowRequest(userID string) bool
}

/**************************
 * Fixed Window (per-user)
 **************************/
type FixedWindowRateLimiter struct {
	maxRequests int
	windowSize  time.Duration

	mu          sync.Mutex
	windowStart time.Time
	count       int
}

func NewFixedWindowRateLimiter(maxRequests int, windowSize time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		maxRequests: maxRequests,
		windowSize:  windowSize,
		windowStart: time.Now(),
	}
}

func (fw *FixedWindowRateLimiter) AllowRequest(_ string) bool {
	now := time.Now()
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = now
		fw.count = 0
	}
	if fw.count < fw.maxRequests {
		fw.count++
		return true
	}
	return false
}

/**************************
 * Sliding Window (per-user)
 **************************/
type SlidingWindowRateLimiter struct {
	maxRequests int
	windowSize  time.Duration

	mu         sync.Mutex
	timestamps []time.Time
}

func NewSlidingWindowRateLimiter(maxRequests int, windowSize time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		maxRequests: maxRequests,
		windowSize:  windowSize,
		timestamps:  make([]time.Time, 0, maxRequests+4),
	}
}

func (s *SlidingWindowRateLimiter) AllowRequest(_ string) bool {
	now := time.Now()
	cutoff := now.Add(-s.windowSize)

	s.mu.Lock()
	defer s.mu.Unlock()

	// evict < cutoff
	i := 0
	for i < len(s.timestamps) && s.timestamps[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		s.timestamps = append([]time.Time{}, s.timestamps[i:]...)
	}

	if len(s.timestamps) < s.maxRequests {
		s.timestamps = append(s.timestamps, now)
		return true
	}
	return false
}

/*************************
 * Token Bucket (per-user)
 *************************/
type TokenBucketRateLimiter struct {
	capacity     int
	refillPerSec float64 // tokens per second

	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
}

func NewTokenBucketRateLimiter(capacity int, refillPerSec float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity:     capacity,
		refillPerSec: refillPerSec,
		tokens:       float64(capacity),
		lastRefill:   time.Now(),
	}
}

func (t *TokenBucketRateLimiter) refill(now time.Time) {
	elapsed := now.Sub(t.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	t.tokens = math.Min(float64(t.capacity), t.tokens+elapsed*t.refillPerSec)
	t.lastRefill = now
}

func (t *TokenBucketRateLimiter) AllowRequest(_ string) bool {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	t.refill(now)
	if t.tokens >= 1.0 {
		t.tokens -= 1.0
		return true
	}
	return false
}

/************************
 * Leaky Bucket (per-user)
 ************************/
type LeakyBucketRateLimiter struct {
	capacity  int
	leakEvery time.Duration

	mu     sync.Mutex
	q      int
	ticker *time.Ticker
	stop   chan struct{}
	once   sync.Once
}

func NewLeakyBucketRateLimiter(capacity int, leakEvery time.Duration) *LeakyBucketRateLimiter {
	l := &LeakyBucketRateLimiter{
		capacity:  capacity,
		leakEvery: leakEvery,
		stop:      make(chan struct{}),
	}
	if l.leakEvery <= 0 {
		l.leakEvery = time.Second
	}
	l.ticker = time.NewTicker(l.leakEvery)
	go l.leakLoop()
	return l
}

func (l *LeakyBucketRateLimiter) leakLoop() {
	for {
		select {
		case <-l.ticker.C:
			l.mu.Lock()
			if l.q > 0 {
				l.q-- // process one request
			}
			l.mu.Unlock()
		case <-l.stop:
			return
		}
	}
}

func (l *LeakyBucketRateLimiter) AllowRequest(_ string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.q < l.capacity {
		l.q++
		return true
	}
	return false
}

func (l *LeakyBucketRateLimiter) Close() {
	l.once.Do(func() {
		close(l.stop)
		l.ticker.Stop()
	})
}

/***********************
 * Factory (Duration in)
 ***********************/
type RateLimiterFactory struct{}

func (f *RateLimiterFactory) CreateRateLimiter(kind string, maxRequests int, window time.Duration) (RateLimiter, error) {
	switch strings.ToUpper(kind) {
	case "FIXED_WINDOW":
		return NewFixedWindowRateLimiter(maxRequests, window), nil
	case "SLIDING_WINDOW":
		return NewSlidingWindowRateLimiter(maxRequests, window), nil
	case "TOKEN_BUCKET":
		// tokens/sec = maxRequests / window(seconds)
		if window <= 0 {
			return nil, errors.New("window must be > 0 for token bucket")
		}
		refillPerSec := float64(maxRequests) / window.Seconds()
		return NewTokenBucketRateLimiter(maxRequests, refillPerSec), nil
	case "LEAKY_BUCKET":
		// leak one request every (window / maxRequests)
		if maxRequests <= 0 || window <= 0 {
			return nil, errors.New("maxRequests and window must be > 0 for leaky bucket")
		}
		leakEvery := time.Duration(float64(window) / float64(maxRequests))
		if leakEvery <= 0 {
			leakEvery = time.Second
		}
		return NewLeakyBucketRateLimiter(maxRequests, leakEvery), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", kind)
	}
}

/****************************
 * Service (Duration in API)
 ****************************/
type RateLimiterService struct {
	mu       sync.RWMutex
	limiters map[string]RateLimiter
}

func NewRateLimiterService() *RateLimiterService {
	return &RateLimiterService{
		limiters: make(map[string]RateLimiter),
	}
}

func (s *RateLimiterService) RegisterUser(userID string, algorithm string, maxRequests int, window time.Duration) error {
	factory := &RateLimiterFactory{}
	limiter, err := factory.CreateRateLimiter(algorithm, maxRequests, window)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	// If replacing a leaky bucket, stop its ticker.
	if old, ok := s.limiters[userID]; ok {
		if lb, ok := old.(*LeakyBucketRateLimiter); ok {
			lb.Close()
		}
	}
	s.limiters[userID] = limiter
	return nil
}

func (s *RateLimiterService) AllowRequest(userID string) (bool, error) {
	s.mu.RLock()
	limiter, ok := s.limiters[userID]
	s.mu.RUnlock()
	if !ok {
		return false, fmt.Errorf("user %s not registered", userID)
	}
	return limiter.AllowRequest(userID), nil
}

func (s *RateLimiterService) CloseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, lim := range s.limiters {
		if lb, ok := lim.(*LeakyBucketRateLimiter); ok {
			lb.Close()
		}
	}
}

/*********
 * Demo
 *********/
func main() {
	svc := NewRateLimiterService()
	defer svc.CloseAll()

	_ = svc.RegisterUser("user_1", "FIXED_WINDOW", 5, 10*time.Second)
	_ = svc.RegisterUser("user_2", "SLIDING_WINDOW", 3, 5*time.Second)
	_ = svc.RegisterUser("user_3", "TOKEN_BUCKET", 5, 10*time.Second) // ~0.5 token/sec
	_ = svc.RegisterUser("user_4", "LEAKY_BUCKET", 3, 4*time.Second)  // leak ~every 1.33s

	for i := 0; i < 7; i++ {
		a1, _ := svc.AllowRequest("user_1")
		a2, _ := svc.AllowRequest("user_2")
		a3, _ := svc.AllowRequest("user_3")
		a4, _ := svc.AllowRequest("user_4")
		fmt.Printf("Tick %d | u1:%v u2:%v u3:%v u4:%v\n", i+1, a1, a2, a3, a4)
		time.Sleep(1 * time.Second)
	}
}
