#ifndef MAIN_WINDOW_H
#define MAIN_WINDOW_H

#include <QMainWindow>

#include "hosts_dock.h"
#include "tasks_dock.h"

class MainWindow : public QMainWindow
{
    Q_OBJECT
public:
    MainWindow(QMainWindow *parent = nullptr);
    ~MainWindow();

private:
    HostsDock *hostsDock;
    TasksDock *tasksDock;

    void CreateMenus();
    void InitFileMenus(QMenu *menu);
    void QuitApplication();
    void InitViewMenus(QMenu *menu);
    void InitHelpMenus(QMenu *menu);
    void ShowAboutDialog();
    void ShowAboutQtDialog();
};

#endif // !MAIN_WINDOW_H
