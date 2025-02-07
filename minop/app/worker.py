import yaml
from PySide6.QtCore import QThread, Signal
from fabric import Connection, Result
import datetime
import os


class FabricWorker(QThread):
    output_signal: Signal = Signal(str)
    finished_signal: Signal = Signal()

    def __init__(self, tasks, servers_file: str, parallel: bool = False) -> None:
        super().__init__()

        self.html_color: dict[str, str] = {
            "blue": "#007bff",
            "green": "#28a745",
            "red": "#dc3545",
            "black": "#212529",
            "orange": "#ffa500",
        }

        self.tasks: list[dict[str, str]] = tasks if tasks is not None else []

        try:
            with open(servers_file, "r", encoding="utf-8") as f:
                self.servers: dict[str, dict[str, str]] = yaml.safe_load(f)
        except Exception as e:
            print(f"Error loading configuration: {str(e)}")
            self.servers = {}

        self.parallel: bool = parallel

    def log(
        self,
        message: str,
        emoji: str,
        color: str,
        server: str = "",
        add_newline: bool = True,
    ) -> None:
        timestamp: str = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        if server:
            formatted_message: str = f"{emoji} [{timestamp}] [{server}] {message}"
        else:
            formatted_message: str = f"{emoji} [{timestamp}] {message}"

        formatted_message: str = formatted_message.replace("\n", "<br>")
        html_message: str = (
            f'<span style="color: {self.html_color.get(color, "#333333")}">{formatted_message}</span>'
            + ("<br />" if add_newline == True else "")
        )
        self.output_signal.emit(html_message)

    def log_output(self, name: str, output: str, color: str, server: str = "") -> None:
        self.log(
            f'{name}:<pre style="padding: 0; margin: 0; color: {self.html_color.get("black")}">{output}</pre>',
            "✨",
            color,
            server,
            False,
        )

    def log_error(self, message: str, server: str = "") -> None:
        self.log(message, "❎", "red", server)

    def execute_tasks(self, conn: Connection, server_name: str) -> None:
        for task in self.tasks:
            task_type: str = task.get("type")

            if task_type == "file":
                local_path: str = task.get("local_path", "")
                local_path: str = os.path.abspath(local_path)

                if not os.path.exists(local_path):
                    error_msg: str = f"Local file not found: {local_path}"
                    self.log_error(error_msg, server_name)
                    raise Exception(error_msg)

                self.log(f"Uploading {task['name']}...", "🚀", "orange", server_name)
                try:
                    conn.put(local_path, task["remote_path"])
                    self.log(
                        f"Successfully uploaded {task['name']} to {task['remote_path']}",
                        "✅",
                        "green",
                        server_name,
                    )
                except Exception as e:
                    error_msg: str = f"Failed to upload {task['name']}: {str(e)}"
                    self.log_error(error_msg, server_name)
                    raise Exception(error_msg)

            elif task_type == "command":
                self.log(f"Executing {task['name']}...", "🚀", "orange", server_name)
                result: Result = conn.run(task["command"], hide=True)
                color: str = "green"
                if result.return_code != 0:
                    color = "red"
                self.log_output(task["name"], result.stdout, color, server_name)

            else:
                error_msg: str = f"Unknown task {task['type']}"
                self.log_error(error_msg, server_name)
                raise Exception(error_msg)

    def run_parallel(self) -> None:
        try:
            for task in self.tasks:
                if task.get("type") == "file":
                    if not os.path.exists(task["local_path"]):
                        raise FileNotFoundError(
                            f"Local file not found: {task['local_path']}"
                        )

            connections: list[Connection] = []
            for server_name, config in self.servers.items():
                conn: Connection = Connection(
                    host=config["host"],
                    user=config["user"],
                    connect_kwargs={
                        "password": config["password"],
                        "look_for_keys": False,
                    },
                )
                connections.append(conn)

            self.log("Starting parallel execution on all servers...", "🖥️", "blue")

            for task in self.tasks:
                task_type: str = task.get("type")

                if task_type == "file":
                    self.log(
                        f"Uploading {task['name']} to all servers...", "🚀", "orange"
                    )
                    for conn in connections:
                        try:
                            conn.put(task["local_path"], task["remote_path"])
                            self.log(
                                f"Successfully uploaded {task['name']} to {conn.host}",
                                "✅",
                                "green",
                                conn.host,
                            )
                        except Exception as e:
                            error_msg: str = (
                                f"Failed to upload {task['name']}: {str(e)}"
                            )
                            self.log_error(error_msg, conn.host)
                            raise Exception(error_msg)

                elif task_type == "command":
                    self.log(
                        f"Executing {task['name']} on all servers...", "🚀", "orange"
                    )
                    for conn in connections:
                        try:
                            result: Result = conn.run(task["command"], hide=True)
                            color: str = "green"
                            if result.return_code != 0:
                                color = "red"
                            self.log_output(
                                task["name"], result.stdout, color, conn.host
                            )
                        except Exception as e:
                            error_msg: str = (
                                f"Failed to execute {task['name']}: {str(e)}"
                            )
                            self.log_error(error_msg, conn.host)
                            raise Exception(error_msg)

                else:
                    error_msg: str = f"Unknown task {task['type']}"
                    self.log_error(error_msg, server_name)
                    raise Exception(error_msg)

            for conn in connections:
                conn.close()
            self.log("All parallel tasks completed successfully", "✅", "blue")

        except FileNotFoundError as e:
            self.log(str(e), "red")
        except Exception as e:
            self.log(f"Parallel execution failed: {str(e)}", "red")

    def run_sequential(self) -> None:
        try:
            for server_name, config in self.servers.items():
                self.log(
                    f"Connecting to {server_name} ({config['host']})...", "🖥️", "blue"
                )

                conn: Connection = Connection(
                    host=config["host"],
                    user=config["user"],
                    connect_kwargs={
                        "password": config["password"],
                        "look_for_keys": False,
                    },
                )

                self.execute_tasks(conn, server_name)

                conn.close()
                self.log(f"Disconnected from {server_name}", "🖥️", "blue")
                self.log("-" * 50, "💠", "black")

        except Exception as e:
            self.log_error(f"Task execution stopped due to error: {str(e)}")

    def run(self) -> None:
        try:
            if self.parallel:
                self.run_parallel()
            else:
                self.run_sequential()
        finally:
            self.finished_signal.emit()
