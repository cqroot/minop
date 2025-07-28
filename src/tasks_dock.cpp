#include "tasks_dock.h"

#include <QVBoxLayout>

TasksDock::TasksDock(QWidget *parent) : QDockWidget(parent)
{
    this->setWindowTitle("Tasks");
    this->setFeatures(QDockWidget::DockWidgetMovable |
                      QDockWidget::DockWidgetFloatable |
                      QDockWidget::DockWidgetClosable);
    this->setAllowedAreas(Qt::LeftDockWidgetArea | Qt::RightDockWidgetArea |
                          Qt::BottomDockWidgetArea);

    QWidget *sidebarContent = new QWidget();
    tasksWidget = new QListWidget();

    QVBoxLayout *layout = new QVBoxLayout(sidebarContent);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->addWidget(tasksWidget);
    sidebarContent->setLayout(layout);

    this->setWidget(sidebarContent);
}

TasksDock::~TasksDock() {}
