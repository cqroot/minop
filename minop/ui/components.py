from PySide6.QtCore import Qt
from PySide6.QtWidgets import (
    QLineEdit,
    QListWidget,
    QPushButton,
    QWidget,
)


class MButton(QPushButton):
    pass


class MCard(QWidget):
    pass


class MInput(QLineEdit):
    pass


class MSidebar(QListWidget):
    def __init__(self, items: list[str]) -> None:
        super().__init__()
        self.addItems(items)

        self.setMinimumWidth(self.sizeHintForColumn(0))
        self.setFixedWidth(self.sizeHintForColumn(0) + 2 * self.frameWidth() + 32)
        self.setHorizontalScrollBarPolicy(Qt.ScrollBarAlwaysOff)
        self.setVerticalScrollBarPolicy(Qt.ScrollBarAlwaysOff)

    def action_toggle(self) -> None:
        if self.isVisible():
            self.setVisible(False)
        else:
            self.setVisible(True)


class MWidget(QWidget):
    pass


def stylesheet(path: str) -> str:
    with open(path, "r") as f:
        qss: str = f.read()

    colors: dict[str, str] = {}
    prop_end: int = -1
    for i, line in enumerate(qss.splitlines()):
        if line == ":root {":
            continue

        if line == "}":
            prop_end = i
            break

        key: str = line.split(":")[0].strip()
        value: str = line.split(":")[1].rstrip(";").strip()
        colors[key] = value

    qss = "\n".join(qss.splitlines()[prop_end + 1 : -1])
    for key, value in colors.items():
        qss = qss.replace("var({})".format(key), value)

    return qss
