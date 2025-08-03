package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// TracerConfig holds configuration for the tracer sidecar
type TracerConfig struct {
	Collector     CollectorConfig `mapstructure:"collector"`
	Tracer        TracerSettings  `mapstructure:"tracer"`
	Proxy         ProxyConfig     `mapstructure:"proxy"`
	Kubernetes    K8sConfig       `mapstructure:"kubernetes"`
	LogLevel      string          `mapstructure:"log_level"`
}

// CollectorConfig holds collector connection settings
type CollectorConfig struct {
	Endpoint    string        `mapstructure:"endpoint"`
	BatchSize   int           `mapstructure:"batch_size"`
	Timeout     time.Duration `mapstructure:"timeout"`
	Insecure    bool          `mapstructure:"insecure"`
}

// TracerSettings holds tracer-specific settings
type TracerSettings struct {
	SamplingRate      float64 `mapstructure:"sampling_rate"`
	MaxSpansPerTrace  int     `mapstructure:"max_spans_per_trace"`
	BufferSize        int     `mapstructure:"buffer_size"`
	FlushInterval     time.Duration `mapstructure:"flush_interval"`
}

// ProxyConfig holds proxy settings
type ProxyConfig struct {
	InboundPort  int      `mapstructure:"inbound_port"`
	OutboundPort int      `mapstructure:"outbound_port"`
	Protocols    []string `mapstructure:"protocols"`
	Ports        []int    `mapstructure:"ports"`
}

// K8sConfig holds Kubernetes-specific settings
type K8sConfig struct {
	PodName     string `mapstructure:"pod_name"`
	Namespace   string `mapstructure:"namespace"`
	ServiceName string `mapstructure:"service_name"`
}

// LoadTracerConfig loads tracer configuration from file
func LoadTracerConfig(configPath string) (*TracerConfig, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	setTracerDefaults()

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DEVMESH")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config TracerConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setTracerDefaults() {
	viper.SetDefault("collector.endpoint", "devmesh-collector:14268")
	viper.SetDefault("collector.batch_size", 100)
	viper.SetDefault("collector.timeout", "5s")
	viper.SetDefault("collector.insecure", true)

	viper.SetDefault("tracer.sampling_rate", 1.0)
	viper.SetDefault("tracer.max_spans_per_trace", 1000)
	viper.SetDefault("tracer.buffer_size", 1000)
	viper.SetDefault("tracer.flush_interval", "5s")

	viper.SetDefault("proxy.inbound_port", 15001)
	viper.SetDefault("proxy.outbound_port", 15002)
	viper.SetDefault("proxy.protocols", []string{"http"})
	viper.SetDefault("proxy.ports", []int{8080})

	viper.SetDefault("log_level", "info")
}