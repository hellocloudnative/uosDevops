package config

import "github.com/spf13/viper"

type Node struct {
	Host   		string
	Port		string
	User        string
	SshKeyDir   string
	BuildExec  	[]string
	Source     	string
	RemoteDir  	string
	DeployExec	[]string
}

func InitNode(cfg *viper.Viper) *Node {

	return &Node{
		Host: 		cfg.GetString("host"),
		Port: 		cfg.GetString("port"),
		User:       cfg.GetString("user"),
		SshKeyDir:  cfg.GetString("sshKeyDir"),
		BuildExec: 	cfg.GetStringSlice("buildExec"),
		Source:    	cfg.GetString("source"),
		RemoteDir:  cfg.GetString("remoteDir"),
		DeployExec:	cfg.GetStringSlice("deployExec"),
	}

}

var NodeConfig = new(Node)