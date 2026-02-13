package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

//go:embed assets/config.json
var defaultConfig []byte

type Config struct {
	Static    bool   `json:"-"` // 是否启用静态配置（不生成外置文件，也不监控）
	Title     string `json:"title"`
	URL       string `json:"url"`
	Icon      string `json:"icon"` // 外置图标路径（空则使用内嵌）
	AutoStart bool   `json:"autoStart"`
	TrayMode  bool   `json:"trayMode"`

	mu       sync.RWMutex `json:"-"`
	saving   bool         `json:"-"` // 防止自循环标记
	watcher  *fsnotify.Watcher
	path     string        // 外置配置文件绝对路径
	dir      string        // 配置文件所在目录
	onChange func(*Config) // 变更回调
}

func LoadConfig(isStatic bool) (*Config, error) {
	c := &Config{
		Static:    isStatic,
		Title:     "WebLauncher",
		URL:       "https://www.example.com",
		AutoStart: false,
		TrayMode:  true,
	}

	// 解析内嵌默认配置
	if err := json.Unmarshal(defaultConfig, c); err != nil {
		return nil, err
	}

	// 静态配置模式
	if c.Static {
		return c, nil
	}

	// 确定外置配置路径
	exe, err := os.Executable()
	if err != nil {
		return c, nil
	}
	c.dir = filepath.Dir(exe)
	c.path = filepath.Join(c.dir, "config.json")

	// 外置配置覆盖
	if _, err := os.Stat(c.path); err == nil {
		data, _ := os.ReadFile(c.path)
		var ext Config
		if err := json.Unmarshal(data, &ext); err == nil {
			if ext.Title != "" {
				c.Title = ext.Title
			}
			if ext.URL != "" {
				c.URL = ext.URL
			}
			if ext.Icon != "" {
				c.Icon = ext.Icon
			}
			c.AutoStart = ext.AutoStart
			c.TrayMode = ext.TrayMode
		}
	} else {
		// 不存在则创建默认外置配置
		c.Save()
	}

	return c, nil
}

func (c *Config) SetStatic(val bool) {
	c.mu.Lock()
	c.Static = val
	c.mu.Unlock()
	if val {
		c.StopWatching()
	} else {
		c.StartWatching()
	}
}

func (c *Config) GetTitle() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Title
}

func (c *Config) GetURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.URL
}

func (c *Config) GetIcon() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Icon
}

func (c *Config) GetAutoStart() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AutoStart
}

func (c *Config) SetAutoStart(val bool) {
	c.mu.Lock()
	c.AutoStart = val
	c.mu.Unlock()
	c.Save()
}

func (c *Config) SetOnChange(fn func(*Config)) {
	c.onChange = fn
}

// Save 保存到外置文件（带防抖标记）
func (c *Config) Save() {
	c.mu.Lock()
	c.saving = true
	data, _ := json.MarshalIndent(c, "", "  ")
	c.mu.Unlock()

	os.WriteFile(c.path, data, 0644)

	// 500ms 后清除标记，避免触发文件监控
	go func() {
		// Windows 下文件写入可能有延迟
		c.mu.Lock()
		c.saving = false
		c.mu.Unlock()
	}()
}

// StartWatching 启动热重载监控
func (c *Config) StartWatching() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	c.watcher = watcher

	// 监控目录（监控文件本身在 Windows 下不可靠）
	err = watcher.Add(c.dir)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) != "config.json" {
					continue
				}
				// 写入或创建事件
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					c.mu.RLock()
					isSaving := c.saving
					c.mu.RUnlock()

					if isSaving {
						continue // 跳过自身触发的保存
					}

					log.Println("Config changed, reloading...")
					// 重新加载
					data, err := os.ReadFile(c.path)
					if err != nil {
						continue
					}

					var newCfg Config
					if err := json.Unmarshal(data, &newCfg); err != nil {
						continue
					}

					c.mu.Lock()
					if newCfg.Title != "" {
						c.Title = newCfg.Title
					}
					if newCfg.URL != "" {
						c.URL = newCfg.URL
					}
					if newCfg.Icon != "" {
						c.Icon = newCfg.Icon
					}
					c.AutoStart = newCfg.AutoStart
					c.TrayMode = newCfg.TrayMode
					c.mu.Unlock()

					if c.onChange != nil {
						c.onChange(c)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()
	return nil
}

func (c *Config) StopWatching() {
	if c.watcher != nil {
		c.watcher.Close()
	}
}
