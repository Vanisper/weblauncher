package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Step 表示一个构建步骤
type Step struct {
	Name string
	Fn   func() error
}

// Run 执行步骤并打印日志
func (s *Step) Run() error {
	fmt.Printf("[Step] %s...\n", s.Name)
	if err := s.Fn(); err != nil {
		return fmt.Errorf("%s 失败: %w", s.Name, err)
	}
	fmt.Printf("[Step] %s 完成\n", s.Name)
	return nil
}

// Builder 构建器
type Builder struct {
	Config      *Config
	ProjectRoot string
}

// NewBuilder 创建构建器
func NewBuilder(cfg *Config, root string) *Builder {
	return &Builder{
		Config:      cfg,
		ProjectRoot: root,
	}
}

// GenerateISS 从模板生成 setup.iss
func (b *Builder) GenerateISS() error {
	templatePath := filepath.Join(b.ProjectRoot, "build", "installer", "setup.template.iss")
	outputPath := filepath.Join(b.ProjectRoot, "build", "installer", "setup.iss")

	template, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取模板失败: %w", err)
	}

	content := string(template)
	replacements := map[string]string{
		"{{APP_ID}}":        b.Config.AppID,
		"{{APP_NAME}}":      b.Config.AppName,
		"{{APP_VERSION}}":   b.Config.AppVersion,
		"{{APP_PUBLISHER}}": b.Config.AppPublisher,
		"{{APP_URL}}":       b.Config.AppURL,
		"{{OUTPUT_NAME}}":   b.Config.OutputName,
	}

	for placeholder, value := range replacements {
		content = strings.ReplaceAll(content, placeholder, value)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入 setup.iss 失败: %w", err)
	}
	return nil
}

// BuildSyso 生成图标资源文件
func (b *Builder) BuildSyso() error {
	iconPath := filepath.Join(b.ProjectRoot, "src", "assets", "icon.ico")
	sysoPath := filepath.Join(b.ProjectRoot, "src", "rsrc.syso")

	// 检查图标是否存在
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		return fmt.Errorf("图标文件不存在: %s", iconPath)
	}

	// 检查 rsrc 是否安装
	rsrc, err := exec.LookPath("rsrc")
	if err != nil {
		fmt.Println("  rsrc 未安装，正在安装...")
		cmd := exec.Command("go", "install", "github.com/akavel/rsrc@latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("安装 rsrc 失败: %w", err)
		}
		// 重新查找
		if rsrc, err = exec.LookPath("rsrc"); err != nil {
			// 尝试 GOPATH
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
			}
			rsrc = filepath.Join(gopath, "bin", "rsrc.exe")
		}
	}

	cmd := exec.Command(rsrc, "-ico", iconPath, "-o", sysoPath)
	cmd.Dir = filepath.Join(b.ProjectRoot, "src")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsrc 执行失败: %w\n%s", err, string(output))
	}
	return nil
}

// BuildExe 构建可执行文件
func (b *Builder) BuildExe() error {
	outputPath := filepath.Join(b.ProjectRoot, ".output", b.Config.OutputName)
	srcDir := filepath.Join(b.ProjectRoot, "src")

	// 确保 .output 目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建 .output 目录失败: %w", err)
	}

	ldflags := "-s -w -H=windowsgui"
	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", outputPath, ".")
	cmd.Dir = srcDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build 失败: %w", err)
	}
	return nil
}

// BuildInstaller 构建安装程序
func (b *Builder) BuildInstaller() error {
	issPath := filepath.Join(b.ProjectRoot, "build", "installer", "setup.iss")

	// 检查 setup.iss 是否存在
	if _, err := os.Stat(issPath); os.IsNotExist(err) {
		return fmt.Errorf("setup.iss 不存在，请先运行生成步骤")
	}

	// 检查程序是否存在
	exePath := filepath.Join(b.ProjectRoot, ".output", b.Config.OutputName)
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("可执行文件不存在: %s", exePath)
	}

	iscc := `C:\Program Files (x86)\Inno Setup 6\ISCC.exe`
	if _, err := os.Stat(iscc); os.IsNotExist(err) {
		return fmt.Errorf("Inno Setup 6 未安装，请从 https://jrsoftware.org/isdl.php 下载")
	}

	cmd := exec.Command(iscc, issPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ISCC 编译失败: %w", err)
	}
	return nil
}

// Clean 清理构建/临时文件
func (b *Builder) Clean(clearOutput bool) error {
	toRemove := []string{
		filepath.Join(b.ProjectRoot, "src", "rsrc.syso"),
		filepath.Join(b.ProjectRoot, "build", "installer", "setup.iss"),
	}

	if clearOutput {
		// 清理 output 目录内容（保留目录）
		outputDir := filepath.Join(b.ProjectRoot, ".output")
		if entries, err := os.ReadDir(outputDir); err == nil {
			for _, entry := range entries {
				toRemove = append(toRemove, filepath.Join(outputDir, entry.Name()))
			}
		}
	}

	for _, f := range toRemove {
		if err := os.RemoveAll(f); err == nil {
			fmt.Printf("  已删除: %s\n", f)
		}
	}
	return nil
}
