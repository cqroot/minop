from PySide6.QtWidgets import QWidget, QHBoxLayout
import qtawesome

from minop.components.mbutton import MButton


class MinopHeader(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self.sidebar_toggle_button: MButton = MButton(qtawesome.icon("fa5s.bars"))
        self.sidebar_toggle_button.setToolTip("Toggle Sidebar")

        self.run_button: MButton = MButton(qtawesome.icon("fa5s.play"))
        self.run_button.setToolTip("Run")

        self.settings_button: MButton = MButton(qtawesome.icon("fa5s.cog"))
        self.settings_button.setToolTip("Settings")

        self.main_layout: QHBoxLayout = QHBoxLayout()
        self.main_layout.addWidget(self.sidebar_toggle_button)
        self.main_layout.addStretch()
        self.main_layout.addWidget(self.run_button)
        self.main_layout.addWidget(self.settings_button)
        self.setLayout(self.main_layout)
