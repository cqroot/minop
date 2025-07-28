#ifndef TASKS_WIDGET_H
#define TASKS_WIDGET_H

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

#endif // !TASKS_WIDGET_H
