package root

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"github.com/billy-playground/registry-load-tester/cmd/internal/runner"
	"github.com/billy-playground/registry-load-tester/cmd/option"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	option.Instance
	option.Pull
}

func pullCmd() *cobra.Command {
	var opts pullOptions

	pullCmd := &cobra.Command{
		Use:   "pull  <num_instances>[=<size>/<duration>] <registry_domain> <token_mode>",
		Short: "Registry Load Tester",
		Long:  "A tool to test the load on a registry by running multiple instances.",
		Args:  cobra.MinimumNArgs(3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Setup arguments
			opts.Instance.SetFlag(args[0])
			opts.Pull.SetFlag(args[1], args[2])

			// Parse options
			if err := opts.Instance.Parse(); err != nil {
				return fmt.Errorf("Error parsing instance option: %v\n", err)
			}
			return opts.Pull.Parse()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull(opts)
		},
	}

	opts.ApplyFlags(pullCmd.Flags())

	return pullCmd
}

func runPull(opts pullOptions) error {
	fmt.Println("json_file,total_size,download_milliseconds,total_count,success_count")
	// Generate a random array of picked JSON files with length numInstances
	files := make([]string, opts.Count)

	allFiles, err := filepath.Glob(filepath.Join("assets/images", "*.json"))
	if err != nil {
		return fmt.Errorf("Error reading image JSON files: %v\n", err)
	}
	if len(allFiles) == 0 {
		return fmt.Errorf("No JSON files found in assets/images\n")
	}
	for i := range opts.Count {
		files[i] = allFiles[rand.Intn(len(allFiles))]
	}

	// Run instanceOption.Count in total
	start := time.Now()
	next := start
	var wg sync.WaitGroup
	testRunner := runner.NewRunner(opts.Token.AccessToken, opts.RegistryDomain)
	batchingEnabled := opts.BatchSize > 0 && opts.BatchInterval > 0
	for i := 0; i < opts.Count; i++ {
		wg.Add(1)
		if batchingEnabled && i%opts.BatchSize == 0 {
			if toWait := next.Sub(time.Now()); toWait > 0 {
				// wait for the next batch
				time.Sleep(toWait)
			}
			next = time.Now().Add(opts.BatchInterval)
		}
		go func() {
			defer wg.Done()
			_ = testRunner.StartNew(files[i])
		}()
	}
	wg.Wait()
	fmt.Printf("Total time taken: %.2f seconds\n", time.Since(start).Seconds())
	return nil
}
