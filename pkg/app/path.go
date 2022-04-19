package app

import (
	"fmt"
	"net/url"
	"gin-api/pkg/config"
	"gin-api/pkg/helpers"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

//buildPath 设置项目路径
func buildPath(path string) {
	app.PathSeparator = string(os.PathSeparator)
	app.LineSeparator = "\n"
	app.Path.Root     = rootPath(path)
	app.Path.App      = makePath("application")
	app.Path.View     = makePath("application/http/view")
	app.Path.Config   = makePath("config")
	app.Path.Route    = makePath("route")
	app.Path.Public   = makePath("public")
	app.Path.Static   = makePath("public/static")
	app.Path.Runtime  = makePath("runtime")
	app.Path.Log      = makePath("runtime/logs")
}

//rootPath 构造项目根目录
func rootPath(path string) string {
	if len(path) == 0 {
		path = GetBuildAbsPath()
	}
	return helpers.StrJoin(strings.TrimRight(path, app.PathSeparator), app.PathSeparator)
}

//GetPathSeparator 获取路径分隔符
func GetPathSeparator() string {
	return app.PathSeparator
}

//GetRootPath 返回项目根目录
func GetRootPath() string {
	return app.Path.Root
}

//GetLogPath 返回日志目录
func GetLogPath() string {
	return app.Path.Log
}

//makePath 构造相关目录
func makePath(name string) string {
	return path.Join(app.Path.Root, name)
}

//RuntimeLogFile 返回系统错误日志路径
func RuntimeLogFile() string {
	return config.GetString("log.runtime_log")
}

//AccessLogFile 返回系统访问日志路径
func AccessLogFile() string {
	fullUrl := GetFullUrl()
	if fullUrl == "" {
		return ""
	}
	u,_ := url.Parse(fullUrl)
	uri := strings.Replace(strings.Trim(u.Path, "/"), "/", "_", -1)
	return fmt.Sprintf(config.GetString("log.access_log"), uri)
}

//UserLogFile 返回用户自定义日志路径
func UserLogFile(filename string) string {
	return fmt.Sprintf(config.GetString("log.user_log"), filename)
}

//TmplPattern 模板文件路径
func TmplPattern() []string {
	return config.GetStringSlice("app.tmpl_pattern")
}

//LogExceptUri 不记录访问日志的uri
func LogExceptUri() []string  {
	return config.GetStringSlice("log.except_uri")
}

//GetBuildAbsPath 获取编译后可执行文件的根目录
//注意： 如果以 go run 运行，则无法获取正确的根目录,因为 go run 生成的可执行文件位于/tmp目录下
func GetBuildAbsPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	return path[:index]
}