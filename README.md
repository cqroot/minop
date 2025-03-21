<div align="center">
  <h1>MinOP</h1>

  <p><i>Batch manage and control your hosts.</i></p>

  <p>
    <a href="https://github.com/cqroot/minop/actions">
      <img src="https://github.com/cqroot/minop/workflows/test/badge.svg" alt="Action Status" />
    </a>
    <a href="https://github.com/cqroot/minop/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/cqroot/minop" />
    </a>
    <a href="https://github.com/cqroot/minop/issues">
      <img src="https://img.shields.io/github/issues/cqroot/minop" />
    </a>
  </p>
</div>

## Usage

Create a file named `host.list` to configure the hosts that need to be operated.

```
root:password@host1:22
user:password@host2
```

Create a file named `minop.yaml` to configure the tasks to be executed.

```yaml
modules:
  - name: command
    command: "ip addr; lscpu"

  - name: file
    src: ./LICENSE
    dst: /root/LICENSE

  - name: script
    script: ./test.sh
```
