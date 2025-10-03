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
      <img src="https://goreportcard.com/badge/github.com/cqroot/minop" alt="Go Report Card" />
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

要从源码安装 `minop`，请确保你已经安装了 Go 然后执行：

```bash
go install github.com/cqroot/minop@latest
```

### 下载编译好的二进制

从 release 界面下载对应平台的二进制文件，并将其所在路径加入到环境变量中。

## 用法

## 创建主机列表文件 hosts

创建一个名为 `hosts` 的文件，它的每一行都是一个远程主机，格式为 `<user>:<password>@<address>:<port>`。示例如下：

```
root:asdf@127.0.0.1:8001

[main]
root:asdf@127.0.0.1:8002
root:asdf@127.0.0.1:8003
```

位于特定章节下的主机，会归属于指定的角色下（如上述示例中，后两个主机的角色为 `main`）。

### 交互式 CLI

直接执行工具会加载当前路径下的 `hosts` 文件，然后进入交互式 CLI 界面。你可以在此执行你想要在远程主机执行的命令。

```bash
minop
```

### 执行任务文件

如果你想通过非交互式方式执行预设好的任务，你可以先创建一个 yaml 文件。比如 `minop.yaml`

```yaml
- name: Copy file to /root on the remote host
  copy: test.txt
  to: /root/test.txt

- name: Copy dir to /root on the remote host
  copy: testdir
  to: /root/testdir

- name: List /root
  shell: ls /root
```

然后执行以下命令来加载当前路径下的 `hosts` 文件和 `minop.yaml` 文件，在远程主机执行我们编排好的任务：

```bash
minop -t minop.yaml
```

## 贡献指南
欢迎贡献！您可以随时提交问题（issue）来报告错误、提出新功能建议，或提交拉取请求（pull request）。

## 许可证

本项目为开源项目，采用 [GPL-3.0 许可证](LICENSE)。