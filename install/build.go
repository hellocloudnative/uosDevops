package install

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	 config2 "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"log"
	"os"
	"path/filepath"
	"uosDevops/config"
	"uosDevops/utils"
)

func Build() {
	ProjectDir := filepath.Join(config.BaseProjectDir, fmt.Sprintf("%s", config.ProjectName))
	config.SetupProject(fmt.Sprintf("%s/config.yaml", ProjectDir))
	WorkBaseDir := filepath.Join("workspace",fmt.Sprintf("%s", config.ProjectName))
	packDir:=filepath.Join(WorkBaseDir,config.NodeConfig.Source)
	deployfile:=config.NodeConfig.DeployExec[0]
	deployDir:=filepath.Join(ProjectDir,deployfile)
	u := &BuildHost{
		Host: config.NodeConfig.Host,
	}
	if(!utils.Exists(ProjectDir)){
		config.Logger.Error(fmt.Sprintf("The project %s is not exit [faild: yes]",config.ProjectName))
		config.Logger.Warning(fmt.Sprintf("Please create your project %s first. ",config.ProjectName))
		config.Logger.Warning(fmt.Sprintf("Please use 'uosDevops create --project %s' ",config.ProjectName))
		os.Exit(1)
	}
	if(utils.Exists(WorkBaseDir)){
		err:=os.RemoveAll(WorkBaseDir)
		utils.CheckIfError(err,fmt.Sprintf("%s",err))
		GitClone()
	}else {
		GitClone()
	}
	config.Logger.Info(fmt.Sprintf("Clone %s successed      [ok]", ProjectDir))
	config.Logger.Info(fmt.Sprintf("Start building the project %s.....", ProjectDir))
	//CI
	_ = utils.Execute(ProjectDir + "/" + "build.sh")
	config.Logger.Info(fmt.Sprintf("End  building the project %s.....", ProjectDir))
	//CD
	if(config.NodeConfig.Source==""){
		config.Logger.Info("the package is empty,please note that")
	}
	u.SendPackage(packDir)
	config.Logger.Info(fmt.Sprintf("[%s:%s] The  Deploy package %s for send  is    successful      [ok]", u.Host,config.NodeConfig.Port,deployfile))
	u.SendPackage(deployDir)
	config.Logger.Info(fmt.Sprintf("[%s:%s] The  Compressed package %s for send  is    successful      [ok]", u.Host,config.NodeConfig.Port,config.NodeConfig.Source))
	cmd:=fmt.Sprintf(`source /etc/profile &&bash -x %s/deploy.sh`,config.NodeConfig.RemoteDir)
	utils.Cmdout(config.NodeConfig.Host,config.NodeConfig.Port,cmd)
	config.Logger.Info(fmt.Sprintf("[%s:%s]  Deploy the project %s is    successful      [ok]", u.Host,config.NodeConfig.Port,config.ProjectName))

}

type BuildHost struct {
	Host      string
}


//SendPackage
func (u  *BuildHost) SendPackage(packName string) {
	utils.SendPackage(u.Host, packName)
}

func GitClone() {
	var auth = &http.BasicAuth{
		Username: config.GitConfig.UserName,
		Password: config.GitConfig.AccessToken,
	}
	directory:= fmt.Sprintf("workspace/%s", config.ProjectName)
	var remote = git.NewRemote(
		filesystem.NewStorage(
			osfs.New(directory),
			cache.NewObjectLRUDefault(),
		),
		&config2.RemoteConfig{
			URLs: []string{config.GitConfig.RepositoryUrl},
		},
	)

	refs, err := remote.List(&git.ListOptions{
		Auth: auth,
	})
	utils.CheckIfError(err,"")

	var ref = plumbing.NewBranchReferenceName(config.GitConfig.Branch)
	if !exists(ref, refs) {
		ref = plumbing.NewTagReferenceName(config.GitConfig.Branch)
		if !exists(ref, refs) {
			config.Logger.Error(config.GitConfig.Branch + " not found   [error]")
		}
	}




	r, err := git.PlainClone( directory,false, &git.CloneOptions{
		URL:               config.GitConfig.RepositoryUrl,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
		ReferenceName:     ref,
	},

)
	utils.CheckIfError(err, fmt.Sprintf("\"%s \", \"07\", \"It is be Cloned: %s\"", err, "yes"))
	head, err := r.Head()
	utils.CheckIfError(err,"")
	config.Logger.Info(fmt.Sprintf("Clone branch is %s",head.Strings()[0]))
	config.Logger.Info(fmt.Sprintf("Clone branch hash is %s",head.Strings()[1]))


}

func exists(ref plumbing.ReferenceName, refs []*plumbing.Reference) bool {
	for _, ref2 := range refs {
		if ref.String() == ref2.Name().String() {
			return true
		}
	}
	return false
}

func handlerr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
