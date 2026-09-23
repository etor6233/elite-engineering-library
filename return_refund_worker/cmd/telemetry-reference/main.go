// AUTHORED finite reference load probe. No business data, provider or refund call.
package main

import (
	"context"
	"elite.local/return-refund-worker/internal/operationaltelemetry"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

func run() error {
	if len(os.Args) == 3 && os.Args[1] == "--certificates" {
		return referenceCertificates(os.Args[2])
	}
	if len(os.Args) != 5 {
		return operationaltelemetry.ErrConfig
	}
	n, e := strconv.Atoi(os.Args[3])
	if e != nil || n < 1 || n > 2000 {
		return operationaltelemetry.ErrConfig
	}
	parallel, e := strconv.Atoi(os.Args[4])
	if e != nil || parallel < 1 || parallel > 16 || n%parallel != 0 {
		return operationaltelemetry.ErrConfig
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	began := time.Now()
	var wg sync.WaitGroup
	errs := make(chan error, parallel)
	for range parallel {
		wg.Go(func() {
			r, e := operationaltelemetry.FromFile(os.Args[1], os.Args[2])
			if e != nil {
				errs <- e
				return
			}
			defer r.Close()
			for range n / parallel {
				if e = r.Observe(ctx, "idle", time.Millisecond); e != nil {
					errs <- e
					cancel()
					return
				}
			}
		})
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(struct {
		Observations int     `json:"observations"`
		Instances    int     `json:"instances"`
		Seconds      float64 `json:"seconds"`
	}{n, parallel, time.Since(began).Seconds()})
}
func main() {
	if run() != nil {
		os.Stderr.WriteString("REFERENCE_TELEMETRY_PROBE_FAILED\n")
		os.Exit(1)
	}
}
