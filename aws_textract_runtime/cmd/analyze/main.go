package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

func main() {
	documentClass := flag.String("document-class", "", "class ID declared REQUIRED in the approved profile")
	profile := flag.String("profile", "", "approved document profile JSON")
	securityReceipt := flag.String("security-receipt", "", "ADMITTED secure local-file receipt JSON")
	input := flag.String("input", "", "local PDF/image already admitted by the security gate")
	output := flag.String("output", "", "new evidence directory")
	maxBytes := flag.Int64("max-bytes", textractruntime.MaxSyncBytes, "approved synchronous byte limit, at most 10 MiB")
	flag.Parse()
	region, err := textractruntime.ApprovedRegion(*profile, *documentClass)
	fatalIf(err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	fatalIf(err)
	client := textract.NewFromConfig(cfg)
	receipt, err := textractruntime.AnalyzeToEvidence(ctx, client, textractruntime.Request{Region: region, DocumentClass: *documentClass, InputPath: *input, OutputDirectory: *output, ProfilePath: *profile, SecurityReceiptPath: *securityReceipt, MaxBytes: *maxBytes})
	fatalIf(err)
	body, _ := json.Marshal(map[string]any{"status": "AWS_TEXTRACT_ANALYSIS_PASS", "receipt": receipt})
	fmt.Println(string(body))
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
