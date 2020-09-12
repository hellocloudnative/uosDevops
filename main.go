package main

import (
	"uosDevops/cmd"
)

/**
  go-git clone all code
  workspace: 需要指定workspace，明确把代码clone到哪里
  url: 要clone的repo url
  referenceName: 可以是分支名称，也可以是tag，也可以是commit id
  auth: 这个是可选的，就要看你有没有授权校验
*/

func main(){

	cmd.Execute()
}
