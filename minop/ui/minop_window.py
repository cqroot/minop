import os

from PySide6.QtWidgets import (
    QWidget,
    QGridLayout,
)

from minop.app.client import Client
from minop.app.hosts import read_hosts, Host
from minop.app.modules import Module
from minop.components.minput import MInput
from minop.ui.minop_container import MinopContainer
from minop.ui.minop_output import MinopOutput
from minop.ui.minop_header import MinopHeader
from minop.ui.minop_sidebar import MinopSidebar


class MinopWindow(QWidget):
    def __init__(self, modules: list[Module], config_path: str) -> None:
        super().__init__()
        self.modules: list[Module] = modules
        self.config_path: str = config_path

        self.__setup_ui()
        self.__connect_all()

    def __setup_ui(self) -> None:
        self.setWindowTitle("MinOP - {}".format(self.config_path))
        self.setStyleSheet(
            """
            MinopWindow {
                background-color: #ffffff
            }
            """
        )

        module_names: list[str] = []
        for module in self.modules:
            module_names.append(module.name)
        self.minop_sidebar: MinopSidebar = MinopSidebar(module_names)
        self.minop_sidebar.setCurrentRow(0)

        self.minop_header: MinopHeader = MinopHeader()

        self.minop_container: MinopContainer = MinopContainer(self.modules)
        self.minop_output: MinopOutput = MinopOutput()

        self.main_layout: QGridLayout = QGridLayout()
        self.main_layout.setContentsMargins(0, 0, 0, 0)
        self.main_layout.addWidget(self.minop_sidebar, 0, 0, 3, 1)
        self.main_layout.addWidget(self.minop_header, 0, 1, 1, 1)
        self.main_layout.addWidget(self.minop_container, 1, 1, 1, 1)
        self.main_layout.addWidget(self.minop_output, 2, 1, 1, 1)
        self.setLayout(self.main_layout)

    def change_current_module(self, index: int) -> None:
        self.minop_container.setCurrentIndex(index)

    def __connect_all(self) -> None:
        self.minop_sidebar.currentRowChanged.connect(self.change_current_module)
        self.minop_header.sidebar_toggle_button.clicked.connect(
            self.minop_sidebar.action_toggle
        )
        self.minop_header.run_button.clicked.connect(self.action_run)

    def get_curr_args(self) -> dict:
        args: dict = {}
        input_widgets: iter[MInput] = self.minop_container.currentWidget().findChildren(
            MInput
        )
        for input_widget in input_widgets:
            args[input_widget.property("name")] = input_widget.text()
        return args

    def action_run(self) -> None:
        args: dict = self.get_curr_args()

        hosts: list[Host] = read_hosts(self.config_path)
        for host in hosts:
            host_info: str = "{}@{}:{}".format(host.user, host.addr, host.port)
            self.minop_output.output_widget.append(
                '<font color="#61afef">{}<br />{}</font>'.format(
                    host_info, "=" * len(host_info)
                ),
            )

            client: Client = Client(host)
            for action in self.modules[self.minop_sidebar.currentRow()].actions:
                output: str = action.run(client, args)
                self.minop_output.output_widget.append(output)
            client.close()
