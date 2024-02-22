// Main FaaS server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func httpHandler(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 1 {
		http.Error(w, "want /{moduleName} prefix", http.StatusBadRequest)
		return
	}
	mod := parts[0]
	log.Printf("module %v requested with query %v", mod, req.URL.Query())

	env := map[string]string{
		"http_path":   req.URL.Path,
		"http_method": req.Method,
		"http_host":   req.Host,
		"http_query":  req.URL.Query().Encode(),
		"remote_addr": req.RemoteAddr,
	}

	modpath := fmt.Sprintf("../target/%v.wasm", mod)
	log.Printf("loading module %v", modpath)
	out, err := InvokeWasmModule(mod, modpath, env)
	if err != nil {
		log.Printf("error loading module %v", modpath)
		http.Error(w, "unable to find module "+modpath, http.StatusNotFound)
		return
	}

	// The module's stdout is written into the response.
	fmt.Fprint(w, out)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", httpHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
