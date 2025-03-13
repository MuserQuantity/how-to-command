package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	envURL   = "HOWTOCOMMAND_URL"
	envModel = "HOWTOCOMMAND_MODEL"
	envToken = "HOWTOCOMMAND_TOKEN"
	cmdName  = "how2cmd"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s <描述你想要执行的操作>\n", cmdName)
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s 如何查看当前目录下的所有文件\n", cmdName)
		fmt.Fprintf(os.Stderr, "  %s 如何压缩当前文件夹\n", cmdName)
		fmt.Fprintf(os.Stderr, "\n环境变量:\n")
		fmt.Fprintf(os.Stderr, "  %s: API地址\n", envURL)
		fmt.Fprintf(os.Stderr, "  %s: 模型名称\n", envModel)
		fmt.Fprintf(os.Stderr, "  %s: API令牌\n", envToken)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// 获取环境变量
	apiURL := os.Getenv(envURL)
	model := os.Getenv(envModel)
	token := os.Getenv(envToken)

	// 验证环境变量
	var missingVars []string
	if apiURL == "" {
		missingVars = append(missingVars, envURL)
	}
	if model == "" {
		missingVars = append(missingVars, envModel)
	}
	if token == "" {
		missingVars = append(missingVars, envToken)
	}

	if len(missingVars) > 0 {
		fmt.Fprintf(os.Stderr, "错误: 以下环境变量未设置:\n")
		for _, v := range missingVars {
			fmt.Fprintf(os.Stderr, "  - %s\n", v)
		}
		fmt.Fprintf(os.Stderr, "\n环境变量设置示例:\n")
		fmt.Fprintf(os.Stderr, "  - Windows (PowerShell):\n")
		fmt.Fprintf(os.Stderr, "    $env:%s = \"https://api.openai.com/v1/chat/completions\"\n", envURL)
		fmt.Fprintf(os.Stderr, "    $env:%s = \"gpt-4o-2024-05-13\"\n", envModel)
		fmt.Fprintf(os.Stderr, "    $env:%s = \"sk-your-api-token\"\n", envToken)
		fmt.Fprintf(os.Stderr, "\n  - Linux/macOS:\n")
		fmt.Fprintf(os.Stderr, "    export %s=\"https://api.openai.com/v1/chat/completions\"\n", envURL)
		fmt.Fprintf(os.Stderr, "    export %s=\"gpt-4o-2024-05-13\"\n", envModel)
		fmt.Fprintf(os.Stderr, "    export %s=\"sk-your-api-token\"\n", envToken)
		os.Exit(1)
	}

	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法获取当前目录: %v\n", err)
		os.Exit(1)
	}

	// 获取当前目录下的文件列表
	files, err := getDirContents(currentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法获取目录内容: %v\n", err)
		os.Exit(1)
	}

	// 构建用户查询
	query := flag.Arg(0)

	// 调用AI接口
	result, err := askAI(apiURL, model, token, currentDir, files, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 调用AI接口失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

// getDirContents 获取指定目录下的文件列表
func getDirContents(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}
