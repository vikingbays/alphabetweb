// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package utils

import (
	"alphabet"
	"alphabet/core/utils"
	"alphabet/env"
	"alphabet/log4go/message"
	"alphabet/web"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Run() {
	ResolveHelperCmdArgs()
	//	fmt.Println("alphabet/cmd/uitls/cmd.go[line29]: ", helperCmdInfoMap)
	//	fmt.Println("alphabet/cmd/uitls/cmd.go[line30]: ", len(helperCmdInfoMapToCmdArgsForNest()))
	//	fmt.Println("alphabet/cmd/uitls/cmd.go[line31]: ", helperCmdInfoMapToCmdArgsForNest())
	//	fmt.Println("alphabet/cmd/uitls/cmd.go[line32]: ", os.Args) //

	if helperCmdInfoMap[1] == nil {
		Helper()
	} else {
		if helperCmdInfoMap[1][0].commandName == "-start" || helperCmdInfoMap[1][0].commandName == "-stop" ||
			helperCmdInfoMap[1][0].commandName == "-genmain" {
			cmdName := helperCmdInfoMap[1][0].commandArgsInfo[0]
			cmdParams := helperCmdInfoMap[1][0].commandArgsInfo[1]

			err := validateCmd(cmdName, cmdParams)
			if err != nil {
				fmt.Println(err)
			} else {
				DoCmd(cmdName, cmdParams)
			}
		} else if helperCmdInfoMap[1][0].commandName == "-build" || helperCmdInfoMap[1][0].commandName == "-build_for_debug" {
			cmdName := helperCmdInfoMap[1][0].commandArgsInfo[0]
			cmdParams := helperCmdInfoMap[1][0].commandArgsInfo[1:]
			if helperCmdInfoMap[1][0].commandArgsInfo[3] == "" {
				cmdParams = cmdParams[:2]
			}
			err := validateCmd(cmdName, cmdParams...)
			if err != nil {
				fmt.Println(err)
			} else {
				DoCmd(cmdName, cmdParams...)
			}
		} else {
			Helper()
		}
	}

	/*
		lenOfArgs := len(os.Args)
			if lenOfArgs < 3 {
				helper()
			} else if lenOfArgs == 3 { // 运行命令，例如： alphabetweb/src/alphabet/cmd/abserver.go -start /xxxx/sample_octopus_frontend_service
				cmdType := os.Args[1]
				cmdParams := os.Args[2]

				err := validateCmd(cmdType, cmdParams)
				if err != nil {
					fmt.Println(err)
				} else {
					DoCmd(cmdType, cmdParams)
				}
			} else {
				cmdType := os.Args[1]
				if (cmdType == "-build" || cmdType == "-build_for_debug") && (lenOfArgs == 4 || lenOfArgs == 5) {
					// 运行命令，例如：    alphabetweb/src/alphabet/cmd/abserver.go -build /xxxx/sample_octopus_frontend_service   /xxxx/target windows
					//  例如： alphabetweb/src/alphabet/cmd/abserver.go -build_for_debug /xxxx/sample_octopus_frontend_service   /xxxx/target windows
					cmdParams := os.Args[2:]
					err := validateCmd(cmdType, cmdParams...)
					if err != nil {
						fmt.Println(err)
					} else {
						DoCmd(cmdType, cmdParams...)
					}
				} else {
					helper()
				}
			}
	*/
}

//
// 生成具体的代码
//
//
func Gen(projectPath string, execProjectPath string, pidPath string, applicationName string, isBuildBin bool, outPath string) string {
	env.Env_Project_Root = projectPath
	srcDirs, err0 := ioutil.ReadDir(projectPath + "/src")
	if err0 != nil {
		//log4go.ErrorLog(message.ERR_CMD0_39063, projectPath+"/src", err0.Error())
		fmt.Println(fmt.Sprintf(message.ERR_CMD0_39063.String(), projectPath+"/src", err0.Error()))
	} else {
		configArgsTag, configArgsStr := GetConfigKeyOrUrlFromCmd()
		if configArgsTag == 1 {
			DoSysconfigPathForKey(configArgsStr)
			//env.SetSysconfigPath(configArgsStr)
		} else if configArgsTag == 2 {
			DoSysconfigPathForUrl(configArgsStr)
		}

		// 根据 /src/[appname]/webconfig.toml文件判断appname是否存在，然后存储到 env.Env_Project_Resource_Apps 中。
		appNamesTemp := make([]string, 0, 1)
		for _, dir := range srcDirs {
			dirName := dir.Name()
			if !strings.HasPrefix(dirName, ".") && dir.IsDir() {
				if utils.ExistFile(projectPath + "/src/" + dirName + "/" + env.GetSysconfigPath() + env.Env_Project_Resource_WebConfig) {
					appNamesTemp = append(appNamesTemp, dir.Name())
				}
			}
		}
		env.Env_Project_Resource_Apps = appNamesTemp
		fmt.Println("alphabet/cmd/uitls/cmd.go[line128]:", env.Env_Project_Resource_Apps, "   from /src/[appname]/", env.Env_Project_Resource_WebConfig)
		alphabet.InitGenApplication(projectPath)

		rootAppNames := make([]string, 0, 1)
		for _, rootAppName := range env.Env_Project_Resource_Apps { //判断定义的appname 是否存在
			_, err := ioutil.ReadDir(env.Env_Project_Resource + "/" + rootAppName)
			if err != nil {
				//log4go.ErrorLog(message.ERR_CMD0_39063, env.Env_Project_Resource+"/"+rootAppName, err)
				fmt.Println(fmt.Sprintf(message.ERR_CMD0_39063.String(), env.Env_Project_Resource+"/"+rootAppName, err))
			} else {
				rootAppNames = append(rootAppNames, rootAppName)
			}
		}

		webappInfo := web.NewWebappInfo(rootAppNames)

		for _, rootAppName := range env.Env_Project_Resource_Apps {
			dirs, err := ioutil.ReadDir(env.Env_Project_Resource + "/" + rootAppName)
			if err != nil {
				//log4go.ErrorLog(message.ERR_CMD0_39063, env.Env_Project_Resource+"/"+rootAppName, err)
				fmt.Println(fmt.Sprintf(message.ERR_CMD0_39063.String(), env.Env_Project_Resource+"/"+rootAppName, err))
			} else {

				for _, dir := range dirs {
					if dir.IsDir() {
						routeFileName := env.Env_Project_Resource + "/" + rootAppName + "/" + dir.Name() + "/" + env.Env_Project_Resource_Apps_WebConfig
						if utils.ExistFile(routeFileName) {
							webappInfo.AddWebapp(routeFileName, dir.Name(), rootAppName)
						}
					}
				}
			}
		}

		projectName := getProjectName(projectPath)

		genFile := ""
		if outPath == "" {
			os.Mkdir(os.TempDir()+"/"+projectName, 0777)
			genFile = os.TempDir() + "/" + projectName + "/" + applicationName
		} else {
			genFile = outPath + "/" + applicationName
		}

		ioutil.WriteFile(genFile, []byte(webappInfo.GenerateWebMainApplication(projectPath, execProjectPath, pidPath, env.Env_Project_Cpunum, env.Env_Web_Mode_IsProduct, isBuildBin)), 0777)
		fmt.Printf("alphabet/cmd/uitls/cmd.go[line173]: path:%s \n", genFile)
		return genFile

	}
	return ""
}

func getProjectName(projectPath string) string {
	pathNames := strings.Split(projectPath, "/")
	projectName := ""
	lenPathNames := len(pathNames)
	for i := 1; i <= lenPathNames; i++ {
		projectName = pathNames[lenPathNames-i]
		if projectName != "" {
			break
		}
	}
	return projectName
}

func DoCmd(cmdName string, params ...string) (err error) {
	if cmdName == "-start" {
		goPathEnv := os.Getenv("GOPATH")
		if runtime.GOOS == "windows" {
			os.Setenv("GOPATH", goPathEnv+";"+params[0])
		} else {
			os.Setenv("GOPATH", goPathEnv+":"+params[0])
		}

		filename := Gen(params[0], params[0], "", "main.go", false, "")
		if utils.ExistFile(filename) {
			//cmd := exec.Command("go", "run", filename, ">", env.Env_Project_Root+"/logs/console.log")

			argParamsNest := helperCmdInfoMapToCmdArgsForNest()
			argsNest := make([]string, 0, 1)
			//argsNest = append(argsNest, "go")
			argsNest = append(argsNest, "run")
			argsNest = append(argsNest, filename)
			for _, argParam := range argParamsNest {
				argsNest = append(argsNest, argParam)
			}
			cmd := exec.Command("go", argsNest...)

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println(err)
				return err
			}

			errStart := cmd.Start()
			if errStart != nil {
				fmt.Printf("start_err: %s  \n", errStart.Error())
			}
			reader := bufio.NewReader(stdout)
			for {
				line, err2 := reader.ReadString('\n')
				if err2 != nil || io.EOF == err2 {
					if err2 != nil {
						fmt.Println(">>>>" + err2.Error())
					}
					break
				}
				fmt.Printf(line)
			}
			cmd.Wait()

		}
	} else if cmdName == "-stop" {
		pidfile := os.TempDir() + "/" + getProjectName(params[0]) + "/pid"
		bytes, err := ioutil.ReadFile(pidfile)
		if err == nil {
			pid, err1 := strconv.Atoi(string(bytes))
			if err1 == nil {
				//log4go.InfoLog(message.INF_CMD0_09014, pid)
				fmt.Println(fmt.Sprintf(message.INF_CMD0_09014.String(), pid))

				if runtime.GOOS == "windows" {
					cmd := exec.Command("tskill", strconv.Itoa(pid))
					cmd.Run()
				} else {
					cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
					cmd.Run()
				}

			} else {
				//log4go.ErrorLog(message.ERR_CMD0_39064, pidfile, err1.Error())
				fmt.Println(fmt.Sprintf(message.ERR_CMD0_39064.String(), pidfile, err1.Error()))
			}
		} else {
			//log4go.ErrorLog(message.ERR_CMD0_39065, pidfile, err.Error())
			fmt.Println(fmt.Sprintf(message.ERR_CMD0_39065.String(), pidfile, err.Error()))
		}

		//pid:=
	} else if cmdName == "-genmain" {
		os.Setenv("GOPATH", params[0])
		filename := Gen(params[0], params[0], "", "main.go", false, "")
		//log4go.InfoLog(message.INF_CMD0_09015, filename)
		fmt.Println(fmt.Sprintf(message.INF_CMD0_09015.String(), filename))
	} else if cmdName == "-build" || cmdName == "-build_for_debug" {

		goPathEnv := os.Getenv("GOPATH")
		if runtime.GOOS == "windows" {
			os.Setenv("GOPATH", goPathEnv+";"+params[0])
		} else {
			os.Setenv("GOPATH", goPathEnv+":"+params[0])
		}

		relProjectPath := "../../" + getProjectName(params[1])

		filepath.Walk(params[0], func(path string, fi os.FileInfo, err error) error {
			if nil == fi {
				return err
			}
			name := fi.Name()

			relPath := path
			if runtime.GOOS == "windows" {
				relPath = strings.Replace(relPath, "\\", "/", -1)
			}
			startPos := strings.Index(relPath, params[0])
			if startPos != -1 {
				relPath = relPath[(len(params[0]) + startPos):]
			}

			if !strings.HasPrefix(name, ".") && !strings.HasSuffix(name, ".go") &&
				!strings.HasSuffix(name, ".sublime-project") && !strings.HasSuffix(name, ".sublime-workspace") &&
				!strings.Contains(name, ".git") && !strings.Contains(relPath, ".git") {
				if fi.IsDir() {
					err := os.Mkdir(params[1]+"/"+relPath, 0777)
					if err != nil {
						fmt.Printf("[err] Create Folder  file [%s] is Error . ErrInfo: %s \n", relPath, err)
					} else {
						fmt.Printf("[ok.] Create Folder  file [%s] is Ok . \n", relPath)
					}

				} else {
					_, err := utils.CopyFile(params[0]+"/"+relPath, params[1]+"/"+relPath)
					if err != nil {
						fmt.Printf("[err] Copy  **File  file [%s] is Error . ErrInfo: %s \n", relPath, err)
					} else {
						fmt.Printf("[ok.] Copy  **File  file [%s] is Ok . \n", relPath)
					}
				}
			}

			return nil
		})

		//filename := gen(params[0], "main.go")
		filename := Gen(params[0], relProjectPath, relProjectPath+"/bin/pid", "main.go", true, "")
		//log4go.InfoLog(message.INF_CMD0_09015, filename)
		fmt.Println(fmt.Sprintf(message.INF_CMD0_09015.String(), filename))

		//projectName := getProjectName(params[0])

		os.Mkdir(params[1], 0777)
		os.Mkdir(params[1]+"/bin", 0777)
		os.Mkdir(params[1]+"/src", 0777)

		fmt.Printf("生成可执行项目路径：%s . \n", params[1])

		if len(params) == 3 {
			//第三个参数的形式是： windows/arm , darwin/amd64 等，如果是多个用分号分割，例如：  windows/arm:darwin/amd64
			{
				if strings.Contains(params[2], ":") {
					goosArchs := strings.Split(params[2], ":")
					for _, goosArch := range goosArchs {
						pos := strings.Index(goosArch, "/")
						if pos != -1 {
							os.Setenv("GOOS", goosArch[0:pos])
							os.Setenv("GOARCH", goosArch[(pos+1):])

							goosStr := os.Getenv("GOOS")
							goarchStr := os.Getenv("GOARCH")
							fmt.Printf("GOOS=%s ，GOARCH=%s . \n", goosStr, goarchStr)
							exeAppname := "server"
							if goosStr != "" {
								exeAppname = exeAppname + "_" + goosStr
							}
							if goarchStr != "" {
								exeAppname = exeAppname + "_" + goarchStr
							}
							if goosStr == "windows" {
								exeAppname = exeAppname + ".exe"
							}
							buildServerCmd(filename, params[1]+"/bin/"+exeAppname, (cmdName == "-build_for_debug"))
							fmt.Printf("生成可执行程序在：%s . \n", "bin/"+exeAppname)
			}
					}
				} else {
					if params[2] != "" {
						pos := strings.Index(params[2], "/")
						if pos != -1 {
							os.Setenv("GOOS", params[2][0:pos])
							os.Setenv("GOARCH", params[2][(pos+1):])

				}
			}
					goosStr := os.Getenv("GOOS")
					goarchStr := os.Getenv("GOARCH")
					fmt.Printf("GOOS=%s ，GOARCH=%s . \n", goosStr, goarchStr)
					exeAppname := "server"
					if goosStr != "" {
						exeAppname = exeAppname + "_" + goosStr
					}
					if goarchStr != "" {
						exeAppname = exeAppname + "_" + goarchStr
		}
					if goosStr == "windows" {
						exeAppname = exeAppname + ".exe"
					}
					buildServerCmd(filename, params[1]+"/bin/"+exeAppname, (cmdName == "-build_for_debug"))
					fmt.Printf("生成可执行程序在：%s . \n", "bin/"+exeAppname)
				}

			}

		}

	}

	return

}

func buildServerCmd(mainFilename string, buildOutFile string, buildForDebug bool) {

	var cmd *exec.Cmd
	if buildForDebug {
		cmd = exec.Command("go", "build", "-gcflags", "-N -l", "-o", buildOutFile, mainFilename)
	} else { //-ldflags "-s -w"
		cmd = exec.Command("go", "build", "-ldflags", "-s -w", "-o", buildOutFile, mainFilename)
	}

	envs := os.Environ()
	inEnvs := make([]string, 0, 1)
	for _, env := range envs {
		if strings.HasPrefix(env, "AB_") && len(env) > 4 {
			inEnvs = append(inEnvs, env[3:])
	} else {
			inEnvs = append(inEnvs, env)
	}
	}
	cmd.Env = inEnvs
	fmt.Printf("cmd.Env=%s \n", cmd.Env)
	fmt.Printf("building.... buildOutFile=%s , mainFilename=%s \n", buildOutFile, mainFilename)
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Running error : ", err)
	}
	fmt.Println(string(bytes))
	/*
		err2 := cmd.Run()
		if err2 != nil {
			fmt.Println(err2)
	}
	*/

}

func validateCmd(cmdName string, params ...string) (err error) {
	err = nil
	if cmdName == "-start" || cmdName == "-stop" || cmdName == "-genmain" {
		if len(params) == 1 && utils.IsFolder(params[0]) {

		} else if len(params) == 1 {
			err = fmt.Errorf("该命令 [%s] 传入的参数不正确，需要传入的是工程根目录文件夹路径。（该路径无效：%s）", cmdName, params[0])
		} else {
			err = fmt.Errorf("该命令 [%s] 只接收1个参数。（当前传入参数个数为%d，不正确）", cmdName, len(params))

		}
	} else if cmdName == "-build" || cmdName == "-build_for_debug" { // 构建二进制程序，接入参数应该有2个或者3个
		if (len(params) == 2 || len(params) == 3) && utils.IsFolder(params[0]) {

		} else if len(params) == 2 || len(params) == 3 {
			err = fmt.Errorf("该命令 [%s] 传入的参数不正确，需要传入的是工程根目录文件夹路径。（该路径无效：%s）", cmdName, params[0])
		} else {
			err = fmt.Errorf("该命令 [%s] 接收参数不正确，一般是2个或者3个参数。（当前传入参数个数为%d，不正确）", cmdName, len(params))

		}
	} else {
		err = fmt.Errorf("指定的命令信息不正确，没有该命令：%s 。", cmdName)
	}
	return
}

type HelperCmdStruct struct {
	commandType     int      // 命令种类，现在有4中，分别对应： 1，2，3，4
	commandName     string   // 命令明
	commandParamNum int      // 命令参数个数
	commandArgsInfo []string // 命令行信息
}

// 定义命令的标准格式
var helperCmdFormatMap map[string]*HelperCmdStruct = map[string]*HelperCmdStruct{
	"-start":           &HelperCmdStruct{1, "-start", 1, []string{}},
	"-stop":            &HelperCmdStruct{1, "-stop", 1, []string{}},
	"-genmain":         &HelperCmdStruct{1, "-genmain", 1, []string{}},
	"-build":           &HelperCmdStruct{1, "-build", 3, []string{}},
	"-build_for_debug": &HelperCmdStruct{1, "-build_for_debug", 3, []string{}},
	"-help":            &HelperCmdStruct{1, "-help", 0, []string{}},
	"-config_key":      &HelperCmdStruct{2, "-config_key", 1, []string{}},
	"-config_url":      &HelperCmdStruct{2, "-config_url", 1, []string{}},
	"-env_":            &HelperCmdStruct{3, "-env_", 1, []string{}},
	"-appsname":        &HelperCmdStruct{4, "-appsname", 1, []string{}},
}

// 存储实际传入的命令信息,key 命令类型： 1，2，3，4
var helperCmdInfoMap map[int][]HelperCmdStruct = nil

func GetHelperCmdInfoMap() map[int][]HelperCmdStruct {
	return helperCmdInfoMap
}

func IsHelpFromCmd() bool {
	if helperCmdInfoMap != nil {
		if helperCmdInfoMap[1] != nil {
			for _, helperCmd := range helperCmdInfoMap[1] {
				if helperCmd.commandName == "-help" {
					return true
				}
			}
		}
	}
	return false
}

/*
 * 获取配置目录信息，返回值是（int, string），如果第一个值是 0 ，表示没有定义配置目录，用默认的。如果第一个值是1，表示使用configkey。如果第一个值是2，表示使用configUrl。
 *
 */
func GetConfigKeyOrUrlFromCmd() (int, string) {
	if helperCmdInfoMap != nil {
		if helperCmdInfoMap[2] != nil {
			if len(helperCmdInfoMap[2]) > 0 {
				if helperCmdInfoMap[2][0].commandName == "-config_key" {
					return 1, helperCmdInfoMap[2][0].commandArgsInfo[1]
				} else if helperCmdInfoMap[2][0].commandName == "-config_url" {
					return 2, helperCmdInfoMap[2][0].commandArgsInfo[1]
				}
			}
		}
	}
	return 0, ""
}

func GetAppnamesFromCmd() []string {
	if helperCmdInfoMap != nil {
		if helperCmdInfoMap[4] != nil {
			if len(helperCmdInfoMap[4]) > 0 {
				if len(helperCmdInfoMap[4][0].commandArgsInfo) == 2 {
					return strings.Split(helperCmdInfoMap[4][0].commandArgsInfo[1], ",")
				}
			}
		}
	}
	return nil
}

func GetEnvsFromCmd() map[string]string {
	kv := make(map[string]string)
	if helperCmdInfoMap != nil {
		if helperCmdInfoMap[3] != nil {
			for _, helperCmd := range helperCmdInfoMap[3] {
				if len(helperCmd.commandArgsInfo) == 2 && strings.HasPrefix(helperCmd.commandArgsInfo[0], "-env_") && len(helperCmd.commandArgsInfo[0]) > 5 {
					kv[helperCmd.commandArgsInfo[0][5:]] = helperCmd.commandArgsInfo[1]
				}
			}
		}
	}
	return kv
}

// 解析参数信息，获取每个命令的参数
func ResolveHelperCmdArgs() {
	helperCmdInfoMap = make(map[int][]HelperCmdStruct)
	lenOfArgs := len(os.Args)
	for i := 1; i < lenOfArgs; i++ {
		key1 := os.Args[i]
		if strings.HasPrefix(os.Args[i], "-env_") {
			key1 = "-env_"
		}
		helperCmdFormat := helperCmdFormatMap[key1]
		if helperCmdFormat != nil {
			if helperCmdFormat.commandType == 1 || helperCmdFormat.commandType == 2 || helperCmdFormat.commandType == 4 {
				helperCmdInfoMap[helperCmdFormat.commandType] = make([]HelperCmdStruct, 0, 1)
			} else if helperCmdFormat.commandType == 3 { // 如果有多个的
				if helperCmdInfoMap[helperCmdFormat.commandType] == nil {
					helperCmdInfoMap[helperCmdFormat.commandType] = make([]HelperCmdStruct, 0, 1)
				}
			}

			if i+helperCmdFormat.commandParamNum >= lenOfArgs { // 参数不够
				fmt.Printf("Error: [%s] 的参数个数不正确。\n", helperCmdFormat.commandName)
				helperCmdInfoMap = nil
				return
			} else {
				hcs := HelperCmdStruct{commandName: key1, commandType: helperCmdFormat.commandType,
					commandParamNum: helperCmdFormat.commandParamNum}
				hcs.commandArgsInfo = make([]string, 0, 1)
				hcs.commandArgsInfo = append(hcs.commandArgsInfo, os.Args[i])
				for j := 1; j <= helperCmdFormat.commandParamNum; j++ {
					helperCmdFormat2 := helperCmdFormatMap[os.Args[i+j]]
					if helperCmdFormat2 != nil { // 说明参数个数不够
						fmt.Printf("Error: [%s] 的第[%d]参数不正确，实际传入数据是[%s]。\n", helperCmdFormat.commandName, j, os.Args[i+j])
						helperCmdInfoMap = nil
						return
					}
					hcs.commandArgsInfo = append(hcs.commandArgsInfo, os.Args[i+j])
				}

				helperCmdInfoMap[helperCmdFormat.commandType] = append(helperCmdInfoMap[helperCmdFormat.commandType], hcs)
				i = i + helperCmdFormat.commandParamNum
			}

		} else {
			fmt.Printf("Error: 当前参数[%s]不是命令。 \n", os.Args[i])
			helperCmdInfoMap = nil
			return
		}
	}
}

// 基于当前传入的args参数转换成 实际运行程序的参数。
func helperCmdInfoMapToCmdArgsForNest() []string {

	argsNest := make([]string, 0, 1)
	if helperCmdInfoMap[2] != nil {
		for _, helperCmdInfo := range helperCmdInfoMap[2] {
			for _, arg := range helperCmdInfo.commandArgsInfo {
				argsNest = append(argsNest, arg)
			}
		}
	}
	if helperCmdInfoMap[3] != nil {
		for _, helperCmdInfo := range helperCmdInfoMap[3] {
			for _, arg := range helperCmdInfo.commandArgsInfo {
				argsNest = append(argsNest, arg)
			}
		}
	}
	if helperCmdInfoMap[4] != nil {
		for _, helperCmdInfo := range helperCmdInfoMap[4] {
			for _, arg := range helperCmdInfo.commandArgsInfo {
				argsNest = append(argsNest, arg)
			}
		}
	}
	return argsNest

}

func Helper() {
	str := `
参考帮助：

命令：  abserver   [command1]  [command2]  [command3]  [command4]

command说明：(所有文件路径的分隔符都必须使用｀／｀)

  (1)、[command1] 参数一：（只能一个有效）

       -start    string
                   启动web服务。
                   参数1：web工程根路径，源码存放在 src/apps下。

       -stop     string
                   停止web服务。
                   参数1：web工程根路径，源码存放在 src/apps下。

       -genmain  string
                   生成web服务启动代码，代码存储到 “${project}/src/apps”下。
                   参数1：web工程根路径，源码存放在 src/apps下。

       -build    string  string  string
                   生成web服务启动二进制程序。
                   参数1：web工程根路径，源码存放在 src/apps下。
                   参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
                   参数3：平台定义，设置生成运行平台的代码，例如：windows/arm , darwin/amd64 等。
                          如果是多个用分号分割，例如：  windows/arm:darwin/amd64
                          具体平台信息可以通过命令查找：go tool dist list
                          补充：如果是android平台，需要设置 CC、CXX、CGO_ENABLED、GO111MODULE参数，在执行命令前定义(用AB_前缀)，例如：
                          AB_CC=/xxxx/android-ndk-r19c/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang        \
                          AB_CXX=/xxxx/android-ndk-r19c/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang++     \
                          AB_CGO_ENABLED=1 AB_GO111MODULE=off                                                                             \
                          go run /xxxx/alphabetweb/src/alphabet/cmd/abserver.go  -build "..." "..." "android/arm64"

       -build_for_debug    string  string  string
                   生成web服务启动二进制程序，支持gdb进行debug（ 在编译时增加参数 -gcflags "-N -l" ）。
                   参数1：web工程根路径，源码存放在 src/apps下。
                   参数2：输出项目文件夹路径，默认可执行文件在此目录的bin目录下。
                   参数3: 参考-build

       -help
                   展现帮助信息。

  (2)、[command2] 参数二：（只能一个有效）
       -config_key    string
                   设置关键字，例如：“base”，找到对应的配置目录，例如：“sysconfig.base” 。
                   该参数值赋给env.Env_Project_Resource_Sysconfig_Folder_Key 。
                   程序启动时，会加载此目录，例如：src/[apps名称]/sysconfig.base 。
                   参数1: 设置这个key值。
       -config_url    string
                   设置一个url，获取远程配置项信息，实现服务配置集中管理。
                   这个url地址会匹配如下目录：
                           [url]/logconfig.toml
                           [url]/logconfig_http.toml
                           [url]/dbsconfig.toml
                           [url]/cachesconfig.toml
                           [url]/webconfig.toml
                           [url]/ms_server_config.toml
                           [url]/ms_client_config_01.toml
                           [url]/...
                   参数1: url地址。

  (3)、[command3] 参数三：（可以多个有效）
       -env_xxx    string
                   设置环境变量，可以定义多个（使用方法os.Setenv）。其中 xxx 就是当前定义的key值。
                   获取环境变量数据，采用 os.Getenv
                   例如： -env_key1  value1    -env_key2  value2 ，那么就是：key1:value1,key2:value2

  (4)、[command4] 参数四：（只能一个有效）
       -appsname    string
                   设置有效的apps名称，如果不设置表示全部有效，如果设置多个就多个有效，多个apps以逗号分割。
                   例如： -appsname    oct_web,oct_service

`
	fmt.Println(".")
	fmt.Printf(str)

}

func HelperForNest() {
	str := `
参考帮助：

命令：  abserver
    默认参数运行

命令：  abserver   -help
    展现帮助信息。

命令：  abserver   [command1]  [command2]  [command3]
    指定参数运行

command说明：(所有文件路径的分隔符都必须使用｀／｀)

  (1)、[command1] 参数一：（只能一个有效）
       -config_key    string
                   设置关键字，例如：“base”，找到对应的配置目录，例如：“sysconfig.base” 。
                   该参数值赋给env.Env_Project_Resource_Sysconfig_Folder_Key 。
                   程序启动时，会加载此目录，例如：src/[apps名称]/sysconfig.base 。
                   参数1: 设置这个key值。
       -config_url    string
                   设置一个url，获取远程配置项信息，实现服务配置集中管理。
                   这个url地址会匹配如下目录：
                           [url]/logconfig.toml
                           [url]/logconfig_http.toml
                           [url]/dbsconfig.toml
                           [url]/cachesconfig.toml
                           [url]/webconfig.toml
                           [url]/ms_server_config.toml
                           [url]/ms_client_config_01.toml
                           [url]/...
                   参数1: url地址。

  (2)、[command2] 参数二：（可以多个有效）
       -env_xxx    string
                   设置环境变量，可以定义多个（使用方法os.Setenv）。其中 xxx 就是当前定义的key值。
                   获取环境变量数据，采用 os.Getenv
                   例如： -env_key1  value1    -env_key2  value2 ，那么就是：key1:value1,key2:value2

  (3)、[command3] 参数三：（只能一个有效）
       -appsname    string
                   设置有效的apps名称，如果不设置表示全部有效，如果设置多个就多个有效，多个apps以逗号分割。
                   例如： -appsname    oct_web,oct_service

`
	fmt.Println(".")
	fmt.Printf(str)

}
