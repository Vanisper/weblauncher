//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

func (c *Config) applyAutoStart() error {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	exe, _ := os.Executable()
	abs, _ := filepath.Abs(exe)
	// Win7 兼容性：确保路径引号正确处理
	val := fmt.Sprintf(`"%s" --tray`, abs)
	name := c.GetTitle()

	if c.GetAutoStart() {
		return k.SetStringValue(name, val)
	} else {
		k.DeleteValue(name)
		return nil
	}
}

// 用于检测 applyAutoStart 是否生效
func isAutoStart() bool {
	return IsAutoStart(config.GetTitle())
}

// 检查是否已设置自启
func IsAutoStart(appName string) bool {
	key, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	_, _, err = key.GetStringValue(appName)
	return err == nil
}

// hideConsole 隐藏控制台窗口
// Win7 兼容：使用 ShowWindow API
func hideConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")

	user32 := syscall.NewLazyDLL("user32.dll")
	showWindow := user32.NewProc("ShowWindow")

	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd != 0 {
		// SW_HIDE = 0
		showWindow.Call(hwnd, 0)
	}
}
