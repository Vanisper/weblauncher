//go:build windows

package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	ipcPipeName = `\\.\pipe\WebLauncher_IPC`
)

// startIPCServer 启动 IPC 服务，监听新实例的请求
func startIPCServer(onOpenURL func()) error {
	// 使用 go-winio 或直接使用 Windows 命名管道
	// 这里使用简单的 TCP 回环地址作为替代方案（Windows 也支持）
	// 或者使用 github.com/microsoft/go-winio 包
	
	// 简化的方案：使用 TCP 127.0.0.1:0 让系统自动分配端口
	// 但这样需要存储端口信息。改用固定端口但带超时检测
	
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	
	// 保存端口到全局变量或文件，让新实例知道如何连接
	// 这里简化为使用固定端口
	_ = listener.Close()
	
	// 使用固定端口
	listener, err = net.Listen("tcp", "127.0.0.1:17896")
	if err != nil {
		return fmt.Errorf("IPC 服务启动失败（可能已有服务在运行）: %w", err)
	}
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleIPCConnection(conn, onOpenURL)
		}
	}()
	
	return nil
}

func handleIPCConnection(conn net.Conn, onOpenURL func()) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	
	reader := bufio.NewReader(conn)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	
	if cmd == "OPEN_URL\n" && onOpenURL != nil {
		onOpenURL()
	}
}

// sendOpenURLCommand 向已运行的实例发送打开 URL 命令
func sendOpenURLCommand() error {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:17896", 2*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接到主实例: %w", err)
	}
	defer conn.Close()
	
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write([]byte("OPEN_URL\n"))
	return err
}
