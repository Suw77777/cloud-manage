package provider

// ECSProvider defines the interface for ECS operations.
type ECSProvider interface {
	DescribeInstances(pageNumber, pageSize int32) ([]ECSInstance, int32, error)
	DescribeInstanceDetail(instanceId string) (*InstanceDetail, error)
	StartInstance(instanceId string) error
	StopInstance(instanceId string, forceStop bool) error
	RebootInstance(instanceId string, forceStop bool) error
}

// CMSProvider defines the interface for CloudMonitor operations.
type CMSProvider interface {
	GetECSMetrics(instanceId string) (*ECSMetricData, error)
}

// SLSProvider defines the interface for SLS operations.
type SLSProvider interface {
	ListLogStores(project string) ([]string, error)
	GetLogs(project, logstore, query string, from, to int64, maxLines int64) ([]LogEntry, int64, error)
}

// OSSProvider defines the interface for OSS operations.
type OSSProvider interface {
	ListBuckets() ([]OSSBucket, error)
	ListObjects(bucket, prefix string, maxKeys int32) ([]OSSObject, bool, error)
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

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp int64             `json:"timestamp"`
	Content   map[string]string `json:"content"`
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

// VPCProvider defines the interface for VPC operations.
type VPCProvider interface {
	ListVPCs() ([]VPC, error)
	GetVPCDetail(vpcId string) (*VPCDetail, error)
	ListVSwitches(vpcId string) ([]VSwitch, error)
}

// VPC represents a VPC network.
type VPC struct {
	VpcId        string `json:"vpcId"`
	VpcName      string `json:"vpcName"`
	CidrBlock    string `json:"cidrBlock"`
	Status       string `json:"status"`
	RegionId     string `json:"regionId"`
	Description  string `json:"description"`
	CreationTime string `json:"creationTime"`
}

// VPCDetail holds detailed information about a VPC.
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

// VSwitch represents a virtual switch in a VPC.
type VSwitch struct {
	VSwitchId    string `json:"vswitchId"`
	VSwitchName  string `json:"vswitchName"`
	CidrBlock    string `json:"cidrBlock"`
	ZoneId       string `json:"zoneId"`
	Status       string `json:"status"`
	VpcId        string `json:"vpcId"`
	CreationTime string `json:"creationTime"`
}

// SLBProvider defines the interface for SLB operations.
type SLBProvider interface {
	ListSLBs() ([]SLB, error)
	GetSLBDetail(slbId string) (*SLBDetail, error)
	ListSLBListeners(slbId string) ([]SLBListener, error)
}

// SLB represents a Server Load Balancer instance.
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

// SLBDetail holds detailed information about an SLB instance.
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

// SLBListener represents a listener on an SLB instance.
type SLBListener struct {
	ListenerPort     int    `json:"listenerPort"`
	ListenerProtocol string `json:"listenerProtocol"`
	Status           string `json:"status"`
	Bandwidth        int    `json:"bandwidth"`
}
