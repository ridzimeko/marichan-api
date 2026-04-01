package httpclient

import (
	"net/http"
	"time"
)

func New() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
	}
}
