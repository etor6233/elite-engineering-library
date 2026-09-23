package finops

import (
	"sort"
	"sync"
)

// Inventory holds every deployed resource so nothing escapes billing. A
// resource is unique per (tenant, provider, id); duplicates are rejected.
type Inventory struct {
	mu        sync.RWMutex
	resources map[string]Resource // key: tenant + "\x00" + provider + "\x00" + id
}

// NewInventory returns an empty inventory.
func NewInventory() *Inventory {
	return &Inventory{resources: make(map[string]Resource)}
}

func invKey(tenant, provider, id string) string {
	return tenant + "\x00" + provider + "\x00" + id
}

// Register adds a validated, tagged resource; duplicates are rejected.
func (i *Inventory) Register(r Resource) error {
	if err := r.Validate(); err != nil {
		return err
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	k := invKey(r.TenantID, r.Provider, r.ID)
	if _, ok := i.resources[k]; ok {
		return ErrDuplicate
	}
	i.resources[k] = r
	return nil
}

// List returns the tenant's resources sorted by (provider, id).
func (i *Inventory) List(tenant string) []Resource {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var out []Resource
	for _, r := range i.resources {
		if r.TenantID == tenant {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].Provider != out[b].Provider {
			return out[a].Provider < out[b].Provider
		}
		return out[a].ID < out[b].ID
	})
	return out
}

// AuditResult reports what escaped the inventory: untagged actual resources,
// unknown (drifted-in) resources, and missing (leaked) inventory resources.
type AuditResult struct {
	Untagged []Resource
	Unknown  []Resource
	Missing  []Resource
}

// Clean reports whether nothing escaped: no untagged, no drift, no leak.
func (a AuditResult) Clean() bool {
	return len(a.Untagged) == 0 && len(a.Unknown) == 0 && len(a.Missing) == 0
}

// Audit compares the tenant's actual cloud resources against the inventory.
func (i *Inventory) Audit(tenant string, actual []Resource) AuditResult {
	i.mu.RLock()
	defer i.mu.RUnlock()
	expected := make(map[string]Resource)
	for _, r := range i.resources {
		if r.TenantID == tenant {
			expected[r.Key()] = r
		}
	}
	seen := make(map[string]bool)
	var res AuditResult
	for _, r := range actual {
		if r.Validate() != nil {
			res.Untagged = append(res.Untagged, r)
			continue
		}
		if _, ok := expected[r.Key()]; !ok {
			res.Unknown = append(res.Unknown, r) // drifted-in, not inventoried
			continue
		}
		seen[r.Key()] = true
	}
	for k, r := range expected {
		if !seen[k] {
			res.Missing = append(res.Missing, r) // leaked out of inventory
		}
	}
	return res
}

// MonthlyCost returns the tenant's total estimated monthly cost.
func (i *Inventory) MonthlyCost(tenant string) int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var total int64
	for _, r := range i.resources {
		if r.TenantID == tenant {
			total += r.EstMonthlyMinorUnits
		}
	}
	return total
}
