// Package ratelimit controls the limit of requests
package ratelimit

import (
	"errors"
	"sync"
	"time"
)

type IP string

type Tokens int

type Bucket struct {
	mu                        sync.Mutex
	IP                        string
	Tokens                    int
	LastTokenRestoreTimeStamp time.Time
}

const (
	BucketRestoreInterval           = 30 * time.Second
	TokenLimit                  int = 0
	AvailableTokensForRenew     int = 15
	AvailableTokensForNewBucket int = 14
)

func CheckLimit(bucket *Bucket) error {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	restoreTokens(bucket)

	if bucket.Tokens <= TokenLimit {
		return errors.New("not enough request tokens")
	}
	removeTokenPerRequest(bucket)
	return nil
}

func removeTokenPerRequest(b *Bucket) {
	b.Tokens--
}

func restoreTokens(b *Bucket) {
	if time.Since(b.LastTokenRestoreTimeStamp) > BucketRestoreInterval {
		b.Tokens = AvailableTokensForRenew
		b.LastTokenRestoreTimeStamp = time.Now()
	}
}
