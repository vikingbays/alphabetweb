// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package web

import (
	"alphabet/core/toml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

//对应的结构是：
//    [[action]]
//    Action1Run1="action1run1"
//其中： name="app1.Action1Run1"  ,   webpath="/app1/action1run1"
//    如果非.开始，匹配当前的appname，假设appname=app1，那么"app1.Action1Run1"
//    如果非／开始，匹配当前的appname，假设appname=app1，那么"/app1/action1run1"
//
type ActionInfo struct {
	Name        string
	Webpath     string
	Appname     string
	RootAppname string
}

//对应的结构是：
//    [[gml]]
//    action1_run1_page1="view/action1_run1_page1.gml"
//其中： name="action1_run1_page1"  ,   webpath="/app1/view/action1_run1_page1.gml"
//    如果非／开始，匹配当前的appname，假设appname=app1，那么"/app1/view/action1_run1_page1.gml"
//
type GmlInfo struct {
	Name        string
	Webpath     string
	Appname     string
	RootAppname string
}

type StaticInfo struct {
	Path        string
	Appname     string
	RootAppname string
}

type FilterInfo struct {
	Path        string
	Appname     string
	RootAppname string
}

type Page404Info struct {
	Path        string
	Appname     string
	RootAppname string
}

type Page500Info struct {
	Path        string
	Appname     string
	RootAppname string
}

type IndexInfo struct {
	Path        string
	Appname     string
	RootAppname string
}

type PackageInfo struct {
	Package     string
	Appname     string
	RootAppname string
}

// 设置启动器信息
type StarterInfo struct {
	Level       int
	Method      string
	Appname     string
	RootAppname string
}

type StarterInfoSort []StarterInfo // StarterInfo

func (sis StarterInfoSort) Len() int           { return len(sis) }
func (sis StarterInfoSort) Less(i, j int) bool { return sis[i].Level < sis[j].Level }
func (sis StarterInfoSort) Swap(i, j int)      { sis[i], sis[j] = sis[j], sis[i] }

//
//封装Webapp的所有路由信息。包括：Action（控制器）、Gml（展现层）、Static（静态资源，例如：js，css等）、Filter（过滤器）、
//404页面、500页面、package（引入的包）。
//
type WebappInfo struct {
	Starters         StarterInfoSort
	Actions          []ActionInfo
	Gmls             []GmlInfo
	Statics          []StaticInfo
	Filter           []FilterInfo
	Page404          []Page404Info
	Page500          []Page500Info
	IndexWebPath     []IndexInfo
	RootIndexWebPath []IndexInfo
	Packages         []PackageInfo
	RootAppnames     []string
	RootAppnameMap   map[string]string
}

//创建WebappInfo对象
//
//@param apps  设置apps的名称。
func NewWebappInfo(apps []string) *WebappInfo {
	webappInfo := new(WebappInfo)
	webappInfo.Statics = make([]StaticInfo, 0)
	webappInfo.Starters = make(StarterInfoSort, 0)
	webappInfo.Actions = make([]ActionInfo, 0)
	webappInfo.Gmls = make([]GmlInfo, 0)
	webappInfo.Packages = make([]PackageInfo, 0)
	webappInfo.Filter = make([]FilterInfo, 0)
	webappInfo.Page404 = make([]Page404Info, 0)
	webappInfo.Page500 = make([]Page500Info, 0)
	webappInfo.IndexWebPath = make([]IndexInfo, 0)
	webappInfo.RootAppnames = apps
	webappInfo.RootAppnameMap = make(map[string]string)
	for idx, rootAppname := range apps {
		webappInfo.RootAppnameMap[rootAppname] = fmt.Sprintf("%d", idx)
	}
	return webappInfo
}

/*
添加app应用

@param path   路由配置文件路径，例如：/home/glassProject/src/glass/app1/route.toml

@param appname   ｀appname｀的名称 ，例如：app1

@param rootAppname  ｀apps｀的名称 ，位置在 src/${rootAppname} ，例如：glass
*/
func (w *WebappInfo) AddWebapp(path string, appname string, rootAppname string) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	} else {

		m, _ := toml.ParseDatas(string(bs))

		w.Packages = append(w.Packages, PackageInfo{rootAppname + "/" + appname, appname, rootAppname})

		if m != nil {

			if m["starter"] != nil {
				starters := m["starter"].([]map[string]interface{})
				for _, startMap := range starters {
					for method, level := range startMap {
						levelStr := ""
						switch vtype := level.(type) {
						case int64:
							levelStr = string(vtype)
						default:
							levelStr = vtype.(string)
						}

						levelInt, err1 := strconv.Atoi(levelStr)
						if err1 != nil {
							levelInt = 99
						} else if levelInt < 0 {
							levelInt = 99
						}

						if !strings.Contains(method, ".") {
							method = w.getPackageAlias(appname, rootAppname) + "." + method
						}
						w.Starters = append(w.Starters, StarterInfo{levelInt, method, appname, rootAppname})
					}
				}
			}

			if m["action"] != nil {
				actions := m["action"].([]map[string]interface{})
				for _, actionMap := range actions {
					for name, webpath := range actionMap {
						webpathstr := webpath.(string)
						if !strings.HasPrefix(webpathstr, "/") {
							webpathstr = "/" + appname + "/" + webpathstr
						}
						if !strings.Contains(name, ".") {
							name = w.getPackageAlias(appname, rootAppname) + "." + name
						}
						w.Actions = append(w.Actions, ActionInfo{name, webpathstr, appname, rootAppname})
					}
				}
			}

			if m["gml"] != nil {
				gmls := m["gml"].([]map[string]interface{})
				for _, gmlMap := range gmls {
					for name, webpath := range gmlMap {
						webpathstr := webpath.(string)
						if strings.HasPrefix(webpathstr, "/") {
							w.Gmls = append(w.Gmls, GmlInfo{name, webpathstr, appname, rootAppname})
						} else {
							w.Gmls = append(w.Gmls, GmlInfo{name, "/" + appname + "/" + webpathstr, appname, rootAppname})
						}

					}
				}
			}

			if m["static"] != nil {
				statics := m["static"].([]map[string]interface{})
				for _, staticMap := range statics {
					for _, webpath := range staticMap {
						webpathstr := webpath.(string)
						webpathArray := strings.Split(webpathstr, ",")
						for _, wp := range webpathArray {
							if strings.HasPrefix(webpathstr, "/") {
								w.Statics = append(w.Statics, StaticInfo{wp, appname, rootAppname})
							} else {
								w.Statics = append(w.Statics, StaticInfo{"/" + appname + "/" + wp, appname, rootAppname})
							}
						}

					}
				}
			}

			if m["filter"] != nil {
				filters := m["filter"].([]map[string]interface{})
				for _, filterMap := range filters {
					for _, classname := range filterMap {
						classnamestr := classname.(string)
						if !strings.Contains(classnamestr, ".") {
							classnamestr = w.getPackageAlias(appname, rootAppname) + "." + classnamestr
						}
						w.Filter = append(w.Filter, FilterInfo{classnamestr, appname, rootAppname})
					}
				}
			}

			if m["404"] != nil {
				page1 := m["404"].([]map[string]interface{})
				for _, pageMap := range page1 {
					for _, p := range pageMap {
						pstr := p.(string)
						if strings.HasPrefix(pstr, "/") {
							w.Page404 = append(w.Page404, Page404Info{pstr, appname, rootAppname})
						} else {
							w.Page404 = append(w.Page404, Page404Info{"/" + appname + "/" + pstr, appname, rootAppname})
						}
					}
				}
			}

			if m["500"] != nil {
				page1 := m["500"].([]map[string]interface{})
				for _, pageMap := range page1 {
					for _, p := range pageMap {
						pstr := p.(string)
						if strings.HasPrefix(pstr, "/") {
							w.Page500 = append(w.Page500, Page500Info{pstr, appname, rootAppname})
						} else {
							w.Page500 = append(w.Page500, Page500Info{"/" + appname + "/" + pstr, appname, rootAppname})
						}
					}
				}
			}

			if m["index"] != nil {
				page1 := m["index"].([]map[string]interface{})
				for _, pageMap := range page1 {
					for _, p := range pageMap {
						pstr := p.(string)
						if strings.HasPrefix(pstr, "/") {
							w.IndexWebPath = append(w.IndexWebPath, IndexInfo{pstr, appname, rootAppname})
						} else {
							w.IndexWebPath = append(w.IndexWebPath, IndexInfo{"/" + appname + "/" + pstr, appname, rootAppname})
						}
					}
				}
			}

			if m["rootindex"] != nil { // 与index不同，rootindex，是根路径指向的index地址，可以理解为 http://ip:port/的响应
				page1 := m["rootindex"].([]map[string]interface{})
				for _, pageMap := range page1 {
					for _, p := range pageMap {
						pstr := p.(string)
						if strings.HasPrefix(pstr, "/") {
							w.RootIndexWebPath = append(w.RootIndexWebPath, IndexInfo{pstr, appname, rootAppname})
						} else {
							w.RootIndexWebPath = append(w.RootIndexWebPath, IndexInfo{"/" + appname + "/" + pstr, appname, rootAppname})
						}
					}
				}
			}

		}
	}
}

var fragment01 string = `
package main

import (
    "runtime"
    "os"
    "os/exec"
    "fmt"
    "strings"
    "path/filepath"
    "io/ioutil"
    "strconv"
	"alphabet"
	"alphabet/env"
	"alphabet/web"
	cmd_utils "alphabet/cmd/utils"

`

var fragment02_val string = `#val1# "#val2#"` //需要定义包引入 例如："apps/app1"

var fragment02_val_optional string = `_ "#val#"` //需要定义包引入 例如：_ "net/http/pprof"

var fragment03_val string = `
)

func main() {
  fmt.Println("abserver.starting.... enter.")
  isBuildBin := #val1#

	runtime.GOMAXPROCS(#val2#)
	env.Env_Project_Resource_Apps = #val3#
	projectPath :="#val4#"

	cmd_utils.ResolveHelperCmdArgs()

	appnamesFromCmd:=cmd_utils.GetAppnamesFromCmd()
  if appnamesFromCmd!=nil{
		env.Env_Project_Resource_Apps = appnamesFromCmd
	}

	envsFromCmd:=cmd_utils.GetEnvsFromCmd()
	for k,v:=range envsFromCmd{
		os.Setenv(k,v)
	}

	fmt.Printf("abserver.starting....  Params ::::: %v \n",os.Args)

	fmt.Printf("abserver.starting....  Pid    ::::: ( %d ) \n",os.Getpid())

	fmt.Printf("abserver.starting....  Envs   ::::: %v  \n",envsFromCmd)

	if cmd_utils.IsHelpFromCmd(){
		cmd_utils.HelperForNest()
		return
	}


	if isBuildBin {
		binPath, _ := exec.LookPath(os.Args[0])
		absBinPath, _ := filepath.Abs(binPath)
		absBinPath = filepath.ToSlash(absBinPath)
		indexNum := strings.LastIndex(absBinPath, "/")
    	absBinPath = string(absBinPath[0 : indexNum])
    	indexNum = strings.LastIndex(absBinPath, "/")
    	absBinPath = string(absBinPath[0 : indexNum])
    	fmt.Println("ProjectPath::::::",absBinPath)

    	projectPath=absBinPath
  }

	env.Env_Project_Root=projectPath

	configArgsTag,configArgsStr:= cmd_utils.GetConfigKeyOrUrlFromCmd()
	if configArgsTag==0{
		fmt.Printf("abserver.starting....  config ::::: none  \n")
	}else if configArgsTag==1{
		fmt.Printf("abserver.starting....  config ::::: configKey = %s  \n",configArgsStr)
		cmd_utils.DoSysconfigPathForKey(configArgsStr)
	}else if configArgsTag==2{
		fmt.Printf("abserver.starting....  config ::::: configUrl = %s  \n",configArgsStr)
		cmd_utils.DoSysconfigPathForUrl(configArgsStr)
	}

	alphabet.Init(projectPath,true)

	fmt.Printf("abserver.starting....  Apps   ::::: %v  \n",env.Env_Project_Resource_Apps)

  for _, rootAppName := range env.Env_Project_Resource_Apps {
	  fmt.Printf("abserver.starting....  Active ::::: rootAppName is  [%s] . Env_Web_Context is : [%s] \n", rootAppName, env.Env_Web_Context[rootAppName])
  }

  ioutil.WriteFile("#val5#", []byte(strconv.Itoa(os.Getpid())), 0777)

	RouteFunc()

	web.RegisterRouteAndStartServer()
}

func RouteFunc() {
`

var fragment04_val string = `web.AddRouteFilter(#val1#,"#val2#")` // 需要定义过滤器 例如：web.AddRouteFilter(ServFilter)

var fragment05_val string = `web.AddRouteActioner("#val1#","#val2#","#val3#").Use(#val4#)`   // web.AddRouteActioner("/app1/action1run1").Use(app1.Action1Run1)
var fragment06_val string = `web.AddRouteTemplate("#val1#","#val2#","#val3#").Tpl("#val4#")` // web.AddRouteTemplate("action1_run1_page1").Tpl("/app1/view/action1_run1_page1.gml")
var fragment07_val string = `web.AddRouteStatics("#val1#","#val2#")`                         // web.AddRouteStatics("/app1/resource/")
var fragment08_val string = `web.AddRoutePage404("#val1#","#val2#")`                         // web.AddRoutePage404("/global/resource/404.html")
var fragment09_val string = `web.AddRoutePage500("#val1#","#val2#")`                         // web.AddRoutePage500("/global/resource/500.html")

var fragment10_val string = `web.AddRouteIndex("#val#","#val2#")` // web.AddRouteIndex("/global/resource/index.html")

var fragment11_val string = `web.AddRouteRootIndex("#val#","#val2#")` // web.AddRouteRootIndex("/global/resource/index.html")

var fragment21_val string = `#val#()` // 设置启动器

var fragment99 string = `
}


`

/*
动态生成当前webapp的可执行程序代码，如果启动该程序代码，那么就运行了该webapp服务。

@param projectPath   设置webapp项目路径，例如：/home/glassProject/

@param execProjectPath   设置执行的webapp项目路径，例如：/home/glassout/

@param pidPath   设置程序运行时进程号文件路径，例如：/home/glassout/bin/pid ，如果为 "" ，表示默认存储到临时目录下

@param cpunum  设置cpu个数

@return string  返回程序代码

*/
func (w *WebappInfo) GenerateWebMainApplication(projectPath string, execProjectPath string, pidPath string, cpunum int, isProd bool, isBuildBin bool) string {
	applString := fragment01
	if !isProd {
		applString = fmt.Sprintf("%s\n%s\n", applString, strings.Replace(fragment02_val_optional, "#val#", "net/http/pprof", -1))
	}

	for _, v := range w.Packages {
		isUsed := false               // fix：用于修复 定义了 app的package，但是没有使用的问题。
		for _, a := range w.Actions { // 查看当前app有没有定义action
			if a.Appname == v.Appname {
				isUsed = true
				break
			}
		}
		if !isUsed {
			for _, f := range w.Filter { // 查看当前app有没有定义filter
				if f.Appname == v.Appname {
					isUsed = true
					break
				}
			}
		}

		if !isUsed {
			for _, starter := range w.Starters {
				if starter.Appname == v.Appname {
					isUsed = true
					break
				}
			}
		}

		if isUsed { // 该 app 有使用
			str := strings.Replace(fragment02_val, "#val1#", w.getPackageAlias(v.Appname, v.RootAppname), -1)
			str = strings.Replace(str, "#val2#", v.Package, -1)

			applString = fmt.Sprintf("%s\n%s\n", applString, str)
		}
	}
	if true {
		str := strings.Replace(fragment03_val, "#val1#", strconv.FormatBool(isBuildBin), -1)
		str = strings.Replace(str, "#val2#", strconv.Itoa(cpunum), -1)
		str = strings.Replace(str, "#val3#", fmt.Sprintf("[]string{\"%s\"}", strings.Join(w.RootAppnames, "\",\"")), -1)
		str = strings.Replace(str, "#val4#", execProjectPath, -1)
		pathNames := strings.Split(projectPath, "/")
		projectName := ""
		lenPathNames := len(pathNames)
		for i := 1; i <= lenPathNames; i++ {
			projectName = pathNames[lenPathNames-i]
			if projectName != "" {
				break
			}
		}

		if pidPath == "" {
			str = strings.Replace(str, "#val5#", strings.Replace(os.TempDir(), "\\", "/", -1)+"/"+projectName+"/pid", -1)

		} else {
			str = strings.Replace(str, "#val5#", pidPath, -1)

		}

		applString = fmt.Sprintf("%s\n%s\n", applString, str)
	}

	for _, filterObject := range w.Filter {
		if filterObject.Path != "" {
			str := strings.Replace(fragment04_val, "#val1#", filterObject.Path, -1)
			str = strings.Replace(str, "#val2#", filterObject.RootAppname, -1)
			applString = fmt.Sprintf("%s\n%s\n", applString, str)
		}
	}

	for _, act := range w.Actions {
		str := strings.Replace(fragment05_val, "#val1#", act.Webpath, -1)
		str = strings.Replace(str, "#val2#", act.Appname, -1)
		str = strings.Replace(str, "#val3#", act.RootAppname, -1)
		str = strings.Replace(str, "#val4#", act.Name, -1)
		applString = fmt.Sprintf("%s\n%s\n", applString, str)
	}

	for _, gml := range w.Gmls {
		str := strings.Replace(fragment06_val, "#val1#", gml.RootAppname, -1)
		str = strings.Replace(str, "#val2#", gml.Appname, -1)
		str = strings.Replace(str, "#val3#", gml.Name, -1)
		str = strings.Replace(str, "#val4#", gml.Webpath, -1)
		applString = fmt.Sprintf("%s\n%s\n", applString, str)
	}

	for _, static := range w.Statics {
		str := strings.Replace(fragment07_val, "#val1#", static.Path, -1)
		str = strings.Replace(str, "#val2#", static.RootAppname, -1)
		applString = fmt.Sprintf("%s\n%s\n", applString, str)
	}

	for _, p404 := range w.Page404 {
		if p404.Path != "" {
			str := strings.Replace(fragment08_val, "#val1#", p404.Path, -1)
			str = strings.Replace(str, "#val2#", p404.RootAppname, -1)
			applString = fmt.Sprintf("%s\n%s\n", applString, str)
		}
	}

	for _, p500 := range w.Page500 {
		if p500.Path != "" {
			str := strings.Replace(fragment09_val, "#val1#", p500.Path, -1)
			str = strings.Replace(str, "#val2#", p500.RootAppname, -1)
			applString = fmt.Sprintf("%s\n%s\n", applString, str)
		}
	}

	for _, indexPage := range w.IndexWebPath {
		if indexPage.Path != "" {
			str := strings.Replace(fragment10_val, "#val#", indexPage.Path, -1)
			str = strings.Replace(str, "#val2#", indexPage.RootAppname, -1)
			applString = fmt.Sprintf("%s\n%s\n", applString, str)
		}
	}

	for _, rootIndexPage := range w.RootIndexWebPath {
		if rootIndexPage.Path != "" {
			str := strings.Replace(fragment11_val, "#val#", rootIndexPage.Path, -1)
			str = strings.Replace(str, "#val2#", rootIndexPage.RootAppname, -1)
			applString = fmt.Sprintf("%s\n%s\n", applString, str)
		}
	}

	sort.Sort(w.Starters)
	for _, starter := range w.Starters {
		applString = fmt.Sprintf("%s\n%s\n", applString, strings.Replace(fragment21_val, "#val#", starter.Method, -1))
	}

	applString = fmt.Sprintf("%s\n%s\n", applString, fragment99)

	return applString

}

func (w *WebappInfo) getPackageAlias(appname string, rootAppname string) string {
	return fmt.Sprintf("%s_%s", appname, w.RootAppnameMap[rootAppname])
}
