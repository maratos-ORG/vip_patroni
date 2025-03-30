package config

import (
	"fmt"
	"os"
	"strings"
	log "vip_patroni/internal/logging"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config represents the configuration of the vip_patroni
type Config struct {
	IP                   string `mapstructure:"ip"`
	Mask                 int    `mapstructure:"netmask"`
	Iface                string `mapstructure:"interface"`
	Interval             int    `mapstructure:"interval"`    //milliseconds
	RetryAfter           int    `mapstructure:"retry-after"` //milliseconds
	RetryNum             int    `mapstructure:"retry-num"`
	LogLevel             string `mapstructure:"log-level"`
	Status               string `mapstructure:"status"`
	PatroniURL           string `mapstructure:"patroni-url"`
	PatroniTimeoutMillis int    `json:"patroni_timeout_millis"` // milliseconds
	APIPort              string `mapstructure:"api-port"`
}

func defineFlags() {
	// When adding new flags here, consider adding them to the Config struct above
	// and then make sure to insert them into the conf instance in NewConfig down below.
	pflag.String("config", "", "Location of the configuration file.")
	pflag.String("ip", "", "Virtual IP address to configure.")
	pflag.String("patroni-url", "http://localhost:8008", "URL of patroni API")
	pflag.String("api-port", "8010", "Port of vip_patroni API")
	pflag.Int("netmask", 32, "The netmask used for the IP address. Defaults to -1 which assigns ipv4 default mask.")
	pflag.String("interface", "", "Network interface to configure on .")
	pflag.String("interval", "1000", "DCS scan interval in milliseconds.")
	pflag.String("log-level", "INFO", "Set log level (trace,debug,info,warning,error,fatal).")
	pflag.CommandLine.SortFlags = false
}

// NewConfig returns a new Config structure
func NewConfig(verApp string) (*Config, error) {
	var err error
	var showHelp, showVersion bool

	if err != nil {
		log.Fatal("unable init logger.NewLogger: %s", err)
	}

	pflag.BoolVar(&showHelp, "help", false, "Show help")
	pflag.BoolVar(&showVersion, "version", false, "Show version of vip_patroni")

	defineFlags()
	pflag.Parse()
	if showHelp {
		pflag.Usage()
		os.Exit(0)
	}
	if showVersion {
		fmt.Println(verApp)
		os.Exit(0)
	}
	// import pflags into viper
	_ = viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("vip_patroni")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetDefault("ip", "postgres://localhost")
	viper.SetDefault("patroni_timeout_millis", 500)
	if viper.IsSet("config") {
		log.Debug("Parsing config: %s", viper.GetString("config"))
		viper.SetConfigFile(viper.GetString("config"))
		err := viper.ReadInConfig()
		if err != nil {
			log.Info("unable to read config file: %s", err)
		}
	} else {
		log.Debug("config file is not specified.")
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Debug("unable to decode viper config into config struct: %s", err)
	}
	return conf, nil
}
