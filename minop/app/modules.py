from typing import NamedTuple
from enum import Enum
import yaml

from minop.actions.action import Action, get_actions


class MArgType(Enum):
    INPUT = "input"


class MArg(NamedTuple):
    name: str
    type: MArgType


class Module(NamedTuple):
    name: str
    args: list[MArg]
    actions: list[Action]


def read_modules(name: str) -> list[Module]:
    modules: list[Module] = []
    with open(name) as f:
        for p in yaml.safe_load(f):
            args: list[MArg] = []
            for arg_obj in p.get("args", []):
                args.append(MArg(arg_obj["name"], MArgType(arg_obj["type"])))

            actions: list[Action] = []
            for action_obj in p.get("actions", []):
                actions.append(get_actions(action_obj))

            modules.append(Module(p["name"], args, actions))

    return modules
