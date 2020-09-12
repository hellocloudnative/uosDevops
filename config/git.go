package config

import "github.com/spf13/viper"

type  Git struct {
	RepositoryType string
	RepositoryUrl  string
	Branch         string
	UserName       string
	AccessToken       string
}
func InitGit(cfg *viper.Viper) *Git {

	return &Git{
		RepositoryType: cfg.GetString("repositoryType"),
		RepositoryUrl:cfg.GetString("repositoryUrl"),
		Branch:cfg.GetString("branch"),
		UserName:cfg.GetString("userName"),
		AccessToken:cfg.GetString("accessToken"),
	}
}

var GitConfig = new(Git)