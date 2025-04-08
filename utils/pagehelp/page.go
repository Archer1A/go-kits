package pagehelp

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type IPage interface {
	GetPage() int
	GetPageSize() int
}

type IResponse[T any] interface {
	GetTotal() int
	GetRows() []T
}
type SyncAllResource[T any] func(ctx context.Context, page, pageSize int) (IResponse[T], error)

type Client[T any] struct {
	MaxRetries    int
	RetryInterval int
	MaxPageSize   int
	StartPage     int
	WorkCount     int
}

func NewClient[T any]() *Client[T] {
	return &Client[T]{
		MaxRetries:    3,
		RetryInterval: 5,
		MaxPageSize:   100,
		StartPage:     1,
		WorkCount:     5,
	}
}

func (c *Client[T]) WithMaxRetries(maxRetries int) *Client[T] {
	c.MaxRetries = maxRetries
	return c
}

func (c *Client[T]) WithRetryInterval(retryInterval int) *Client[T] {
	c.RetryInterval = retryInterval
	return c
}

func (c *Client[T]) WithMaxPageSize(maxPageSize int) *Client[T] {
	c.MaxPageSize = maxPageSize
	return c
}

func (c *Client[T]) WithStartPage(startPage int) *Client[T] {
	c.StartPage = startPage
	return c
}

func (c *Client[T]) WithWorkCount(workCount int) *Client[T] {
	c.WorkCount = workCount
	return c
}

func (c *Client[T]) SyncAll(ctx context.Context, f SyncAllResource[T]) ([]T, error) {
	firstPageBody, err := f(ctx, c.StartPage, c.MaxPageSize)
	if err != nil {
		return nil, err
	}
	var allResources []T
	allResources = append(allResources, firstPageBody.GetRows()...)
	if firstPageBody.GetTotal() <= c.MaxPageSize {
		return allResources, nil
	}
	totalPage := (firstPageBody.GetTotal() + c.MaxPageSize - 1) / c.MaxPageSize
	errChan := make(chan error, totalPage-1)
	pageChan := make(chan int, totalPage-1)
	resultChan := make(chan []T, totalPage-1)
	wg := sync.WaitGroup{}
	for i := 0; i < c.WorkCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range pageChan {
				var lastErr error
				var result []T
				for t := 0; t < c.MaxRetries; t++ {
					resp, err := f(ctx, page, c.MaxPageSize)
					if err == nil {
						result = append(result, resp.GetRows()...)
						lastErr = nil
						break
					}
					lastErr = err
					time.Sleep(time.Duration(c.RetryInterval) * time.Second)
				}
				if lastErr != nil {
					errChan <- fmt.Errorf("failed to get page %d after %d retries: %w",
						page, c.MaxRetries, lastErr)
					return
				}
				resultChan <- result
			}

		}()
	}

	go func() {
		for page := c.StartPage + 1; page <= totalPage; page++ {
			pageChan <- page
		}
		close(pageChan)
	}()

	go func() {
		wg.Wait()
		close(errChan)
		close(resultChan)
	}()
	for {
		select {
		case err := <-errChan:
			return allResources, err
		case res, ok := <-resultChan:
			if !ok {
				return allResources, nil
			}
			allResources = append(allResources, res...)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

}
