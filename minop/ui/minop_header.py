from PySide6.QtCore import Qt
from PySide6.QtWidgets import QWidget, QHBoxLayout
import qtawesome

from minop.ui.components import MButton


class MinopHeader(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self.sidebar_toggle_button: MButton = MButton(
            qtawesome.icon("fa5s.bars", color="#2e3235"), ""
        )
        self.sidebar_toggle_button.setToolTip("Toggle Sidebar")
        self.sidebar_toggle_button.setFocusPolicy(Qt.FocusPolicy.NoFocus)

        self.toggle_parallel_button: MButton = MButton(
            qtawesome.icon("fa5s.toggle-off", color="#2e3235"), ""
        )
        self.toggle_parallel_button.setToolTip("Toggle Parallel")
        self.toggle_parallel_button.setCheckable(True)
        self.toggle_parallel_button.setFocusPolicy(Qt.FocusPolicy.NoFocus)
        self.toggle_parallel_button.toggled.connect(
            lambda state: self.toggle_parallel_button.setIcon(
                qtawesome.icon(
                    "fa5s.toggle-on" if state else "fa5s.toggle-off", color="#2e3235"
                )
            )
        )

        self.run_button: MButton = MButton(
            qtawesome.icon("fa5s.play", color="#2e3235"), ""
        )
        self.run_button.setToolTip("Run")
        self.run_button.setFocusPolicy(Qt.FocusPolicy.NoFocus)

        self.main_layout: QHBoxLayout = QHBoxLayout()
        self.main_layout.addWidget(self.sidebar_toggle_button)
        self.main_layout.addStretch()
        self.main_layout.addWidget(self.toggle_parallel_button)
        self.main_layout.addWidget(self.run_button)
        self.setLayout(self.main_layout)
