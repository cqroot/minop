from typing import NamedTuple
from enum import Enum
import yaml


class MArgType(Enum):
    INPUT = "input"


class MArg(NamedTuple):
    name: str
    type: MArgType


class Module(NamedTuple):
    name: str
    args: list[MArg]
    actions: list[dict[str, str]]


def read_modules(name: str) -> list[Module]:
    modules: list[Module] = []
    with open(name) as f:
        for p in yaml.safe_load(f):
            args: list[MArg] = []
            for arg_obj in p.get("args", []):
                args.append(MArg(arg_obj["name"], MArgType(arg_obj["type"])))

            modules.append(Module(p["name"], args, p.get("actions", [])))

    return modules
