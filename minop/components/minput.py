#!/usr/bin/env python3

from PySide6.QtWidgets import QLineEdit


class MInput(QLineEdit):
    def __init__(self) -> None:
        super().__init__()

        self.setStyleSheet(
            """
            {name} {{
                background-color: #ffffff;
                border: 1px solid #d9d9d9;
                border-radius: 5px;
                padding: 3px;
            }}

            {name}:hover {{
                border-bottom: 1px solid #5b5b5d;
            }}

            {name}:focus {{
                background-color: #f6f6f6;
                border-bottom: 1px solid #2e3235;
            }}
            """.format(
                name=type(self).__name__
            ),
        )
