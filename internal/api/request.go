package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ReqEffector func(ctx context.Context) (map[string]interface{}, error)

type request struct {
	target string
	body   map[string]interface{}
	resp   map[string]interface{}
}

type ApiReq interface {
	MakeGetReq(ctx context.Context) (map[string]interface{}, error)
	Repeater(effector ReqEffector, retries int, delay time.Duration) ReqEffector
}

func CreateApiReq(t string, b map[string]interface{}) ApiReq {
	return &request{
		target: t,
		body:   b,
	}
}

func (req *request) Repeater(effector ReqEffector, retries int, delay time.Duration) ReqEffector {
	return func(ctx context.Context) (map[string]interface{}, error) {
		for i := 0; ; i++ {

			resp, err := effector(ctx)
			if err == nil || i >= retries {
				return resp, err
			}

			fmt.Printf("Attempt %d failed; retrying in %v\n", i+1, delay)

			delay += time.Second

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}

func (req *request) MakeGetReq(ctx context.Context) (map[string]interface{}, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(req.target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return nil, err
	}

	return j, nil
}
