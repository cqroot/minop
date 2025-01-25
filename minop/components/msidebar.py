#!/usr/bin/env python3

from PySide6.QtCore import Qt
from PySide6.QtWidgets import QListWidget


class MSidebar(QListWidget):
    def __init__(self, items: list[str]) -> None:
        super().__init__()
        self.addItems(items)

        self.setMinimumWidth(self.sizeHintForColumn(0))
        self.setFixedWidth(self.sizeHintForColumn(0) + 2 * self.frameWidth() + 32)
        self.setHorizontalScrollBarPolicy(Qt.ScrollBarAlwaysOff)
        self.setVerticalScrollBarPolicy(Qt.ScrollBarAlwaysOff)
        self.setStyleSheet(
            """
            {name} {{
                outline: none;
                background-color: #2e3235;
                padding: 8px;
            }}

            {name}::item {{
                color: #dddddd;
                background-color: #2e3235;
                border: transparent;
                border-radius: 5px;
                padding: 8px;
            }}

            {name}::item:selected {{
                color: #ffffff;
                background-color: #464646;
            }}
            """.format(
                name=type(self).__name__
            )
        )

    def action_toggle(self) -> None:
        if self.isVisible():
            self.setVisible(False)
        else:
            self.setVisible(True)
