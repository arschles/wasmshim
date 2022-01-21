# WASM Shim Host

This is the process that the ContainerD shim starts up to run one or more WASM modules. It is started as a standard process with command line arguments, and knows how to run a module in response to an HTTP request, but its value comes from the functionality it provides before and after it runs the module:

1. It defines a standard WASM module interface for sending HTTP requests to WASM, and a standard interface for receiving them back from the module
2. It runs a highly scalable HTTP server to receive HTTP requests from the network, send them to the module, and send them back to the client
3. (TODO) It serves prometheus metrics about loaded modules, in-flight HTTP requests, and more
4. (TODO) It exposes an "admin" API that allows one to modify the running state of modules, load new modules, and more
