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

### Create the Hosts File

Create a file named `hosts`. Each line in this file should represent a remote host in the format `<user>:<password>@<address>:<port>`. Example:

```
root:asdf@127.0.0.1:8001

[main]
root:asdf@127.0.0.1:8002
root:asdf@127.0.0.1:8003
```

Hosts listed under a specific section header (like `[main]` in the example) will be assigned to that role.

### Interactive CLI

Running the tool directly will load the `hosts` file from the current directory and start an interactive CLI. Here, you can execute commands you wish to run on the remote hosts.

```bash
minop
```

### Execute Task Files

To execute predefined tasks non-interactively, first create a YAML file, for example `minop.yaml`:

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

Then, run the following command to load the `hosts` file and the `minop.yaml` file from the current directory, executing the orchestrated tasks on the remote hosts:

```bash
minop -t minop.yaml
```

## Contributing

Contributions are welcome! Feel free to open an issue to report bugs, suggest new features, or submit a pull request.

## License

This project is open source, licensed under the [GPL-3.0 License](LICENSE).