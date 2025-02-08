import os

import yaml
from PySide6.QtCore import QThread
from PySide6.QtWidgets import QGridLayout
from jinja2 import Template

from minop.app.modules import Module
from minop.app.worker import FabricWorker
from minop.ui.components import MWidget, stylesheet, MSidebar, MInput
from minop.ui.minop_container import MinopContainer
from minop.ui.minop_output import MinopOutput
from minop.ui.minop_header import MinopHeader
from minop.ui.stylesheet import get_stylesheet


class MinopWindow(MWidget):
    def __init__(self, modules: list[Module], config_path: str) -> None:
        super().__init__()
        self.modules: list[Module] = modules
        self.config_path: str = config_path
        self.thread: QThread = QThread()

        self.__setup_ui()
        self.__connect_all()
        self.setStyleSheet(get_stylesheet())

    def __setup_ui(self) -> None:
        self.setWindowTitle("MinOP - {}".format(self.config_path))

        module_names: list[str] = []
        for module in self.modules:
            module_names.append(module.name)
        self.minop_sidebar: MSidebar = MSidebar(module_names)
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
        self.minop_header.run_button.clicked.connect(self.start_task)

    def get_curr_args(self) -> dict:
        args: dict = {}
        input_widgets: iter[MInput] = self.minop_container.currentWidget().findChildren(
            MInput
        )
        for input_widget in input_widgets:
            args[input_widget.property("name")] = input_widget.text()
        return args

    def start_task(self) -> None:
        self.minop_header.run_button.setEnabled(False)
        self.minop_output.output_widget.clear()

        args: dict = self.get_curr_args()
        tasks: list[dict[str, str]] = []
        for action in self.modules[self.minop_sidebar.currentRow()].actions:
            task = {}
            for k, v in action.items():
                templator: Template = Template(v)
                task[k] = templator.render(args)
            tasks.append(task)

        self.worker = FabricWorker(
            tasks=tasks,
            servers_file=self.config_path,
            parallel=self.minop_header.toggle_parallel_button.isChecked(),
        )
        self.worker.output_signal.connect(self.update_output)
        self.worker.finished_signal.connect(self.task_finished)
        self.worker.start()

    def update_output(self, text):
        self.minop_output.output_widget.insertHtml(text)
        self.minop_output.output_widget.verticalScrollBar().setValue(
            self.minop_output.output_widget.verticalScrollBar().maximum()
        )

    def task_finished(self):
        self.minop_header.run_button.setEnabled(True)
        print(self.minop_output.output_widget.toHtml())
