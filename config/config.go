package config

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var AppConfig Config

var Logger *logrus.Logger

func processError(err error) {
	Logger.Fatalf("Loading config failed with exception: '%s'\n", err)
}

type Config struct {
	ContentPath struct {
		ApiSpecYamlFile                string `yaml:"apiSpecYamlFile"`
		ApiSpecDefaultYamlFile         string `yaml:"apiSpecDefaultYamlFile"`
		ProductLayoutFile              string `yaml:"productLayoutFile"`
		DocumentExplorerDefinitionFile string `yaml:"documentExplorerDefinitionFile"`
		TenantProviderApiFile          string `yaml:"tenantProviderApiFile"`
	} `yaml:"contentpath"`

	GitHub struct {
		GitHubRawContentHost  string `yaml:"gitHubRawContentHost"`
		GitHubSourceOwner     string `yaml:"gitHubSourceOwner"`
		GitHubSourceRepo      string `yaml:"gitHubSourceRepo" envconfig:"GITHUB_TENANT_CONTENT_REPO"`
		GitHubContentBranch   string `yaml:"gitHubContentBranch" envconfig:"GITHUB_TENANT_CONTENT_BRANCH"`
		GitHubUserName        string `yaml:"gitHubUser" envconfig:"GITHUB_TENANT_REPO_USER_NAME"`
		GitHubAuthToken       string `yaml:"gitHubAuthToken" envconfig:"GITHUB_TENANT_REPO_AUTH_TOKEN"`
		GitHubContentFullPath string
	} `yaml:"github"`
}

func ReadFile(cfg *Config) {
	// f, err := os.Open("config.yml")
	f, err := ioutil.ReadFile("resources/config.yml")

	if err != nil {
		processError(err)
	}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		processError(err)
	}
}

func ReadEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

func AddLogFields(logger *logrus.Logger) *logrus.Entry {
	pc, file, _, _ := runtime.Caller(1)

	filename := file[strings.LastIndex(file, "/")+1:]
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	pid := os.Getpid()
	return logger.WithField("file", filename).WithField("function", fn).WithField("processID", pid)
}
