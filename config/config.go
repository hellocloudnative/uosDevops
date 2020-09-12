package config

import go_logger "github.com/phachon/go-logger"

var (
	ProjectName    string
	BaseProjectDir string = "project"
	Logger                = go_logger.NewLogger()
)

const TemplateConfigyaml = string(`projectName: helloworld
git:
  repositoryType: gitlab 
  repositoryUrl:  https://gitlab.deepin.io/container/helloworld.git
  branch: master
  userName: zhangpengxuan
  accessToken: 19YmsZQ2VbsRrVTG8dPB
node:
  host: 47.94.104.204
  port: 22
  user: root
  sshKeyDir: /Users/zhangpengxuan/.ssh/id_rsa
  buildExec: ["build.sh"]
  source: deployment.yaml
  remoteDir: /data
  deployExec: ["deploy.sh"]
`)
const Templatebuildshell = `#!/bin/bash
#The bash file  to  build  for you
`
const Templatedeployshell = `#!/bin/bash
#The bash file  to  deploy  for you
`

func init() {

	_ = Logger.Detach("console")

	// console adapter config
	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true, // Does the text display the color
		JsonFormat: false, // Whether or not formatted into a JSON string
		Format:     "",   // JsonFormat is false, logger message output to console format string
	}
	// add output to the console
	_ = Logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)

	// file adapter config
	fileConfig := &go_logger.FileConfig{
		Filename: "logs/uosDevops.log", // The file name of the logger output, does not exist automatically
		// If you want to separate separate logs into files, configure LevelFileName parameters.
		LevelFileName: map[int]string{
			Logger.LoggerLevel("error"): "logs/error.log", // The error level log is written to the error.log file.
			Logger.LoggerLevel("info"):  "logs/info.log",  // The info level log is written to the info.log file.
			Logger.LoggerLevel("debug"): "logs/debug.log", // The debug level log is written to the debug.log file.
		},
		MaxSize:    1024 * 1024, // File maximum (KB), default 0 is not limited
		MaxLine:    100000,      // The maximum number of lines in the file, the default 0 is not limited
		DateSlice:  "d",         // Cut the document by date, support "Y" (year), "m" (month), "d" (day), "H" (hour), default "no".
		JsonFormat: false,       // Whether the file data is written to JSON formatting
		Format:     "",          // JsonFormat is false, logger message written to file format string
	}
	// add output to the file
	_ = Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
}
