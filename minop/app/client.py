import paramiko
from paramiko.channel import ChannelStdinFile, ChannelFile, ChannelStderrFile

from minop.app.hosts import Host


class Client:
    def __init__(self, host: Host) -> None:
        self.__client: paramiko.SSHClient = paramiko.SSHClient()
        self.__client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        self.__client.connect(
            host.addr, port=host.port, username=host.user, password=host.password
        )

    def close(self) -> None:
        self.__client.close()

    def exec_command(
        self, cmd: str
    ) -> tuple[ChannelStdinFile, ChannelFile, ChannelStderrFile]:
        return self.__client.exec_command(cmd)
