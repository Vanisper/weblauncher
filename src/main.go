package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/energye/systray"
)

//go:embed assets/icon.ico
var embeddedIcon []byte

var (
	trayMode = flag.Bool("tray", false, "强制托盘模式")
	openOnce = flag.Bool("open", false, "仅打开浏览器并退出")
)

var (
	config   *Config
	menuAuto *systray.MenuItem
)

func main() {
	flag.Parse()

	var err error
	config, err = LoadConfig()
	if err != nil {
		fmt.Println("加载配置失败:", err)
		os.Exit(1)
	}

	// 命令行覆盖
	if *trayMode {
		config.TrayMode = true
	}
	if *openOnce {
		config.TrayMode = false
	}

	// 非托盘模式
	if !config.TrayMode {
		openBrowser(config.GetURL())
		return
	}

	// 应用自启设置
	config.applyAutoStart()

	// 启动配置热重载
	config.StartWatching()
	defer config.StopWatching()

	systray.Run(onReady, onExit)
}

func onReady() {
	// 设置图标（读取外置或内嵌）
	// 实际项目中应将内嵌图标转为 []byte 传入
	systray.SetIcon(getIconData())
	systray.SetTitle(config.GetTitle())
	systray.SetTooltip(config.GetTitle())
	// systray.SetTemplateIcon(getIconData(), getIconData()) // 模板图标支持
	// 双击托盘图标打开网页
	systray.SetOnDClick(func(menu systray.IMenu) {
		openBrowser(config.GetURL())
	})

	// 菜单
	menuOpen := systray.AddMenuItem("打开网页", "Open URL")
	menuOpen.Click(func() {
		openBrowser(config.GetURL())
	})

	menuAuto = systray.AddMenuItemCheckbox("开机自启", "Auto start on boot", config.GetAutoStart())
	menuAuto.Click(func() {
		newState := !config.GetAutoStart()
		config.SetAutoStart(newState)
		config.applyAutoStart()
		if newState {
			menuAuto.Check()
		} else {
			menuAuto.Uncheck()
		}
	})

	systray.AddSeparator()
	menuReload := systray.AddMenuItem("重载配置", "Reload config.json")
	menuReload.Click(func() {
		// 手动触发重载逻辑，实际已由 fsnotify 处理，可能也有需求
		// TODO: 可能需要增加防抖，避免重复触发
	})

	menuQuit := systray.AddMenuItem("退出", "Quit")
	menuQuit.Click(func() {
		systray.Quit()
	})

	// 配置变更回调（热重载后更新 UI）
	config.SetOnChange(func(c *Config) {
		systray.SetTitle(c.GetTitle())
		systray.SetTooltip(c.GetTitle())
		if c.GetAutoStart() {
			menuAuto.Check()
		} else {
			menuAuto.Uncheck()
		}
	})

	// 启动时自动打开浏览器
	openBrowser(config.GetURL())
}

func onExit() {
	config.StopWatching()
}

// getIconData 优先读取外置图标，否则返回内嵌字节
func getIconData() []byte {
	// 如果 config.Icon 指定了外置路径，尝试读取
	if config.GetIcon() != "" {
		data, err := os.ReadFile(config.GetIcon())
		if err == nil {
			return data
		}
		fmt.Fprintf(os.Stderr, "警告: 无法读取外置图标 %s: %v，使用内嵌图标\n", config.GetIcon(), err)
	}
	return embeddedIcon
}
