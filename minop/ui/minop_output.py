from PySide6.QtWidgets import QWidget, QVBoxLayout

from minop.ui.components import MOutput


class MinopOutput(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self.output_widget: MOutput = MOutput()
        self.output_widget.setReadOnly(True)

        self.main_layout: QVBoxLayout = QVBoxLayout()
        self.main_layout.addWidget(self.output_widget)
        self.setLayout(self.main_layout)
