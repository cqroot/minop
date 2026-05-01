<div align="center">
  <h1>MINOP</h1>
  <i>A simple tool for remote task orchestration and batch execution.</i>
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

## Installation

### Install from Source

To install `minop` from source, ensure you have Go installed and run:

```bash
go install github.com/cqroot/minop@latest
```

### Download Pre-compiled Binaries

Download the binary for your platform from the releases page and add its directory to your system's PATH.

## Usage

### Create the Config File

Create a file named `minop.yaml` in YAML format. This file contains both hosts and tasks.

#### Hosts Section

The file should contain groups of hosts under the `hosts` key, where each group is a list of host strings in the format `<user>:<password>@<address>:<port>`. Example:

```yaml
hosts:
  all:
    - root:asdf@127.0.0.1:8001

  main:
    - root:asdf@127.0.0.1:8002
    - root:asdf@127.0.0.1:8003
```

Hosts listed under a specific section header (like `main` in the example) will be assigned to that role.

#### Tasks Section

Add your tasks under the `tasks` key:

```yaml
tasks:
  - name: Copy file to /root on the remote host
    copy: test.txt
    to: /root/test.txt

  - name: Copy dir to /root on the remote host
    copy: testdir
    to: /root/testdir

  - name: List /root
    shell: ls /root
```

### Execute Tasks

Run the following command to execute tasks on the remote hosts:

```bash
minop
```

This will load `./minop.yaml` by default. You can specify a different config file:

```bash
minop -c /path/to/config.yaml
```

### Interactive CLI

Start an interactive CLI mode to execute commands on remote hosts:

```bash
minop cli
```

You can also specify a different config file:

```bash
minop cli -c /path/to/config.yaml
```

## Contributing

Contributions are welcome! Feel free to open an issue to report bugs, suggest new features, or submit a pull request.

## License

This project is open source, licensed under the [GPL-3.0 License](LICENSE).