#ifndef TASKS_DOCK_H
#define TASKS_DOCK_H

#include <QDockWidget>
#include <QListWidget>

class TasksDock : public QDockWidget
{
    Q_OBJECT
public:
    TasksDock(QWidget *parent = nullptr);
    ~TasksDock();

private:
    QListWidget *tasksWidget;
};

#endif // !TASKS_DOCK_H
