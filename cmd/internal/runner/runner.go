package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"context"

	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/billy-playground/registry-load-tester/cmd/internal/image"
)

// Runner can be used to start a new test instance to download blobs and manifests.
type Runner struct {
	accessToken string
}

// It takes a JSON file as input and downloads the blobs and manifests specified in the file.
func NewRunner(accessToken string) *Runner {
	return &Runner{
		accessToken: accessToken,
	}
}

// StartNew starts a new test instance to download blobs and manifests.
func (r *Runner) StartNew(fileName string) error {
	// Parse JSON file
	data, err := parseJSON(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Record start time
	startTime := time.Now()

	// Set up repository client
	ctx := context.Background()
	repo, err := remote.NewRepository(data.Manifest)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}
	repo.Client = &auth.Client{
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", r.accessToken)},
		},
	}

	// Download manifest and blobs concurrently
	var wg sync.WaitGroup
	var totalCount = 1 + len(data.Blobs)
	var successCount atomic.Int32

	if data.Manifest != "" {
		wg.Add(1)
		go func(manifest string) {
			defer wg.Done()
			// Fetch the manifest
			_, rc, err := repo.Manifests().FetchReference(ctx, manifest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading manifest: %v\n", err)
				return
			}
			defer rc.Close()
			if _, err := io.Copy(io.Discard, rc); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading manifest response: %v\n", err)
			}
			successCount.Add(1)
		}(data.Manifest)
	}

	for _, blob := range data.Blobs {
		wg.Add(1)
		go func(blob string) {
			defer wg.Done()
			// Fetch the blob
			_, rc, err := repo.Blobs().FetchReference(ctx, blob)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading blob: %v\n", err)
				return
			}
			defer rc.Close()
			if _, err := io.Copy(io.Discard, rc); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading blob response: %v\n", err)
			}
			successCount.Add(1)
		}(blob)
	}

	wg.Wait()

	// Record end time and calculate elapsed time
	endTime := time.Now()
	downloadMilliseconds := endTime.Sub(startTime).Milliseconds()

	// Output results
	fmt.Printf("%s,%d,%d,%d,%d\n", fileName, data.Size, downloadMilliseconds, totalCount, successCount.Load())
	return nil
}

func parseJSON(filePath string) (*image.Data, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var blobData image.Data
	if err := json.Unmarshal(data, &blobData); err != nil {
		return nil, err
	}
	return &blobData, nil
}
