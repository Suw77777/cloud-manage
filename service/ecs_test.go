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
