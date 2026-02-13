//go:build windows

package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Singleton struct {
	name   string
	handle windows.Handle
}

// NewSingleton 创建单例锁，如果程序已运行则返回错误
func NewSingleton(name string) (*Singleton, error) {
	mutexName, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}

	// 创建安全属性（允许所有用户访问）
	sa := &windows.SecurityAttributes{
		Length:             uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		InheritHandle:      0,
		SecurityDescriptor: nil,
	}

	// 创建命名互斥量
	handle, err := windows.CreateMutex(sa, false, mutexName)
	if err != nil {
		// 检查是否是已存在的错误
		if err == windows.ERROR_ALREADY_EXISTS {
			return nil, fmt.Errorf("程序已在运行")
		}
		return nil, fmt.Errorf("创建互斥量失败: %w", err)
	}

	// 再次检查 GetLastError（CreateMutex 成功时可能返回 ERROR_ALREADY_EXISTS）
	if windows.GetLastError() == windows.ERROR_ALREADY_EXISTS {
		windows.CloseHandle(handle)
		return nil, fmt.Errorf("程序已在运行")
	}

	return &Singleton{
		name:   name,
		handle: handle,
	}, nil
}

// Release 释放单例锁
func (s *Singleton) Release() {
	if s.handle != 0 {
		windows.ReleaseMutex(s.handle)
		windows.CloseHandle(s.handle)
		s.handle = 0
	}
}

// IsRunning 检查程序是否已在运行
func IsRunning(name string) bool {
	mutexName, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return false
	}

	sa := &windows.SecurityAttributes{
		Length:             uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		InheritHandle:      0,
		SecurityDescriptor: nil,
	}

	handle, err := windows.CreateMutex(sa, false, mutexName)
	if err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			return true
		}
		return false
	}
	defer windows.CloseHandle(handle)

	if windows.GetLastError() == windows.ERROR_ALREADY_EXISTS {
		return true
	}
	return false
}
