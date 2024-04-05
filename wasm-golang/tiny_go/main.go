package main

//go:wasm-module main
//export addInt
func addInt(x, y int) int {
	return x + y
}

//go:wasm-module main
//export addFloat
func addFloat(x, y float64) float64 {
	return x + y
}

// main is required for the `wasi` target, even if it isn't used.
func main() {}
