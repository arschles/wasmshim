package main

import (
	"context"
	"log"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	ctx := namespaces.WithNamespace(context.TODO(), "default")

	// Create containerd client
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}

	// Get the image ref to create the container for
	ctx = namespaces.WithNamespace(ctx, "default")
	img, err := client.Pull(ctx, "docker.io/library/redis:alpine", containerd.WithPullUnpack) //"docker.io/library/busybox:latest")
	if err != nil {
		panic(err)
	}

	// img, err = client.GetImage(ctx, "docker.io/library/redis:alpine")
	// if err != nil {
	// 	panic(err)
	// }

	// set options we will pass to the shim (not really setting anything here, but we could)
	// var opts v1opts.Options

	// Create a container object in containerd
	cntr, err := client.NewContainer(ctx, "myContainer",
		// All the basic things needed to create the container
		containerd.WithSnapshotter("overlayfs"),
		containerd.WithNewSnapshot("redis-snapshot", img),
		// containerd.WithImage(img),
		containerd.WithNewSpec(oci.WithImageConfig(img)),

		// Set the option for the shim we want
		// the shim name ends up being
		// shim binary should be containerd-shim-wasm-v1
		containerd.WithRuntime("io.containerd.wasm.v1", nil), //, &opts),
	)
	if err != nil {
		panic(err)
	}
	defer cntr.Delete(ctx)
	task, err := cntr.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		panic(err)
	}
	defer task.Delete(ctx)
	exitStatusCh, err := task.Wait(ctx)
	if err != nil {
		panic(err)
	}
	if err := task.Start(ctx); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
		panic(err)
	}

	select {
	case exitStatus := <-exitStatusCh:
		status, exitedAt, err := exitStatus.Result()
		if err != nil {
			panic(err)
		}
		log.Println("Exit status:", status, "exited at:", exitedAt)
	case <-time.After(3 * time.Second):
		log.Fatal("didn't receive termination signal within 3 seconds. I DIED")
	}
}
