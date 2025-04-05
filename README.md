# how-to-command (how2cmd)

一个基于大语言模型的命令行工具，帮助用户快速找到所需的命令。

## 安装

```bash
git clone https://github.com/MuserQuantity/how-to-command.git
# windows下
go build -o $env:GOPATH\bin\how2cmd.exe
# linux下
go build -o $GOPATH/bin/how2cmd
# 非root用户，需要使用sudo
sudo $(which go) build -o $GOPATH/bin/how2cmd
```

## 环境变量设置

使用前需要设置以下环境变量：

- `HOWTOCOMMAND_URL`: API 接口地址（例如：https://api.openai.com/v1/chat/completions）
- `HOWTOCOMMAND_MODEL`: 使用的模型名称（例如：gpt-4o-2024-05-13）
- `HOWTOCOMMAND_TOKEN`: API 访问令牌

### windows下的设置
```
set HOWTOCOMMAND_URL=https://api.openai.com/v1/chat/completions
set HOWTOCOMMAND_MODEL=gpt-4o-2024-05-13
set HOWTOCOMMAND_TOKEN=sk-proj-1234567890
```

### linux下的设置
```
export HOWTOCOMMAND_URL=https://api.openai.com/v1/chat/completions
export HOWTOCOMMAND_MODEL=gpt-4o-2024-05-13
export HOWTOCOMMAND_TOKEN=sk-proj-1234567890
```

## 使用方法

```bash
# 查看帮助
how2cmd --help

# 示例：查询如何列出当前目录下的所有文件
how2cmd 如何查看当前目录下的所有文件

# 示例：查询如何压缩文件夹
how2cmd 如何压缩当前文件夹
```

## 特点

1. 自动获取当前工作目录和文件列表作为上下文
2. 支持中文自然语言描述
3. 只输出命令，简洁直观
4. 支持任何兼容 OpenAI API 的大语言模型服务

## 注意事项

1. 确保已正确设置所有必需的环境变量
2. 建议在执行命令前仔细检查输出的命令是否符合预期
3. 该工具仅提供建议，请在理解命令的作用后再执行
