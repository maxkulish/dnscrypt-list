package config

import (
	"flag"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

// SourceType defines the type of the Source
type SourceType int8

const (
	// BlackList - domains to block
	blackList SourceType = iota + 1
	// WhiteList - domains to allow
	whiteList
)

var (
	// dnscrypt-list --conf=/etc/dnscrypt-list/config.yml
	configFile = flag.String("conf", "./config.yml", "Path to the config file")
)

// Update keeps timing params. How often to update different sources
type Update struct {
	Sources   time.Duration `yaml:"sources"`
	Blacklist time.Duration `yaml:"blacklist"`
	Whitelist time.Duration `yaml:"whitelist"`
}

// Output keeps paths to the output files
// these two params are defined in the dnscrypt-proxy config
// [blocked_names] and [allowed_names]
type Output struct {
	Whitelist string `yaml:"whitelist_path"`
	Blacklist string `yaml:"blacklist_path"`
}

// RawTarget keeps params of sources
// URL - "https://example.com/hosts.txt"
// Format: host, domain, bind, url
// Notes - additional information about target
type RawTarget struct {
	URL    string `yaml:"url,omitempty"`
	File   string `yaml:"file,omitempty"`
	Format string `yaml:"format"`
	Type   string `yaml:"type"`
	Notes  string `yaml:"notes,omitempty"`
}

// Sources keeps the list of RawTarget
type Sources struct {
	Targets []*RawTarget `yaml:"targets"`
}

// SourcesLink keeps file path or url to the sources list
type SourcesLink struct {
	FilePath string `yaml:"file,omitempty"`
	URL      string `yaml:"url,omitempty"`
}

// Config represents dnscrypt-list configuration parsed from the file --conf
type Config struct {
	Timeout     time.Duration `yaml:"timeout"`
	TempDir     string        `yaml:"temp_dir"`
	WhiteListDB string        `yaml:"whitelist_db"`
	BlackListDB string        `yaml:"blacklist_db"`
	Update      `yaml:"update"`
	SourcesLink `yaml:"sources"`
	Output      `yaml:"output"`
	SourceList  *Sources
}

// configFromYAML reads yaml configuration file
func configFromYAML(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("error closing file: %s\n", zap.Error(err))
		}
	}()
	config := &Config{}
	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return &Config{}, err
	}
	return config, nil
}

// sourceFromYAML reads yaml source file
func sourceFromYAML(path string) (*Sources, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("error closing file: %s\n", zap.Error(err))
		}
	}()

	sources := &Sources{}
	err = yaml.NewDecoder(f).Decode(sources)
	if err != nil {
		return &Sources{}, err
	}
	return sources, nil
}

// Get read config file
func Get() (*Config, error) {

	flag.Parse()

	logger.SetLogger()

	var conf *Config

	// Read config.yml file
	conf, err := configFromYAML(*configFile)
	if err != nil {
		return nil, err
	}

	//Read local source file and add it to the config
	if conf.SourcesLink.FilePath != "" {
		conf.SourceList, err = sourceFromYAML(conf.SourcesLink.FilePath)
		if err != nil {
			return nil, err
		}
	}

	// Read remote source file
	//if conf.SourcesLink.URL != "" {
	//	// TODO: get remote source file
	//	conf, err = load(conf.SourcesLink.FilePath)
	//}

	return conf, nil
}
