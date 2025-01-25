#!/usr/bin/env python3

from PySide6.QtGui import QIcon, QPixmap
from PySide6.QtWidgets import QPushButton


class MButton(QPushButton):
    def __init__(self, content: QIcon | QPixmap | str) -> None:
        if isinstance(content, QIcon):
            super().__init__(content, "")
        else:
            super().__init__(content)

        self.setStyleSheet(
            """
            {name} {{
                color: #0a0a0a;
                background-color: #ffffff;
                padding: 5px;
                border: none;
                border-radius: 5px;
            }}

            {name}:hover {{
                background-color: #cacbcd;
            }}

            {name}:pressed,
            {name}:checked {{
                background-color: #babbbc;
            }}
            """.format(
                name=type(self).__name__
            )
        )
