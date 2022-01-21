package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	wasmtime "github.com/bytecodealliance/wasmtime-go"
)

func main() {
	moduleLoc := flag.String("module", "", "module location")
	isWat := flag.Bool("wat", false, "whether the module is WAT (Web Assembly Text Format)")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()
	store := wasmtime.NewStore(wasmtime.NewEngine())
	rawWasm, err := os.ReadFile(*moduleLoc)
	if err != nil {
		log.Fatalf("error reading WASM file: %s", err)
	}
	if *isWat {
		r, err := wasmtime.Wat2Wasm(string(rawWasm))
		if err != nil {
			log.Fatalf("Error converting from WAT: %s", err)
		}
		rawWasm = r
	}
	module, err := wasmtime.NewModule(store.Engine, rawWasm)
	if err != nil {
		log.Fatalf("Error creating module: %s", err)
	}
	instance, err := wasmtime.NewInstance(store, module, nil)
	if err != nil {
		log.Fatalf("Error creating instance: %s", err)
	}
	run := instance.GetExport(store, "run").Func()

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting host server on %s", addr)
	http.ListenAndServe(addr, handler(store, run))
}
