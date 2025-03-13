package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func askAI(apiURL, model, token, currentDir string, files []string, query string) (string, error) {
	// 获取当前操作系统
	osType := runtime.GOOS

	// 获取系统架构
	osArch := runtime.GOARCH

	// 获取操作系统详细描述
	var osDescription string
	switch osType {
	case "windows":
		osDescription = "Windows系统"
	case "darwin":
		osDescription = "macOS系统"
	case "linux":
		osDescription = "Linux系统"
	default:
		osDescription = fmt.Sprintf("未知系统(%s)", osType)
	}

	systemPrompt := fmt.Sprintf(`你是一个专业的命令行助手。

系统信息:
- 操作系统类型: %s (%s)
- 系统架构: %s
- 当前工作目录: %s
- 目录下的文件: %s

请根据用户的描述和当前系统环境，给出最合适的命令。
要求：
1. 只输出命令本身，不要解释
2. 根据不同系统给出对应命令:
   - Windows: 优先使用PowerShell命令，如需多个命令用分号(;)连接
   - Linux/macOS: 使用Shell命令，如需多个命令用 && 连接
3. 命令应该是安全的，不会造成数据丢失
4. 尽量使用系统原生命令，避免依赖可能不存在的工具
5. 如果用户的描述不够清晰，只输出"请提供更详细的描述"`,
		osType, osDescription, osArch, currentDir, strings.Join(files, ", "))

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: query},
	}

	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}

	// 打印请求信息，帮助调试
	fmt.Fprintf(os.Stderr, "发送请求到: %s\n", apiURL)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("服务器返回错误代码: %d %s", resp.StatusCode, resp.Status)
	}

	// 检查Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		// 读取响应内容用于调试
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)
		// 截断过长的响应
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return "", fmt.Errorf("服务器返回了非JSON内容 (Content-Type: %s): %s", contentType, bodyStr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// 如果JSON解析失败，打印部分响应内容以便调试
		bodyStr := string(body)
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		return "", fmt.Errorf("JSON解码失败: %v\n响应内容: %s", err, bodyStr)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API错误: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return result.Choices[0].Message.Content, nil
}
