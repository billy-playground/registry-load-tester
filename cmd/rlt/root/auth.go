package root

import (
	"fmt"
	"sync"
	"time"

	"github.com/billy-playground/registry-load-tester/cmd/internal/runner"
	"github.com/billy-playground/registry-load-tester/cmd/option"
	"github.com/billy-playground/registry-load-tester/internal/auth"
	"github.com/spf13/cobra"
)

type authOptions struct {
	option.Instance
	option.Registry
	refreshToken string
}

func authCmd() *cobra.Command {
	var opts authOptions

	authCmd := &cobra.Command{
		Use:   "auth  <num_instances>[=<size>/<duration>] <registry_domain>",
		Short: "authenticate to a registry",
		Long: `run authentication workloads simultaneously with customized options

Example - authenticate 10 images against registry.example.com without using any token.
  rlt auth 10 registry.example.com

Example - authenticate 100 images against registry.example.com, starting 10 instances every 500 milliseconds using the specified token.
  rlt auth 100=10/500ms registry.example.com --refresh-token=$registry_token

Example - authenticate 20 images against registry.example.com via a custom endpoint -e cus.fe.example.com.
  rlt auth 20 registry.example.com none -e cus.fe.example.com
`,
		Args: cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Setup arguments
			opts.Instance.SetFlag(args[0])
			opts.Registry.SetFlag(args[1])

			// Parse options
			if err := opts.Instance.Parse(); err != nil {
				return fmt.Errorf("Error parsing instance option: %v\n", err)
			}
			return opts.Registry.Parse()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAuth(opts)
		},
	}

	opts.Registry.ApplyFlags(authCmd.Flags())
	authCmd.Flags().StringVarP(&opts.refreshToken, "refresh-token", "r", "", "Token used for refreshing")

	return authCmd
}

func runAuth(opts authOptions) error {
	fmt.Println("timestamp,is_success")

	// Run instanceOption.Count in total
	// TODO: isolate this
	start := time.Now()
	next := start
	var wg sync.WaitGroup

	authHeader, err := auth.GetAuthHeader(opts.RegistryDomain)
	if err != nil {
		return err
	}

	if authHeader == "" {
		// no way to run auth at all since the authorization header is empty
		return nil
	}

	realm, service, err := auth.ParseRealmAndService(authHeader)
	if err != nil {
		return fmt.Errorf("failed to parse auth header: %v", err)
	}

	testRunner := runner.NewAuthRunner(realm, service, opts.RegistryDomain, opts.refreshToken)
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
			err = testRunner.StartNew()
			fmt.Printf("%s,%t\n", time.Now().Format(time.RFC3339), err == nil)
		}()
	}
	wg.Wait()
	fmt.Printf("Total time taken: %.2f seconds\n", time.Since(start).Seconds())
	return nil
}
