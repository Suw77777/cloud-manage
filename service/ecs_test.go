package service

import "testing"

func TestListInstances_EmptyCredentials(t *testing.T) {
	svc := NewECSService()

	_, err := svc.ListInstances("", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeyId")
	}

	_, err = svc.ListInstances("key", "", "cn-hangzhou")
	if err == nil {
		t.Error("expected error for empty accessKeySecret")
	}

	_, err = svc.ListInstances("key", "secret", "")
	if err == nil {
		t.Error("expected error for empty region")
	}
}

func TestNewECSService(t *testing.T) {
	svc := NewECSService()
	if svc == nil {
		t.Error("expected non-nil ECSService")
	}
}

func TestListInstancesMultiRegion_EmptyRegions(t *testing.T) {
	svc := NewECSService()
	results := svc.ListInstancesMultiRegion("key", "secret", []string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty regions, got %d", len(results))
	}
}

func TestListInstancesMultiRegion_MultipleRegions(t *testing.T) {
	svc := NewECSService()
	// With fake credentials, this will error per-region but should not panic.
	regions := []string{"cn-hangzhou", "cn-beijing", "cn-shanghai"}
	results := svc.ListInstancesMultiRegion("fake-key", "fake-secret", regions)

	if len(results) != len(regions) {
		t.Errorf("expected %d results, got %d", len(regions), len(results))
	}

	for i, r := range results {
		if r.Region != regions[i] {
			t.Errorf("result[%d].Region = %q, want %q", i, r.Region, regions[i])
		}
		// With fake credentials, we expect an error.
		if r.Error == "" {
			t.Errorf("result[%d] expected error for fake credentials", i)
		}
	}
}
