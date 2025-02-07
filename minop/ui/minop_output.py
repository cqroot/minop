from PySide6.QtWidgets import QWidget, QVBoxLayout, QTextEdit, QSizePolicy

from minop.ui.components import MOutput


class MinopOutput(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self.setSizePolicy(QSizePolicy.Expanding, QSizePolicy.Expanding)

        self.output_widget: MOutput = MOutput()
        self.output_widget.setReadOnly(True)
        # self.output_widget.setFontFamily("Courier New")
        self.output_widget.setStyleSheet(
            # color:  # ffffff;
            # background - color:  # 5b5b5d;
            """
            MinopOutput QTextEdit {
                border: 1px solid #d9d9d9;
                border-radius: 5px;
                padding: 10px;
                font-family: "Cascadia Code", Consolas;
            }
            """
        )

        self.main_layout: QVBoxLayout = QVBoxLayout()
        self.main_layout.addWidget(self.output_widget)
        self.setLayout(self.main_layout)
