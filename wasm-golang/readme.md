# GoLang WebAssembly

### Prerequisites

- [GoLang](https://go.dev/doc/install)
- [TinyGo](https://tinygo.org/getting-started/install/) Optional*

## Getting Started

1. Clone this repository
    ```shell
    git clone https://github.com/GabiBizdoc/random
    cd wasm-golang
    ```
2. Copy the wasm_exec.js file to your project directory.

   You must use same version of wasm_exec.js as the version of go/tinygo you are using to compile
   ```shell
   cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./www/
   ```
   #### TinyGO
   ```shell
   cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js ./www_tinygo/
   ```
3. Build the Wasm file using the following command
   ```shell
   GOARCH=wasm GOOS=js go build -o www/main.wasm -ldflags '-s'
   ```
   #### TinyGO
   ```shell
   tinygo build -o www_tinygo/main.wasm -target=wasi tiny_go/main.go
   ```
   Optimized for speed and size.
   ```shell
   tinygo build -o www_tinygo/main.wasm -opt=2 -no-debug -target=wasi tiny_go/main.go
   ```
   - Optimize for speed `-opt=2`. By default, tinygo will optimize for size `-opt=z`
   - Remove debug symbols `-no-debug`
   - Disable goroutines `-scheduler=none`
   - Don’t print panic messages `-panic=trap`
   - Disable the GC `-gc=leaking` (for very short-lived programs only)
   
   For more details, please visit https://tinygo.org/docs/guides/optimizing-binaries/
   ```shell
   tinygo build -o www_tinygo/main.wasm -no-debug -scheduler=none -panic=trap -gc=leaking -target=wasi tiny_go/main.go  
   ```
4. Lastly, serve the www directory using an HTTP file server
   ```shell
   go run ./cmd/server.go www
   ```
   ```shell
   go run ./cmd/server.go www_tinygo
   ```
   
## DEMO
   - [Wasm TinyGO Demo](https://gabibizdoc.github.io/random/wasm-tinygo-demo/)
   - [Wasm GO Demo](https://gabibizdoc.github.io/random/wasm-go-demo/)