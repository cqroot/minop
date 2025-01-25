import html
from jinja2 import Template

from minop.app.client import Client


class Command:
    def __init__(self, cmd: str):
        self.cmd: str = cmd
        pass

    def run(self, client: Client, args: dict) -> str:
        templator: Template = Template(self.cmd)
        self.cmd = templator.render(args)

        _, stdout, stderr = client.exec_command(self.cmd)
        stdout_str: str = stdout.read().decode()
        stderr_str: str = stderr.read().decode()
        result_str: str = '<font color="#89ca78">===> COMMAND:</font> {}<br />'.format(
            self.cmd
        )

        if stdout_str:
            result_str += (
                '<font color="#89ca78">===> STDOUT:</font><pre>{}</pre>'.format(
                    html.escape(stdout_str)
                )
            )
        if stderr_str:
            result_str += (
                '<font color="#89ca78">===> STDERR:</font><pre>{}</pre>'.format(
                    html.escape(stderr_str)
                )
            )

        return result_str
