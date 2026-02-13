//go:build darwin

package main

import (
	"syscall"
	"os"
)

// flock 使用 syscall 实现文件锁
func flock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

// unflock 释放文件锁
func unflock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
