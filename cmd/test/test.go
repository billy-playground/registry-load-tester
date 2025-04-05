package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/billy-playground/registry-load-tester/cmd/internal/runner"
	"github.com/billy-playground/registry-load-tester/cmd/option"
)

func main() {
	// Input arguments:
	// 0. Input check and help
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <num_instances>[=<size>/<duration>] <registry_domain> <token_mode> [<registry_endpoint>]")
		fmt.Println("num_instances: Number of instances to run")
		fmt.Println("registry_domain: Domain of the registry")
		fmt.Println("token_mode: Token mode (none, anonymous, token=<token>)")
		fmt.Println("registry_endpoint: Endpoint of the registry (default: registry)")
		os.Exit(1)
	}
	// 1. Number of instances to run and how to run it
	instanceOption, err := option.ParseInstanceOption(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing instance option: %v\n", err)
		os.Exit(1)
	}

	// 2. Registry domain
	var registry = os.Args[2]
	token, err := option.ParseTokenOption(os.Args[3], registry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing token option: %v\n", err)
		os.Exit(1)
	}
	// 4. Registry endpoint (default to registry)
	endpoint := registry
	if len(os.Args) > 4 {
		endpoint = os.Args[4]
		// Resolve registry to endpoint
		http.DefaultClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if after, found := strings.CutPrefix(addr, registry); found {
					return net.Dial(network, endpoint+after)
				}
				return net.Dial(network, addr)
			},
		}
	}

	fmt.Println("json_file,total_size,download_milliseconds,total_count,success_count")
	// Generate a random array of picked JSON files with length numInstances
	files := make([]string, instanceOption.Count)

	allFiles, err := filepath.Glob(filepath.Join("assets/images", "*.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading image JSON files: %v\n", err)
		os.Exit(1)
	}
	if len(allFiles) == 0 {
		fmt.Fprintf(os.Stderr, "No JSON files found in assets/images\n")
		os.Exit(1)
	}
	for i := range instanceOption.Count {
		files[i] = allFiles[rand.Intn(len(allFiles))]
	}

	// Run instanceOption.Count in total
	start := time.Now()
	next := start
	var wg sync.WaitGroup
	testRunner := runner.NewRunner(token, registry)
	batchingEnabled := instanceOption.BatchSize > 0 && instanceOption.BatchInterval > 0
	for i := 0; i < instanceOption.Count; i++ {
		wg.Add(1)
		// start at most instanceOption.BatchSize in instanceOption.BatchInterval
		if batchingEnabled && i%instanceOption.BatchSize == 0 {
			now := time.Now()
			toWait := next.Sub(time.Now())
			if toWait <= 0 {
				next = now.Add(instanceOption.BatchInterval)
			} else {
				time.Sleep(toWait)
			}
		}
		fmt.Printf("[%s] Starting instance %d/%d with file %s\n", time.Now().Format("0000-00-00 00:00:00"), i+1, instanceOption.Count, files[i])
		go func() {
			defer wg.Done()
			_ = testRunner.StartNew(files[i])
		}()
	}
	wg.Wait()
	fmt.Printf("Total time taken: %.2f seconds\n", time.Since(start).Seconds())
}
