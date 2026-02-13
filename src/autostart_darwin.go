//go:build darwin

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c *Config) applyAutoStart() error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "Library", "LaunchAgents")
	os.MkdirAll(dir, 0755)

	label := "com.example." + c.GetTitle() // 简单处理，实际应用需清理特殊字符
	path := filepath.Join(dir, label+".plist")

	if !c.GetAutoStart() {
		os.Remove(path)
		return nil
	}

	exe, _ := os.Executable()
	abs, _ := filepath.Abs(exe)

	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>--tray</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
`, label, abs)

	return os.WriteFile(path, []byte(content), 0644)
}
