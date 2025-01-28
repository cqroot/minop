from PySide6.QtWidgets import (
    QLabel,
    QStackedWidget,
    QGridLayout,
    QWidget,
    QSizePolicy,
)

from minop.app.modules import Module, MArgType
from minop.ui.components import MInput, MLabel


class MinopContainer(QStackedWidget):
    def __init__(self, modules: list[Module]) -> None:
        super().__init__()

        self.setSizePolicy(QSizePolicy(QSizePolicy.Expanding, QSizePolicy.Fixed))
        self.setContentsMargins(10, 0, 10, 0)

        for i, module in enumerate(modules):
            page_widget: QWidget = QWidget()
            page_widget.setObjectName("page_widget")
            page_widget.setStyleSheet(
                """
                #page_widget {
                    border: 1px solid #d9d9d9;
                    border-radius: 5px;
                }
                """
            )

            page_layout: QGridLayout = QGridLayout()
            page_layout.setContentsMargins(50, 20, 50, 20)
            page_layout.setSpacing(10)

            for j, arg in enumerate(module.args):
                page_layout.addWidget(MLabel("{}:".format(arg.name)), j, 0)
                match arg.type:
                    case MArgType.INPUT:
                        arg_widget: MInput = MInput()
                        arg_widget.setProperty("name", arg.name)
                        page_layout.addWidget(arg_widget, j, 1)

            page_widget.setLayout(page_layout)
            self.insertWidget(i, page_widget)
