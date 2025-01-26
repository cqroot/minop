from PySide6.QtWidgets import QWidget, QVBoxLayout, QTextEdit, QSizePolicy


class MinopOutput(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self.setSizePolicy(QSizePolicy.Expanding, QSizePolicy.Expanding)

        self.output_widget: QTextEdit = QTextEdit()
        self.output_widget.setReadOnly(True)
        self.output_widget.setStyleSheet(
            """
            MinopOutput QTextEdit {
                color: #ffffff;
                background-color: #5b5b5d;
                border-radius: 5px;
                padding: 10px;
                font-family: Consolas;
            }
            """
        )

        self.main_layout: QVBoxLayout = QVBoxLayout()
        self.main_layout.addWidget(self.output_widget)
        self.setLayout(self.main_layout)
