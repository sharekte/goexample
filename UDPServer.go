package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
)

var commandChan = make(chan string, 0) //开辟无缓冲通道，存放客户端发来的命令

func main() {
	var tempPath string
	var app string

	udpServer, err := net.ListenPacket("udp", ":8080") //监听UDP 8080端口
	if err != nil {
		fmt.Println(err)
		return
	}
	defer udpServer.Close()
	osVersion := runtime.GOOS //获取操作系统
	switch osVersion {
	case `darwin`:
		tempPath = `/applications/%s.app`
	case `windows`:
		tempPath = `c:\program files\%s.exe`
	}
	for {
		go udpHandle(udpServer)                       //另开一个协程读客户端数据
		app = <-commandChan                           //阻塞读通道数据
		path := fmt.Sprintf(tempPath, app)            //拼接app实际路径
		_, err := exec.Command("open", path).Output() //在终端终端执行命令
		if err != nil {
			fmt.Errorf("%s", err)
		}
	}
}

func udpHandle(c net.PacketConn) {
	var buffer = make([]byte, 1024)     //开辟1024字节缓冲区
	count, _, err := c.ReadFrom(buffer) //从客户端读数据到缓冲区
	if err != nil {
		log.Fatalln(err)
	}
	command := string(buffer[:count]) //读缓冲区
	commandChan <- command            //将数据塞进通道
}
