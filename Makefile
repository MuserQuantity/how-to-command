.PHONY: build install clean

# 默认目标
all: build

# 编译为"how2cmd"二进制文件
build:
	go build -o how2cmd

# 安装到GOPATH/bin
install:
	go install -v

# 使用自定义输出名称安装
install-custom:
	go build -o $(GOPATH)/bin/how2cmd

# 清理生成的文件
clean:
	rm -f how2cmd

# 帮助
help:
	@echo "可用目标:"
	@echo "  build         - 编译为how2cmd可执行文件"
	@echo "  install       - 安装到GOPATH/bin (保持原名how-to-command)"
	@echo "  install-custom - 安装到GOPATH/bin (使用how2cmd名称)"
	@echo "  clean         - 删除生成的文件" 