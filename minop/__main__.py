import os
import sys
from PySide6.QtWidgets import QApplication

from minop.ui.minop_window import MinopWindow
from minop.app.modules import Module, read_modules


def main() -> None:
    modules: list[Module] = read_modules(
        os.path.join(os.path.dirname(os.path.realpath(__file__)), "modules.yaml")
    )

    cwd: str = os.getcwd()
    config_path = os.path.join(cwd, "minop.yaml")

    app: QApplication = QApplication([])
    minop_window: MinopWindow = MinopWindow(modules, config_path)
    minop_window.resize(800, 600)
    minop_window.show()
    sys.exit(app.exec())


main()
