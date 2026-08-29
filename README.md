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

## Installation

### Install from Source

To install `minop` from source, ensure you have Go installed and run:

```bash
go install github.com/cqroot/minop@latest
```

### Download Pre-compiled Binaries

Download the binary for your platform from the releases page and add its directory to your system's PATH.

## Usage

minop reads two configuration files:

| File          | Flag              | Purpose                                                                  |
| ------------- | ----------------- | ------------------------------------------------------------------------ |
| `hosts.yaml`  | `-H` / `--hosts-file` | Hosts grouped by role. Stable across runs; usually version-controlled. |
| `minop.yaml`  | `-t` / `--task`       | Tasks to execute. Changes from run to run.                             |

Both files are looked up in the current working directory by default.

### Hosts File (`hosts.yaml`)

A flat YAML map from role name to a list of host strings in the
format `<user>:<password>@<address>:<port>`. The role named `all` is
special: every task targets it unless the task specifies a different
role.

```yaml
all:
  - root:PASSWORD@127.0.0.1:8001

main:
  - root:PASSWORD@127.0.0.1:8002
  - root:PASSWORD@127.0.0.1:8003
```

### Tasks File (`minop.yaml`)

A list of tasks under the `tasks` key. Each task supports one of three
operation types: `copy`, `shell`, or `local`.

```yaml
tasks:
  - name: Copy file to the remote host
    copy: test.txt
    to: /tmp/test.txt

  - name: Copy directory to the remote host
    copy: testdir
    to: /tmp/testdir

  - name: List /tmp on the remote host
    shell: ls /tmp
```

### Run Tasks

```bash
minop                       # use ./hosts.yaml and ./minop.yaml
minop -H ./prod-hosts.yaml  # custom hosts file
minop -t ./deploy.yaml      # custom tasks file
minop -p 20                 # run with concurrency 20 (default 10)
```

### Inspect Configuration

```bash
minop host                  # show the parsed hosts as a tree
minop task                  # list the tasks defined in minop.yaml
minop check                 # validate minop.yaml and list its tasks
```

### Interactive CLI

```bash
minop cli                   # REPL: type a command, hit Enter, run on every host
```

Type `exit` (or `quit`) to leave the CLI, and `help` to see built-in
shortcuts.

## Contributing

Contributions are welcome! Feel free to open an issue to report bugs, suggest new features, or submit a pull request.

## License

This project is open source, licensed under the [GPL-3.0 License](LICENSE).
