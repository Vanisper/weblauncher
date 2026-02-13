//go:build linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c *Config) applyAutoStart() error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "autostart")
	os.MkdirAll(dir, 0755)

	path := filepath.Join(dir, c.GetTitle()+".desktop")

	if !c.GetAutoStart() {
		os.Remove(path)
		return nil
	}

	exe, _ := os.Executable()
	abs, _ := filepath.Abs(exe)

	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s --tray
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`, c.GetTitle(), abs)

	return os.WriteFile(path, []byte(content), 0644)
}
