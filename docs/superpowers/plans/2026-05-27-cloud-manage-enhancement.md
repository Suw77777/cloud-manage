# Cloud Manage 增强实现计划

> **对于自动化工作代理:** 必须使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 来逐任务执行此计划。步骤使用 checkbox (`- [ ]`) 语法进行跟踪。

**目标:** 修复现有问题，添加 VPC/SLB 管理，增强现有功能，提升工程质量。

**架构:** 遵循已建立的三层模式：`provider/interfaces.go`（接口）→ `provider/aliyun/*.go`（SDK 实现）→ `service/*.go`（业务逻辑）→ `app.go`/`main.go`（CLI/GUI 暴露）。每个新服务（VPC、SLB）将遵循此模式。

**技术栈:** Go 1.26.3, 阿里云 SDK (vpc-20160428, slb-20140515), Bubble Tea (TUI), Wails (GUI), Vue 3 (前端)

---

## 文件结构

### 新建文件
- `provider/mock_vpc.go` - VPC Mock 测试桩
- `provider/mock_slb.go` - SLB Mock 测试桩
- `provider/aliyun/vpc.go` - VPC SDK 实现
- `provider/aliyun/slb.go` - SLB SDK 实现
- `provider/aliyun/vpc_test.go` - VPC Provider 测试
- `provider/aliyun/slb_test.go` - SLB Provider 测试
- `service/vpc.go` - VPC 服务层
- `service/slb.go` - SLB 服务层
- `service/vpc_test.go` - VPC 服务测试
- `service/slb_test.go` - SLB 服务测试
- `internal/tui/views/vpc.go` - TUI VPC 视图
- `internal/tui/views/slb.go` - TUI SLB 视图

### 修改文件
- `go.mod` - 添加 VPC 和 SLB SDK 依赖
- `provider/interfaces.go` - 添加 VPC/SLB 接口和类型
- `app.go` - 添加 VPC/SLB 方法供 GUI 使用
- `main.go` - 添加 VPC/SLB CLI 处理器，修复版本常量
- `cmd/cli/cli_test.go` - 修复失败的测试
- `internal/tui/app.go` - 添加 VPC/SLB 标签页
- `frontend/src/App.vue` - 添加 VPC/SLB 标签页
- `frontend/src/components/` - 添加 VPC/SLB 组件

---

## 任务 1: 修复失败的 CLI 测试

**文件:**
- 修改: `cmd/cli/cli_test.go`
- 修改: `main.go:23`（版本常量）

测试 `TestCLI_CMSProducts` 失败是因为 CLI 的 `cms` 处理器不支持 `products` 子命令。测试期望 `cms products` 列出支持的产品，但 `handleCMS()` 只处理 `metrics`。

- [ ] **步骤 1: 在 main.go 的 handleCMS 中添加 `products` 动作**

```go
// 在 main.go 中，修改 handleCMS 函数（第 365-390 行）：
func handleCMS(action string, args []string) {
	svc := service.NewCMSService()

	switch action {
	case "products":
		products := svc.GetSupportedProducts()
		if outputJSON {
			printJSON(products)
		} else {
			printProducts(products)
		}

	case "metrics":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "用法: cloud-manage cms metrics <实例ID>\n")
			os.Exit(1)
		}
		instanceId := args[0]
		result, err := svc.GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printMetrics(*result)
		}

	default:
		fmt.Println(`CMS 操作:
  products          列出支持的云产品
  metrics <id>      查询实例监控指标`)
	}
}
```

- [ ] **步骤 2: 添加 `printProducts` 辅助函数**

```go
// 在 main.go 的 printMetrics 函数之后添加：
func printProducts(products []service.CloudProduct) {
	fmt.Printf("\n支持的云产品:\n")
	for _, p := range products {
		fmt.Printf("\n  [%s] %s (namespace: %s)\n", p.ID, p.Name, p.Namespace)
		fmt.Println("  监控指标:")
		for _, m := range p.Metrics {
			fmt.Printf("    - %-30s %s (%s)\n", m.Name, m.Unit, m.Description)
		}
	}
}
```

- [ ] **步骤 3: 运行失败的测试验证通过**

运行: `cd /home/shanxi/cloud-manage && go test ./cmd/cli/ -v -run TestCLI_CMSProducts`
预期: PASS

- [ ] **步骤 4: 运行所有 CLI 测试**

运行: `cd /home/shanxi/cloud-manage && go test ./cmd/cli/ -v`
预期: 全部 3 个测试 PASS

- [ ] **步骤 5: 提交**

```bash
git add main.go cmd/cli/cli_test.go
git commit -m "修复: 添加 cms products 子命令以修复失败的 CLI 测试"
```

---

## 任务 2: 添加 VPC SDK 依赖

**文件:**
- 修改: `go.mod`
- 修改: `go.sum`（自动更新）

- [ ] **步骤 1: 添加 VPC SDK 依赖**

运行: `cd /home/shanxi/cloud-manage && go get github.com/alibabacloud-go/vpc-20160428/v3@latest`
预期: go.mod 更新了 vpc 依赖

- [ ] **步骤 2: 验证依赖已添加**

运行: `grep vpc /home/shanxi/cloud-manage/go.mod`
预期: 显示 vpc 依赖行

- [ ] **步骤 3: 运行 go mod tidy**

运行: `cd /home/shanxi/cloud-manage && go mod tidy`
预期: 正常退出

- [ ] **步骤 4: 提交**

```bash
git add go.mod go.sum
git commit -m "依赖: 添加阿里云 VPC SDK"
```

---

## 任务 3: 添加 SLB SDK 依赖

**文件:**
- 修改: `go.mod`
- 修改: `go.sum`（自动更新）

- [ ] **步骤 1: 添加 SLB SDK 依赖**

运行: `cd /home/shanxi/cloud-manage && go get github.com/alibabacloud-go/slb-20140515/v3@latest`
预期: go.mod 更新了 slb 依赖

- [ ] **步骤 2: 验证依赖已添加**

运行: `grep slb /home/shanxi/cloud-manage/go.mod`
预期: 显示 slb 依赖行（与现有的 sls 区分）

- [ ] **步骤 3: 运行 go mod tidy**

运行: `cd /home/shanxi/cloud-manage && go mod tidy`
预期: 正常退出

- [ ] **步骤 4: 提交**

```bash
git add go.mod go.sum
git commit -m "依赖: 添加阿里云 SLB SDK"
```

---

## 任务 4: 添加 VPC Provider 接口和类型

**文件:**
- 修改: `provider/interfaces.go`
- 新建: `provider/mock_vpc.go`

- [ ] **步骤 1: 在 interfaces.go 中添加 VPC 类型**

```go
// 在 provider/interfaces.go 的 OSSObject 结构体之后添加：

// VPCProvider 定义 VPC 操作接口。
type VPCProvider interface {
	ListVPCs() ([]VPC, error)
	GetVPCDetail(vpcId string) (*VPCDetail, error)
	ListVSwitches(vpcId string) ([]VSwitch, error)
}

// VPC 表示一个 VPC 网络。
type VPC struct {
	VpcId        string `json:"vpcId"`
	VpcName      string `json:"vpcName"`
	CidrBlock    string `json:"cidrBlock"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	Description  string `json:"description"`
	CreationTime string `json:"creationTime"`
}

// VPCDetail 保存 VPC 的详细信息。
type VPCDetail struct {
	VpcId          string   `json:"vpcId"`
	VpcName        string   `json:"vpcName"`
	CidrBlock      string   `json:"cidrBlock"`
	Status         string   `json:"status"`
	RegionId       string   `json:"regionId"`
	Description    string   `json:"description"`
	CreationTime   string   `json:"creationTime"`
	VSwitchIds     []string `json:"vswitchIds"`
	NatGatewayIds  []string `json:"natGatewayIds"`
	RouterTableIds []string `json:"routerTableIds"`
}

// VSwitch 表示 VPC 中的虚拟交换机。
type VSwitch struct {
	VSwitchId    string `json:"vswitchId"`
	VSwitchName  string `json:"vswitchName"`
	CidrBlock    string `json:"cidrBlock"`
	ZoneId       string `json:"zoneId"`
	Status       string `json:"status"`
	VpcId        string `json:"vpcId"`
	CreationTime string `json:"creationTime"`
}
```

- [ ] **步骤 2: 创建 mock_vpc.go**

```go
// provider/mock_vpc.go
package provider

// MockVPCProvider 实现 VPCProvider 用于测试。
type MockVPCProvider struct {
	VPCs      []VPC
	VPCDetail *VPCDetail
	VSwitches []VSwitch
	Err       error
}

func (m *MockVPCProvider) ListVPCs() ([]VPC, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VPCs, nil
}

func (m *MockVPCProvider) GetVPCDetail(vpcId string) (*VPCDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VPCDetail, nil
}

func (m *MockVPCProvider) ListVSwitches(vpcId string) ([]VSwitch, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.VSwitches, nil
}
```

- [ ] **步骤 3: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build ./provider/...`
预期: 正常退出

- [ ] **步骤 4: 提交**

```bash
git add provider/interfaces.go provider/mock_vpc.go
git commit -m "feat(provider): 添加 VPC 接口和 Mock"
```

---

## 任务 5: 添加 SLB Provider 接口和类型

**文件:**
- 修改: `provider/interfaces.go`
- 新建: `provider/mock_slb.go`

- [ ] **步骤 1: 在 interfaces.go 中添加 SLB 类型**

```go
// 在 provider/interfaces.go 的 VSwitch 结构体之后添加：

// SLBProvider 定义 SLB 操作接口。
type SLBProvider interface {
	ListSLBs() ([]SLB, error)
	GetSLBDetail(slbId string) (*SLBDetail, error)
	ListSLBListeners(slbId string) ([]SLBListener, error)
}

// SLB 表示一个负载均衡实例。
type SLB struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	CreationTime     string `json:"creationTime"`
}

// SLBDetail 保存 SLB 的详细信息。
type SLBDetail struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	VSwitchId        string `json:"vswitchId"`
	CreationTime     string `json:"creationTime"`
	ListenerCount    int    `json:"listenerCount"`
	Bandwidth        int    `json:"bandwidth"`
}

// SLBListener 表示 SLB 上的监听器。
type SLBListener struct {
	ListenerPort     int    `json:"listenerPort"`
	ListenerProtocol string `json:"listenerProtocol"`
	Status           string `json:"status"`
	Bandwidth        int    `json:"bandwidth"`
}
```

- [ ] **步骤 2: 创建 mock_slb.go**

```go
// provider/mock_slb.go
package provider

// MockSLBProvider 实现 SLBProvider 用于测试。
type MockSLBProvider struct {
	SLBs      []SLB
	SLBDetail *SLBDetail
	Listeners []SLBListener
	Err       error
}

func (m *MockSLBProvider) ListSLBs() ([]SLB, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.SLBs, nil
}

func (m *MockSLBProvider) GetSLBDetail(slbId string) (*SLBDetail, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.SLBDetail, nil
}

func (m *MockSLBProvider) ListSLBListeners(slbId string) ([]SLBListener, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Listeners, nil
}
```

- [ ] **步骤 3: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build ./provider/...`
预期: 正常退出

- [ ] **步骤 4: 提交**

```bash
git add provider/interfaces.go provider/mock_slb.go
git commit -m "feat(provider): 添加 SLB 接口和 Mock"
```

---

## 任务 6: 实现 VPC Provider（阿里云 SDK）

**文件:**
- 新建: `provider/aliyun/vpc.go`

- [ ] **步骤 1: 创建 provider/aliyun/vpc.go**

```go
package aliyun

import (
	"cloud-manage/provider"
	"fmt"

	vpc "github.com/alibabacloud-go/vpc-20160428/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// VPCProvider 封装阿里云 VPC SDK 客户端。
type VPCProvider struct {
	client *vpc.Client
}

// NewVPCProvider 创建新的 VPCProvider。
func NewVPCProvider(accessKeyId, accessKeySecret, region string) (*VPCProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("vpc.%s.aliyuncs.com", region))

	client, err := vpc.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建 VPC 客户端失败: %w", err)
	}
	return &VPCProvider{client: client}, nil
}

// ListVPCs 列出区域内的所有 VPC。
func (p *VPCProvider) ListVPCs() ([]provider.VPC, error) {
	request := &vpc.DescribeVpcsRequest{
		RegionId: p.client.RegionId,
	}

	response, err := p.client.DescribeVpcs(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeVpcs 失败: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("空响应体")
	}

	vpcs := make([]provider.VPC, 0)
	if body.Vpcs != nil && body.Vpcs.Vpc != nil {
		for _, v := range body.Vpcs.Vpc {
			vpcs = append(vpcs, provider.VPC{
				VpcId:        tea.StringValue(v.VpcId),
				VpcName:      tea.StringValue(v.VpcName),
				CidrBlock:    tea.StringValue(v.CidrBlock),
				Status:       tea.StringValue(v.Status),
				RegionId:     tea.StringValue(v.RegionId),
				Description:  tea.StringValue(v.Description),
				CreationTime: tea.StringValue(v.CreationTime),
			})
		}
	}

	return vpcs, nil
}

// GetVPCDetail 查询单个 VPC 的详细信息。
func (p *VPCProvider) GetVPCDetail(vpcId string) (*provider.VPCDetail, error) {
	request := &vpc.DescribeVpcAttributeRequest{
		VpcId: tea.String(vpcId),
	}

	response, err := p.client.DescribeVpcAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeVpcAttribute 失败: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("空响应体")
	}

	vswitchIds := make([]string, 0)
	if body.VSwitchIds != nil && body.VSwitchIds.VSwitchId != nil {
		for _, id := range body.VSwitchIds.VSwitchId {
			vswitchIds = append(vswitchIds, tea.StringValue(id))
		}
	}

	natGatewayIds := make([]string, 0)
	if body.NatGatewayIds != nil && body.NatGatewayIds.NatGatewayId != nil {
		for _, id := range body.NatGatewayIds.NatGatewayId {
			natGatewayIds = append(natGatewayIds, tea.StringValue(id))
		}
	}

	routerTableIds := make([]string, 0)
	if body.RouterTableIds != nil && body.RouterTableIds.RouterTableId != nil {
		for _, id := range body.RouterTableIds.RouterTableId {
			routerTableIds = append(routerTableIds, tea.StringValue(id))
		}
	}

	return &provider.VPCDetail{
		VpcId:          tea.StringValue(body.VpcId),
		VpcName:        tea.StringValue(body.VpcName),
		CidrBlock:      tea.StringValue(body.CidrBlock),
		Status:         tea.StringValue(body.Status),
		RegionId:       tea.StringValue(body.RegionId),
		Description:    tea.StringValue(body.Description),
		CreationTime:   tea.StringValue(body.CreationTime),
		VSwitchIds:     vswitchIds,
		NatGatewayIds:  natGatewayIds,
		RouterTableIds: routerTableIds,
	}, nil
}

// ListVSwitches 列出 VPC 中的所有虚拟交换机。
func (p *VPCProvider) ListVSwitches(vpcId string) ([]provider.VSwitch, error) {
	request := &vpc.DescribeVSwitchesRequest{
		VpcId:    tea.String(vpcId),
		RegionId: p.client.RegionId,
	}

	response, err := p.client.DescribeVSwitches(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeVSwitches 失败: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("空响应体")
	}

	vswitches := make([]provider.VSwitch, 0)
	if body.VSwitchs != nil {
		for _, vs := range body.VSwitchs {
			vswitches = append(vswitches, provider.VSwitch{
				VSwitchId:    tea.StringValue(vs.VSwitchId),
				VSwitchName:  tea.StringValue(vs.VSwitchName),
				CidrBlock:    tea.StringValue(vs.CidrBlock),
				ZoneId:       tea.StringValue(vs.ZoneId),
				Status:       tea.StringValue(vs.Status),
				VpcId:        tea.StringValue(vs.VpcId),
				CreationTime: tea.StringValue(vs.CreationTime),
			})
		}
	}

	return vswitches, nil
}
```

- [ ] **步骤 2: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build ./provider/aliyun/...`
预期: 正常退出

- [ ] **步骤 3: 提交**

```bash
git add provider/aliyun/vpc.go
git commit -m "feat(provider/aliyun): 实现 VPC Provider"
```

---

## 任务 7: 实现 SLB Provider（阿里云 SDK）

**文件:**
- 新建: `provider/aliyun/slb.go`

- [ ] **步骤 1: 创建 provider/aliyun/slb.go**

```go
package aliyun

import (
	"cloud-manage/provider"
	"fmt"

	slb "github.com/alibabacloud-go/slb-20140515/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// SLBProvider 封装阿里云 SLB SDK 客户端。
type SLBProvider struct {
	client *slb.Client
}

// NewSLBProvider 创建新的 SLBProvider。
func NewSLBProvider(accessKeyId, accessKeySecret, region string) (*SLBProvider, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	config.Endpoint = tea.String(fmt.Sprintf("slb.%s.aliyuncs.com", region))

	client, err := slb.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建 SLB 客户端失败: %w", err)
	}
	return &SLBProvider{client: client}, nil
}

// ListSLBs 列出区域内的所有 SLB 实例。
func (p *SLBProvider) ListSLBs() ([]provider.SLB, error) {
	region := tea.StringValue(p.client.RegionId)
	request := &slb.DescribeLoadBalancersRequest{
		RegionId: tea.String(region),
	}

	response, err := p.client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeLoadBalancers 失败: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("空响应体")
	}

	slbs := make([]provider.SLB, 0)
	if body.LoadBalancers != nil && body.LoadBalancers.LoadBalancer != nil {
		for _, lb := range body.LoadBalancers.LoadBalancer {
			slbs = append(slbs, provider.SLB{
				LoadBalancerId:   tea.StringValue(lb.LoadBalancerId),
				LoadBalancerName: tea.StringValue(lb.LoadBalancerName),
				Address:          tea.StringValue(lb.Address),
				AddressType:      tea.StringValue(lb.AddressType),
				Status:           tea.StringValue(lb.Status),
				RegionId:         tea.StringValue(lb.RegionId),
				VpcId:            tea.StringValue(lb.VpcId),
				CreationTime:     tea.StringValue(lb.CreateTime),
			})
		}
	}

	return slbs, nil
}

// GetSLBDetail 查询单个 SLB 实例的详细信息。
func (p *SLBProvider) GetSLBDetail(slbId string) (*provider.SLBDetail, error) {
	request := &slb.DescribeLoadBalancerAttributeRequest{
		LoadBalancerId: tea.String(slbId),
	}

	response, err := p.client.DescribeLoadBalancerAttribute(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeLoadBalancerAttribute 失败: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("空响应体")
	}

	return &provider.SLBDetail{
		LoadBalancerId:   tea.StringValue(body.LoadBalancerId),
		LoadBalancerName: tea.StringValue(body.LoadBalancerName),
		Address:          tea.StringValue(body.Address),
		AddressType:      tea.StringValue(body.AddressType),
		Status:           tea.StringValue(body.Status),
		RegionId:         tea.StringValue(body.RegionId),
		VpcId:            tea.StringValue(body.VpcId),
		VSwitchId:        tea.StringValue(body.VSwitchId),
		CreationTime:     tea.StringValue(body.CreateTime),
		ListenerCount:    int(tea.Int32Value(body.ListenerPorts.ListenerPort)),
		Bandwidth:        int(tea.Int32Value(body.Bandwidth)),
	}, nil
}

// ListSLBListeners 列出 SLB 实例上的所有监听器。
func (p *SLBProvider) ListSLBListeners(slbId string) ([]provider.SLBListener, error) {
	request := &slb.DescribeLoadBalancerListenersRequest{
		LoadBalancerId: tea.String(slbId),
	}

	response, err := p.client.DescribeLoadBalancerListeners(request)
	if err != nil {
		return nil, fmt.Errorf("DescribeLoadBalancerListeners 失败: %w", err)
	}

	body := response.Body
	if body == nil {
		return nil, fmt.Errorf("空响应体")
	}

	listeners := make([]provider.SLBListener, 0)
	if body.Listeners != nil {
		for _, l := range body.Listeners {
			listeners = append(listeners, provider.SLBListener{
				ListenerPort:     int(tea.Int32Value(l.ListenerPort)),
				ListenerProtocol: tea.StringValue(l.ListenerProtocol),
				Status:           tea.StringValue(l.Status),
				Bandwidth:        int(tea.Int32Value(l.Bandwidth)),
			})
		}
	}

	return listeners, nil
}
```

- [ ] **步骤 2: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build ./provider/aliyun/...`
预期: 正常退出

- [ ] **步骤 3: 提交**

```bash
git add provider/aliyun/slb.go
git commit -m "feat(provider/aliyun): 实现 SLB Provider"
```

---

## 任务 8: 实现 VPC 服务层

**文件:**
- 新建: `service/vpc.go`
- 新建: `service/vpc_test.go`

- [ ] **步骤 1: 创建 service/vpc.go**

```go
package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// VPCProviderFactory 根据凭证和区域创建 VPCProvider。
type VPCProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.VPCProvider, error)

// VPCService 处理 VPC 业务逻辑。
type VPCService struct {
	providerFactory VPCProviderFactory
}

// NewVPCService 创建使用默认 Provider 工厂的 VPCService。
func NewVPCService() *VPCService {
	return &VPCService{
		providerFactory: func(accessKeyId, accessKeySecret, region string) (provider.VPCProvider, error) {
			return aliyun.NewVPCProvider(accessKeyId, accessKeySecret, region)
		},
	}
}

// NewVPCServiceWithProvider 创建使用自定义 Provider 工厂的 VPCService（用于测试）。
func NewVPCServiceWithProvider(factory VPCProviderFactory) *VPCService {
	return &VPCService{providerFactory: factory}
}

// VPCAdapter 是与 Provider 无关的 VPC 表示。
type VPCAdapter struct {
	VpcId        string `json:"vpcId"`
	VpcName      string `json:"vpcName"`
	CidrBlock    string `json:"cidrBlock"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	Description  string `json:"description"`
	CreationTime string `json:"creationTime"`
}

// VPCDetailAdapter 是与 Provider 无关的 VPC 详情表示。
type VPCDetailAdapter struct {
	VpcId          string   `json:"vpcId"`
	VpcName        string   `json:"vpcName"`
	CidrBlock      string   `json:"cidrBlock"`
	Status         string   `json:"status"`
	RegionId       string   `json:"regionId"`
	Description    string   `json:"description"`
	CreationTime   string   `json:"creationTime"`
	VSwitchIds     []string `json:"vswitchIds"`
	NatGatewayIds  []string `json:"natGatewayIds"`
	RouterTableIds []string `json:"routerTableIds"`
}

// VSwitchAdapter 是与 Provider 无关的 VSwitch 表示。
type VSwitchAdapter struct {
	VSwitchId    string `json:"vswitchId"`
	VSwitchName  string `json:"vswitchName"`
	CidrBlock    string `json:"cidrBlock"`
	ZoneId       string `json:"zoneId"`
	Status       string `json:"status"`
	VpcId        string `json:"vpcId"`
	CreationTime string `json:"creationTime"`
}

// ListVPCsResult 保存列出 VPC 的结果。
type ListVPCsResult struct {
	VPCs []VPCAdapter `json:"vpcs"`
}

// ListVSwitchesResult 保存列出 VSwitch 的结果。
type ListVSwitchesResult struct {
	VSwitches []VSwitchAdapter `json:"vswitches"`
}

// ListVPCs 列出区域内的所有 VPC。
func (s *VPCService) ListVPCs(accessKeyId, accessKeySecret, region string) (*ListVPCsResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret 和 region 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 VPC Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	vpcs, err := p.ListVPCs()
	if err != nil {
		return nil, fmt.Errorf("列出 VPC 失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]VPCAdapter, 0, len(vpcs))
	for _, v := range vpcs {
		adapters = append(adapters, VPCAdapter{
			VpcId:        v.VpcId,
			VpcName:      v.VpcName,
			CidrBlock:    v.CidrBlock,
			Status:       v.Status,
			RegionId:     v.RegionId,
			Description:  v.Description,
			CreationTime: v.CreationTime,
		})
	}

	return &ListVPCsResult{VPCs: adapters}, nil
}

// GetVPCDetail 查询单个 VPC 的详细信息。
func (s *VPCService) GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId string) (*VPCDetailAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || vpcId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 vpcId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 VPC Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	detail, err := p.GetVPCDetail(vpcId)
	if err != nil {
		return nil, fmt.Errorf("获取 VPC 详情失败: %s", security.SanitizeErrorMessage(err))
	}

	return &VPCDetailAdapter{
		VpcId:          detail.VpcId,
		VpcName:        detail.VpcName,
		CidrBlock:      detail.CidrBlock,
		Status:         detail.Status,
		RegionId:       detail.RegionId,
		Description:    detail.Description,
		CreationTime:   detail.CreationTime,
		VSwitchIds:     detail.VSwitchIds,
		NatGatewayIds:  detail.NatGatewayIds,
		RouterTableIds: detail.RouterTableIds,
	}, nil
}

// ListVSwitches 列出 VPC 中的所有虚拟交换机。
func (s *VPCService) ListVSwitches(accessKeyId, accessKeySecret, region, vpcId string) (*ListVSwitchesResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || vpcId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 vpcId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 VPC Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	vswitches, err := p.ListVSwitches(vpcId)
	if err != nil {
		return nil, fmt.Errorf("列出 VSwitch 失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]VSwitchAdapter, 0, len(vswitches))
	for _, vs := range vswitches {
		adapters = append(adapters, VSwitchAdapter{
			VSwitchId:    vs.VSwitchId,
			VSwitchName:  vs.VSwitchName,
			CidrBlock:    vs.CidrBlock,
			ZoneId:       vs.ZoneId,
			Status:       vs.Status,
			VpcId:        vs.VpcId,
			CreationTime: vs.CreationTime,
		})
	}

	return &ListVSwitchesResult{VSwitches: adapters}, nil
}
```

- [ ] **步骤 2: 创建 service/vpc_test.go**

```go
package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListVPCs_MockProvider(t *testing.T) {
	mock := &provider.MockVPCProvider{
		VPCs: []provider.VPC{
			{VpcId: "vpc-001", VpcName: "test-vpc-1", CidrBlock: "10.0.0.0/16", Status: "Available"},
			{VpcId: "vpc-002", VpcName: "test-vpc-2", CidrBlock: "172.16.0.0/12", Status: "Available"},
		},
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	result, err := svc.ListVPCs("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("意外错误: %v", err)
	}
	if len(result.VPCs) != 2 {
		t.Errorf("期望 2 个 VPC，实际 %d 个", len(result.VPCs))
	}
	if result.VPCs[0].VpcId != "vpc-001" {
		t.Errorf("期望第一个 VPC ID 为 'vpc-001'，实际为 '%s'", result.VPCs[0].VpcId)
	}
}

func TestListVPCs_ProviderError(t *testing.T) {
	mock := &provider.MockVPCProvider{
		Err: errors.New("API 错误"),
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	_, err := svc.ListVPCs("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("期望 Provider 返回错误")
	}
}

func TestListVPCs_EmptyCredentials(t *testing.T) {
	svc := NewVPCService()

	_, err := svc.ListVPCs("", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("期望空 accessKeyId 时报错")
	}
}

func TestNewVPCService(t *testing.T) {
	svc := NewVPCService()
	if svc == nil {
		t.Error("期望非空 VPCService")
	}
}

func TestListVSwitches_MockProvider(t *testing.T) {
	mock := &provider.MockVPCProvider{
		VSwitches: []provider.VSwitch{
			{VSwitchId: "vsw-001", VSwitchName: "test-vsw-1", CidrBlock: "10.0.1.0/24", ZoneId: "cn-hangzhou-a"},
			{VSwitchId: "vsw-002", VSwitchName: "test-vsw-2", CidrBlock: "10.0.2.0/24", ZoneId: "cn-hangzhou-b"},
		},
	}

	svc := NewVPCServiceWithProvider(func(a, b, c string) (provider.VPCProvider, error) {
		return mock, nil
	})

	result, err := svc.ListVSwitches("key", "secret", "cn-hangzhou", "vpc-001")
	if err != nil {
		t.Errorf("意外错误: %v", err)
	}
	if len(result.VSwitches) != 2 {
		t.Errorf("期望 2 个 VSwitch，实际 %d 个", len(result.VSwitches))
	}
}
```

- [ ] **步骤 3: 运行 VPC 服务测试**

运行: `cd /home/shanxi/cloud-manage && go test ./service/ -v -run TestVPC -run TestListVPC -run TestNewVPC -run TestListVSwitch`
预期: 所有 VPC 测试 PASS

- [ ] **步骤 4: 提交**

```bash
git add service/vpc.go service/vpc_test.go
git commit -m "feat(service): 实现 VPC 服务层及测试"
```

---

## 任务 9: 实现 SLB 服务层

**文件:**
- 新建: `service/slb.go`
- 新建: `service/slb_test.go`

- [ ] **步骤 1: 创建 service/slb.go**

```go
package service

import (
	"cloud-manage/provider"
	"cloud-manage/provider/aliyun"
	"cloud-manage/security"
	"fmt"
)

// SLBProviderFactory 根据凭证和区域创建 SLBProvider。
type SLBProviderFactory func(accessKeyId, accessKeySecret, region string) (provider.SLBProvider, error)

// SLBService 处理 SLB 业务逻辑。
type SLBService struct {
	providerFactory SLBProviderFactory
}

// NewSLBService 创建使用默认 Provider 工厂的 SLBService。
func NewSLBService() *SLBService {
	return &SLBService{
		providerFactory: func(accessKeyId, accessKeySecret, region string) (provider.SLBProvider, error) {
			return aliyun.NewSLBProvider(accessKeyId, accessKeySecret, region)
		},
	}
}

// NewSLBServiceWithProvider 创建使用自定义 Provider 工厂的 SLBService（用于测试）。
func NewSLBServiceWithProvider(factory SLBProviderFactory) *SLBService {
	return &SLBService{providerFactory: factory}
}

// SLBAdapter 是与 Provider 无关的 SLB 表示。
type SLBAdapter struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	CreationTime     string `json:"creationTime"`
}

// SLBDetailAdapter 是与 Provider 无关的 SLB 详情表示。
type SLBDetailAdapter struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	VSwitchId        string `json:"vswitchId"`
	CreationTime     string `json:"creationTime"`
	ListenerCount    int    `json:"listenerCount"`
	Bandwidth        int    `json:"bandwidth"`
}

// SLBListenerAdapter 是与 Provider 无关的 SLB 监听器表示。
type SLBListenerAdapter struct {
	ListenerPort     int    `json:"listenerPort"`
	ListenerProtocol string `json:"listenerProtocol"`
	Status           string `json:"status"`
	Bandwidth        int    `json:"bandwidth"`
}

// ListSLBsResult 保存列出 SLB 的结果。
type ListSLBsResult struct {
	SLBs []SLBAdapter `json:"slbs"`
}

// ListSLBListenersResult 保存列出 SLB 监听器的结果。
type ListSLBListenersResult struct {
	Listeners []SLBListenerAdapter `json:"listeners"`
}

// ListSLBs 列出区域内的所有 SLB 实例。
func (s *SLBService) ListSLBs(accessKeyId, accessKeySecret, region string) (*ListSLBsResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret 和 region 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 SLB Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	slbs, err := p.ListSLBs()
	if err != nil {
		return nil, fmt.Errorf("列出 SLB 失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]SLBAdapter, 0, len(slbs))
	for _, lb := range slbs {
		adapters = append(adapters, SLBAdapter{
			LoadBalancerId:   lb.LoadBalancerId,
			LoadBalancerName: lb.LoadBalancerName,
			Address:          lb.Address,
			AddressType:      lb.AddressType,
			Status:           lb.Status,
			RegionId:         lb.RegionId,
			VpcId:            lb.VpcId,
			CreationTime:     lb.CreationTime,
		})
	}

	return &ListSLBsResult{SLBs: adapters}, nil
}

// GetSLBDetail 查询单个 SLB 实例的详细信息。
func (s *SLBService) GetSLBDetail(accessKeyId, accessKeySecret, region, slbId string) (*SLBDetailAdapter, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || slbId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 slbId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 SLB Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	detail, err := p.GetSLBDetail(slbId)
	if err != nil {
		return nil, fmt.Errorf("获取 SLB 详情失败: %s", security.SanitizeErrorMessage(err))
	}

	return &SLBDetailAdapter{
		LoadBalancerId:   detail.LoadBalancerId,
		LoadBalancerName: detail.LoadBalancerName,
		Address:          detail.Address,
		AddressType:      detail.AddressType,
		Status:           detail.Status,
		RegionId:         detail.RegionId,
		VpcId:            detail.VpcId,
		VSwitchId:        detail.VSwitchId,
		CreationTime:     detail.CreationTime,
		ListenerCount:    detail.ListenerCount,
		Bandwidth:        detail.Bandwidth,
	}, nil
}

// ListSLBListeners 列出 SLB 实例上的所有监听器。
func (s *SLBService) ListSLBListeners(accessKeyId, accessKeySecret, region, slbId string) (*ListSLBListenersResult, error) {
	if accessKeyId == "" || accessKeySecret == "" || region == "" || slbId == "" {
		return nil, fmt.Errorf("accessKeyId、accessKeySecret、region 和 slbId 为必填项")
	}

	p, err := s.providerFactory(accessKeyId, accessKeySecret, region)
	if err != nil {
		return nil, fmt.Errorf("初始化 SLB Provider 失败: %s", security.SanitizeErrorMessage(err))
	}

	listeners, err := p.ListSLBListeners(slbId)
	if err != nil {
		return nil, fmt.Errorf("列出 SLB 监听器失败: %s", security.SanitizeErrorMessage(err))
	}

	adapters := make([]SLBListenerAdapter, 0, len(listeners))
	for _, l := range listeners {
		adapters = append(adapters, SLBListenerAdapter{
			ListenerPort:     l.ListenerPort,
			ListenerProtocol: l.ListenerProtocol,
			Status:           l.Status,
			Bandwidth:        l.Bandwidth,
		})
	}

	return &ListSLBListenersResult{Listeners: adapters}, nil
}
```

- [ ] **步骤 2: 创建 service/slb_test.go**

```go
package service

import (
	"cloud-manage/provider"
	"errors"
	"testing"
)

func TestListSLBs_MockProvider(t *testing.T) {
	mock := &provider.MockSLBProvider{
		SLBs: []provider.SLB{
			{LoadBalancerId: "lb-001", LoadBalancerName: "test-slb-1", Address: "1.2.3.4", Status: "active"},
			{LoadBalancerId: "lb-002", LoadBalancerName: "test-slb-2", Address: "5.6.7.8", Status: "active"},
		},
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	result, err := svc.ListSLBs("key", "secret", "cn-hangzhou")
	if err != nil {
		t.Errorf("意外错误: %v", err)
	}
	if len(result.SLBs) != 2 {
		t.Errorf("期望 2 个 SLB，实际 %d 个", len(result.SLBs))
	}
	if result.SLBs[0].LoadBalancerId != "lb-001" {
		t.Errorf("期望第一个 SLB ID 为 'lb-001'，实际为 '%s'", result.SLBs[0].LoadBalancerId)
	}
}

func TestListSLBs_ProviderError(t *testing.T) {
	mock := &provider.MockSLBProvider{
		Err: errors.New("API 错误"),
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	_, err := svc.ListSLBs("key", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("期望 Provider 返回错误")
	}
}

func TestListSLBs_EmptyCredentials(t *testing.T) {
	svc := NewSLBService()

	_, err := svc.ListSLBs("", "secret", "cn-hangzhou")
	if err == nil {
		t.Error("期望空 accessKeyId 时报错")
	}
}

func TestNewSLBService(t *testing.T) {
	svc := NewSLBService()
	if svc == nil {
		t.Error("期望非空 SLBService")
	}
}

func TestListSLBListeners_MockProvider(t *testing.T) {
	mock := &provider.MockSLBProvider{
		Listeners: []provider.SLBListener{
			{ListenerPort: 80, ListenerProtocol: "HTTP", Status: "running"},
			{ListenerPort: 443, ListenerProtocol: "HTTPS", Status: "running"},
		},
	}

	svc := NewSLBServiceWithProvider(func(a, b, c string) (provider.SLBProvider, error) {
		return mock, nil
	})

	result, err := svc.ListSLBListeners("key", "secret", "cn-hangzhou", "lb-001")
	if err != nil {
		t.Errorf("意外错误: %v", err)
	}
	if len(result.Listeners) != 2 {
		t.Errorf("期望 2 个监听器，实际 %d 个", len(result.Listeners))
	}
}
```

- [ ] **步骤 3: 运行 SLB 服务测试**

运行: `cd /home/shanxi/cloud-manage && go test ./service/ -v -run TestSLB -run TestListSLB -run TestNewSLB`
预期: 所有 SLB 测试 PASS

- [ ] **步骤 4: 提交**

```bash
git add service/slb.go service/slb_test.go
git commit -m "feat(service): 实现 SLB 服务层及测试"
```

---

## 任务 10: 添加 VPC/SLB 到应用层（GUI 支持）

**文件:**
- 修改: `app.go`

- [ ] **步骤 1: 在 App 结构体中添加 VPC/SLB 服务**

```go
// 在 app.go 中，修改 App 结构体（第 18-24 行）：
type App struct {
	ctx     context.Context
	ecsSvc  *service.ECSService
	cmsSvc  *service.CMSService
	slsSvc  *service.SLSService
	ossSvc  *service.OSSService
	vpcSvc  *service.VPCService
	slbSvc  *service.SLBService
}
```

- [ ] **步骤 2: 更新 NewApp 构造函数**

```go
// 在 app.go 中，修改 NewApp（第 27-34 行）：
func NewApp() *App {
	return &App{
		ecsSvc: service.NewECSService(),
		cmsSvc: service.NewCMSService(),
		slsSvc: service.NewSLSService(),
		ossSvc: service.NewOSSService(),
		vpcSvc: service.NewVPCService(),
		slbSvc: service.NewSLBService(),
	}
}
```

- [ ] **步骤 3: 添加 VPC 视图类型和方法**

```go
// 在 app.go 的 OSSObjectResult 结构体之后添加：

// VPCView 是前端友好的 VPC 表示。
type VPCView struct {
	VpcId        string `json:"vpcId"`
	VpcName      string `json:"vpcName"`
	CidrBlock    string `json:"cidrBlock"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	Description  string `json:"description"`
	CreationTime string `json:"creationTime"`
}

// VPCDetailResult 是 VPC 详情查询的响应类型。
type VPCDetailResult struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Detail  *VPCDetailView `json:"detail,omitempty"`
}

// VPCDetailView 是前端友好的 VPC 详情表示。
type VPCDetailView struct {
	VpcId          string   `json:"vpcId"`
	VpcName        string   `json:"vpcName"`
	CidrBlock      string   `json:"cidrBlock"`
	Status         string   `json:"status"`
	RegionId       string   `json:"regionId"`
	Description    string   `json:"description"`
	CreationTime   string   `json:"creationTime"`
	VSwitchIds     []string `json:"vswitchIds"`
	NatGatewayIds  []string `json:"natGatewayIds"`
	RouterTableIds []string `json:"routerTableIds"`
}

// VSwitchView 是前端友好的 VSwitch 表示。
type VSwitchView struct {
	VSwitchId    string `json:"vswitchId"`
	VSwitchName  string `json:"vswitchName"`
	CidrBlock    string `json:"cidrBlock"`
	ZoneId       string `json:"zoneId"`
	Status       string `json:"status"`
	VpcId        string `json:"vpcId"`
	CreationTime string `json:"creationTime"`
}

// VPCListResult 是 VPC 列表查询的响应类型。
type VPCListResult struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	VPCs    []VPCView  `json:"vpcs"`
}

// VSwitchListResult 是 VSwitch 列表查询的响应类型。
type VSwitchListResult struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	VSwitches []VSwitchView `json:"vswitches"`
}

// ListVPCs 列出区域内的所有 VPC。
func (a *App) ListVPCs(accessKeyId, accessKeySecret, region string) VPCListResult {
	result, err := a.vpcSvc.ListVPCs(accessKeyId, accessKeySecret, region)
	if err != nil {
		return VPCListResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	vpcs := make([]VPCView, 0, len(result.VPCs))
	for _, v := range result.VPCs {
		vpcs = append(vpcs, VPCView{
			VpcId:        v.VpcId,
			VpcName:      v.VpcName,
			CidrBlock:    v.CidrBlock,
			Status:       v.Status,
			RegionId:     v.RegionId,
			Description:  v.Description,
			CreationTime: v.CreationTime,
		})
	}

	return VPCListResult{
		Success: true,
		Message: fmt.Sprintf("找到 %d 个 VPC", len(vpcs)),
		VPCs:    vpcs,
	}
}

// GetVPCDetail 查询单个 VPC 的详细信息。
func (a *App) GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId string) VPCDetailResult {
	detail, err := a.vpcSvc.GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId)
	if err != nil {
		return VPCDetailResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	return VPCDetailResult{
		Success: true,
		Message: "VPC 详情获取成功",
		Detail: &VPCDetailView{
			VpcId:          detail.VpcId,
			VpcName:        detail.VpcName,
			CidrBlock:      detail.CidrBlock,
			Status:         detail.Status,
			RegionId:       detail.RegionId,
			Description:    detail.Description,
			CreationTime:   detail.CreationTime,
			VSwitchIds:     detail.VSwitchIds,
			NatGatewayIds:  detail.NatGatewayIds,
			RouterTableIds: detail.RouterTableIds,
		},
	}
}

// ListVSwitches 列出 VPC 中的所有虚拟交换机。
func (a *App) ListVSwitches(accessKeyId, accessKeySecret, region, vpcId string) VSwitchListResult {
	result, err := a.vpcSvc.ListVSwitches(accessKeyId, accessKeySecret, region, vpcId)
	if err != nil {
		return VSwitchListResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	vswitches := make([]VSwitchView, 0, len(result.VSwitches))
	for _, vs := range result.VSwitches {
		vswitches = append(vswitches, VSwitchView{
			VSwitchId:    vs.VSwitchId,
			VSwitchName:  vs.VSwitchName,
			CidrBlock:    vs.CidrBlock,
			ZoneId:       vs.ZoneId,
			Status:       vs.Status,
			VpcId:        vs.VpcId,
			CreationTime: vs.CreationTime,
		})
	}

	return VSwitchListResult{
		Success:   true,
		Message:   fmt.Sprintf("找到 %d 个 VSwitch", len(vswitches)),
		VSwitches: vswitches,
	}
}
```

- [ ] **步骤 4: 添加 SLB 视图类型和方法**

```go
// 在 app.go 的 VSwitchListResult 之后添加：

// SLBView 是前端友好的 SLB 表示。
type SLBView struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	CreationTime     string `json:"creationTime"`
}

// SLBDetailResult 是 SLB 详情查询的响应类型。
type SLBDetailResult struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Detail  *SLBDetailView  `json:"detail,omitempty"`
}

// SLBDetailView 是前端友好的 SLB 详情表示。
type SLBDetailView struct {
	LoadBalancerId   string `json:"loadBalancerId"`
	LoadBalancerName string `json:"loadBalancerName"`
	Address          string `json:"address"`
	AddressType      string `json:"addressType"`
	Status           string `json:"status"`
	RegionId         string `json:"regionId"`
	VpcId            string `json:"vpcId"`
	VSwitchId        string `json:"vswitchId"`
	CreationTime     string `json:"creationTime"`
	ListenerCount    int    `json:"listenerCount"`
	Bandwidth        int    `json:"bandwidth"`
}

// SLBListenerView 是前端友好的 SLB 监听器表示。
type SLBListenerView struct {
	ListenerPort     int    `json:"listenerPort"`
	ListenerProtocol string `json:"listenerProtocol"`
	Status           string `json:"status"`
	Bandwidth        int    `json:"bandwidth"`
}

// SLBListResult 是 SLB 列表查询的响应类型。
type SLBListResult struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	SLBs    []SLBView `json:"slbs"`
}

// SLBListenerListResult 是 SLB 监听器列表查询的响应类型。
type SLBListenerListResult struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	Listeners []SLBListenerView `json:"listeners"`
}

// ListSLBs 列出区域内的所有 SLB 实例。
func (a *App) ListSLBs(accessKeyId, accessKeySecret, region string) SLBListResult {
	result, err := a.slbSvc.ListSLBs(accessKeyId, accessKeySecret, region)
	if err != nil {
		return SLBListResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	slbs := make([]SLBView, 0, len(result.SLBs))
	for _, lb := range result.SLBs {
		slbs = append(slbs, SLBView{
			LoadBalancerId:   lb.LoadBalancerId,
			LoadBalancerName: lb.LoadBalancerName,
			Address:          lb.Address,
			AddressType:      lb.AddressType,
			Status:           lb.Status,
			RegionId:         lb.RegionId,
			VpcId:            lb.VpcId,
			CreationTime:     lb.CreationTime,
		})
	}

	return SLBListResult{
		Success: true,
		Message: fmt.Sprintf("找到 %d 个 SLB", len(slbs)),
		SLBs:    slbs,
	}
}

// GetSLBDetail 查询单个 SLB 实例的详细信息。
func (a *App) GetSLBDetail(accessKeyId, accessKeySecret, region, slbId string) SLBDetailResult {
	detail, err := a.slbSvc.GetSLBDetail(accessKeyId, accessKeySecret, region, slbId)
	if err != nil {
		return SLBDetailResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	return SLBDetailResult{
		Success: true,
		Message: "SLB 详情获取成功",
		Detail: &SLBDetailView{
			LoadBalancerId:   detail.LoadBalancerId,
			LoadBalancerName: detail.LoadBalancerName,
			Address:          detail.Address,
			AddressType:      detail.AddressType,
			Status:           detail.Status,
			RegionId:         detail.RegionId,
			VpcId:            detail.VpcId,
			VSwitchId:        detail.VSwitchId,
			CreationTime:     detail.CreationTime,
			ListenerCount:    detail.ListenerCount,
			Bandwidth:        detail.Bandwidth,
		},
	}
}

// ListSLBListeners 列出 SLB 实例上的所有监听器。
func (a *App) ListSLBListeners(accessKeyId, accessKeySecret, region, slbId string) SLBListenerListResult {
	result, err := a.slbSvc.ListSLBListeners(accessKeyId, accessKeySecret, region, slbId)
	if err != nil {
		return SLBListenerListResult{
			Success: false,
			Message: security.SanitizeErrorMessage(err),
		}
	}

	listeners := make([]SLBListenerView, 0, len(result.Listeners))
	for _, l := range result.Listeners {
		listeners = append(listeners, SLBListenerView{
			ListenerPort:     l.ListenerPort,
			ListenerProtocol: l.ListenerProtocol,
			Status:           l.Status,
			Bandwidth:        l.Bandwidth,
		})
	}

	return SLBListenerListResult{
		Success:   true,
		Message:   fmt.Sprintf("找到 %d 个监听器", len(listeners)),
		Listeners: listeners,
	}
}
```

- [ ] **步骤 5: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build .`
预期: 正常退出

- [ ] **步骤 6: 提交**

```bash
git add app.go
git commit -m "feat(app): 添加 VPC 和 SLB GUI 支持"
```

---

## 任务 11: 添加 VPC/SLB CLI 命令

**文件:**
- 修改: `main.go`

- [ ] **步骤 1: 在 knownServices 中添加 vpc 和 slb**

```go
// 在 main.go 中，修改 knownServices（第 46-53 行）：
var knownServices = map[string]bool{
	"ecs":     true,
	"cms":     true,
	"sls":     true,
	"oss":     true,
	"vpc":     true,
	"slb":     true,
	"help":    true,
	"version": true,
}
```

- [ ] **步骤 2: 在 runCLI switch 中添加 VPC/SLB case**

```go
// 在 main.go 中，在 handleOSS case 之后添加（第 179 行）：
	case "vpc":
		handleVPC(action, remainingArgs)
	case "slb":
		handleSLB(action, remainingArgs)
```

- [ ] **步骤 3: 添加 handleVPC 函数**

```go
// 在 main.go 的 handleOSS 函数之后添加：

// ========== VPC ==========

func handleVPC(action string, args []string) {
	svc := service.NewVPCService()

	switch action {
	case "list":
		regions := []string{region}
		if region == "all" {
			regions = getAllRegions()
		}
		for _, rgn := range regions {
			result, err := svc.ListVPCs(accessKeyId, accessKeySecret, rgn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[%s] 错误: %v\n", rgn, err)
				continue
			}
			if outputJSON {
				printJSON(result)
			} else {
				printVPCs(result.VPCs, rgn)
			}
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "用法: cloud-manage vpc detail <vpc-id>\n")
			os.Exit(1)
		}
		vpcId := args[0]
		result, err := svc.GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printVPCDetail(result)
		}

	case "vswitches":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "用法: cloud-manage vpc vswitches <vpc-id>\n")
			os.Exit(1)
		}
		vpcId := args[0]
		result, err := svc.ListVSwitches(accessKeyId, accessKeySecret, region, vpcId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printVSwitches(result.VSwitches, vpcId)
		}

	default:
		fmt.Println(`VPC 操作:
  list                    列出 VPC
  detail <vpc-id>         查看 VPC 详情
  vswitches <vpc-id>      列出 VSwitch`)
	}
}
```

- [ ] **步骤 4: 添加 handleSLB 函数**

```go
// 在 handleVPC 函数之后添加：

// ========== SLB ==========

func handleSLB(action string, args []string) {
	svc := service.NewSLBService()

	switch action {
	case "list":
		regions := []string{region}
		if region == "all" {
			regions = getAllRegions()
		}
		for _, rgn := range regions {
			result, err := svc.ListSLBs(accessKeyId, accessKeySecret, rgn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[%s] 错误: %v\n", rgn, err)
				continue
			}
			if outputJSON {
				printJSON(result)
			} else {
				printSLBs(result.SLBs, rgn)
			}
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "用法: cloud-manage slb detail <slb-id>\n")
			os.Exit(1)
		}
		slbId := args[0]
		result, err := svc.GetSLBDetail(accessKeyId, accessKeySecret, region, slbId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printSLBDetail(result)
		}

	case "listeners":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "用法: cloud-manage slb listeners <slb-id>\n")
			os.Exit(1)
		}
		slbId := args[0]
		result, err := svc.ListSLBListeners(accessKeyId, accessKeySecret, region, slbId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printSLBListeners(result.Listeners, slbId)
		}

	default:
		fmt.Println(`SLB 操作:
  list                      列出 SLB 实例
  detail <slb-id>           查看 SLB 详情
  listeners <slb-id>        列出监听端口`)
	}
}
```

- [ ] **步骤 5: 添加 VPC/SLB 打印辅助函数**

```go
// 在 main.go 的 printObjects 函数之后添加：

func printVPCs(vpcs []service.VPCAdapter, region string) {
	fmt.Printf("\n=== %s 的 VPC ===\n", region)
	if len(vpcs) == 0 {
		fmt.Println("未找到 VPC")
		return
	}
	for _, v := range vpcs {
		fmt.Printf("  [%s] %-20s CIDR: %-18s %s\n", v.Status, v.VpcId, v.CidrBlock, v.VpcName)
	}
}

func printVPCDetail(d *service.VPCDetailAdapter) {
	fmt.Printf(`
VPC 详情:
  ID:               %s
  名称:             %s
  CIDR:             %s
  状态:             %s
  区域:             %s
  描述:             %s
  创建时间:         %s
  VSwitch:          %v
  NAT 网关:         %v
  路由表:           %v
`,
		d.VpcId, d.VpcName, d.CidrBlock, d.Status, d.RegionId,
		d.Description, d.CreationTime, d.VSwitchIds, d.NatGatewayIds, d.RouterTableIds)
}

func printVSwitches(vswitches []service.VSwitchAdapter, vpcId string) {
	fmt.Printf("\n=== %s 的 VSwitch ===\n", vpcId)
	if len(vswitches) == 0 {
		fmt.Println("未找到 VSwitch")
		return
	}
	for _, vs := range vswitches {
		fmt.Printf("  [%s] %-20s CIDR: %-18s 可用区: %s\n", vs.Status, vs.VSwitchId, vs.CidrBlock, vs.ZoneId)
	}
}

func printSLBs(slbs []service.SLBAdapter, region string) {
	fmt.Printf("\n=== %s 的 SLB ===\n", region)
	if len(slbs) == 0 {
		fmt.Println("未找到 SLB")
		return
	}
	for _, lb := range slbs {
		fmt.Printf("  [%s] %-20s %-15s %-10s %s\n", lb.Status, lb.LoadBalancerId, lb.Address, lb.AddressType, lb.LoadBalancerName)
	}
}

func printSLBDetail(d *service.SLBDetailAdapter) {
	fmt.Printf(`
SLB 详情:
  ID:               %s
  名称:             %s
  地址:             %s
  类型:             %s
  状态:             %s
  区域:             %s
  VPC:              %s
  VSwitch:          %s
  监听器数:         %d
  带宽:             %d Mbps
  创建时间:         %s
`,
		d.LoadBalancerId, d.LoadBalancerName, d.Address, d.AddressType,
		d.Status, d.RegionId, d.VpcId, d.VSwitchId,
		d.ListenerCount, d.Bandwidth, d.CreationTime)
}

func printSLBListeners(listeners []service.SLBListenerAdapter, slbId string) {
	fmt.Printf("\n=== %s 的监听器 ===\n", slbId)
	if len(listeners) == 0 {
		fmt.Println("未找到监听器")
		return
	}
	for _, l := range listeners {
		fmt.Printf("  [%s] %s:%d (带宽: %d)\n", l.Status, l.ListenerProtocol, l.ListenerPort, l.Bandwidth)
	}
}
```

- [ ] **步骤 6: 更新帮助文本**

```go
// 在 printUsage 函数中，在 oss 服务之后添加 vpc 和 slb：
	vpc               VPC 网络管理
	slb               负载均衡管理
```

- [ ] **步骤 7: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build .`
预期: 正常退出

- [ ] **步骤 8: 提交**

```bash
git add main.go
git commit -m "feat(cli): 添加 VPC 和 SLB CLI 命令"
```

---

## 任务 12: 添加 VPC/SLB 到 TUI

**文件:**
- 新建: `internal/tui/views/vpc.go`
- 新建: `internal/tui/views/slb.go`
- 修改: `internal/tui/app.go`

- [ ] **步骤 1: 创建 internal/tui/views/vpc.go**

```go
package views

import (
	"cloud-manage/service"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type VPCMsg struct {
	VPCs []service.VPCAdapter
	Err  error
}

type VPCModel struct {
	table   table.Model
	vpcs    []service.VPCAdapter
	loading bool
	err     error
}

func NewVPCModel() VPCModel {
	columns := []table.Column{
		{Title: "VPC ID", Width: 20},
		{Title: "名称", Width: 20},
		{Title: "CIDR", Width: 18},
		{Title: "状态", Width: 10},
		{Title: "区域", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true).Foreground(lipgloss.Color("6"))
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	t.SetStyles(s)

	return VPCModel{
		table: t,
	}
}

func (m VPCModel) Init() tea.Cmd {
	return nil
}

func (m VPCModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case VPCMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.vpcs = msg.VPCs
		m.updateTable()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			m.err = nil
			return m, func() tea.Msg {
				return VPCMsg{}
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *VPCModel) updateTable() {
	rows := make([]table.Row, 0, len(m.vpcs))
	for _, v := range m.vpcs {
		rows = append(rows, table.Row{
			v.VpcId,
			v.VpcName,
			v.CidrBlock,
			v.Status,
			v.RegionId,
		})
	}
	m.table.SetRows(rows)
}

func (m VPCModel) View() string {
	if m.loading {
		return "正在加载 VPC..."
	}
	if m.err != nil {
		return fmt.Sprintf("错误: %v", m.err)
	}
	if len(m.vpcs) == 0 {
		return "未找到 VPC。按 'r' 刷新。"
	}
	return m.table.View()
}

func (m VPCModel) Title() string {
	return "VPC 网络"
}

func (m VPCModel) HelpText() string {
	return strings.Join([]string{
		"r: 刷新",
		"↑/↓: 导航",
		"q: 退出",
	}, " | ")
}
```

- [ ] **步骤 2: 创建 internal/tui/views/slb.go**

```go
package views

import (
	"cloud-manage/service"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SLBMsg struct {
	SLBs []service.SLBAdapter
	Err  error
}

type SLBModel struct {
	table   table.Model
	slbs    []service.SLBAdapter
	loading bool
	err     error
}

func NewSLBModel() SLBModel {
	columns := []table.Column{
		{Title: "SLB ID", Width: 20},
		{Title: "名称", Width: 20},
		{Title: "地址", Width: 15},
		{Title: "类型", Width: 10},
		{Title: "状态", Width: 10},
		{Title: "VPC", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true).Foreground(lipgloss.Color("6"))
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	t.SetStyles(s)

	return SLBModel{
		table: t,
	}
}

func (m SLBModel) Init() tea.Cmd {
	return nil
}

func (m SLBModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SLBMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.slbs = msg.SLBs
		m.updateTable()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			m.err = nil
			return m, func() tea.Msg {
				return SLBMsg{}
			}
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *SLBModel) updateTable() {
	rows := make([]table.Row, 0, len(m.slbs))
	for _, lb := range m.slbs {
		rows = append(rows, table.Row{
			lb.LoadBalancerId,
			lb.LoadBalancerName,
			lb.Address,
			lb.AddressType,
			lb.Status,
			lb.VpcId,
		})
	}
	m.table.SetRows(rows)
}

func (m SLBModel) View() string {
	if m.loading {
		return "正在加载 SLB..."
	}
	if m.err != nil {
		return fmt.Sprintf("错误: %v", m.err)
	}
	if len(m.slbs) == 0 {
		return "未找到 SLB。按 'r' 刷新。"
	}
	return m.table.View()
}

func (m SLBModel) Title() string {
	return "负载均衡"
}

func (m SLBModel) HelpText() string {
	return strings.Join([]string{
		"r: 刷新",
		"↑/↓: 导航",
		"q: 退出",
	}, " | ")
}
```

- [ ] **步骤 3: 更新 internal/tui/app.go 以包含 VPC/SLB 标签页**

```go
// 首先读取 internal/tui/app.go，然后将 VPC 和 SLB 模型添加到 App 结构体中
// 并更新标签页处理以包含新视图。
// 具体更改取决于 app.go 的当前结构。
```

- [ ] **步骤 4: 验证编译**

运行: `cd /home/shanxi/cloud-manage && go build ./internal/tui/...`
预期: 正常退出

- [ ] **步骤 5: 提交**

```bash
git add internal/tui/views/vpc.go internal/tui/views/slb.go internal/tui/app.go
git commit -m "feat(tui): 添加 VPC 和 SLB 视图"
```

---

## 任务 13: 运行所有测试并验证

**文件:**
- 无（仅验证）

- [ ] **步骤 1: 运行所有 Go 测试**

运行: `cd /home/shanxi/cloud-manage && go test ./... -v 2>&1 | tail -50`
预期: 所有测试 PASS

- [ ] **步骤 2: 运行 go vet**

运行: `cd /home/shanxi/cloud-manage && go vet ./...`
预期: 正常退出（无问题）

- [ ] **步骤 3: 构建二进制文件**

运行: `cd /home/shanxi/cloud-manage && go build -o cloud-manage .`
预期: 二进制文件创建成功

- [ ] **步骤 4: 测试 VPC CLI 帮助**

运行: `cd /home/shanxi/cloud-manage && ./cloud-manage vpc`
预期: 显示 VPC 操作帮助文本

- [ ] **步骤 5: 测试 SLB CLI 帮助**

运行: `cd /home/shanxi/cloud-manage && ./cloud-manage slb`
预期: 显示 SLB 操作帮助文本

- [ ] **步骤 6: 测试 CMS products 命令**

运行: `cd /home/shanxi/cloud-manage && ./cloud-manage cms products`
预期: 显示支持的云产品列表

- [ ] **步骤 7: 最终提交并更新版本号**

```bash
git add -A
git commit -m "feat: v0.1.0 - 添加 VPC/SLB 管理，修复测试，增强功能"
```

---

## 总结

此计划涵盖：
1. **修复现有问题**: 添加缺失的 `cms products` CLI 子命令
2. **添加 VPC 管理**: 从 Provider 到 CLI/TUI/GUI 的完整实现
3. **添加 SLB 管理**: 从 Provider 到 CLI/TUI/GUI 的完整实现
4. **工程质量**: 所有更改包含测试，遵循已建立的模式

总计：13 个任务，约 50 个步骤
