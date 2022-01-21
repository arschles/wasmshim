package main

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	v1opts "github.com/containerd/containerd/pkg/runtimeoptions/v1"
)

func main() {
	ctx := namespaces.WithNamespace(context.TODO(), "default")

	// Create containerd client
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}

	// Get the image ref to create the container for
	img, err := client.GetImage(ctx, "docker.io/library/busybox:latest")
	if err != nil {
		panic(err)
	}

	// set options we will pass to the shim (not really setting anything here, but we could)
	var opts v1opts.Options

	// Create a container object in containerd
	cntr, err := client.NewContainer(ctx, "myContainer",
		// All the basic things needed to create the container
		containerd.WithSnapshotter("overlayfs"),
		containerd.WithNewSnapshot("myContainer-snapshot", img),
		containerd.WithImage(img),
		containerd.WithNewSpec(oci.WithImageConfig(img)),

		// Set the option for the shim we want
		// the shim name ends up being
		// containerd-shim-wasm-v1
		containerd.WithRuntime("io.containerd.wasm.v1", &opts),
	)
	if err != nil {
		panic(err)
	}

	// cleanup
	cntr.Delete(ctx)
}
