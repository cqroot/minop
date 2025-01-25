from typing import NamedTuple
import yaml


class Host(NamedTuple):
    addr: str
    port: int
    user: str
    password: str


def read_hosts(path: str) -> list[Host]:
    hosts: list[Host] = []
    with open(path, "r") as f:
        for host_obj in yaml.safe_load(f):
            hosts.append(
                Host(
                    host_obj.get("addr", "127.0.0.1"),
                    host_obj.get("port", 22),
                    host_obj.get("user", "root"),
                    host_obj.get("password"),
                )
            )

    return hosts
