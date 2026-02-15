package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"weblauncher/internal/build"
)

func main() {
	var (
		envFile   = flag.String("env", ".env", "环境配置文件路径")
		onlyClean = flag.Bool("only-clean", false, "仅执行清理工作（临时、构建文件）")
		clean     = flag.Bool("clean", false, "流程结束后清理临时文件")
		release   = flag.Bool("release", false, "构建完整发布包")
	)
	flag.Parse()

	// 获取项目根目录（兼容 go run 和编译后的可执行文件）
	projectRoot := getProjectRoot()

	// 如果 .env 路径是相对路径，基于项目根目录解析
	envPath := *envFile
	if !filepath.IsAbs(envPath) {
		envPath = filepath.Join(projectRoot, envPath)
	}

	// 加载配置
	cfg, err := build.LoadConfig(envPath, projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		fmt.Fprintln(os.Stderr, "提示: 运行 'task init' 或手动复制 .env.example 到 .env")
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "配置无效: %v\n", err)
		os.Exit(1)
	}

	b := build.NewBuilder(cfg, projectRoot)

	// 仅清理
	if *onlyClean {
		fmt.Println("清理构建文件...")
		if err := b.Clean(true); err != nil {
			fmt.Fprintf(os.Stderr, "清理失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("清理完成")
		return
	}

	// 构建步骤
	var steps []build.Step

	if *release {
		// 完整发布流程
		steps = []build.Step{
			{Name: "生成安装脚本", Fn: b.GenerateISS},
			{Name: "生成图标资源", Fn: b.BuildSyso},
			{Name: "构建可执行文件", Fn: b.BuildExe},
			{Name: "构建安装程序", Fn: b.BuildInstaller},
		}
	} else {
		// 仅构建程序
		steps = []build.Step{
			{Name: "生成图标资源", Fn: b.BuildSyso},
			{Name: "构建可执行文件", Fn: b.BuildExe},
		}
	}

	// 执行构建
	fmt.Printf("开始构建: %s v%s\n", cfg.AppName, cfg.AppVersion)
	fmt.Printf("项目目录: %s\n", projectRoot)
	fmt.Println()

	for _, step := range steps {
		if err := step.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\n构建失败: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("构建完成!")
	fmt.Println("========================================")

	if *clean {
		fmt.Println()
		fmt.Println("清理临时文件...")
		if err := b.Clean(false); err != nil {
			fmt.Fprintf(os.Stderr, "清理失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("清理完成")
		fmt.Println()
	}

	outputDir := filepath.Join(projectRoot, ".output")
	if *release {
		fmt.Printf("程序:   %s\n", filepath.Join(outputDir, cfg.OutputName))
		fmt.Printf("安装包: %s_Setup.exe\n", filepath.Join(outputDir, cfg.AppName))
	} else {
		fmt.Printf("程序: %s\n", filepath.Join(outputDir, cfg.OutputName))
	}
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 获取当前文件路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		// 退而求其次，使用工作目录
		wd, _ := os.Getwd()
		return wd
	}
	// cmd/build/main.go -> 项目根目录
	return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}
