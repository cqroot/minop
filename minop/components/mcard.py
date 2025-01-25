from PySide6.QtWidgets import QWidget


class MCard(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self.setStyleSheet(
            """
            {name} {{
                background-color: #000000;
                border: 1px solid #d9d9d9;
                border-radius: 5px;
            }}
            """.format(
                name=type(self).__name__
            )
        )
