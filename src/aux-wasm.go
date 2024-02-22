// Auxiliary functions for initializing and interacting with
// wasm modules.
package main

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log"
	"os"
)

// initRuntime instantiates the wasm runtime, setting up exported functions from the host
// that the wasm module can use for logging purposes.
func initRuntime(r wazero.Runtime, modname string, ctx context.Context) error {
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(v uint32) {
			log.Printf("[%v]: %v", modname, v)
		}).
		Export("log_i32").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, ptr uint32, len uint32) {
			// Read the string from the module's exported memory.
			if bytes, ok := mod.Memory().Read(ptr, len); ok {
				log.Printf("[%v]: %v", modname, string(bytes))
			} else {
				log.Printf("[%v]: log_string: unable to read wasm memory", modname)
			}
		}).
		Export("log_string").
		Instantiate(ctx)
	return err
}

// InvokeWasmModule invokes the given WASM module (given as a file path),
// setting its env vars according to env. Returns the module's stdout.
func InvokeWasmModule(modname string, wasmPath string, env map[string]string) (string, error) {
	ctx := context.Background()

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	err := initRuntime(r, modname, ctx)
	if err != nil {
		return "", err
	}

	wasmObj, err := os.ReadFile(wasmPath)
	if err != nil {
		return "", err
	}

	// Set up stdout redirection and env vars for the module.
	var stdoutBuf bytes.Buffer
	config := wazero.NewModuleConfig().WithStdout(&stdoutBuf)

	for k, v := range env {
		config = config.WithEnv(k, v)
	}

	// Instantiate the module. This invokes the _start function by default.
	_, err = r.InstantiateWithConfig(ctx, wasmObj, config)
	if err != nil {
		return "", err
	}

	return stdoutBuf.String(), nil
}
