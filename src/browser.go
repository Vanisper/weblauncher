//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

func openBrowser(url string) error {
	cmd := exec.Command("cmd", "/c", "start", "", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		// TODO: 待验证
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW - Win7 支持
	}
	return cmd.Start()
}
