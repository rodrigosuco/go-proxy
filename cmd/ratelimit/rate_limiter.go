// Package ratelimit controls the limit of requests
package ratelimit

import (
	"errors"
	"time"
)

type IP string

type Tokens int

type Bucket struct {
	IP                        string
	Tokens                    int
	LastTokenRestoreTimeStamp time.Time
}

const (
	BucketRestoreInterval     = 30 * time.Second
	TokenLimit            int = 0
)

func CheckLimit(bucket *Bucket) error {
	restoreTokens(bucket)

	if bucket.Tokens <= TokenLimit {
		return errors.New("not enough request tokens")
	}
	removeTokenPerRequest(bucket)
	return nil
}

func removeTokenPerRequest(b *Bucket) {
	b.Tokens = b.Tokens - 1
}

func restoreTokens(b *Bucket) {
	if time.Since(b.LastTokenRestoreTimeStamp) > BucketRestoreInterval {
		b.Tokens = 15
		b.LastTokenRestoreTimeStamp = time.Now()
	}
}
