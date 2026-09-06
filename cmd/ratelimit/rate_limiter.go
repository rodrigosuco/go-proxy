// Package ratelimit controls the limit of requests
package ratelimit

type IP string

type Tokens int

type Bucket map[IP]Tokens

func CheckLimit(requestIP string, buckets []Bucket) error {
	return nil
}
