package main

import (
	"cloud-manage/internal/config"
	"cloud-manage/internal/consts"
	"cloud-manage/internal/tui"
	"cloud-manage/service"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	accessKeyId     string
	accessKeySecret string
	region          string
	outputJSON      bool
	forceGUI        bool
	forceCLI        bool
	forceTUI        bool
)

func init() {
	flag.StringVar(&accessKeyId, "id", "", "AccessKey ID (or env CLOUD_ACCESS_KEY_ID)")
	flag.StringVar(&accessKeySecret, "secret", "", "AccessKey Secret (or env CLOUD_ACCESS_KEY_SECRET)")
	flag.StringVar(&region, "region", "cn-hangzhou", "Region ID (default: cn-hangzhou)")
	flag.BoolVar(&outputJSON, "json", false, "Output in JSON format")
	flag.BoolVar(&forceGUI, "gui", false, "Force GUI mode (requires display)")
	flag.BoolVar(&forceCLI, "cli", false, "Force CLI mode")
	flag.BoolVar(&forceTUI, "tui", false, "Force TUI mode (terminal UI)")
}

// knownServices are the valid CLI subcommands.
var knownServices = map[string]bool{
	"ecs":     true,
	"cms":     true,
	"sls":     true,
	"oss":     true,
	"vpc":     true,
	"slb":     true,
	"config":  true,
	"help":    true,
	"version": true,
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	fmt.Printf("\n  Cloud 管理小助手 %s\n\n", consts.Version)

	mode := detectMode()

	switch mode {
	case "gui":
		runGUI()
	case "tui":
		runTUI()
	case "cli":
		runCLI()
	}
}

// detectMode determines whether to run in GUI, TUI, or CLI mode.
func detectMode() string {
	// 1. Explicit flags take priority
	if forceGUI {
		return "gui"
	}
	if forceTUI {
		return "tui"
	}
	if forceCLI {
		return "cli"
	}

	// 2. If subcommand is present, use CLI mode
	args := flag.Args()
	if len(args) > 0 && knownServices[args[0]] {
		return "cli"
	}

	// 3. Check for graphical environment
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("WAYLAND_SOCKET") != "" {
		return "gui"
	}

	// 4. No display found, default to TUI
	return "tui"
}

// runGUI starts the Wails desktop application.
func runGUI() {
	fmt.Println("  检测到图形环境，启动 GUI 模式...")
	fmt.Println()

	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Cloud 管理小助手",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runTUI starts the terminal UI application.
func runTUI() {
	p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// loadCredentials loads credentials with priority: command line > env > config file > default
func loadCredentials() {
	// Priority 1: Command line flags (already set)
	// Priority 2: Environment variables
	envAK := os.Getenv("CLOUD_ACCESS_KEY_ID")
	envSK := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	envRegion := os.Getenv("CLOUD_REGION")

	// Priority 3: Config file
	var cfgAK, cfgSK, cfgRegion string
	if config.HasConfig() {
		cfg, err := config.Load()
		if err == nil && cfg.CurrentProfile != "" {
			if profile, ok := cfg.Profiles[cfg.CurrentProfile]; ok {
				cfgAK = profile.AccessKeyID
				cfgSK = profile.AccessKeySecret
				cfgRegion = profile.Region
			}
		}
	}

	// Apply priority and warn on conflicts
	if accessKeyId == "" {
		if envAK != "" {
			accessKeyId = envAK
			if cfgAK != "" && envAK != cfgAK {
				fmt.Fprintf(os.Stderr, "警告: 使用环境变量 CLOUD_ACCESS_KEY_ID，与配置文件不一致\n")
			}
		} else if cfgAK != "" {
			accessKeyId = cfgAK
		}
	} else {
		if envAK != "" && accessKeyId != envAK {
			fmt.Fprintf(os.Stderr, "警告: 使用命令行参数 -id，与环境变量不一致\n")
		}
		if cfgAK != "" && accessKeyId != cfgAK {
			fmt.Fprintf(os.Stderr, "警告: 使用命令行参数 -id，与配置文件不一致\n")
		}
	}

	if accessKeySecret == "" {
		if envSK != "" {
			accessKeySecret = envSK
			if cfgSK != "" && envSK != cfgSK {
				fmt.Fprintf(os.Stderr, "警告: 使用环境变量 CLOUD_ACCESS_KEY_SECRET，与配置文件不一致\n")
			}
		} else if cfgSK != "" {
			accessKeySecret = cfgSK
		}
	} else {
		if envSK != "" && accessKeySecret != envSK {
			fmt.Fprintf(os.Stderr, "警告: 使用命令行参数 -secret，与环境变量不一致\n")
		}
		if cfgSK != "" && accessKeySecret != cfgSK {
			fmt.Fprintf(os.Stderr, "警告: 使用命令行参数 -secret，与配置文件不一致\n")
		}
	}

	// Handle region with default
	defaultRegion := "cn-hangzhou"
	if region == defaultRegion {
		// User didn't set -region flag
		if envRegion != "" {
			region = envRegion
			if cfgRegion != "" && envRegion != cfgRegion {
				fmt.Fprintf(os.Stderr, "警告: 使用环境变量 CLOUD_REGION，与配置文件不一致\n")
			}
		} else if cfgRegion != "" {
			region = cfgRegion
		}
	} else {
		// User set -region flag
		if envRegion != "" && region != envRegion {
			fmt.Fprintf(os.Stderr, "警告: 使用命令行参数 -region，与环境变量不一致\n")
		}
		if cfgRegion != "" && region != cfgRegion {
			fmt.Fprintf(os.Stderr, "警告: 使用命令行参数 -region，与配置文件不一致\n")
		}
	}
}

// runCLI processes command-line operations.
func runCLI() {
	// Load credentials with priority handling
	loadCredentials()

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

	// Validate credentials (not needed for help display or product listing)
	needsCredentials := action != "" && !(serviceName == "cms" && action == "products") && serviceName != "config"
	if needsCredentials && (accessKeyId == "" || accessKeySecret == "") {
		fmt.Fprintf(os.Stderr, "Error: AccessKey ID and Secret are required.\n")
		fmt.Fprintf(os.Stderr, "Use -id/-secret flags, set CLOUD_ACCESS_KEY_ID/CLOUD_ACCESS_KEY_SECRET environment variables,\n")
		fmt.Fprintf(os.Stderr, "or use 'cloud-manage config add <profile>' to configure credentials.\n")
		os.Exit(1)
	}

	remainingArgs := []string{}
	if len(args) > 2 {
		remainingArgs = args[2:]
	}

	switch serviceName {
	case "help":
		printUsage()
	case "version":
		fmt.Printf("  cloud-manage %s\n", consts.Version)
	case "config":
		handleConfig(action, remainingArgs)
	case "ecs":
		handleECS(action, remainingArgs)
	case "cms":
		handleCMS(action, remainingArgs)
	case "sls":
		handleSLS(action, remainingArgs)
	case "oss":
		handleOSS(action, remainingArgs)
	case "vpc":
		handleVPC(action, remainingArgs)
	case "slb":
		handleSLB(action, remainingArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown service: %s\n", serviceName)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage:
  cloud-manage [flags] <service> <action> [args...]

Modes:
  --gui             Force GUI mode (requires display)
  --tui             Force TUI mode (terminal UI)
  --cli             Force CLI mode
  (auto)            Auto-detect based on display environment

Services:
  ecs               ECS 实例管理
  cms               云监控指标查询
  sls               日志服务查询
  oss               对象存储管理
  vpc               VPC 网络管理
  slb               负载均衡管理

Commands:
  help              显示帮助信息
  version           显示版本号

Flags:`)
	flag.PrintDefaults()
	fmt.Println(`
Examples:
  # Auto-detect mode (TUI if no display, GUI if display available)
  cloud-manage

  # Force TUI mode
  cloud-manage --tui

  # Force GUI mode
  cloud-manage --gui

  # CLI: List ECS instances
  cloud-manage -id LTAI4xxx -secret xxx ecs list

  # CLI: View ECS instance detail
  cloud-manage ecs detail i-xxx

  # CLI: Start ECS instance
  cloud-manage ecs start i-xxx

  # CLI: Query monitoring metrics
  cloud-manage cms metrics i-xxx

  # CLI: List OSS buckets
  cloud-manage oss buckets

  # CLI: List OSS objects
  cloud-manage oss objects my-bucket --prefix logs/

  # CLI: Query SLS logs
  cloud-manage sls logs my-project my-logstore --query "level: ERROR"

  # CLI: List VPCs
  cloud-manage vpc list

  # CLI: List VSwitches in a VPC
  cloud-manage vpc vswitches vpc-xxx

  # CLI: List SLBs
  cloud-manage slb list

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
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs detail <instance-id>\n")
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
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs start <instance-id>\n")
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
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs stop <instance-id> [--force]\n")
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
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs reboot <instance-id> [--force]\n")
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
	case "products":
		products := svc.GetSupportedProducts()
		if outputJSON {
			printJSON(products)
		} else {
			printProducts(products)
		}

	case "metrics":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage cms metrics <instance-id>\n")
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
  products          列出支持的云产品
  metrics <id>      查询实例监控指标`)
	}
}

// ========== SLS ==========

func handleSLS(action string, args []string) {
	svc := service.NewSLSService()

	switch action {
	case "logstores":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage sls logstores <project>\n")
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
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage sls logs <project> <logstore> [--query <query>] [--from <timestamp>] [--to <timestamp>] [--max <lines>]\n")
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
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage oss objects <bucket> [--prefix <prefix>] [--max <count>]\n")
			os.Exit(1)
		}
		bucket := args[0]
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

func handleVPC(action string, args []string) {
	svc := service.NewVPCService()

	switch action {
	case "list":
		result, err := svc.ListVPCs(accessKeyId, accessKeySecret, region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printVPCs(result.VPCs)
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage vpc detail <vpc-id>\n")
			os.Exit(1)
		}
		vpcId := args[0]
		result, err := svc.GetVPCDetail(accessKeyId, accessKeySecret, region, vpcId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printVPCDetail(result)
		}

	case "vswitches":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage vpc vswitches <vpc-id>\n")
			os.Exit(1)
		}
		vpcId := args[0]
		result, err := svc.ListVSwitches(accessKeyId, accessKeySecret, region, vpcId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printVSwitches(result.VSwitches)
		}

	default:
		fmt.Println(`VPC Actions:
  list                        列出 VPC
  detail <vpc-id>             查看 VPC 详情
  vswitches <vpc-id>          列出虚拟交换机`)
	}
}

func handleSLB(action string, args []string) {
	svc := service.NewSLBService()

	switch action {
	case "list":
		result, err := svc.ListSLBs(accessKeyId, accessKeySecret, region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printSLBs(result.SLBs)
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage slb detail <slb-id>\n")
			os.Exit(1)
		}
		slbId := args[0]
		result, err := svc.GetSLBDetail(accessKeyId, accessKeySecret, region, slbId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printSLBDetail(result)
		}

	case "listeners":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage slb listeners <slb-id>\n")
			os.Exit(1)
		}
		slbId := args[0]
		result, err := svc.ListSLBListeners(accessKeyId, accessKeySecret, region, slbId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			printJSON(result)
		} else {
			printSLBListeners(result.Listeners)
		}

	default:
		fmt.Println(`SLB Actions:
  list                          列出 SLB 实例
  detail <slb-id>               查看 SLB 详情
  listeners <slb-id>            列出监听器`)
	}
}

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

func printVPCs(vpcs []service.VPCAdapter) {
	if len(vpcs) == 0 {
		fmt.Println("No VPCs found")
		return
	}
	fmt.Printf("\nVPCs:\n")
	for _, v := range vpcs {
		fmt.Printf("  %-20s %-20s %-18s %s\n", v.VpcId, v.VpcName, v.CidrBlock, v.Status)
	}
}

func printVPCDetail(d *service.VPCDetailAdapter) {
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

func printVSwitches(vswitches []service.VSwitchAdapter) {
	if len(vswitches) == 0 {
		fmt.Println("No VSwitches found")
		return
	}
	fmt.Printf("\nVSwitches:\n")
	for _, vs := range vswitches {
		fmt.Printf("  %-20s %-20s %-18s %s\n", vs.VSwitchId, vs.VSwitchName, vs.CidrBlock, vs.ZoneId)
	}
}

func printSLBs(slbs []service.SLBAdapter) {
	if len(slbs) == 0 {
		fmt.Println("No SLBs found")
		return
	}
	fmt.Printf("\nSLBs:\n")
	for _, lb := range slbs {
		fmt.Printf("  %-20s %-20s %-15s %-10s %s\n", lb.LoadBalancerId, lb.LoadBalancerName, lb.Address, lb.AddressType, lb.Status)
	}
}

func printSLBDetail(d *service.SLBDetailAdapter) {
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

func printSLBListeners(listeners []service.SLBListenerAdapter) {
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

func getAllRegions() []string {
	return []string{
		"cn-hangzhou", "cn-shanghai", "cn-beijing", "cn-shenzhen",
		"cn-guangzhou", "cn-chengdu", "cn-hongkong",
		"ap-southeast-1", "ap-southeast-2", "ap-southeast-3",
		"us-east-1", "us-west-1", "eu-central-1", "eu-west-1",
	}
}

// ========== Config ==========

func handleConfig(action string, args []string) {
	switch action {
	case "init":
		handleConfigInit(args)
	case "add":
		handleConfigAdd(args)
	case "remove":
		handleConfigRemove(args)
	case "list":
		handleConfigList(args)
	case "switch":
		handleConfigSwitch(args)
	case "show":
		handleConfigShow(args)
	case "reset":
		handleConfigReset(args)
	default:
		fmt.Println(`Config Actions:
  init              初始化配置文件
  add <profile>     添加账号
  remove <profile>  删除账号
  list              列出所有账号
  switch <profile>  切换默认账号
  show              显示当前账号
  reset             重置配置文件（删除所有凭证）`)
	}
}

func handleConfigInit(args []string) {
	// Check if config already exists
	if config.HasConfig() {
		fmt.Fprintf(os.Stderr, "配置文件已存在: ")
		path, _ := config.GetConfigPath()
		fmt.Fprintf(os.Stderr, "%s\n", path)
		fmt.Fprintf(os.Stderr, "使用 --force 参数强制重新生成\n")
		os.Exit(1)
	}

	// Init config
	if err := config.InitConfig(false); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	path, _ := config.GetConfigPath()
	fmt.Printf("配置文件已生成: %s\n", path)
	fmt.Println("请编辑配置文件，填入您的云账号凭证。")
}

func handleConfigAdd(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: cloud-manage config add <profile>\n")
		os.Exit(1)
	}

	profileName := args[0]

	// Check if config exists
	if !config.HasConfig() {
		fmt.Println("配置文件不存在，正在初始化...")
		if err := config.InitConfig(false); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Read credentials from flags or environment
	akId := accessKeyId
	if akId == "" {
		akId = os.Getenv("CLOUD_ACCESS_KEY_ID")
	}
	akSecret := accessKeySecret
	if akSecret == "" {
		akSecret = os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	}
	profileRegion := region

	if akId == "" || akSecret == "" {
		fmt.Fprintf(os.Stderr, "Error: AccessKey ID and Secret are required.\n")
		fmt.Fprintf(os.Stderr, "Use -id/-secret flags or set CLOUD_ACCESS_KEY_ID/CLOUD_ACCESS_KEY_SECRET environment variables.\n")
		os.Exit(1)
	}

	// Create profile
	profile := &config.Profile{
		AccessKeyID:     akId,
		AccessKeySecret: akSecret,
		Region:          profileRegion,
	}

	// Add profile
	if err := config.AddProfile(profileName, profile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("账号 '%s' 已添加。\n", profileName)
}

func handleConfigRemove(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: cloud-manage config remove <profile>\n")
		os.Exit(1)
	}

	profileName := args[0]

	if err := config.RemoveProfile(profileName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("账号 '%s' 已删除。\n", profileName)
}

func handleConfigList(args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("没有配置任何账号。")
		fmt.Println("使用 'cloud-manage config add <profile>' 添加账号。")
		return
	}

	fmt.Println("已配置的账号:")
	fmt.Println()

	// Sort profile names
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		profile := cfg.Profiles[name]
		marker := " "
		if name == cfg.CurrentProfile {
			marker = "*"
		}
		fmt.Printf("  %s %-15s %-20s %s\n", marker, name, profile.AccessKeyID, profile.Region)
	}

	fmt.Println()
	fmt.Println("* 表示当前使用的账号")
}

func handleConfigSwitch(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: cloud-manage config switch <profile>\n")
		os.Exit(1)
	}

	profileName := args[0]

	if err := config.SwitchProfile(profileName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已切换到账号 '%s'。\n", profileName)
}

func handleConfigShow(args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cfg.CurrentProfile == "" {
		fmt.Println("当前没有设置默认账号。")
		fmt.Println("使用 'cloud-manage config switch <profile>' 切换账号。")
		return
	}

	profile, ok := cfg.Profiles[cfg.CurrentProfile]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: 当前账号 '%s' 不存在。\n", cfg.CurrentProfile)
		os.Exit(1)
	}

	fmt.Printf("当前账号: %s\n", cfg.CurrentProfile)
	fmt.Printf("  AccessKey ID: %s\n", profile.AccessKeyID)
	fmt.Printf("  Region:       %s\n", profile.Region)
	if profile.Endpoint != "" {
		fmt.Printf("  Endpoint:     %s\n", profile.Endpoint)
	}
}

func handleConfigReset(args []string) {
	if err := config.ResetConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("配置文件已重置。")
}
