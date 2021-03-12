package config

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"os"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

type Config struct {
	ContentPath struct {
		ApiSpecYamlFile       string `yaml:"tenantProviderApiFile"`
		ProductLayoutFile     string `yaml:"productLayoutFile"`
		TenantProviderApiFile string `yaml:"apiSpecYamlFile"`
	} `yaml:"contentpath"`

	GitHub struct {
		GitHubRawContentHost  string `yaml:"gitHubRawContentHost"`
		GitHubSourceOwner     string `yaml:"gitHubSourceOwner"`
		GitHubSourceRepo      string `yaml:"gitHubSourceRepo" envconfig:"GITHUB_SOURCE_REPO"`
		GitHubContentBranch   string `yaml:"gitHubContentBranch" envconfig:"GITHUB_CONTENT_BRANCH"`
		GitHubUserName        string `yaml:"gitHubUser" envconfig:"GITHUB_USER_NAME"`
		GitHubAuthToken       string `yaml:"gitHubAuthToken" envconfig:"GITHUB_AUTH_TOKEN"`
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
