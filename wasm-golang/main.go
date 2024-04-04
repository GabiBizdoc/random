//go:build wasm

package main

import (
	"sync"
	"syscall/js"
)

func addIntsFunction(this js.Value, p []js.Value) interface{} {
	var result int
	for _, value := range p {
		result += value.Int()
	}
	return js.ValueOf(result)
}
func addFloatFunction(this js.Value, p []js.Value) interface{} {
	var result float64
	for _, value := range p {
		result += value.Float()
	}
	return js.ValueOf(result)
}
func main() {
	js.Global().Set("addInt", js.FuncOf(addIntsFunction))
	js.Global().Set("addFloat", js.FuncOf(addFloatFunction))
	select {}
}

// this doesn't work, wasm is single-threaded;
// multiple goroutines won't make things faster in CPU-bound operations
func testConcurrency(n int) {
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			fibo(38)
			wg.Done()
		}()
	}
	wg.Wait()
}
func fibo(n int) int {
	if n <= 1 {
		return n
	}
	return fibo(n-1) + fibo(n-2)
}
