package lib

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// ExecProgram 执行程序
func ExecProgram(name string) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd.exe", "/C", "start", name)
		if err := cmd.Start(); err != nil {
			fmt.Println("启动", name, "错误", err)
		}
		break
	case "linux":
		err := os.Chmod(name, 0777)
		if err != nil {
			fmt.Println("修改文件", name, "权限错误", err)
		}
		//cmd := exec.Command("/bin/bash", "-c", name)
		cmd := exec.Command(name)
		//if err = cmd.Run(); err != nil {
		if err = cmd.Start(); err != nil {
			fmt.Println("启动", name, "错误!", err)
		}
	}

	time.Sleep(2 * time.Second)
}
