//go:build e2e
// +build e2e

package test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetComments(t *testing.T) {
	client := resty.New()
	resp, err := client.R().GET(BASE_URL + "/api/comments")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestPostComment(t *testing.T) {
	client := resty.New()
	resp, err := client
		.R()
		.SetBody(`{"slug": "/", "author": "john doe", "body": "hello world"}`)
		.Post(BASE_URL + "/api/comments")
	
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}