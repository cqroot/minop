#include "main_window.h"

#include <QMenuBar>
#include <QVBoxLayout>

MainWindow::MainWindow(QMainWindow *parent) : QMainWindow(parent)
{
    this->resize(800, 500);

    QWidget *centralWidget = new QWidget(this);
    setCentralWidget(centralWidget);

    hostsDock = new HostsDock(this);
    addDockWidget(Qt::LeftDockWidgetArea, hostsDock);

    tasksDock = new TasksDock(this);
    addDockWidget(Qt::LeftDockWidgetArea, tasksDock);

    CreateMenus();
}

MainWindow::~MainWindow() {}

void MainWindow::CreateMenus()
{
    QMenuBar *menuBar = this->menuBar();

    QMenu *viewMenu = menuBar->addMenu(tr("&View"));
    InitViewMenus(viewMenu);
}

void MainWindow::InitViewMenus(QMenu *menu)
{
    QAction *showHostsAction = new QAction(tr("Show Hosts"), this);
    showHostsAction->setCheckable(true);
    showHostsAction->setChecked(true);
    connect(showHostsAction, &QAction::toggled, hostsDock,
            &HostsDock::setVisible);
    connect(hostsDock, &QDockWidget::visibilityChanged, showHostsAction,
            &QAction::setChecked);
    menu->addAction(showHostsAction);

    QAction *showTasksAction = new QAction(tr("Show Tasks"), this);
    showTasksAction->setCheckable(true);
    showTasksAction->setChecked(true);
    connect(showTasksAction, &QAction::toggled, tasksDock,
            &TasksDock::setVisible);
    connect(tasksDock, &QDockWidget::visibilityChanged, showTasksAction,
            &QAction::setChecked);
    menu->addAction(showTasksAction);
}
