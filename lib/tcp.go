package lib

import (
	"fmt"
	"net"
	"time"
)

const (
	upHost  = "127.0.0.1:7778" //Upgrade.exe Port
	proHost = "127.0.0.1:7777" //Client.exe Port
)

var (
	ServerMsg    string //Tcp Server 返回信息
	tcpStatus    string //Tcp Server 运行状态
	clientStatus string //客户端运行状态
)

// StartTcpServer 后台启动 Tcp Server
func StartTcpServer() {
	if tcpStatus == "running" { //是否已经运行
		fmt.Println("Server has running!")
		return
	}
	go startTcpServer() //后台启动
}

// CheckClientUpdateTcp 通过 TCP 验证 Client是否正在更新
func CheckClientUpdateTcp() {
	fmt.Println("[Upgrade] Check program update status...")
	for {
		if checkUpgradeRun("client update") {
			fmt.Println("[Upgrade] program status is update.")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// WaitClientStopTcp 通过 TCP 等待 Client.exe 停止
func WaitClientStopTcp() {
	fmt.Println("[Upgrade] Wait old program stop...")
	for {
		err := scanPort(proHost)
		if err != nil {
			fmt.Println("[Upgrade] old program has stopped!")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// WaitNewRunningTcp 通过 TCP 等待 Client.exe 启动
func WaitNewRunningTcp() {
	fmt.Println("[Upgrade] Wait new client run...")
	for {
		if clientStatus == "running" {
			fmt.Println("[Upgrade] new client running!")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// check Client
func checkUpgradeRun(cStr string) bool {
	var str string
	var err error
	var st = false
	for i := 0; i < 3; i++ {
		str, err = startTcpClient("by ws_upgrade")
		if err == nil {
			if str == cStr {
				fmt.Println("[TCP:Client] return Msg:", str)
				st = true
				break
			} else {
				fmt.Println("[TCP:Client] 端口开启,不是program!")
			}
		}
		time.Sleep(1 * time.Second)
	}
	return st
}

/* TCP 相关 */
// Tcp Server 启动
func startTcpServer() {
	fmt.Println("[TCP:Server] at:", upHost)
	listen, err := net.Listen("tcp", upHost) //开启端口
	if err != nil {
		fmt.Println("[TCP:Server] listen Error:", err)
	}
	tcpStatus = "running"
	fmt.Println("[TCP:Server] started")
	for {
		if tcpStatus == "stop" { //停止
			break
		}
		// 接受数据
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("[TCP:Server] accept Error:", err)
		}
		go acceptClient(conn)
	}
	err = listen.Close() //关闭端口
	if err != nil {
		fmt.Println("[TCP:Server] Stop Error:", err)
		return
	}
	tcpStatus = "stopped"
}

// Tcp Server 接受消息
func acceptClient(conn net.Conn) {
	fmt.Println("[TCP:Server] set resMsg:", ServerMsg)
	//创建消息缓冲区
	buffer := make([]byte, 512)
	//先读取消息
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("[TCP:Server] read Msg Error:", err)
		return
	}
	clientMsg := string(buffer[0:n])
	fmt.Println("[TCP:Server] accept Msg:", clientMsg)
	//如果接受到 running信息 就退出
	if clientMsg == "running" {
		clientStatus = "running"
		tcpStatus = "stop"
	}
	//再发送消息
	conn.Write([]byte(ServerMsg))
	conn.Close()
}

// Tcp Client 启动 发送和读取 消息
func startTcpClient(msg string) (string, error) {
	var str string
	connTimeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", proHost, connTimeout)
	if err != nil {
		//fmt.Println("Port none connect,Error:", err)
		return "", err
	}
	// 发送消息
	fmt.Println("[TCP:Client] send Msg:", msg)
	conn.Write([]byte(msg))
	// 接受消息
	for {
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			//fmt.Println("Receive Server Error:", err)
			return "", err
		}
		str = string(buf[0:n])
		if str != "" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	//关闭
	err = conn.Close()
	if err != nil {
		//fmt.Println("Stop tcp client Error:", err)
		return str, err
	}
	//fmt.Println("Client Exit")
	return str, nil
}

/* 端口扫描 */

// 端口扫描
func scanPort(host string) error {
	connTimeout := 1 * time.Second
	//conn, err := net.Dial("tcp", host+":"+port)
	conn, err := net.DialTimeout("tcp", host, connTimeout)
	if err != nil {
		fmt.Println("Port none connect,Error:", err)
		return err
	}

	//关闭
	err = conn.Close()
	if err != nil {
		fmt.Println("Stop tcp client Error:", err)
		return err
	}
	return nil
}
