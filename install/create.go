package install

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"uosDevops/config"
	"uosDevops/utils"
)
func CreateFile(projectDir string,fileName string,template string){
	mask := syscall.Umask(0)    // 改为 0000 八进制
	defer syscall.Umask(mask)   // 改为原来的 umask
	file,err:=os.Create(fmt.Sprintf("%s/%s",projectDir,fileName))
	os.Chmod(projectDir+"/"+fileName, 0766)
	defer file.Close()
	utils.CheckIfError(err,fmt.Sprintf("%s",err))
	_,err=file.Write([]byte(template))
	utils.CheckIfError(err,fmt.Sprintf("%s",err))
}

func Create(){
	ProjectDir := filepath.Join(config.BaseProjectDir,fmt.Sprintf("%s",config.ProjectName))
	err:=os.Mkdir(ProjectDir,os.ModePerm)
	utils.CheckIfError(err,fmt.Sprintf("%s",err))
	CreateFile(ProjectDir,"config.yaml",config.TemplateConfigyaml)
	CreateFile(ProjectDir,"build.sh",config.Templatebuildshell)
	CreateFile(ProjectDir,"deploy.sh",config.Templatedeployshell)

}
