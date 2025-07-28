#ifndef HOSTS_WIDGET_H
#define HOSTS_WIDGET_H

#include <QDockWidget>
#include <QListWidget>

class HostsDock : public QDockWidget
{
    Q_OBJECT
public:
    HostsDock(QWidget *parent = nullptr);
    ~HostsDock();

private:
    QListWidget *hostsWidget;
};

#endif // !HOSTS_WIDGET_H
