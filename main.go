package main

import (
	"encoding/base64"
	"github.com/alexfaker/Pantasy/config"
	"github.com/alexfaker/Pantasy/router"
	"log"
	"strings"
)

var (
	AppName      string // 应用名称
	AppVersion   string // 应用版本
	BuildVersion string // 编译版本
	BuildTime    string // 编译时间
	GitRevision  string // Git版本
	GitBranch    string // Git分支
	GoVersion    string // Golang信息
)

func outputInformation() {
	//v, _ := base64.StdEncoding.DecodeString(AppVersion)
	//AppVersion = strings.TrimSpace(string(v))
	v, _ := base64.StdEncoding.DecodeString(BuildVersion)
	BuildVersion = strings.TrimSpace(string(v))
	v, _ = base64.StdEncoding.DecodeString(BuildTime)
	BuildTime = strings.TrimSpace(string(v))
	v, _ = base64.StdEncoding.DecodeString(GitRevision)
	GitRevision = strings.TrimSpace(string(v))
	v, _ = base64.StdEncoding.DecodeString(GitBranch)
	GitBranch = strings.TrimSpace(string(v))
	v, _ = base64.StdEncoding.DecodeString(GoVersion)
	GoVersion = strings.TrimSpace(string(v))

	log.Println("App Name:", AppName)
	log.Println("App Version:", AppVersion)
	log.Println("Build version:", BuildVersion)
	log.Println("Build time:", BuildTime)
	log.Println("Git revision:", GitRevision)
	log.Println("Git branch:", GitBranch)
	log.Println("Golang Version:", GoVersion)
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	outputInformation()
	config.Initialize()

	router.Run()
}

// NodeCode 获取节点随机唯一编号
