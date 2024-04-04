# GoLang WebAssembly

### Prerequisites

- [GoLang](https://go.dev/doc/install)

## Getting Started

1. Clone this repository
    ```shell
    git clone https://github.com/GabiBizdoc/random
    cd wasm-golang
    ```
2. Copy the wasm_exec.js file to your project directory
   ```shell
   cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./www/
   ```
3. Build the Wasm file using the following command
   ```shell
   GOARCH=wasm GOOS=js go build -o www/main.wasm -ldflags '-s'
   ```
4. Lastly, serve the www directory using an HTTP file server
   ```shell
   go run ./cmd/server.go
   ```