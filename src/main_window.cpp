#include "main_window.h"

#include <QApplication>
#include <QMenuBar>
#include <QMessageBox>
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

    QMenu *fileMenu = menuBar->addMenu(tr("&File"));
    InitFileMenus(fileMenu);

    QMenu *viewMenu = menuBar->addMenu(tr("&View"));
    InitViewMenus(viewMenu);

    QMenu *helpMenu = menuBar->addMenu(tr("&Help"));
    InitHelpMenus(helpMenu);
}

void MainWindow::InitFileMenus(QMenu *menu)
{
    QAction *quitAction = new QAction(tr("&Quit"), this);
    quitAction->setShortcut(QKeySequence(Qt::CTRL | Qt::Key_Q));
    menu->addAction(quitAction);
    connect(quitAction, &QAction::triggered, this,
            &MainWindow::QuitApplication);
}

void MainWindow::QuitApplication() { QApplication::quit(); }

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

void MainWindow::InitHelpMenus(QMenu *menu)
{
    QAction *aboutAction = new QAction(tr("&About Minop"), this);
    menu->addAction(aboutAction);
    connect(aboutAction, &QAction::triggered, this,
            &MainWindow::ShowAboutDialog);

    QAction *aboutQtAction = new QAction(tr("About &Qt"), this);
    menu->addAction(aboutQtAction);
    connect(aboutQtAction, &QAction::triggered, this,
            &MainWindow::ShowAboutQtDialog);
}

void MainWindow::ShowAboutDialog()
{
    QMessageBox::about(this, tr("About Minop"),
                       tr("<b>Minop</b><br>"
                          "Version 0.1<br><br>"
                          "Copyright © 2025 The Minop Authors"));
}

void MainWindow::ShowAboutQtDialog()
{
    QMessageBox::aboutQt(this, tr("About Qt"));
}
