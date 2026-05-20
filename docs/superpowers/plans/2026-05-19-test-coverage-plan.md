# 测试覆盖率提升实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ []`) syntax for tracking.

**Goal:** 为 Cloud Manage 项目建立完整的分层测试体系，覆盖 Provider、Service、CLI 三层，支持 Mock 和集成测试两种模式。

**Architecture:** 采用接口抽象 + Mock 模式，Provider 层定义接口，Service 层依赖接口而非具体实现，测试时注入 Mock 对象。集成测试使用 build tag 隔离。

**Tech Stack:** Go testing, testify/assert (可选), build tags

---

## 文件结构

### 新建文件

| 文件 | 职责 |
|------|------|
| `provider/interfaces.go` | 定义 Provider 接口 |
| `provider/mock_ecs.go` | ECS Mock 实现 |
| `provider/mock_cms.go` | CMS Mock 实现 |
| `provider/mock_sls.go` | SLS Mock 实现 |
| `provider/mock_oss.go` | OSS Mock 实现 |
| `provider/aliyun/ecs_test.go` | ECS Provider 测试 |
| `provider/aliyun/cms_test.go` | CMS Provider 测试 |
| `provider/aliyun/sls_test.go` | SLS Provider 测试 |
| `provider/aliyun/oss_test.go` | OSS Provider 测试 |
| `service/cms_test.go` | CMS Service 测试 |
| `service/sls_test.go` | SLS Service 测试 |
| `service/oss_test.go` | OSS Service 测试 |
| `cmd/cli/cli_test.go` | CLI 端到端测试 |

### 修改文件

| 文件 | 修改内容 |
|------|----------|
| `service/ecs.go` | 使用接口替代具体类型 |
| `service/cms.go` | 使用接口替代具体类型 |
| `service/sls.go` | 使用接口替代具体类型 |
| `service/oss.go` | 使用接口替代具体类型 |
| `app.go` | 使用接口替代具体类型 |

---

## Task 1: 定义 Provider 接口

**Files:**
- Create: `provider/interfaces.go`

- [ ] **Step 1: 创建 ECS Provider 接口**

```go
package provider

// ECSProvider defines the interface for ECS operations.
type ECSProvider interface {
	DescribeInstances(pageNumber, pageSize int32) ([]ECSInstance, int32, error)
	DescribeInstanceDetail(instanceId string) (*InstanceDetail, error)
	StartInstance(instanceId string) error
	StopInstance(instanceId string, forceStop bool) error
	RebootInstance(instanceId string, forceStop bool) error
}

// ECSInstance represents an ECS instance.
type ECSInstance struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	ZoneId       string `json:"zoneId"`
	PublicIp     string `json:"publicIp"`
	PrivateIp    string `json:"privateIp"`
	CreationTime string `json:"creationTime"`
}

// InstanceDetail holds detailed information about an ECS instance.
type InstanceDetail struct {
	InstanceId         string   `json:"instanceId"`
	InstanceName       string   `json:"instanceName"`
	Description        string   `json:"description"`
	HostName           string   `json:"hostName"`
	Status             string   `json:"status"`
	RegionId           string   `json:"regionId"`
	ZoneId             string   `json:"zoneId"`
	InstanceType       string   `json:"instanceType"`
	Cpu                int32    `json:"cpu"`
	Memory             int32    `json:"memory"`
	ImageId            string   `json:"imageId"`
	InternetChargeType string   `json:"internetChargeType"`
	CreationTime       string   `json:"creationTime"`
	ExpiredTime        string   `json:"expiredTime"`
	StoppedMode        string   `json:"stoppedMode"`
	PublicIp           []string `json:"publicIp"`
	PrivateIp          []string `json:"privateIp"`
	SecurityGroupIds   []string `json:"securityGroupIds"`
}
```

- [ ] **Step 2: 添加 CMS Provider 接口**

在 `provider/interfaces.go` 中追加：

```go
// CMSProvider defines the interface for CloudMonitor operations.
type CMSProvider interface {
	GetECSMetrics(instanceId string) (*ECSMetricData, error)
}

// ECSMetricData holds metric data for an ECS instance.
type ECSMetricData struct {
	InstanceId         string   `json:"instanceId"`
	CPUUtilization     *float64 `json:"cpuUtilization,omitempty"`
	MemoryUtilization  *float64 `json:"memoryUtilization,omitempty"`
	DiskReadBPS        *float64 `json:"diskReadBps,omitempty"`
	DiskWriteBPS       *float64 `json:"diskWriteBps,omitempty"`
	InternetRX         *float64 `json:"internetRx,omitempty"`
	InternetTX         *float64 `json:"internetTx,omitempty"`
	UpdateTime         string   `json:"updateTime"`
}
```

- [ ] **Step 3: 添加 SLS Provider 接口**

在 `provider/interfaces.go` 中追加：

```go
// SLSProvider defines the interface for SLS operations.
type SLSProvider interface {
	ListLogStores(project string) ([]string, error)
	GetLogs(project, logstore, query string, from, to int64, maxLines int64) ([]LogEntry, int64, error)
}

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp int64             `json:"timestamp"`
	Content   map[string]string `json:"content"`
}
```

- [ ] **Step 4: 添加 OSS Provider 接口**

在 `provider/interfaces.go` 中追加：

```go
// OSSProvider defines the interface for OSS operations.
type OSSProvider interface {
	ListBuckets() ([]OSSBucket, error)
	ListObjects(bucket, prefix string, maxKeys int32) ([]OSSObject, bool, error)
}

// OSSBucket represents an OSS bucket.
type OSSBucket struct {
	Name             string `json:"name"`
	Location         string `json:"location"`
	CreationDate     string `json:"creationDate"`
	StorageClass     string `json:"storageClass"`
	ExtranetEndpoint string `json:"extranetEndpoint"`
	IntranetEndpoint string `json:"intranetEndpoint"`
}

// OSSObject represents an OSS object.
type OSSObject struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	ETag         string `json:"etag"`
	Type         string `json:"type"`
	StorageClass string `json:"storageClass"`
}
```

- [ ] **Step 5: 验证编译**

```bash
go build ./...
```

- [ ] **Step 6: 提交**

```bash
git add provider/interfaces.go
git commit -m "feat: define provider interfaces for testability"
```

---

## Task 2: 创建 Mock 实现

**Files:**
- Create: `provider/mock_ecs.go`
- Create: `provider/mock_cms.go`
- Create: `provider/mock_sls.go`
- Create: `provider/mock_oss.go`

- [ ] **Step 1: 创建 ECS Mock**

```go
// provider/mock_ecs.go
package provider

// MockECSProvider implements ECSProvider for testing.
type MockECSProvider struct {
	Instances       []ECSInstance
	TotalCount      int32
	InstanceDetail  *InstanceDetail
	Err             error
	StartErr        error
	StopErr         error
	RebootErr       error
}

func (m *MockECSProvider) DescribeInstances(pageNumber, pageSize int32) ([]ECSInstance, int32, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	return m.Instances, m.TotalCount, nil
}

func (m *MockECSProvider) DescribeInstanceDetail(instanceId string) (*InstanceDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.InstanceDetail, nil
}

func (m *MockECSProvider) StartInstance(instanceId string) error {
	return m.StartErr
}

func (m *MockECSProvider) StopInstance(instanceId string, forceStop bool) error {
	return m.StopErr
}

func (m *MockECSProvider) RebootInstance(instanceId string, forceStop bool) error {
	return m.RebootErr
}
```

- [ ] **Step 2: 创建 CMS Mock**

```go
// provider/mock_cms.go
package provider

// MockCMSProvider implements CMSProvider for testing.
type MockCMSProvider struct {
	Metrics *ECSMetricData
	Err     error
}

func (m *MockCMSProvider) GetECSMetrics(instanceId string) (*ECSMetricData, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Metrics, nil
}
```

- [ ] **Step 3: 创建 SLS Mock**

```go
// provider/mock_sls.go
package provider

// MockSLSProvider implements SLSProvider for testing.
type MockSLSProvider struct {
	LogStores []string
	Logs      []LogEntry
	Count     int64
	Err       error
}

func (m *MockSLSProvider) ListLogStores(project string) ([]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.LogStores, nil
}

func (m *MockSLSProvider) GetLogs(project, logstore, query string, from, to int64, maxLines int64) ([]LogEntry, int64, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	return m.Logs, m.Count, nil
}
```

- [ ] **Step 4: 创建 OSS Mock**

```go
// provider/mock_oss.go
package provider

// MockOSSProvider implements OSSProvider for testing.
type MockOSSProvider struct {
	Buckets     []OSSBucket
	Objects     []OSSObject
	IsTruncated bool
	Err         error
}

func (m *MockOSSProvider) ListBuckets() ([]OSSBucket, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Buckets, nil
}

func (m *MockOSSProvider) ListObjects(bucket, prefix string, maxKeys int32) ([]OSSObject, bool, error) {
	if m.Err != nil {
		return nil, false, m.Err
	}
	return m.Objects, m.IsTruncated, nil
}
```

- [ ] **Step 5: 验证编译**

```bash
go build ./provider/...
```

- [ ] **Step 6: 提交**

```bash
git add provider/mock_*.go
git commit -m "feat: add mock providers for testing"
```

---

## Task 3: 重构 Service 层使用接口

**Files:**
- Modify: `service/ecs.go`
- Modify: `service/cms.go`
- Modify: `service/sls.go`
- Modify: `service/oss.go`

- [ ] **Step 1: 修改 ECSService 使用接口**

在 `service/ecs.go` 中，将 provider 创建改为接口注入：

```go
package service

import (
	"cloud-manage/provider"
	"fmt"
	"sync"
)

// ECSService handles ECS business logic.
type ECSService struct {
	providerFactory func(accessKeyId, accessKeySecret, region string) (provider.ECSProvider, error)
}

// NewECSService creates a new ECSService with default provider factory.
func NewECSService() *ECSService {
	return &ECSService{
		providerFactory: defaultECSProviderFactory,
	}
}

// NewECSServiceWithProvider creates a new ECSService with custom provider factory (for testing).
func NewECSServiceWithProvider(factory func(string, string, string) (provider.ECSProvider, error)) *ECSService {
	return &ECSService{providerFactory: factory}
}

func defaultECSProviderFactory(accessKeyId, accessKeySecret, region string) (provider.ECSProvider, error) {
	// This will be updated to use actual aliyun provider
	return nil, fmt.Errorf("not implemented - use NewECSServiceWithProvider or set default factory")
}
```

- [ ] **Step 2: 更新 ListInstances 方法**

在 `service/ecs.go` 中更新：

```go
// ListInstances lists ECS instances for a single region.
func (s *ECSService) ListInstances(accessKeyId, accessKeySecret, region string) (*ListInstancesResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId, accessKeySecret and region are required")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ECS provider: %w", err)
	}

	instances, total, err := p.DescribeInstances(1, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	adapters := make([]ECSInstanceAdapter, 0, len(instances))
	for _, inst := range instances {
		adapters = append(adapters, ECSInstanceAdapter{
			InstanceId:   inst.InstanceId,
			InstanceName: inst.InstanceName,
			Status:       inst.Status,
			RegionId:     inst.RegionId,
			ZoneId:       inst.ZoneId,
			PublicIp:     inst.PublicIp,
			PrivateIp:    inst.PrivateIp,
			CreationTime: inst.CreationTime,
		})
	}

	return &ListInstancesResult{
		Instances:  adapters,
		TotalCount: total,
	}, nil
}
```

- [ ] **Step 3: 类似方式重构 CMS/SLS/OSS Service**

为每个 Service 添加 `providerFactory` 字段和 `NewXxxServiceWithProvider` 构造函数。

- [ ] **Step 4: 验证编译**

```bash
go build ./service/...
```

- [ ] **Step 5: 运行现有测试**

```bash
go test ./service/...
```

- [ ] **Step 6: 提交**

```bash
git add service/*.go
git commit -m "refactor: service layer uses interfaces for testability"
```

---

## Task 4: Service 层单元测试

**Files:**
- Create: `service/cms_test.go`
- Create: `service/sls_test.go`
- Create: `service/oss_test.go`
- Update: `service/ecs_test.go`

- [ ] **Step 1: 创建 CMS Service 测试**

```go
// service/cms_test.go
package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestGetSupportedProducts(t *testing.T) {
	svc := NewCMSService()
	products := svc.GetSupportedProducts()

	if len(products) == 0 {
		t.Error("expected non-empty products list")
	}

	// Verify ECS product exists
	found := false
	for _, p := range products {
		if p.ID == "ecs" {
			found = true
			if len(p.Metrics) == 0 {
				t.Error("expected ECS to have metrics")
			}
			break
		}
	}
	if !found {
		t.Error("expected ECS product in list")
	}
}

func TestGetProductMetrics_ValidProduct(t *testing.T) {
	svc := NewCMSService()
	product, err := svc.GetProductMetrics("ecs")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if product == nil {
		t.Fatal("expected non-nil product")
	}
	if product.ID != "ecs" {
		t.Errorf("expected product ID 'ecs', got '%s'", product.ID)
	}
}

func TestGetProductMetrics_InvalidProduct(t *testing.T) {
	svc := NewCMSService()
	_, err := svc.GetProductMetrics("invalid")

	if err == nil {
		t.Error("expected error for invalid product")
	}
}

func TestGetInstanceMetrics_MockProvider(t *testing.T) {
	cpuVal := 75.5
	mock := &provider.MockCMSProvider{
		Metrics: &provider.ECSMetricData{
			InstanceId:     "i-test",
			CPUUtilization: &cpuVal,
			UpdateTime:     "2026-05-19T00:00:00Z",
		},
	}

	svc := NewCMSServiceWithProvider(func(a, b, c string) (provider.CMSProvider, error) {
		return mock, nil
	})

	result, err := svc.GetInstanceMetrics("key", "secret", "cn-hangzhou", "i-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.InstanceId != "i-test" {
		t.Errorf("expected instanceId 'i-test', got '%s'", result.InstanceId)
	}
	if result.CPUUtilization == nil || *result.CPUUtilization != 75.5 {
		t.Error("expected CPU utilization 75.5")
	}
}

func TestGetInstanceMetrics_ProviderError(t *testing.T) {
	mock := &provider.MockCMSProvider{
		Err: errors.New("API error"),
	}

	svc := NewCMSServiceWithProvider(func(a, b, c string) (provider.CMSProvider, error) {
		return mock, nil
	})

	_, err := svc.GetInstanceMetrics("key", "secret", "cn-hangzhou", "i-test")
	if err == nil {
		t.Error("expected error from provider")
	}
}
```

- [ ] **Step 2: 创建 SLS Service 测试**

```go
// service/sls_test.go
package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListLogStores_MockProvider(t *testing.T) {
	mock := &provider.MockSLSProvider{
		LogStores: []string{"logstore1", "logstore2"},
	}

	svc := NewSLSServiceWithProvider(func(a, b, c string) (provider.SLSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListLogStores("key", "secret", "cn-hangzhou", "test-project")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 logstores, got %d", len(result))
	}
}

func TestQueryLogs_MockProvider(t *testing.T) {
	mock := &provider.MockSLSProvider{
		Logs: []provider.LogEntry{
			{Timestamp: 1000, Content: map[string]string{"level": "INFO"}},
			{Timestamp: 2000, Content: map[string]string{"level": "ERROR"}},
		},
		Count: 2,
	}

	svc := NewSLSServiceWithProvider(func(a, b, c string) (provider.SLSProvider, error) {
		return mock, nil
	})

	result, err := svc.QueryLogs("key", "secret", "cn-hangzhou", "project", "logstore", "", 0, 100, 100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}
}

func TestQueryLogs_EmptyCredentials(t *testing.T) {
	svc := NewSLSService()
	_, err := svc.QueryLogs("", "", "", "", "", "", 0, 0, 0)
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}
```

- [ ] **Step 3: 创建 OSS Service 测试**

```go
// service/oss_test.go
package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListBuckets_MockProvider(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Buckets: []provider.OSSBucket{
			{Name: "bucket1", Location: "oss-cn-hangzhou"},
			{Name: "bucket2", Location: "oss-cn-shenzhen"},
		},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListBuckets("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(result.Buckets))
	}
}

func TestListObjects_MockProvider(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Objects: []provider.OSSObject{
			{Key: "file1.txt", Size: 100},
			{Key: "file2.txt", Size: 200},
		},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListObjects("key", "secret", "cn-hangzhou", "test-bucket", "", 100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Objects) != 2 {
		t.Errorf("expected 2 objects, got %d", len(result.Objects))
	}
}

func TestDetectBucketRegion(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Buckets: []provider.OSSBucket{
			{Name: "test-bucket", Location: "oss-cn-shenzhen"},
		},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	region, err := svc.DetectBucketRegion("key", "secret", "test-bucket")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if region != "cn-shenzhen" {
		t.Errorf("expected region 'cn-shenzhen', got '%s'", region)
	}
}

func TestDetectBucketRegion_NotFound(t *testing.T) {
	mock := &provider.MockOSSProvider{
		Buckets: []provider.OSSBucket{},
	}

	svc := NewOSSServiceWithProvider(func(a, b, c string) (provider.OSSProvider, error) {
		return mock, nil
	})

	_, err := svc.DetectBucketRegion("key", "secret", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent bucket")
	}
}
```

- [ ] **Step 4: 更新 ECS 测试使用 Mock**

```go
// service/ecs_test.go
package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListInstances_MockProvider(t *testing.T) {
	mock := &provider.MockECSProvider{
		Instances: []provider.ECSInstance{
			{InstanceId: "i-001", InstanceName: "test-1", Status: "Running"},
			{InstanceId: "i-002", InstanceName: "test-2", Status: "Stopped"},
		},
		TotalCount: 2,
	}

	svc := NewECSServiceWithProvider(func(a, b, c string) (provider.ECSProvider, error) {
		return mock, nil
	})

	result, err := svc.ListInstances("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result.Instances) != 2 {
		t.Errorf("expected 2 instances, got %d", len(result.Instances))
	}
	if result.Instances[0].InstanceId != "i-001" {
		t.Errorf("expected first instance ID 'i-001', got '%s'", result.Instances[0].InstanceId)
	}
}

func TestListInstances_ProviderError(t *testing.T) {
	mock := &provider.MockECSProvider{
		Err: errors.New("API error"),
	}

	svc := NewECSServiceWithProvider(func(a, b, c string) (provider.ECSProvider, error) {
		return mock, nil
	})

	_, err := svc.ListInstances("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("expected error from provider")
	}
}

// Keep existing tests...
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
```

- [ ] **Step 5: 运行所有 Service 测试**

```bash
go test ./service/... -v
```

Expected: All tests PASS

- [ ] **Step 6: 提交**

```bash
git add service/*_test.go
git commit -m "test: add service layer unit tests with mocks"
```

---

## Task 5: Provider 层集成测试

**Files:**
- Create: `provider/aliyun/ecs_test.go`
- Create: `provider/aliyun/cms_test.go`
- Create: `provider/aliyun/sls_test.go`
- Create: `provider/aliyun/oss_test.go`

- [ ] **Step 1: 创建 ECS Provider 集成测试**

```go
// provider/aliyun/ecs_test.go
//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func getTestCredentials(t *testing.T) (string, string, string) {
	t.Helper()
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	region := os.Getenv("CLOUD_REGION")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: CLOUD_ACCESS_KEY_ID or CLOUD_ACCESS_KEY_SECRET not set")
	}
	if region == "" {
		region = "cn-hangzhou"
	}
	return accessKeyId, accessKeySecret, region
}

func TestECSProvider_Integration(t *testing.T) {
	accessKeyId, accessKeySecret, region := getTestCredentials(t)

	provider, err := NewECSProvider(accessKeyId, accessKeySecret, region)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	instances, total, err := provider.DescribeInstances(1, 10)
	if err != nil {
		t.Fatalf("DescribeInstances failed: %v", err)
	}

	t.Logf("Found %d instances (total: %d)", len(instances), total)
	for _, inst := range instances {
		t.Logf("  %s: %s (%s)", inst.InstanceId, inst.InstanceName, inst.Status)
	}
}
```

- [ ] **Step 2: 创建 CMS Provider 集成测试**

```go
// provider/aliyun/cms_test.go
//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func TestCMSProvider_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	instanceId := os.Getenv("TEST_INSTANCE_ID")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}
	if instanceId == "" {
		t.Skip("Skipping integration test: TEST_INSTANCE_ID not set")
	}

	provider, err := NewCMSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	metrics, err := provider.GetECSMetrics(instanceId)
	if err != nil {
		t.Fatalf("GetECSMetrics failed: %v", err)
	}

	t.Logf("Metrics for %s:", instanceId)
	if metrics.CPUUtilization != nil {
		t.Logf("  CPU: %.2f%%", *metrics.CPUUtilization)
	}
	if metrics.MemoryUtilization != nil {
		t.Logf("  Memory: %.2f%%", *metrics.MemoryUtilization)
	}
}
```

- [ ] **Step 3: 创建 SLS Provider 集成测试**

```go
// provider/aliyun/sls_test.go
//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func TestSLSProvider_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	project := os.Getenv("TEST_SLS_PROJECT")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}
	if project == "" {
		t.Skip("Skipping integration test: TEST_SLS_PROJECT not set")
	}

	provider, err := NewSLSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	logstores, err := provider.ListLogStores(project)
	if err != nil {
		t.Fatalf("ListLogStores failed: %v", err)
	}

	t.Logf("Found %d logstores in project %s:", len(logstores), project)
	for _, ls := range logstores {
		t.Logf("  - %s", ls)
	}
}
```

- [ ] **Step 4: 创建 OSS Provider 集成测试**

```go
// provider/aliyun/oss_test.go
//go:build integration

package aliyun

import (
	"os"
	"testing"
)

func TestOSSProvider_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}

	provider, err := NewOSSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	buckets, err := provider.ListBuckets()
	if err != nil {
		t.Fatalf("ListBuckets failed: %v", err)
	}

	t.Logf("Found %d buckets:", len(buckets))
	for _, b := range buckets {
		t.Logf("  - %s (%s)", b.Name, b.Location)
	}
}

func TestOSSProvider_ListObjects_Integration(t *testing.T) {
	accessKeyId := os.Getenv("CLOUD_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	bucket := os.Getenv("TEST_OSS_BUCKET")

	if accessKeyId == "" || accessKeySecret == "" {
		t.Skip("Skipping integration test: credentials not set")
	}
	if bucket == "" {
		t.Skip("Skipping integration test: TEST_OSS_BUCKET not set")
	}

	provider, err := NewOSSProvider(accessKeyId, accessKeySecret, "cn-hangzhou")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	objects, _, err := provider.ListObjects(bucket, "", 10)
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}

	t.Logf("Found %d objects in bucket %s:", len(objects), bucket)
	for _, obj := range objects {
		t.Logf("  - %s (%d bytes)", obj.Key, obj.Size)
	}
}
```

- [ ] **Step 5: 运行单元测试（排除集成测试）**

```bash
go test ./provider/aliyun/... -v
```

Expected: SKIP (因为没有 integration tag)

- [ ] **Step 6: 提交**

```bash
git add provider/aliyun/*_test.go
git commit -m "test: add provider layer integration tests with build tag"
```

---

## Task 6: CLI 端到端测试

**Files:**
- Create: `cmd/cli/cli_test.go`

- [ ] **Step 1: 创建 CLI 测试**

```go
// cmd/cli/cli_test.go
package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestCLI_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/cli/main.go", "help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !contains(outputStr, "Usage:") {
		t.Error("Expected help output to contain 'Usage:'")
	}
}

func TestCLI_Version(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/cli/main.go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !contains(outputStr, "v0.0.12") {
		t.Error("Expected version output to contain version number")
	}
}

func TestCLI_MissingCredentials(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("CLOUD_ACCESS_KEY_ID")
	os.Unsetenv("CLOUD_ACCESS_KEY_SECRET")

	cmd := exec.Command("go", "run", "../../cmd/cli/main.go", "ecs", "list")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error for missing credentials")
	}

	outputStr := string(output)
	if !contains(outputStr, "AccessKey") {
		t.Error("Expected error message about AccessKey")
	}
}

func TestCLI_CMSProducts(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/cli/main.go", "cms", "products")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error, got: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !contains(outputStr, "ecs") {
		t.Error("Expected products output to contain 'ecs'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: 运行 CLI 测试**

```bash
go test ./cmd/cli/... -v -run TestCLI_Help
```

Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add cmd/cli/cli_test.go
git commit -m "test: add CLI end-to-end tests"
```

---

## Task 7: 测试覆盖率报告

**Files:**
- Create: `scripts/test.sh` (update)

- [ ] **Step 1: 更新测试脚本**

```bash
#!/bin/bash
# scripts/test.sh - Run tests with coverage

set -e

echo "Running unit tests..."
go test ./... -v -coverprofile=coverage.out -covermode=atomic

echo ""
echo "Coverage report:"
go tool cover -func=coverage.out

echo ""
echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "Total coverage:"
go tool cover -func=coverage.out | tail -1
```

- [ ] **Step 2: 运行测试并查看覆盖率**

```bash
chmod +x scripts/test.sh
./scripts/test.sh
```

- [ ] **Step 3: 提交**

```bash
git add scripts/test.sh coverage.out coverage.html
git commit -m "test: add coverage reporting"
```

---

## 验收标准

1. **单元测试覆盖率** > 60%
   - Service 层: 100% 函数覆盖
   - Provider 接口: 100% Mock 实现

2. **集成测试** 可选运行
   - 使用 `go test -tags=integration` 运行
   - 需要环境变量配置凭证

3. **CI/CD 就绪**
   - 所有单元测试可在无凭证环境运行
   - 测试执行时间 < 30 秒

4. **文档完整**
   - 测试说明在 README.md
   - Mock 使用示例
