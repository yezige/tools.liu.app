package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var defaultConf = []byte(`
core:
  enabled: true # enable httpd server
  address: "" # ip address to bind (default: any)
  port: "8088" # ignore this port number if auto_tls is enabled (listen 443).
  worker_num: 0 # default worker number is runtime.NumCPU()
  mode: "release" # release mode or debug mode
  shutdown_timeout: 30 # default is 30 second
  pid:
    enabled: true
    path: "app.pid"
    override: true

api:
  stat_go_uri: /api/stat/go
  google_api_key: ""

log:
  format: "string" # string or json
  access_log: "stdout" # stdout: output to console, or define log path like "log/access_log"
  access_level: "debug"
  error_log: "stderr" # stderr: output to console, or define log path like "log/error_log"
  error_level: "debug"
  hide_token: true

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

site:
  name: "YouTube 视频下载"
  logo: "Tools"
  author: "Yezi"
  hostname: "tools.liu.app"
  menu:
    - name: "menu-youtube"
      url: "/youtube"
    - name: "menu-about"
      url: "/about"
    - name: "menu-tags"
      url: "/tags"
  twitter:
    image: "https://tools.liu.app/static/images/twitter_image.webp"
  google_analytics: "G-xxx"
  google_site_verification: ""
  baidu_site_verification: "code-xxx"
  baidu_analytics: ""
  google_adsense: ""
  disqus:
    shortname: "liu-tools"
    lazyload: true
    site: https://tools.liu.app
    api: 
    mode: 1
    badge: 叶子
    timeout: 3000
    apihost:  # disqusapi.swig文件使用
    count: true # 不使用disqus自带的count方式，防止引disqus的js时404
    emoji_path: https://cdnjs.cloudflare.com/ajax/libs/emojione/2.2.7/assets/png/
  fancybox: true
  sponsors:
    - sponsor1
    - sponsor2
`)

// ConfYaml is config structure.
type ConfYaml struct {
	Core  SectionCore  `yaml:"core" mapstructure:"core" json:"core"`
	API   SectionAPI   `yaml:"api" mapstructure:"api" json:"api"`
	Log   SectionLog   `yaml:"log" mapstructure:"log" json:"log"`
	Site  SectionSite  `yaml:"site" mapstructure:"site" json:"site"`
	Redis SectionRedis `yaml:"redis" mapstructure:"redis" json:"redis"`
}

// SectionCore is sub section of config.
type SectionCore struct {
	Enabled         bool       `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	Address         string     `yaml:"address" mapstructure:"address" json:"address"`
	Port            string     `yaml:"port" mapstructure:"port" json:"port"`
	WorkerNum       int64      `yaml:"worker_num" mapstructure:"worker_num" json:"worker_num"`
	Mode            string     `yaml:"mode" mapstructure:"mode" json:"mode"`
	ShutdownTimeout int64      `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout" json:"shutdown_timeout"`
	PID             SectionPID `yaml:"pid" mapstructure:"pid" json:"pid"`
}

type SectionPID struct {
	Enabled  bool   `yaml:"enabled"`
	Path     string `yaml:"path"`
	Override bool   `yaml:"override"`
}

// SectionAPI is sub section of config.
type SectionAPI struct {
	StatGoURI    string `yaml:"stat_go_uri" mapstructure:"stat_go_uri" json:"stat_go_uri"`
	GoogleAPIKey string `yaml:"google_api_key" mapstructure:"google_api_key" json:"google_api_key"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Format      string `yaml:"format" mapstructure:"format" json:"format"`
	AccessLog   string `yaml:"access_log" mapstructure:"access_log" json:"access_log"`
	AccessLevel string `yaml:"access_level" mapstructure:"access_level" json:"access_level"`
	ErrorLog    string `yaml:"error_log" mapstructure:"error_log" json:"error_log"`
	ErrorLevel  string `yaml:"error_level" mapstructure:"error_level" json:"error_level"`
	HideToken   bool   `yaml:"hide_token" mapstructure:"hide_token" json:"hide_token"`
}

type SectionRedis struct {
	Addr     string `yaml:"addr" mapstructure:"addr" json:"addr"`
	Password string `yaml:"password" mapstructure:"password" json:"password"`
	DB       int    `yaml:"db" mapstructure:"db" json:"db"`
}

type SectionSite struct {
	Name                   string              `yaml:"name" mapstructure:"name" json:"name"`
	Logo                   string              `yaml:"logo" mapstructure:"logo" json:"logo"`
	Author                 string              `yaml:"author" mapstructure:"author" json:"author"`
	Hostname               string              `yaml:"hostname" mapstructure:"hostname" json:"hostname"`
	Menu                   []map[string]string `yaml:"menu" mapstructure:"menu" json:"menu"`
	Twitter                SectionTwitter      `yaml:"twitter" mapstructure:"twitter" json:"twitter"`
	GoogleAnalytics        string              `yaml:"google_analytics" mapstructure:"google_analytics" json:"google_analytics"`
	GoogleSiteVerification string              `yaml:"google_site_verification" mapstructure:"google_site_verification" json:"google_site_verification"`
	BaiduSiteVerification  string              `yaml:"baidu_site_verification" mapstructure:"baidu_site_verification" json:"baidu_site_verification"`
	BaiduAnalytics         string              `yaml:"baidu_analytics" mapstructure:"baidu_analytics" json:"baidu_analytics"`
	GoogleAdsense          string              `yaml:"google_adsense" mapstructure:"google_adsense" json:"google_adsense"`
	Disqus                 SectionDisqus       `yaml:"disqus" mapstructure:"disqus" json:"disqus"`
	Fancybox               bool                `yaml:"fancybox" mapstructure:"fancybox" json:"fancybox"`
	Sponsors               []string            `yaml:"sponsors" mapstructure:"sponsors" json:"sponsors"`
	PayAmount              float64             `yaml:"pay_amount" mapstructure:"pay_amount" json:"pay_amount"`
	CrossedAmount          float64             `yaml:"crossed_amount" mapstructure:"crossed_amount" json:"crossed_amount"`
}

type SectionTwitter struct {
	Image string `yaml:"image" mapstructure:"image" json:"image"`
}

type SectionDisqus struct {
	Shortname string `yaml:"shortname" mapstructure:"shortname" json:"shortname"`
	Lazyload  bool   `yaml:"lazyload" mapstructure:"lazyload" json:"lazyload"`
	Site      string `yaml:"site" mapstructure:"site" json:"site"`
	API       string `yaml:"api" mapstructure:"api" json:"api"`
	Mode      int    `yaml:"mode" mapstructure:"mode" json:"mode"`
	Badge     string `yaml:"badge" mapstructure:"badge" json:"badge"`
	Timeout   int    `yaml:"timeout" mapstructure:"timeout" json:"timeout"`
	Apihost   string `yaml:"apihost" mapstructure:"apihost" json:"apihost"`
	Count     bool   `yaml:"count" mapstructure:"count" json:"count"`
	EmojiPath string `yaml:"emoji_path" mapstructure:"emoji_path" json:"emoji_path"`
}

func setDefault() {
}

var conf *ConfYaml

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath ...string) (*ConfYaml, error) {
	conf = &ConfYaml{}

	// load default values
	setDefault()

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()        // read in environment variables that match
	viper.SetEnvPrefix("TOOLS") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if len(confPath) > 0 && confPath[0] != "" {
		content, err := ioutil.ReadFile(confPath[0])
		if err != nil {
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
		fmt.Println("Using config file:", confPath)
	} else {
		// Search config in home directory with name ".gorush" (without extension).
		viper.AddConfigPath("/etc/tools.liu.app/")
		viper.AddConfigPath("$HOME/.tools.liu.app")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
			// load default config
			return conf, err
		}
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}

	if conf.Core.WorkerNum == int64(0) {
		conf.Core.WorkerNum = int64(runtime.NumCPU())
	}

	return conf, nil
}

func GetConfig() (*ConfYaml, error) {
	return conf, nil
}
