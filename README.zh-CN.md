<div align="center">
  <h1>MINOP</h1>
  <i>一个简单的远程任务编排和批量执行工具。</i>
  <p align="center">
    <a href="README.md">English</a>
    ·
    <a href="README.zh-CN.md">简体中文</a>
    <br />
  </p>
  <p>
    <a href="https://github.com/cqroot/minop/actions">
      <img src="https://github.com/cqroot/minop/workflows/test/badge.svg" alt="Action Status" />
    </a>
    <a href="https://codecov.io/gh/cqroot/minop">
      <img src="https://codecov.io/gh/cqroot/minop/branch/main/graph/badge.svg" alt="Codecov" />
    </a>
    <a href="https://goreportcard.com/report/github.com/cqroot/minop">
      <img src="https://goreportcard.com/badge/github.com/cqroot/minop.svg" alt="Go Report Card" />
    </a>
    <a href="https://pkg.go.dev/github.com/cqroot/minop">
      <img src="https://pkg.go.dev/badge/github.com/cqroot/minop.svg" alt="Go Reference" />
    </a>
    <a href="https://github.com/cqroot/minop/tags">
      <img src="https://img.shields.io/github/v/tag/cqroot/minop" alt="Git tag" />
    </a>
    <a href="https://github.com/cqroot/minop/blob/main/go.mod">
      <img src="https://img.shields.io/github/go-mod/go-version/cqroot/minop" alt="Go Version" />
    </a>
    <a href="https://github.com/cqroot/minop/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/cqroot/minop" />
    </a>
    <a href="https://github.com/cqroot/minop/issues">
      <img src="https://img.shields.io/github/issues/cqroot/minop" />
    </a>
    <a href="https://github.com/cqroot/minop/releases">
      <img src="https://img.shields.io/github/downloads/cqroot/minop/total?label=github%20downloads" />
    </a>
  </p>
  <hr>
</div>

## 安装

### 从源码安装

要安装 `minop`，请确保你已经安装了 Go 然后执行：

```bash
go install github.com/cqroot/minop@latest
```

### 下载编译好的二进制

从 release 界面下载对应平台的二进制文件，并将其所在路径加入到环境变量中。

## 特性

- **批量 SSH 执行** — 在 `hosts.yaml` 的每台主机上并发运行 shell 命令（默认并发 10，可通过 `--max-procs` 调整）。
- **文件与目录上传** — 内置 `copy` 任务，可选 `backup: true` 自动保留远端已存在的文件。
- **本地命令** — 在 minop 主机自身执行 shell 片段，适合编排胶水代码。
- **交互式 REPL** — `minop cli` 启动 TUI，输入回车即分发到所有目标；输出区在上方独立滚动，输入框固定在底部，并显示快捷键提示。
- **兼容 ansible 任务语法** — `minop.yaml` 采用 ansible 用户熟悉的"动作 key 即类型"约定：`shell:`、`local:`、`copy:` 各自的动作 key 即决定操作类型，不需要显式 `type:` 字段。
- **检查命令** — `minop host` 以树状结构展示已解析的主机；`minop task` 列出 `minop.yaml` 中的所有任务但不执行。

## 用法

minop 读取两个配置文件：

| 文件          | 参数                | 作用                                                                 |
| ------------- | ------------------- | -------------------------------------------------------------------- |
| `hosts.yaml`  | `-H` / `--hosts-file` | 按角色分组的机器清单。跨多次运行基本不变，通常纳入版本控制。       |
| `minop.yaml`  | `-t` / `--task`       | 任务编排。每次运行可能不一样。                                       |

默认在当前工作目录下查找这两个文件。

### 主机文件（`hosts.yaml`）

扁平的 YAML，键为角色名，值为该角色下的主机字符串列表，格式为
`<user>:<password>@<address>:<port>`。`all` 是特殊角色：所有未指定
`role` 的任务都会作用于它。

```yaml
all:
  - root:PASSWORD@127.0.0.1:8001

main:
  - root:PASSWORD@127.0.0.1:8002
  - root:PASSWORD@127.0.0.1:8003
```

### 任务文件（`minop.yaml`）

`tasks` 键下是一个任务列表，每个任务声明**恰好一个**动作 —— 动作
key 本身决定操作类型，不需要额外的 `type:` 字段。

| 动作 key | 操作 | 值 |
|---|---|---|
| `shell: <cmd>` | 在每台目标主机上执行 | 命令字符串 |
| `local: <cmd>` | 在 minop 主机上执行 | 命令字符串 |
| `copy: {...}` | 上传文件或目录 | 嵌套对象 |

```yaml
tasks:
  - name: Copy a file to the remote host
    copy:
      src: test.txt
      dest: /tmp/test.txt

  - name: Copy a directory to the remote host
    copy:
      src: testdir
      dest: /tmp/testdir

  - name: List /tmp on the remote host
    shell: ls /tmp
```

`copy` 可选 `backup: true`——上传前将远端已存在的文件重命名为
`dest.minop_bak`。

### 执行任务

```bash
minop run                       # 使用 ./hosts.yaml 和 ./minop.yaml
minop run -H ./prod-hosts.yaml  # 自定义主机文件
minop run -t ./deploy.yaml      # 自定义任务文件
minop run -p 20                 # 并发度调到 20（默认 10）
```

### 检查配置

```bash
minop host                  # 以树状形式展示解析后的主机
minop task                  # 列出 minop.yaml 中定义的任务
```

### 交互式 CLI

```bash
minop cli                   # 进入交互式命令行：输入一条命令，回车即可在所有机器上执行
```

输入 `exit`（或 `quit`）退出，输入 `help` 查看内置快捷键。

## 贡献指南

欢迎贡献！您可以随时提交问题（issue）来报告错误、提出新功能建议，或提交拉取请求（pull request）。

## 许可证

本项目为开源项目，采用 [GPL-3.0 许可证](LICENSE)。
