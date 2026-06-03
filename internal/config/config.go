package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigDir is the configuration directory name.
	ConfigDir = "cloud-manage"
	// ConfigFile is the configuration file name.
	ConfigFile = "config.yaml"
	// CurrentVersion is the current config file version.
	CurrentVersion = 1
)

// Config represents the application configuration.
type Config struct {
	// Version is the config file version for migration.
	Version int `yaml:"version"`
	// CurrentProfile is the name of the active profile.
	CurrentProfile string `yaml:"current_profile"`
	// Theme is the UI theme (light, dark, auto).
	Theme string `yaml:"theme"`
	// MemoryLimit is the memory limit in MB.
	MemoryLimit int `yaml:"memory_limit"`
	// Concurrency is the number of concurrent requests for multi-region queries.
	Concurrency int `yaml:"concurrency"`
	// PasswordPolicy defines the password complexity requirements.
	PasswordPolicy PasswordPolicy `yaml:"password_policy"`
	// Profiles is a map of profile name to profile configuration.
	Profiles map[string]*Profile `yaml:"profiles"`
}

// PasswordPolicy defines the password complexity requirements.
type PasswordPolicy struct {
	MinLength        int  `yaml:"min_length"`
	RequireUppercase bool `yaml:"require_uppercase"`
	RequireLowercase bool `yaml:"require_lowercase"`
	RequireDigit     bool `yaml:"require_digit"`
	RequireSpecial   bool `yaml:"require_special"`
}

// Profile represents a cloud account profile.
type Profile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint,omitempty"`
}

var (
	// globalConfig is the cached configuration.
	globalConfig *Config
	// configPath is the path to the configuration file.
	configPath string
	// mu protects globalConfig and configPath.
	mu sync.RWMutex
)

// GetConfigPath returns the path to the configuration file.
func GetConfigPath() (string, error) {
	mu.RLock()
	if configPath != "" {
		path := configPath
		mu.RUnlock()
		return path, nil
	}
	mu.RUnlock()

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	configDir := filepath.Join(userConfigDir, ConfigDir)
	path := filepath.Join(configDir, ConfigFile)

	mu.Lock()
	configPath = path
	mu.Unlock()

	return path, nil
}

// SetConfigPath sets a custom configuration file path (for testing).
func SetConfigPath(path string) {
	mu.Lock()
	defer mu.Unlock()
	configPath = path
	globalConfig = nil
}

// Load loads the configuration from the file.
// If the file does not exist, it returns a default configuration.
func Load() (*Config, error) {
	mu.RLock()
	if globalConfig != nil {
		cfg := globalConfig
		mu.RUnlock()
		return cfg, nil
	}
	mu.RUnlock()

	path, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config
		cfg := DefaultConfig()
		return cfg, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Migrate if needed
	if cfg.Version < CurrentVersion {
		if err := migrate(cfg); err != nil {
			return nil, fmt.Errorf("failed to migrate config: %w", err)
		}
	}

	// Validate
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	mu.Lock()
	globalConfig = cfg
	mu.Unlock()

	return cfg, nil
}

// Save saves the configuration to the file.
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := []byte(`# Cloud 管理小助手配置文件
# 文档: https://github.com/anthropics/cloud-manage#配置说明

`)

	// Write file with restricted permissions
	fullData := append(header, data...)
	if err := os.WriteFile(path, fullData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	mu.Lock()
	globalConfig = cfg
	mu.Unlock()

	return nil
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	return &Config{
		Version:        CurrentVersion,
		CurrentProfile: "",
		Theme:          "auto",
		MemoryLimit:    256,
		Concurrency:    3,
		PasswordPolicy: PasswordPolicy{
			MinLength:        8,
			RequireUppercase: true,
			RequireLowercase: true,
			RequireDigit:     true,
			RequireSpecial:   false,
		},
		Profiles: make(map[string]*Profile),
	}
}

// DemoConfig returns a demo configuration with comments.
func DemoConfig() *Config {
	cfg := DefaultConfig()
	cfg.CurrentProfile = "prod"
	cfg.Profiles["prod"] = &Profile{
		AccessKeyID:     "LTAI4xxx",
		AccessKeySecret: "encrypted:xxxx",
		Region:          "cn-hangzhou",
	}
	cfg.Profiles["dev"] = &Profile{
		AccessKeyID:     "LTAI4xxx",
		AccessKeySecret: "encrypted:xxxx",
		Region:          "cn-shanghai",
	}
	return cfg
}

// generateDemoYAML generates a demo YAML with comments.
func generateDemoYAML() []byte {
	return []byte(`# Cloud 管理小助手配置文件
# 文档: https://github.com/anthropics/cloud-manage#配置说明

# 配置文件版本（自动迁移，请勿手动修改）
version: 1

# 当前使用的账号 profile
# 可选值: prod, dev, test 等（对应下方 profiles 中的 key）
current_profile: prod

# 主题设置
# 可选值: light (浅色), dark (暗色), auto (跟随系统)
theme: auto

# 内存限制 (可选)
# 默认: 256MB
# 可通过环境变量 CLOUD_MEMORY_LIMIT 覆盖
memory_limit: 256

# 多 Region 查询并发数 (可选)
# 默认: 3
# 可通过环境变量 CLOUD_CONCURRENCY 覆盖
concurrency: 3

# 密码策略 (可选)
password_policy:
  # 最小密码长度
  min_length: 8
  # 是否要求大写字母
  require_uppercase: true
  # 是否要求小写字母
  require_lowercase: true
  # 是否要求数字
  require_digit: true
  # 是否要求特殊字符
  require_special: false

# 账号配置
# 每个 profile 代表一个云账号
profiles:
  # 生产环境账号
  prod:
    # 阿里云 AccessKey ID
    access_key_id: "LTAI4xxx"
    # 阿里云 AccessKey Secret (加密存储)
    access_key_secret: "encrypted:xxxx"
    # 默认区域
    # 可选值: cn-hangzhou, cn-shanghai, cn-beijing, cn-shenzhen 等
    region: "cn-hangzhou"
    # 自定义 Endpoint (可选)
    # 一般不需要配置，SDK 会自动选择
    # endpoint: "ecs.cn-hangzhou.aliyuncs.com"

  # 开发环境账号
  dev:
    access_key_id: "LTAI4xxx"
    access_key_secret: "encrypted:xxxx"
    region: "cn-shanghai"

  # 测试环境账号
  test:
    access_key_id: "LTAI4xxx"
    access_key_secret: "encrypted:xxxx"
    region: "cn-beijing"
`)
}

// InitConfig initializes the configuration file.
// If the file already exists, it returns an error unless force is true.
func InitConfig(force bool) error {
	path, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("config file already exists: %s", path)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write demo config
	if err := os.WriteFile(path, generateDemoYAML(), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// migrate migrates the configuration to the current version.
func migrate(cfg *Config) error {
	// Future migrations go here
	cfg.Version = CurrentVersion
	return nil
}

// validate validates the configuration.
func validate(cfg *Config) error {
	if cfg.Version < 1 {
		return fmt.Errorf("invalid config version: %d", cfg.Version)
	}

	if cfg.MemoryLimit < 0 {
		return fmt.Errorf("invalid memory limit: %d", cfg.MemoryLimit)
	}

	if cfg.Concurrency < 1 {
		return fmt.Errorf("invalid concurrency: %d", cfg.Concurrency)
	}

	// Validate theme
	validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
	if cfg.Theme != "" && !validThemes[cfg.Theme] {
		return fmt.Errorf("invalid theme: %s", cfg.Theme)
	}

	return nil
}

// HasConfig checks if the configuration file exists.
func HasConfig() bool {
	path, err := GetConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// GetProfile returns the current profile.
func GetProfile() (*Profile, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	if cfg.CurrentProfile == "" {
		return nil, fmt.Errorf("no current profile set")
	}

	profile, ok := cfg.Profiles[cfg.CurrentProfile]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", cfg.CurrentProfile)
	}

	return profile, nil
}

// GetProfileByName returns a profile by name.
func GetProfileByName(name string) (*Profile, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	profile, ok := cfg.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", name)
	}

	return profile, nil
}

// AddProfile adds or updates a profile.
func AddProfile(name string, profile *Profile) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]*Profile)
	}

	cfg.Profiles[name] = profile

	// If this is the first profile, set it as current
	if cfg.CurrentProfile == "" {
		cfg.CurrentProfile = name
	}

	return Save(cfg)
}

// RemoveProfile removes a profile.
func RemoveProfile(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if _, ok := cfg.Profiles[name]; !ok {
		return fmt.Errorf("profile not found: %s", name)
	}

	delete(cfg.Profiles, name)

	// If we removed the current profile, switch to another one
	if cfg.CurrentProfile == name {
		cfg.CurrentProfile = ""
		for k := range cfg.Profiles {
			cfg.CurrentProfile = k
			break
		}
	}

	return Save(cfg)
}

// SwitchProfile switches the current profile.
func SwitchProfile(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if _, ok := cfg.Profiles[name]; !ok {
		return fmt.Errorf("profile not found: %s", name)
	}

	cfg.CurrentProfile = name
	return Save(cfg)
}

// ListProfiles returns a list of profile names.
func ListProfiles() []string {
	cfg, err := Load()
	if err != nil {
		return nil
	}

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	return names
}

// ResetConfig deletes the configuration file.
func ResetConfig() error {
	path, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	mu.Lock()
	globalConfig = nil
	mu.Unlock()

	return nil
}
