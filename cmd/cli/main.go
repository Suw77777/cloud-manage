package main

import (
	"cloud-manage/service"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const version = "v0.0.11"

var (
	accessKeyId     string
	accessKeySecret string
	region          string
	outputJSON      bool
)

func init() {
	flag.StringVar(&accessKeyId, "id", "", "AccessKey ID (or env CLOUD_ACCESS_KEY_ID)")
	flag.StringVar(&accessKeySecret, "secret", "", "AccessKey Secret (or env CLOUD_ACCESS_KEY_SECRET)")
	flag.StringVar(&region, "region", "cn-hangzhou", "Region ID (default: cn-hangzhou)")
	flag.BoolVar(&outputJSON, "json", false, "Output in JSON format")
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	fmt.Printf("\n  Cloud 管理小助手 CLI %s\n\n", version)

	// Get credentials from flags or environment
	if accessKeyId == "" {
		accessKeyId = os.Getenv("CLOUD_ACCESS_KEY_ID")
	}
	if accessKeySecret == "" {
		accessKeySecret = os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	}

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	serviceName := args[0]
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	// Validate credentials (not needed for help/version)
	if serviceName != "help" && serviceName != "version" {
		if action != "" && (accessKeyId == "" || accessKeySecret == "") {
			fmt.Fprintf(os.Stderr, "Error: AccessKey ID and Secret are required.\n")
			fmt.Fprintf(os.Stderr, "Use -id/-secret flags or set CLOUD_ACCESS_KEY_ID/CLOUD_ACCESS_KEY_SECRET environment variables.\n")
			os.Exit(1)
		}
	}

	remainingArgs := []string{}
	if len(args) > 2 {
		remainingArgs = args[2:]
	}

	switch serviceName {
	case "help":
		printUsage()
	case "version":
		fmt.Printf("  cloud-cli %s\n", version)
	case "ecs":
		handleECS(action, remainingArgs)
	case "cms":
		handleCMS(action, remainingArgs)
	case "sls":
		handleSLS(action, remainingArgs)
	case "oss":
		handleOSS(action, remainingArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown service: %s\n", serviceName)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage:
  cloud-cli [flags] <service> <action> [args...]

Services:
  ecs               ECS 实例管理
  cms               云监控指标查询
  sls               日志服务查询
  oss               对象存储管理

Commands:
  help              显示帮助信息
  version           显示版本号

Flags:`)
	flag.PrintDefaults()
	fmt.Println(`
Examples:
  # List ECS instances
  cloud-cli -id LTAI4xxx -secret xxx ecs list

  # View ECS instance detail
  cloud-cli ecs detail i-xxx

  # Start ECS instance
  cloud-cli ecs start i-xxx

  # Query monitoring metrics
  cloud-cli cms metrics i-xxx

  # List OSS buckets
  cloud-cli oss buckets

  # List OSS objects
  cloud-cli oss objects my-bucket --prefix logs/

  # Query SLS logs
  cloud-cli sls logs my-project my-logstore --query "level: ERROR"

Environment Variables:
  CLOUD_ACCESS_KEY_ID      AccessKey ID
  CLOUD_ACCESS_KEY_SECRET  AccessKey Secret`)
}

// ========== ECS ==========

func handleECS(action string, args []string) {
	svc := service.NewECSService()

	switch action {
	case "list":
		regions := []string{region}
		if region == "all" {
			regions = getAllRegions()
		}
		if len(regions) == 1 {
			result, err := svc.ListInstances(accessKeyId, accessKeySecret, regions[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if outputJSON {
				printJSON(result)
			} else {
				printInstances(result.Instances, regions[0])
			}
		} else {
			results := svc.ListInstancesMultiRegion(accessKeyId, accessKeySecret, regions)
			if outputJSON {
				printJSON(results)
			} else {
				for _, r := range results {
					if r.Error != "" {
						fmt.Printf("\n[%s] Error: %s\n", r.Region, r.Error)
					} else {
						printInstances(r.Instances, r.Region)
					}
				}
			}
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli ecs detail <instance-id>\n")
			os.Exit(1)
		}
		instanceId := args[0]
		result, err := svc.GetInstanceDetail(accessKeyId, accessKeySecret, region, instanceId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printInstanceDetail(result)
		}

	case "start":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli ecs start <instance-id>\n")
			os.Exit(1)
		}
		instanceId := args[0]
		err := svc.StartInstance(accessKeyId, accessKeySecret, region, instanceId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance %s start command sent successfully\n", instanceId)

	case "stop":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli ecs stop <instance-id> [--force]\n")
			os.Exit(1)
		}
		instanceId := args[0]
		force := false
		for _, arg := range args[1:] {
			if arg == "--force" || arg == "-f" {
				force = true
			}
		}
		err := svc.StopInstance(accessKeyId, accessKeySecret, region, instanceId, force)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance %s stop command sent successfully\n", instanceId)

	case "reboot":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli ecs reboot <instance-id> [--force]\n")
			os.Exit(1)
		}
		instanceId := args[0]
		force := false
		for _, arg := range args[1:] {
			if arg == "--force" || arg == "-f" {
				force = true
			}
		}
		err := svc.RebootInstance(accessKeyId, accessKeySecret, region, instanceId, force)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance %s reboot command sent successfully\n", instanceId)

	default:
		fmt.Println(`ECS Actions:
  list              列出实例
  detail <id>       查看实例详情
  start <id>        启动实例
  stop <id>         停止实例 [--force]
  reboot <id>       重启实例 [--force]`)
	}
}

// ========== CMS ==========

func handleCMS(action string, args []string) {
	svc := service.NewCMSService()

	switch action {
	case "metrics":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli cms metrics <instance-id>\n")
			os.Exit(1)
		}
		instanceId := args[0]
		result, err := svc.GetInstanceMetrics(accessKeyId, accessKeySecret, region, instanceId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printMetrics(*result)
		}

	default:
		fmt.Println(`CMS Actions:
  metrics <id>      查询实例监控指标`)
	}
}

// ========== SLS ==========

func handleSLS(action string, args []string) {
	svc := service.NewSLSService()

	switch action {
	case "logstores":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli sls logstores <project>\n")
			os.Exit(1)
		}
		project := args[0]
		logstores, err := svc.ListLogStores(accessKeyId, accessKeySecret, region, project)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(logstores)
		} else {
			fmt.Printf("LogStores in %s:\n", project)
			for _, ls := range logstores {
				fmt.Printf("  - %s\n", ls)
			}
		}

	case "logs":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli sls logs <project> <logstore> [--query <query>] [--from <timestamp>] [--to <timestamp>] [--max <lines>]\n")
			os.Exit(1)
		}
		project := args[0]
		logstore := args[1]

		query := ""
		from := time.Now().Add(-1 * time.Hour).Unix()
		to := time.Now().Unix()
		maxLines := int64(100)

		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--query", "-q":
				if i+1 < len(args) {
					query = args[i+1]
					i++
				}
			case "--from":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &from)
					i++
				}
			case "--to":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &to)
					i++
				}
			case "--max":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &maxLines)
					i++
				}
			}
		}

		result, err := svc.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, maxLines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printLogs(result)
		}

	default:
		fmt.Println(`SLS Actions:
  logstores <project>                    列出 Logstore
  logs <project> <logstore> [options]    查询日志
    --query, -q <query>    查询表达式
    --from <timestamp>     开始时间 (Unix timestamp)
    --to <timestamp>       结束时间 (Unix timestamp)
    --max <lines>          最大返回行数 (default: 100)`)
	}
}

// ========== OSS ==========

func handleOSS(action string, args []string) {
	svc := service.NewOSSService()

	switch action {
	case "buckets":
		result, err := svc.ListBuckets(accessKeyId, accessKeySecret, region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printBuckets(result.Buckets)
		}

	case "objects":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-cli oss objects <bucket> [--prefix <prefix>] [--max <count>]\n")
			os.Exit(1)
		}
		// Clean bucket name: remove trailing :// or other URL artifacts
		bucket := strings.TrimSuffix(args[0], "://")
		bucket = strings.TrimSuffix(bucket, "/")
		prefix := ""
		maxKeys := int32(100)

		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--prefix", "-p":
				if i+1 < len(args) {
					prefix = args[i+1]
					i++
				}
			case "--max":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &maxKeys)
					i++
				}
			}
		}

		result, err := svc.ListObjects(accessKeyId, accessKeySecret, region, bucket, prefix, maxKeys)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printObjects(result.Objects, bucket)
		}

	default:
		fmt.Println(`OSS Actions:
  buckets                           列出 Bucket
  objects <bucket> [options]        列出对象
    --prefix, -p <prefix>   前缀过滤
    --max <count>           最大返回数量 (default: 100)`)
	}
}

// ========== Output Formatting ==========

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func printInstances(instances []service.ECSInstanceAdapter, region string) {
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

func printInstanceDetail(d *service.InstanceDetailAdapter) {
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

func printMetrics(m service.ECSMetricAdapter) {
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

func printLogs(result *service.LogQueryResult) {
	fmt.Printf("\nFound %d logs (hasMore: %v):\n", result.Count, result.HasMore)
	for _, entry := range result.Entries {
		ts := time.Unix(entry.Timestamp/1000, 0)
		fmt.Printf("\n[%s]\n", ts.Format("2006-01-02 15:04:05"))
		for k, v := range entry.Content {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
}

func printBuckets(buckets []service.BucketAdapter) {
	if len(buckets) == 0 {
		fmt.Println("No buckets found")
		return
	}
	fmt.Printf("\nBuckets:\n")
	for _, b := range buckets {
		fmt.Printf("  %-30s %-15s %s\n", b.Name, b.Location, b.CreationDate)
	}
}

func printObjects(objects []service.ObjectAdapter, bucket string) {
	fmt.Printf("\nObjects in %s:\n", bucket)
	if len(objects) == 0 {
		fmt.Println("No objects found")
		return
	}
	for _, obj := range objects {
		size := formatSize(obj.Size)
		fmt.Printf("  %-50s %10s %s\n", obj.Key, size, obj.LastModified)
	}
}

func formatSize(bytes int64) string {
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

func getAllRegions() []string {
	return []string{
		"cn-hangzhou", "cn-shanghai", "cn-beijing", "cn-shenzhen",
		"cn-guangzhou", "cn-chengdu", "cn-hongkong",
		"ap-southeast-1", "ap-southeast-2", "ap-southeast-3",
		"us-east-1", "us-west-1", "eu-central-1", "eu-west-1",
	}
}
