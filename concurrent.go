package main

import (
	"context"
	"sync"
)

func concurrentExecute[T comparable, R any](f func(T, context.Context) (R, error), items []T, concurrency int) (map[T]R, error) {
	var wg sync.WaitGroup
	results := make(map[T]R, len(items))
	var mu sync.Mutex
	sem := make(chan struct{}, concurrency) // limit concurrent executions
	// Use context to handle cancellation on error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	for _, item := range items {
		wg.Add(1)
		go func(param T){
			sem <- struct{}{} // acquire a slot
			defer func() { <-sem }() // release the slot
			defer wg.Done()
			select {
				case <-ctx.Done():
					return
				default:
			}
			res, err := f(param, ctx)
			if err != nil {
				errCh <- err
				cancel()
				return
			}
			mu.Lock()
			results[param] = res
			mu.Unlock()
		}(item)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
		case err := <-errCh:
			return nil, err
		case <-done:
			return results, nil
	}
}

func concurrentExecuteAny[T any, R any](f func(T, context.Context) (R, error), items []T, concurrency int) ([]R, error) {
	var wg sync.WaitGroup
	results := make([]R, 0, len(items))
	var mu sync.Mutex
	sem := make(chan struct{}, concurrency) // limit concurrent executions
	// Use context to handle cancellation on error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	for _, item := range items {
		wg.Add(1)
		go func(param T){
			sem <- struct{}{} // acquire a slot
			defer func() { <-sem }() // release the slot
			defer wg.Done()
			select {
				case <-ctx.Done():
					return
				default:
			}
			res, err := f(param, ctx)
			if err != nil {
				errCh <- err
				cancel()
				return
			}
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(item)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
		case err := <-errCh:
			return nil, err
		case <-done:
			return results, nil
	}
}

func concurrentExecuteAnyWithContext[T any, R any](f func(T, context.Context) (R, error), items []T, concurrency int, ctx context.Context) ([]R, error) {
	var wg sync.WaitGroup
	results := make([]R, 0, len(items))
	var mu sync.Mutex
	sem := make(chan struct{}, concurrency) // limit concurrent executions
	// Use context to handle cancellation on error
	errCh := make(chan error, 1)
	for _, item := range items {
		wg.Add(1)
		go func(param T){
			sem <- struct{}{} // acquire a slot
			defer func() { <-sem }() // release the slot
			defer wg.Done()
			select {
				case <-ctx.Done():
					return
				default:
			}
			res, err := f(param, ctx)
			if err != nil {
				errCh <- err
				return
			}
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(item)
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
		case err := <-errCh:
			return nil, err
		case <-done:
			return results, nil
	}
}