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

## Features

- **Batch SSH execution** — run shell commands on every host in `hosts.yaml` concurrently (default 10, tunable via `--max-procs`).
- **File & directory uploads** — built-in `copy` task with optional `backup: true` that preserves any pre-existing remote file.
- **Local commands** — run shell snippets on the minop host itself, useful for orchestration glue.
- **Interactive REPL** — `minop cli` opens a TUI where every input line is dispatched to all targets; output scrolls above a pinned input box with hint shortcuts.
- **Ansible-style task syntax** — `minop.yaml` follows the "action key as type" convention familiar to ansible users: `shell:`, `local:`, and `copy:` each pick the operation type by their key, no explicit `type:` field required.
- **Inspection commands** — `minop host` prints the parsed host tree; `minop task` lists every task in `minop.yaml` without executing.

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

A list of tasks under the `tasks` key. Each task declares exactly one
action — the action key itself picks the operation type, no separate
`type:` field is needed.

| Action key | Operation | Body |
|---|---|---|
| `shell: <cmd>` | run on every target host | command string |
| `local: <cmd>` | run on the minop host | command string |
| `copy: {...}` | upload a file or directory | nested object |

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

`copy` accepts an optional `backup: true` flag — when set, any pre-existing
remote file at `dest` is renamed to `dest.minop_bak` before the upload.

### Run Tasks

```bash
minop run                   # use ./hosts.yaml and ./minop.yaml
minop run -H ./prod-hosts.yaml  # custom hosts file
minop run -t ./deploy.yaml      # custom tasks file
minop run -p 20                 # run with concurrency 20 (default 10)
```

### Inspect Configuration

```bash
minop host                  # show the parsed hosts as a tree
minop task                  # list the tasks defined in minop.yaml
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
