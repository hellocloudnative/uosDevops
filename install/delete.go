package install

import (
	"fmt"
	"os"
	"path/filepath"
	"uosDevops/config"
	"uosDevops/utils"
)

func Delete(){
	ProjectDir := filepath.Join(config.BaseProjectDir, fmt.Sprintf("%s", config.ProjectName))
	WorkBaseDir := filepath.Join("workspace",fmt.Sprintf("%s", config.ProjectName))
	err:=os.RemoveAll(ProjectDir)
	utils.CheckIfError(err,fmt.Sprintf("%s",err))
	err=os.RemoveAll(WorkBaseDir)
	utils.CheckIfError(err,fmt.Sprintf("%s",err))
}