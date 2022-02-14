package web

import (
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Get(string) (*http.Response, error)
}

type Hasher struct {
	client HttpClient
	hasher hash.Hash
}

// NewHasher
// client - web client for making GET requests
// hasher - calculates hash of the response by chunks
func NewHasher(client HttpClient, hasher hash.Hash) *Hasher {
	return &Hasher{client: client, hasher: hasher}
}

// Calc gets web address, make GET request, calculate md5 hash of response and write it to the result slice
func (wh *Hasher) Calc(address string) ([]byte, error) {
	parsedUrl, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address %s: %v", address, err)
	}
	if parsedUrl == nil {
		return nil, fmt.Errorf("invalid address %s", address)
	}
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	resp, err := wh.client.Get(parsedUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get response from %s: %v", address, err)
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("body from %s is nil", address)
	}
	defer resp.Body.Close()

	_, err = io.Copy(wh.hasher, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy body from %s to hasher: %v", address, err)
	}

	h := wh.hasher.Sum(nil)
	// later for performance reason we can take it as argument and allocate it somewhere else once
	result := make([]byte, len(h))
	copy(result, h)

	wh.hasher.Reset()

	return result, nil
}
