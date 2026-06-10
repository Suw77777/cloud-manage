package main

import (
	"cloud-manage/internal/cli"
	"cloud-manage/internal/config"
	"cloud-manage/internal/consts"
	"cloud-manage/internal/tui"
	"embed"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
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
	profileName     string
	outputJSON      bool
	forceGUI        bool
	forceCLI        bool
	forceTUI        bool
)

func init() {
	flag.StringVar(&accessKeyId, "id", "", "AccessKey ID (or env CLOUD_ACCESS_KEY_ID)")
	flag.StringVar(&accessKeySecret, "secret", "", "AccessKey Secret (or env CLOUD_ACCESS_KEY_SECRET)")
	flag.StringVar(&region, "region", "cn-hangzhou", "Region ID (default: cn-hangzhou)")
	flag.StringVar(&profileName, "profile", "", "Use specified profile (overrides current_profile)")
	flag.BoolVar(&outputJSON, "json", false, "Output in JSON format")
	flag.BoolVar(&forceGUI, "gui", false, "Force GUI mode (requires display)")
	flag.BoolVar(&forceCLI, "cli", false, "Force CLI mode")
	flag.BoolVar(&forceTUI, "tui", false, "Force TUI mode (terminal UI)")
}

var knownServices = map[string]bool{
	"ecs": true, "cms": true, "sls": true, "oss": true,
	"vpc": true, "slb": true, "config": true, "help": true, "version": true,
}

func main() {
	flag.Usage = printUsage
	flag.Parse()
	setMemoryLimit()

	fmt.Printf("\n  Cloud 管理小助手 %s\n\n", consts.Version)

	switch detectMode() {
	case "gui":
		runGUI()
	case "tui":
		runTUI()
	case "cli":
		runCLI()
	}
}

// setMemoryLimit sets the memory limit from config, env var, or default.
func setMemoryLimit() {
	limitMB := 256
	if config.HasConfig() {
		if cfg, err := config.Load(); err == nil && cfg.MemoryLimit > 0 {
			limitMB = cfg.MemoryLimit
		}
	}
	if envLimit := os.Getenv("CLOUD_MEMORY_LIMIT"); envLimit != "" {
		if parsed, err := strconv.Atoi(envLimit); err == nil && parsed > 0 {
			limitMB = parsed
		}
	}
	debug.SetMemoryLimit(int64(limitMB) * 1024 * 1024)
}

// detectMode determines whether to run in GUI, TUI, or CLI mode.
func detectMode() string {
	if forceGUI {
		return "gui"
	}
	if forceTUI {
		return "tui"
	}
	if forceCLI {
		return "cli"
	}
	args := flag.Args()
	if len(args) > 0 && knownServices[args[0]] {
		return "cli"
	}
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("WAYLAND_SOCKET") != "" {
		return "gui"
	}
	return "tui"
}

func runGUI() {
	fmt.Println("  检测到图形环境，启动 GUI 模式...")
	fmt.Println()
	app := NewApp()
	if err := wails.Run(&options.App{
		Title:  "Cloud 管理小助手",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup:   app.startup,
		Bind:        []interface{}{app},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runTUI() {
	p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// loadCredentials loads credentials with priority: command line > env > --profile > current_profile > default.
func loadCredentials() {
	envAK := os.Getenv("CLOUD_ACCESS_KEY_ID")
	envSK := os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	envRegion := os.Getenv("CLOUD_REGION")

	var cfgAK, cfgSK, cfgRegion string
	if config.HasConfig() {
		if cfg, err := config.Load(); err == nil {
			targetProfile := profileName
			if targetProfile == "" {
				targetProfile = cfg.CurrentProfile
			}
			if targetProfile != "" {
				if profile, ok := cfg.Profiles[targetProfile]; ok {
					cfgAK = profile.AccessKeyID
					cfgSK = profile.AccessKeySecret
					cfgRegion = profile.Region
					if config.IsEncrypted(cfgSK) {
						if decrypted, err := config.GetProfileWithCredentials(targetProfile); err != nil {
							fmt.Fprintf(os.Stderr, "警告: 无法解密 profile '%s' 的凭证: %v\n", targetProfile, err)
						} else {
							cfgSK = decrypted.AccessKeySecret
						}
					}
				} else if profileName != "" {
					fmt.Fprintf(os.Stderr, "警告: profile '%s' 不存在\n", profileName)
				}
			}
		}
	}

	// Apply priority: CLI flag > env > config
	if accessKeyId == "" {
		if envAK != "" {
			accessKeyId = envAK
		} else if cfgAK != "" {
			accessKeyId = cfgAK
		}
	}
	if accessKeySecret == "" {
		if envSK != "" {
			accessKeySecret = envSK
		} else if cfgSK != "" {
			accessKeySecret = cfgSK
		}
	}
	defaultRegion := "cn-hangzhou"
	if region == defaultRegion {
		if envRegion != "" {
			region = envRegion
		} else if cfgRegion != "" {
			region = cfgRegion
		}
	}
}

func getConcurrency() int {
	concurrency := 3
	if config.HasConfig() {
		if cfg, err := config.Load(); err == nil && cfg.Concurrency > 0 {
			concurrency = cfg.Concurrency
		}
	}
	if envConcurrency := os.Getenv("CLOUD_CONCURRENCY"); envConcurrency != "" {
		if parsed, err := strconv.Atoi(envConcurrency); err == nil && parsed > 0 {
			concurrency = parsed
		}
	}
	return concurrency
}

func runCLI() {
	loadCredentials()
	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	serviceName, action := args[0], ""
	if len(args) > 1 {
		action = args[1]
	}
	remainingArgs := args[2:]

	needsCredentials := action != "" && !(serviceName == "cms" && action == "products") && serviceName != "config"
	if needsCredentials && (accessKeyId == "" || accessKeySecret == "") {
		fmt.Fprintf(os.Stderr, "Error: AccessKey ID and Secret are required.\n")
		fmt.Fprintf(os.Stderr, "Use -id/-secret flags, set CLOUD_ACCESS_KEY_ID/CLOUD_ACCESS_KEY_SECRET environment variables,\n")
		fmt.Fprintf(os.Stderr, "or use 'cloud-manage config add <profile>' to configure credentials.\n")
		os.Exit(1)
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

// ========== ECS ==========

func handleECS(action string, args []string) {
	h := cli.NewECSHandler()

	switch action {
	case "list":
		regions := []string{region}
		if region == "all" {
			regions = consts.RegionIDs()
		}
		if len(regions) == 1 {
			result, err := h.ListInstances(accessKeyId, accessKeySecret, regions[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if outputJSON {
				cli.PrintJSON(result)
			} else {
				cli.PrintInstances(result.Instances, regions[0])
			}
		} else {
			results := h.ListInstancesMultiRegion(accessKeyId, accessKeySecret, regions, getConcurrency())
			if outputJSON {
				cli.PrintJSON(results)
			} else {
				for _, r := range results {
					if r.Error != "" {
						fmt.Printf("\n[%s] Error: %s\n", r.Region, r.Error)
					} else {
						cli.PrintInstances(r.Instances, r.Region)
					}
				}
			}
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs detail <instance-id>\n")
			os.Exit(1)
		}
		result, err := h.GetInstanceDetail(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintInstanceDetail(result)
		}

	case "start":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs start <instance-id>\n")
			os.Exit(1)
		}
		if err := h.StartInstance(accessKeyId, accessKeySecret, region, args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance %s start command sent successfully\n", args[0])

	case "stop":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs stop <instance-id> [--force]\n")
			os.Exit(1)
		}
		force := false
		for _, arg := range args[1:] {
			if arg == "--force" || arg == "-f" {
				force = true
			}
		}
		if err := h.StopInstance(accessKeyId, accessKeySecret, region, args[0], force); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance %s stop command sent successfully\n", args[0])

	case "reboot":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage ecs reboot <instance-id> [--force]\n")
			os.Exit(1)
		}
		force := false
		for _, arg := range args[1:] {
			if arg == "--force" || arg == "-f" {
				force = true
			}
		}
		if err := h.RebootInstance(accessKeyId, accessKeySecret, region, args[0], force); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Instance %s reboot command sent successfully\n", args[0])

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
	h := cli.NewCMSHandler()

	switch action {
	case "products":
		if outputJSON {
			cli.PrintJSON(h.GetSupportedProducts())
		} else {
			cli.PrintProducts(h.GetSupportedProducts())
		}

	case "metrics":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage cms metrics <instance-id>\n")
			os.Exit(1)
		}
		result, err := h.GetInstanceMetrics(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintMetrics(*result)
		}

	default:
		fmt.Println(`CMS Actions:
  products          列出支持的云产品
  metrics <id>      查询实例监控指标`)
	}
}

// ========== SLS ==========

func handleSLS(action string, args []string) {
	h := cli.NewSLSHandler()

	switch action {
	case "logstores":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage sls logstores <project>\n")
			os.Exit(1)
		}
		logstores, err := h.ListLogStores(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(logstores)
		} else {
			fmt.Printf("LogStores in %s:\n", args[0])
			for _, ls := range logstores {
				fmt.Printf("  - %s\n", ls)
			}
		}

	case "export":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage sls export <project> <logstore> [--format csv|json] [--output <file>] [--query <query>] [--from <time>] [--to <time>] [--max <lines>]\n")
			os.Exit(1)
		}
		project, logstore := args[0], args[1]
		query, format, outputPath := "", "csv", ""
		fromStr, toStr := "", ""
		maxLines := int64(1000)
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--query", "-q":
				if i+1 < len(args) {
					query = args[i+1]
					i++
				}
			case "--format", "-f":
				if i+1 < len(args) {
					format = args[i+1]
					i++
				}
			case "--output", "-o":
				if i+1 < len(args) {
					outputPath = args[i+1]
					i++
				}
			case "--from":
				if i+1 < len(args) {
					fromStr = args[i+1]
					i++
				}
			case "--to":
				if i+1 < len(args) {
					toStr = args[i+1]
					i++
				}
			case "--max", "-m":
				if i+1 < len(args) {
					fmt.Sscanf(args[i+1], "%d", &maxLines)
					i++
				}
			}
		}
		result, err := h.ExportLogs(accessKeyId, accessKeySecret, region, project, logstore, query, fromStr, toStr, maxLines, format, outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("导出成功!\n  文件: %s\n  数量: %d 条\n  格式: %s\n", result.FilePath, result.Count, result.Format)

	case "logs":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage sls logs <project> <logstore> [--query <query>] [--from <timestamp>] [--to <timestamp>] [--max <lines>]\n")
			os.Exit(1)
		}
		project, logstore := args[0], args[1]
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
		result, err := h.QueryLogs(accessKeyId, accessKeySecret, region, project, logstore, query, from, to, maxLines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintLogs(result)
		}

	default:
		fmt.Println(`SLS Actions:
  logstores <project>                    列出 Logstore
  logs <project> <logstore> [options]    查询日志
  export <project> <logstore> [options]  导出日志到文件

Query Options:
    --query, -q <query>    查询表达式
    --from <time>          开始时间 (ISO 8601, 相对时间如 1h/30m/7d, Unix 时间戳)
    --to <time>            结束时间
    --max <lines>          最大返回行数 (default: 100, 最大 5000)

Export Options:
    --format <format>      导出格式 (csv, json, default: csv)
    --output, -o <file>    输出文件名 (自动生成)`)
	}
}

// ========== OSS ==========

func handleOSS(action string, args []string) {
	h := cli.NewOSSHandler()

	switch action {
	case "buckets":
		result, err := h.ListBuckets(accessKeyId, accessKeySecret, region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintBuckets(result.Buckets)
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
		result, err := h.ListObjects(accessKeyId, accessKeySecret, region, bucket, prefix, maxKeys)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintObjects(result.Objects, bucket)
		}

	default:
		fmt.Println(`OSS Actions:
  buckets                           列出 Bucket
  objects <bucket> [options]        列出对象
    --prefix, -p <prefix>   前缀过滤
    --max <count>           最大返回数量 (default: 100)`)
	}
}

// ========== VPC ==========

func handleVPC(action string, args []string) {
	h := cli.NewVPCHandler()

	switch action {
	case "list":
		result, err := h.ListVPCs(accessKeyId, accessKeySecret, region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintVPCs(result.VPCs)
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage vpc detail <vpc-id>\n")
			os.Exit(1)
		}
		result, err := h.GetVPCDetail(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintVPCDetail(result)
		}

	case "vswitches":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage vpc vswitches <vpc-id>\n")
			os.Exit(1)
		}
		result, err := h.ListVSwitches(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintVSwitches(result.VSwitches)
		}

	default:
		fmt.Println(`VPC Actions:
  list                        列出 VPC
  detail <vpc-id>             查看 VPC 详情
  vswitches <vpc-id>          列出虚拟交换机`)
	}
}

// ========== SLB ==========

func handleSLB(action string, args []string) {
	h := cli.NewSLBHandler()

	switch action {
	case "list":
		result, err := h.ListSLBs(accessKeyId, accessKeySecret, region)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintSLBs(result.SLBs)
		}

	case "detail":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage slb detail <slb-id>\n")
			os.Exit(1)
		}
		result, err := h.GetSLBDetail(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintSLBDetail(result)
		}

	case "listeners":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "Usage: cloud-manage slb listeners <slb-id>\n")
			os.Exit(1)
		}
		result, err := h.ListSLBListeners(accessKeyId, accessKeySecret, region, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if outputJSON {
			cli.PrintJSON(result)
		} else {
			cli.PrintSLBListeners(result.Listeners)
		}

	default:
		fmt.Println(`SLB Actions:
  list                          列出 SLB 实例
  detail <slb-id>               查看 SLB 详情
  listeners <slb-id>            列出监听器`)
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
	if config.HasConfig() {
		path, _ := config.GetConfigPath()
		fmt.Fprintf(os.Stderr, "配置文件已存在: %s\n", path)
		fmt.Fprintf(os.Stderr, "使用 --force 参数强制重新生成\n")
		os.Exit(1)
	}
	if err := config.InitConfig(false); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	path, _ := config.GetConfigPath()
	fmt.Printf("配置文件已生成: %s\n请编辑配置文件，填入您的云账号凭证。\n", path)
}

func handleConfigAdd(args []string) {
	saveCredentials := false
	filteredArgs := []string{}
	for _, arg := range args {
		if arg == "--save" {
			saveCredentials = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	args = filteredArgs

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: cloud-manage config add <profile> [--save]\n")
		os.Exit(1)
	}

	profileName := args[0]
	if !config.HasConfig() {
		fmt.Println("配置文件不存在，正在初始化...")
		if err := config.InitConfig(false); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	akId := accessKeyId
	if akId == "" {
		akId = os.Getenv("CLOUD_ACCESS_KEY_ID")
	}
	akSecret := accessKeySecret
	if akSecret == "" {
		akSecret = os.Getenv("CLOUD_ACCESS_KEY_SECRET")
	}
	if akId == "" {
		fmt.Fprintf(os.Stderr, "Error: AccessKey ID is required.\n")
		os.Exit(1)
	}

	profile := &config.Profile{
		AccessKeyID:     akId,
		AccessKeySecret: akSecret,
		Region:          region,
	}
	if err := config.AddProfile(profileName, profile, saveCredentials); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if saveCredentials {
		fmt.Printf("账号 '%s' 已添加（凭证已加密保存）。\n", profileName)
	} else {
		fmt.Printf("账号 '%s' 已添加（仅保存 AccessKey ID 和 Region）。\n", profileName)
	}
}

func handleConfigRemove(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: cloud-manage config remove <profile>\n")
		os.Exit(1)
	}
	if err := config.RemoveProfile(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("账号 '%s' 已删除。\n", args[0])
}

func handleConfigList(args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if len(cfg.Profiles) == 0 {
		fmt.Println("没有配置任何账号。使用 'cloud-manage config add <profile>' 添加账号。")
		return
	}
	fmt.Println("已配置的账号:")
	fmt.Println()
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
	if err := config.SwitchProfile(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("已切换到账号 '%s'。\n", args[0])
}

func handleConfigShow(args []string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if cfg.CurrentProfile == "" {
		fmt.Println("当前没有设置默认账号。使用 'cloud-manage config switch <profile>' 切换账号。")
		return
	}
	profile, ok := cfg.Profiles[cfg.CurrentProfile]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: 当前账号 '%s' 不存在。\n", cfg.CurrentProfile)
		os.Exit(1)
	}
	fmt.Printf("当前账号: %s\n  AccessKey ID: %s\n  Region: %s\n", cfg.CurrentProfile, profile.AccessKeyID, profile.Region)
	if profile.Endpoint != "" {
		fmt.Printf("  Endpoint: %s\n", profile.Endpoint)
	}
}

func handleConfigReset(args []string) {
	if err := config.ResetConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("配置文件已重置。")
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
  cloud-manage                          # Auto-detect mode
  cloud-manage --tui                    # Force TUI mode
  cloud-manage ecs list                 # CLI: List ECS instances
  cloud-manage ecs detail i-xxx         # CLI: View ECS detail
  cloud-manage sls logs proj logstore   # CLI: Query SLS logs
  cloud-manage oss buckets              # CLI: List OSS buckets
  cloud-manage vpc list                 # CLI: List VPCs
  cloud-manage slb list                 # CLI: List SLBs

Environment Variables:
  CLOUD_ACCESS_KEY_ID      AccessKey ID
  CLOUD_ACCESS_KEY_SECRET  AccessKey Secret`)
}
