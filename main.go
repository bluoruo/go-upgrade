package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"upgrade/lib"
)

const (
	newPath = "/new_version/" //新程序目录
	appVer  = "2.1"
)

var (
	sysOs       string //操作系统
	basePath    string //程序所在目录
	programName string //程序名
)

// 初始化
func init() {
	sysOs = runtime.GOOS
	basePath, _ = os.Getwd()
}

// 版本信息
func verInfo() {
	fmt.Println("######################################")
	fmt.Println("# Auto upgrade for WorkStation Client#")
	fmt.Println("# Ver", appVer, "at 2023-03-05              #")
	fmt.Println("# By Comanche Lab.                   #")
	fmt.Println("######################################")
	fmt.Println(" Option:")
	fmt.Println("	the [old program] must be in the same directory as [ws_upgrade].")
	fmt.Println("	the [new program] must be in the </new_version/> temp directory.")
	fmt.Println("	-name <program name>	your program name.")
	fmt.Println("")
	fmt.Println(" Manual:")
	fmt.Println("	1. [old program] run update => listen port 7777 => exec [upgrade] => check port 7778 serverMsg <start update> (ok) => exit [old program]")
	fmt.Println("	2. [upgrade] start => listen port 7778 serverMsg <start update> and check port 7777 (down) => upgrade [new program] and exec [new program] => change serverMsg(wait client).")
	fmt.Println("	3. [new program] Start => check port 7778 (ok) and <wait client> => Send Msg <running> to port 7778")
	fmt.Println("	4. [upgrade] receive msg <running> => exit [upgrade]")
	fmt.Println("----------------------------------------------------------------------------")
}

// 显示更新基础信息
func echoUpdateBaseInfo() {

	fmt.Println("OS:", sysOs)
	fmt.Println("Path:", basePath)
	fmt.Println("program:", programName)
	fmt.Println("new program path:", basePath+newPath)
	fmt.Println("----------------------------------------------------------------------------")
}

// 格式化用户名
func formatProgramName() {
	if sysOs == "windows" {
		if !strings.Contains(programName, ".exe") {
			programName = programName + ".exe"
		}
	}
}

// 通过 tcp 方式更新
func upgradeByTcp() {
	//更新前
	formatProgramName() //格式化程序名
	echoUpdateBaseInfo()
	//lib.UpgradeProgress = "wait client stop" //更新进度 -- 等待客户端停止
	//监听 7778 And resMsg = "start update"
	lib.ServerMsg = "start update"                        //tpc server返回信息
	lib.StartTcpServer()                                  //监听 7778
	time.Sleep(1 * time.Second)                           //等待1s
	lib.CheckClientUpdateTcp()                            //确认客户端是否处在更新中
	fmt.Println("[Upgrade] wait Program stop.")           //进度提示
	lib.WaitClientStopTcp()                               //等待客户端结束
	fmt.Println("[Upgrade] Check new Program.")           //进度提示
	if lib.FileExists(basePath + newPath + programName) { //新程序存在
		//开始更新
		fmt.Println("[Upgrade] Start update....") //进度提示
		// 老程序更名
		fmt.Println("[Upgrade] Rename old program.")
		if lib.FileReName(basePath+"/"+programName,
			basePath+"/"+fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), programName)) {
			// 复制新版本到当前目录
			fmt.Println("[Upgrade] Copy new program.")
			err := lib.FileCopy(basePath+newPath+programName, basePath+"/"+programName)
			if err == nil {
				_ = os.Remove(basePath + newPath + programName)      //删除之前的更新文件
				fmt.Println("[Upgrade] Starting new program.......") //提示
				go lib.ExecProgram(basePath + "/" + programName)     //运行新程序
				lib.ServerMsg = "wait client"                        //tpc server返回信息
				lib.WaitNewRunningTcp()                              //等待新程序启动
				time.Sleep(1 * time.Second)                          //1s后退出
				os.Exit(0)
			}
		}
	}
	//更新阶段出错 就运行老程序
	fmt.Println("[Upgrade] Error: start old program!")
	lib.ExecProgram(basePath + "/" + programName) //运行旧程序
	os.Exit(0)                                    //退出更新程序
}

func main() {
	verInfo() //显示版本信息
	//加载参数
	flag.StringVar(&programName, "n", "client", "program name")
	flag.Parse()
	//执行更新
	upgradeByTcp()
}
