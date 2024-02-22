# WebAssembly FaaS

## Overview

A WebAssembly Function as a Service (FaaS) implemented as an HTTP server.
The service allows users to register wasm modules compiled from any supported programming language
as functions on the server. These functions can then be invoked on behalf of a client via HTTP requests.

## Build and Run the Server
### Registering Modules
Registering a module as a function with the service is a matter of adding a `module.wasm` file in
the target directory. The modules can be built be any means for example the `miller_rabin.go` in
the `examples/primality-tester` can be built running the following in the project directory:
```
GOOS=wasip1 GOARCH=wasm gotip build -o target/miller-rabin.wasm examples/primality-tester/miller_rabin.go
```
### Running the Server
From the `src` directory run 
```
go run .
```

### Making Requests as a Client
A FaaS client issue an HTTP GET request with a path consisting of both a module name
and a query string (as appropriate for invoked function).

As an example if running the server on your local machine a client can invoke the `miller-rabin`
module to test the primality of the natural number 170141183460469231731687303715884105727
by issuing such a request:
```
curl "localhost:8080/miller-rabin?number=170141183460469231731687303715884105727"
```
