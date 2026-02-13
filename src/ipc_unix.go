//go:build !windows

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

func getIPCPath() string {
	tmpDir := os.TempDir()
	return filepath.Join(tmpDir, "weblauncher.ipc")
}

// startIPCServer 启动 IPC 服务（Unix Domain Socket）
func startIPCServer(onOpenURL func()) error {
	ipcPath := getIPCPath()
	
	// 清理旧的 socket 文件
	os.Remove(ipcPath)
	
	listener, err := net.Listen("unix", ipcPath)
	if err != nil {
		return fmt.Errorf("IPC 服务启动失败: %w", err)
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
	ipcPath := getIPCPath()
	
	conn, err := net.DialTimeout("unix", ipcPath, 2*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接到主实例: %w", err)
	}
	defer conn.Close()
	
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write([]byte("OPEN_URL\n"))
	return err
}
