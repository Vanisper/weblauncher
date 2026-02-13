//go:build !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Singleton 使用文件锁实现单例（适用于 Linux/macOS）
type Singleton struct {
	name string
	file *os.File
}

// NewSingleton 创建单例锁
func NewSingleton(name string) (*Singleton, error) {
	// 使用系统临时目录或用户目录存储锁文件
	lockDir := os.TempDir()
	lockPath := filepath.Join(lockDir, name+".lock")

	// 尝试创建并锁定文件
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("无法创建锁文件: %w", err)
	}

	// 尝试获取文件锁（非阻塞）
	if err := flock(file); err != nil {
		file.Close()
		return nil, fmt.Errorf("程序已在运行")
	}

	// 写入 PID 便于调试
	fmt.Fprintf(file, "%d", os.Getpid())

	return &Singleton{
		name: name,
		file: file,
	}, nil
}

// Release 释放单例锁
func (s *Singleton) Release() {
	if s.file != nil {
		unflock(s.file)
		s.file.Close()
		// 尝试删除锁文件（可选）
		lockDir := os.TempDir()
		lockPath := filepath.Join(lockDir, s.name+".lock")
		os.Remove(lockPath)
		s.file = nil
	}
}

// IsRunning 检查程序是否已在运行
func IsRunning(name string) bool {
	s, err := NewSingleton(name)
	if err != nil {
		return true
	}
	s.Release()
	return false
}
