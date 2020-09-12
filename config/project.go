package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

//git配置项
var cfgGit  *viper.Viper

//node配置项

var cfgNode *viper.Viper

//载入配置文件
func SetupProject(path string)  {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		Logger.Error(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}
	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		Logger.Error(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}
	cfgGit = viper.Sub("git")
	if cfgGit == nil {
		panic("No found git in the configuration")
	}
	GitConfig = InitGit(cfgGit)

	cfgNode = viper.Sub("node")
	if cfgNode == nil {
		panic("No found node in the configuration")
	}
	NodeConfig= InitNode(cfgNode)
}