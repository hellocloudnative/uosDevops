package utils

import (
	"fmt"
	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
	"uosDevops/config"
)
//检查错误
func CheckIfError(err error, message string) {
	if err != nil {
		config.Logger.Error(message)
		os.Exit(1)
	}
}

//检查debug
func CheckIfDeBug(err error, message string) {
	if err != nil {
		config.Logger.Debug(message)
	}
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

//传包
func SendPackage(host string, packName string) {
	//本地包
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey(fmt.Sprintf("%s",config.NodeConfig.User), fmt.Sprintf("%s",config.NodeConfig.SshKeyDir), ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod
	// Create a new SCP client
	client := scp.NewClient(fmt.Sprintf("%s",host+":"+config.NodeConfig.Port), &clientConfig)
	// Connect to the remote server
	err := client.Connect()
	CheckIfError(err,fmt.Sprintf("[%s:%s] Couldn't establish a connection to the remote server",config.NodeConfig.Host,config.NodeConfig.Port))

	config.Logger.Info(fmt.Sprintf("[%s:%s] The ssh connection is successful      [ok]",config.NodeConfig.Host,config.NodeConfig.Port))

	// Open a file
	f, err := os.Open(packName)
	config.Logger.Info(fmt.Sprintf("[%s:%s]  Read this  file    [ok]",config.NodeConfig.Host,config.NodeConfig.Port))
	CheckIfError(err,fmt.Sprintf("[%s:%s] open file is faild :[%s]",config.NodeConfig.Host,config.NodeConfig.Port,packName))
	// Close client connection after the file has been copied
	defer client.Close()
	// Close the file after it has been copied
	defer f.Close()
	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)
	remoteFileName :=path.Base(packName)
	err = client.CopyFile(f, fmt.Sprintf("%s",filepath.Join(config.NodeConfig.RemoteDir,remoteFileName)), "0655")
	CheckIfError(err,fmt.Sprintf("%s,Error while copying file",err))
}
func readFile(name string)(string){
	content,err :=ioutil.ReadFile(name)
	CheckIfError(err,fmt.Sprintf("[globals] read file err is : %s", err))
	return string(content)
}
func AddReformat(host string,port string)(string){
	if strings.Index(host,":") == -1{
		host = fmt.Sprintf("%s:%s",host,port)
	}
	return host
}

func sshAuthMethod( pkFile string) ssh.AuthMethod {
	var am ssh.AuthMethod
		pkData:=readFile(pkFile)
		pk,_ := ssh.ParsePrivateKey([]byte(pkData))
		am = ssh.PublicKeys(pk)
	return am
}
//ssh connect
func Connect(user,pkFile,host string,port string) (*ssh.Session, error){
	auth :=[]ssh.AuthMethod{sshAuthMethod(pkFile)}
	config :=ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256"},
	}
	clientConfig :=&ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: time.Duration(1) * time.Minute,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) (err error) {
			return nil
		},
	}
	addr := AddReformat(host,port)
	client,err :=ssh.Dial("tcp",addr,clientConfig)
	if err!=nil{
		return nil,err
	}
	session,err :=client.NewSession()
	if err!=nil{
		return nil,err
	}
	modes :=ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err:=session.RequestPty("xterm", 80, 40, modes);err!=nil{
		return nil,err
	}
	return session,nil
}

//shell命令
func Cmdout(host string,port string,cmd string)([]byte){
	session,err :=Connect(config.NodeConfig.User,config.NodeConfig.SshKeyDir,host,port)
	defer func() {
		if r := recover(); r != nil {
			config.Logger.Error(fmt.Sprintf("[%s]Error create ssh session failed,%s", host, err))
		}
	}()
	if err !=nil{
		panic(1)
	}
	defer session.Close()
	b,err :=session.CombinedOutput(cmd)
	config.Logger.Info(fmt.Sprintf("[%s]command result running ...: %s",host,b))
	config.Logger.Info("command is completed     [ok]")
	defer func() {
		if r := recover(); r != nil {
			config.Logger.Error(fmt.Sprintf("[%s]Error exec command failed: %s", host, err))
			os.Exit(1)
		}
	}()
	if err !=nil{
		panic(1)

	}
	return b
}
func Execute(bash string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s",bash))

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		config.Logger.Error(fmt.Sprintf("Error starting command: %s......", err.Error()))
		return err
	}

	go asyncLog(stdout)
	go asyncLog(stderr)

	if err := cmd.Wait(); err != nil {
		config.Logger.Error(fmt.Sprintf("Error starting command: %s......", err.Error()))
		return err
	}
	return nil
}

//shell 实时输出日志
func asyncLog(reader io.ReadCloser) error {
	cache := "" //缓存不足一行的日志信息
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err!=io.EOF{
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
			config.Logger.Info(fmt.Sprintf("%s%s", cache, line))
			cache = s[len(s)-1]
		}
	}
	return nil
}
