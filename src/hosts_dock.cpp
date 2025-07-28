#include "hosts_dock.h"

#include <QVBoxLayout>

HostsDock::HostsDock(QWidget *parent) : QDockWidget(parent)
{
    this->setWindowTitle("Hosts");
    this->setFeatures(QDockWidget::DockWidgetMovable |
                      QDockWidget::DockWidgetFloatable |
                      QDockWidget::DockWidgetClosable);
    this->setAllowedAreas(Qt::LeftDockWidgetArea | Qt::RightDockWidgetArea |
                          Qt::BottomDockWidgetArea);

    QWidget *sidebarContent = new QWidget();
    hostsWidget = new QListWidget();

    QVBoxLayout *layout = new QVBoxLayout(sidebarContent);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->addWidget(hostsWidget);
    sidebarContent->setLayout(layout);

    this->setWidget(sidebarContent);
}

HostsDock::~HostsDock() {}
