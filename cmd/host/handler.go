package main

import (
	"io"
	"log"
	"net/http"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func handler(
	store *wasmtime.Store,
	runner *wasmtime.Func,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		args := map[string]interface{}{
			"method":  r.Method,
			"url":     r.URL.String(),
			"headers": r.Header,
			"body":    bodyBytes,
		}
		ret, err := runner.Call(store, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		retMap, ok := ret.(map[string]interface{})
		if !ok {
			log.Printf("invalid return from module: %v", ret)
			http.Error(w, "invalid return from module", http.StatusInternalServerError)
			return
		}
		retHeaders, ok := retMap["headers"].(map[string][]string)
		if !ok {
			log.Printf("invalid headers returned from module: %v", retHeaders)
			http.Error(w, "invalid headers returned from module", http.StatusInternalServerError)
			return
		}
		retBody, ok := retMap["body"].([]byte)
		if !ok {
			log.Printf("invalid body returned from module: %v", retBody)
			http.Error(w, "invalid body returned from module", http.StatusInternalServerError)
			return
		}
		retStatus, ok := retMap["status"].(int)
		if !ok {
			log.Printf("invalid status returned from module: %v", retStatus)
			http.Error(w, "invalid status returned from module", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(retStatus)
		for key, vals := range retHeaders {
			for _, val := range vals {
				w.Header().Add(key, val)
			}
		}
		w.Write(retBody)
	})
}
