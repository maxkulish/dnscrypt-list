package config

import (
	"flag"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
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
	AllowedNames string `yaml:"allowed_names"`
	BlockedNames string `yaml:"blocked_names"`
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

// load reads yaml configuration file
func load(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// loadSource reads sources yaml file
func loadSource(file string) (*Sources, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	s := &Sources{}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Get read config file
func Get() (*Config, error) {

	flag.Parse()

	logger.SetLogger()

	var conf *Config

	// Read config.yml file
	conf, err := load(*configFile)
	if err != nil {
		return nil, err
	}

	//Read local source file
	if conf.SourcesLink.FilePath != "" {
		conf.SourceList, err = loadSource(conf.SourcesLink.FilePath)
		if err != nil {
			return nil, err
		}
	}

	// Read remote source file
	//if conf.SourcesLink.URL != "" {
	//	// TODO: get remote source file
	//	conf, err = load(conf.SourcesLink.FilePath)
	//}

	//fmt.Printf("Config: %v", conf)

	return conf, nil
}
