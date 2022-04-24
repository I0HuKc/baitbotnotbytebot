package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
)

type ReqEffector[R any] func(ctx context.Context, method string) (R, error)

type request[B, R any] struct {
	target  string
	timeout time.Duration
	body    B
	resp    R
}

type ApiReq[R any] interface {
	Repeater(effector ReqEffector[R], retries int, delay time.Duration) ReqEffector[R]
	NewRequest(ctx context.Context, method string) (R, error)
}

func CreateApiReq[B, R any](target string, body B, to time.Duration) ApiReq[R] {
	return &request[B, R]{
		target:  target,
		body:    body,
		timeout: to,
	}
}

func (req *request[B, R]) Repeater(effector ReqEffector[R], retries int, delay time.Duration) ReqEffector[R] {
	return func(ctx context.Context, method string) (R, error) {
		for i := 0; ; i++ {

			resp, err := effector(ctx, method)
			if err == nil || i >= retries {
				return resp, err
			}

			if os.Getenv("APP_ENV") == core.LocalEnv {
				fmt.Printf("Attempt %d failed; retrying in %v\n", i+1, delay)
			}

			delay += time.Second

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return *new(R), ctx.Err()
			}
		}
	}
}

func (req *request[B, R]) NewRequest(ctx context.Context, method string) (R, error) {
	client := http.Client{
		Timeout: req.timeout,
	}

	body, err := json.Marshal(req.body)
	if err != nil {
		return *new(R), err
	}

	nreq, err := http.NewRequest(method, req.target, bytes.NewBuffer(body))
	if err != nil {
		return *new(R), err
	}

	resp, err := client.Do(nreq)
	if err != nil {
		return *new(R), err
	}

	defer resp.Body.Close()

	var respBody R
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return *new(R), err
	}
	return respBody, nil
}
