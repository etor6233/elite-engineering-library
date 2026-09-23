// AUTHORED profile materialization glue. The same Go loader used by the host
// validates all source/override data before creating an absent destination.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	wc "elite.local/enterprise/internal/warrantyclaim"
)

type activation struct {
	Schema                string            `json:"schema"`
	Scope                 string            `json:"scope"`
	ProfileFile           string            `json:"profile_file"`
	ProfileSHA256         string            `json:"profile_sha256"`
	SourceSHA256          string            `json:"source_sha256"`
	TenantID              string            `json:"tenant_id"`
	OrganizationID        string            `json:"organization_id"`
	FactoryOrganizationID string            `json:"factory_organization_id"`
	Overrides             map[string]string `json:"overrides"`
	ProductionAuthorized  bool              `json:"production_authorized"`
}

func materialize(args []string) (activation, error) {
	var empty activation
	flags := flag.NewFlagSet("warranty-profile", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	input := flags.String("input", "", "existing profile JSON")
	expected := flags.String("sha256", "", "exact source profile SHA256")
	output := flags.String("output", "", "absent output directory")
	tenant := flags.String("tenant-id", "", "explicit tenant override")
	org := flags.String("organization-id", "", "explicit store override")
	factory := flags.String("factory-organization-id", "", "explicit factory override")
	if err := flags.Parse(args); err != nil {
		return empty, err
	}
	if flags.NArg() != 0 || *input == "" || *output == "" || !wc.ValidSHA(*expected) {
		return empty, errors.New("input, sha256 and absent output are required")
	}
	f, err := os.Open(*input)
	if err != nil {
		return empty, err
	}
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() < 1 || info.Size() > 32768 {
		f.Close()
		return empty, errors.New("source profile must be a bounded regular file")
	}
	raw, err := io.ReadAll(io.LimitReader(f, 32769))
	closeErr := f.Close()
	if err != nil {
		return empty, err
	}
	if closeErr != nil {
		return empty, closeErr
	}
	source, err := wc.LoadProfile(raw, *expected)
	if err != nil {
		return empty, err
	}
	doc := source.Document()
	overrides := map[string]string{}
	if *tenant != "" {
		doc.TenantID = *tenant
		overrides["tenant_id"] = *tenant
	}
	if *org != "" {
		doc.OrganizationID = *org
		overrides["organization_id"] = *org
	}
	if *factory != "" {
		doc.FactoryOrganizationID = *factory
		overrides["factory_organization_id"] = *factory
	}
	profileBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return empty, err
	}
	profileBytes = append(profileBytes, '\n')
	profileHash := wc.SHA(profileBytes)
	selected, err := wc.LoadProfile(profileBytes, profileHash)
	if err != nil {
		return empty, err
	}
	boundTenant, boundOrg := selected.Scope()
	value := activation{Schema: "elite-warranty-activation/v1", Scope: "LIBRARY_INFRASTRUCTURE", ProfileFile: "profile.json", ProfileSHA256: profileHash, SourceSHA256: *expected, TenantID: boundTenant, OrganizationID: boundOrg, FactoryOrganizationID: doc.FactoryOrganizationID, Overrides: overrides, ProductionAuthorized: false}
	activationBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return empty, err
	}
	activationBytes = append(activationBytes, '\n')
	dest, err := filepath.Abs(*output)
	if err != nil {
		return empty, err
	}
	// Mkdir is exclusive; existing directories/files are never reused. A partial
	// failure is preserved for inspection and cannot be mistaken for a new run.
	if err = os.Mkdir(dest, 0700); err != nil {
		return empty, err
	}
	for _, artifact := range []struct {
		name string
		data []byte
	}{{"profile.json", profileBytes}, {"activation.json", activationBytes}} {
		file, err := os.OpenFile(filepath.Join(dest, artifact.name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return empty, err
		}
		_, writeErr := file.Write(artifact.data)
		syncErr := file.Sync()
		closeErr := file.Close()
		if writeErr != nil {
			return empty, writeErr
		}
		if syncErr != nil {
			return empty, syncErr
		}
		if closeErr != nil {
			return empty, closeErr
		}
	}
	return value, nil
}
func main() {
	result, err := materialize(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "warranty profile materialization failed:", err)
		os.Exit(1)
	}
	if err = json.NewEncoder(os.Stdout).Encode(result); err != nil {
		os.Exit(1)
	}
}
