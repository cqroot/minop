from abc import abstractmethod

from minop.actions.command import Command
from minop.app.client import Client


class Action:
    @abstractmethod
    def run(self, client: Client, args: dict):
        pass


def get_actions(params: dict):
    if params.get("action") == "command":
        return Command(params.get("command"))
    return None
