# WebLauncher

一个将网站打包成 Windows 桌面应用程序的工具。支持系统托盘运行、开机自启、配置热重载等功能。

## 功能特性

- 🖥️ **系统托盘模式** - 最小化到系统托盘，双击打开网页
- 🔧 **配置热重载** - 修改配置文件无需重启程序
- 🚀 **开机自启** - 支持设置开机自动启动
- 🚫 **单例模式** - 防止程序重复运行
- 📦 **一键打包** - 内置构建系统，支持生成安装程序（Inno Setup）
- ⚙️ **灵活配置** - 支持外部配置文件覆盖内置默认值

## 项目结构

```
weblauncher/
├── cmd/build/          # 构建工具
│   └── main.go         # 构建程序入口
├── internal/build/     # 构建逻辑
│   ├── config.go       # 构建配置
│   └── steps.go        # 构建步骤
├── src/                # 主程序源码
│   ├── assets/         # 嵌入资源
│   │   ├── config.json # 默认配置（会被嵌入）
│   │   └── icon.ico    # 应用图标
│   ├── main.go         # 程序入口
│   ├── config.go       # 配置管理
│   ├── browser.go      # 浏览器调用
│   └── ...             # 平台相关实现
├── build/              # 构建相关文件
│   └── installer/      # 安装程序模板
├── .env.example        # 环境配置示例
└── Taskfile.yml        # 任务定义
```

## 快速开始

### 环境要求

- Go 1.20+
- Windows（主要支持平台）
- [Task](https://taskfile.dev/)（可选，用于运行任务）
- [Inno Setup 6](https://jrsoftware.org/isdl.php)（可选，用于构建安装程序）

### 初始化项目

```bash
# 复制环境配置
cp .env.example .env

# 或使用 Task
task init
```

编辑 `.env` 文件，设置你的应用信息：

```env
APP_ID={{你的GUID}
APP_NAME=YourAppName
APP_PUBLISHER=YourName
APP_URL=YourAppUrl
OUTPUT_NAME=带后缀的输出名称
```

> 💡 提示：生成 GUID 可以使用 `task new-guid` 命令

### 配置默认网页

编辑 `src/assets/config.json`：

```json
{
  "title": "我的应用",
  "url": "https://www.example.com",
  "icon": "",
  "autoStart": false,
  "trayMode": true
}
```

### 构建

```bash
# 仅构建可执行文件
task build

# 构建完整发布包（包含安装程序）
task release

# 开发模式运行
task run
```

构建输出位于 `.output/` 目录。

## 配置说明

### 嵌入式配置（`src/assets/config.json`）

打包时嵌入到程序中，作为默认配置。用户首次运行时，会自动在数据目录生成可修改的外部配置。

### 外部配置（`config.json`）

程序运行时会自动创建，优先级高于嵌入式配置：

- **Windows**: 程序所在目录（可写时）或 `%APPDATA%/{title}/`
- **热重载**: 修改后自动生效，无需重启

配置项说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| `title` | string | 应用标题（显示在托盘菜单） |
| `url` | string | 要打开的网页地址 |
| `icon` | string | 外置图标路径（空则使用内嵌图标） |
| `autoStart` | bool | 是否开机自启 |
| `trayMode` | bool | 是否启用托盘模式 |

### 静态配置模式

使用 `-static` 参数启动，程序将不会生成外部配置文件：

```bash
weblauncher.exe -static
```

## 命令行参数

| 参数 | 说明 |
|------|------|
| `-tray` | 强制启用托盘模式 |
| `-open` | 仅打开浏览器并退出 |
| `-static` | 启用静态配置模式 |

## 技术栈

- **GUI**: [systray](https://github.com/energye/systray) - 跨平台系统托盘库
- **文件监控**: [fsnotify](https://github.com/fsnotify/fsnotify) - 配置文件热重载
- **环境配置**: [godotenv](https://github.com/joho/godotenv) - 构建配置管理
- **安装程序**: [Inno Setup](https://jrsoftware.org/isinfo.php) - Windows 安装包生成

## 跨平台支持

项目代码结构支持跨平台，目前主要实现和测试在 **Windows** 平台。

已预留其他平台实现文件：
- `autostart_darwin.go` / `autostart_linux.go`
- `ipc_unix.go`
- `singleton_unix.go`
