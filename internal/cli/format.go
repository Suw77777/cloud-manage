package cli

import (
	"cloud-manage/service"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// PrintJSON marshals a value to JSON and prints it.
func PrintJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

// PrintInstances prints a list of ECS instances for a region.
func PrintInstances(instances []service.ECSInstanceAdapter, region string) {
	fmt.Printf("\n=== %s ===\n", region)
	if len(instances) == 0 {
		fmt.Println("No instances found")
		return
	}
	for _, inst := range instances {
		status := "●"
		if inst.Status != "Running" {
			status = "○"
		}
		fmt.Printf("  %s [%s] %-20s %-15s | Public: %-15s Private: %-15s\n",
			status, inst.Status, inst.InstanceId, inst.InstanceName, inst.PublicIp, inst.PrivateIp)
	}
}

// PrintInstanceDetail prints detailed information about an ECS instance.
func PrintInstanceDetail(d *service.InstanceDetailAdapter) {
	fmt.Printf(`
Instance Detail:
  ID:               %s
  Name:             %s
  Description:      %s
  Hostname:         %s
  Status:           %s
  Region:           %s
  Zone:             %s
  Type:             %s
  CPU:              %d cores
  Memory:           %d MB
  Image:            %s
  Charge Type:      %s
  Created:          %s
  Expired:          %s
  Stopped Mode:     %s
  Public IPs:       %s
  Private IPs:      %s
  Security Groups:  %s
`,
		d.InstanceId, d.InstanceName, d.Description, d.HostName,
		d.Status, d.RegionId, d.ZoneId, d.InstanceType,
		d.Cpu, d.Memory, d.ImageId, d.InternetChargeType,
		d.CreationTime, d.ExpiredTime, d.StoppedMode,
		strings.Join(d.PublicIp, ", "), strings.Join(d.PrivateIp, ", "),
		strings.Join(d.SecurityGroupIds, ", "))
}

// PrintMetrics prints ECS monitoring metrics.
func PrintMetrics(m service.ECSMetricAdapter) {
	fmt.Printf("\nMetrics for %s:\n", m.InstanceId)
	if m.CPUUtilization != nil {
		fmt.Printf("  CPU Utilization:     %.2f%%\n", *m.CPUUtilization)
	}
	if m.MemoryUtilization != nil {
		fmt.Printf("  Memory Utilization:  %.2f%%\n", *m.MemoryUtilization)
	}
	if m.DiskReadBPS != nil {
		fmt.Printf("  Disk Read BPS:       %.2f B/s\n", *m.DiskReadBPS)
	}
	if m.DiskWriteBPS != nil {
		fmt.Printf("  Disk Write BPS:      %.2f B/s\n", *m.DiskWriteBPS)
	}
	if m.InternetRX != nil {
		fmt.Printf("  Internet RX:         %.2f bps\n", *m.InternetRX)
	}
	if m.InternetTX != nil {
		fmt.Printf("  Internet TX:         %.2f bps\n", *m.InternetTX)
	}
	fmt.Printf("  Update Time:         %s\n", m.UpdateTime)
}

// PrintProducts prints supported cloud products and their metrics.
func PrintProducts(products []service.CloudProduct) {
	fmt.Printf("\n支持的云产品:\n")
	for _, p := range products {
		fmt.Printf("\n  [%s] %s (namespace: %s)\n", p.ID, p.Name, p.Namespace)
		fmt.Println("  监控指标:")
		for _, m := range p.Metrics {
			fmt.Printf("    - %-30s %s (%s)\n", m.Name, m.Unit, m.Description)
		}
	}
}

// PrintLogs prints SLS log query results.
func PrintLogs(result *service.LogQueryResult) {
	fmt.Printf("\nFound %d logs (hasMore: %v):\n", result.Count, result.HasMore)
	for _, entry := range result.Entries {
		ts := time.Unix(entry.Timestamp/1000, 0)
		fmt.Printf("\n[%s]\n", ts.Format("2006-01-02 15:04:05"))
		for k, v := range entry.Content {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
}

// PrintBuckets prints a list of OSS buckets.
func PrintBuckets(buckets []service.BucketAdapter) {
	if len(buckets) == 0 {
		fmt.Println("No buckets found")
		return
	}
	fmt.Printf("\nBuckets:\n")
	for _, b := range buckets {
		fmt.Printf("  %-30s %-15s %s\n", b.Name, b.Location, b.CreationDate)
	}
}

// PrintObjects prints objects in an OSS bucket.
func PrintObjects(objects []service.ObjectAdapter, bucket string) {
	fmt.Printf("\nObjects in %s:\n", bucket)
	if len(objects) == 0 {
		fmt.Println("No objects found")
		return
	}
	for _, obj := range objects {
		size := FormatSize(obj.Size)
		fmt.Printf("  %-50s %10s %s\n", obj.Key, size, obj.LastModified)
	}
}

// FormatSize formats bytes into a human-readable string.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// PrintVPCs prints a list of VPCs.
func PrintVPCs(vpcs []service.VPCAdapter) {
	if len(vpcs) == 0 {
		fmt.Println("No VPCs found")
		return
	}
	fmt.Printf("\nVPCs:\n")
	for _, v := range vpcs {
		fmt.Printf("  %-20s %-20s %-18s %s\n", v.VpcId, v.VpcName, v.CidrBlock, v.Status)
	}
}

// PrintVPCDetail prints detailed information about a VPC.
func PrintVPCDetail(d *service.VPCDetailAdapter) {
	fmt.Printf(`
VPC Detail:
  ID:               %s
  Name:             %s
  CIDR:             %s
  Status:           %s
  Region:           %s
  Description:      %s
  Created:          %s
  VSwitch IDs:      %s
`,
		d.VpcId, d.VpcName, d.CidrBlock, d.Status, d.RegionId,
		d.Description, d.CreationTime, strings.Join(d.VSwitchIds, ", "))
}

// PrintVSwitches prints a list of VSwitches.
func PrintVSwitches(vswitches []service.VSwitchAdapter) {
	if len(vswitches) == 0 {
		fmt.Println("No VSwitches found")
		return
	}
	fmt.Printf("\nVSwitches:\n")
	for _, vs := range vswitches {
		fmt.Printf("  %-20s %-20s %-18s %s\n", vs.VSwitchId, vs.VSwitchName, vs.CidrBlock, vs.ZoneId)
	}
}

// PrintSLBs prints a list of SLB instances.
func PrintSLBs(slbs []service.SLBAdapter) {
	if len(slbs) == 0 {
		fmt.Println("No SLBs found")
		return
	}
	fmt.Printf("\nSLBs:\n")
	for _, lb := range slbs {
		fmt.Printf("  %-20s %-20s %-15s %-10s %s\n", lb.LoadBalancerId, lb.LoadBalancerName, lb.Address, lb.AddressType, lb.Status)
	}
}

// PrintSLBDetail prints detailed information about an SLB instance.
func PrintSLBDetail(d *service.SLBDetailAdapter) {
	fmt.Printf(`
SLB Detail:
  ID:               %s
  Name:             %s
  Address:          %s
  Address Type:     %s
  Status:           %s
  Region:           %s
  VPC ID:           %s
  VSwitch ID:       %s
  Created:          %s
  Listeners:        %d
  Bandwidth:        %d Mbps
`,
		d.LoadBalancerId, d.LoadBalancerName, d.Address, d.AddressType,
		d.Status, d.RegionId, d.VpcId, d.VSwitchId,
		d.CreationTime, d.ListenerCount, d.Bandwidth)
}

// PrintSLBListeners prints listeners of an SLB instance.
func PrintSLBListeners(listeners []service.SLBListenerAdapter) {
	if len(listeners) == 0 {
		fmt.Println("No listeners found")
		return
	}
	fmt.Printf("\nSLB Listeners:\n")
	for _, l := range listeners {
		fmt.Printf("  Port: %-6d Protocol: %-8s Status: %-10s Bandwidth: %d Mbps\n",
			l.ListenerPort, l.ListenerProtocol, l.Status, l.Bandwidth)
	}
}
