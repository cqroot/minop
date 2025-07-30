#ifndef HOSTS_DOCK_H
#define HOSTS_DOCK_H

#include <QDockWidget>
#include <QTreeWidget>

class HostsTree : public QTreeWidget
{
    Q_OBJECT
public:
    HostsTree(QWidget *parent = nullptr);
    ~HostsTree();

public slots:
    void CreateHostGroup();
    void CreateHost();

protected:
    void contextMenuEvent(QContextMenuEvent *event) override;

private:
    void LoadData();
};

class HostsDock : public QDockWidget
{
    Q_OBJECT
public:
    HostsDock(QWidget *parent = nullptr);
    ~HostsDock();

public slots:
    void CreateHostGroup();
    void CreateHost();

private:
    HostsTree *hostsWidget;
};

#endif // !HOSTS_DOCK_H
