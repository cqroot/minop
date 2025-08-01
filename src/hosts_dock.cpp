#include "hosts_dock.h"

#include "dbmanager.h"
#include "exception.h"
#include <QContextMenuEvent>
#include <QDebug>
#include <QFormLayout>
#include <QHeaderView>
#include <QInputDialog>
#include <QMenu>
#include <QMessageBox>
#include <QVBoxLayout>

HostsTree::HostsTree(QWidget *parent) : QTreeWidget(parent) { LoadData(); }

HostsTree::~HostsTree() {}

void HostsTree::contextMenuEvent(QContextMenuEvent *event)
{
    QMenu menu(this);
    QAction *newHostGroupAction = menu.addAction("New Host Group");
    connect(newHostGroupAction, &QAction::triggered, this,
            &HostsTree::CreateHostGroup);

    QAction *newHostAction = menu.addAction("New Host");
    connect(newHostAction, &QAction::triggered, this, &HostsTree::CreateHost);

    menu.exec(event->globalPos());
}

void HostsTree::CreateHostGroup()
{
    bool ok;
    QString name = QInputDialog::getText(
        this, "New Host Group", "Host group name:", QLineEdit::Normal, "", &ok);
    if (!ok || name.isEmpty()) {
        return;
    }

    DbManager &dbmgr = DbManager::Instance();
    int id = 0;
    try {
        id = dbmgr.CreateHostGroup(DbManager::HostGroup(0, name));
    } catch (MinopException &e) {
        QMessageBox::critical(this, "Error", e.Message());
        return;
    }

    QTreeWidgetItem *item = new QTreeWidgetItem();
    item->setText(0, name);
    item->setData(0, Qt::UserRole, id);

    addTopLevelItem(item);
    expandItem(item);
}

void HostsTree::CreateHost()
{
    QTreeWidgetItem *parentItem = currentItem();
    if (!parentItem) {
        QMessageBox::warning(this, "Warning",
                             "Please select a host group first");
        return;
    }

    if (parentItem->data(0, Qt::UserRole + 1).toBool()) {
        parentItem = parentItem->parent();
    }

    if (!parentItem) {
        QMessageBox::warning(this, "Warning", "Invalid host group selection");
        return;
    }

    QDialog dialog(this);
    dialog.setWindowTitle("New Host");

    QFormLayout form(&dialog);
    QLineEdit nameEdit, ipEdit, userEdit, passEdit;
    passEdit.setEchoMode(QLineEdit::Password);

    form.addRow("Name:", &nameEdit);
    form.addRow("IP:", &ipEdit);
    form.addRow("Username:", &userEdit);
    form.addRow("Password:", &passEdit);

    QDialogButtonBox buttons(QDialogButtonBox::Ok | QDialogButtonBox::Cancel,
                             Qt::Horizontal, &dialog);
    form.addRow(&buttons);

    connect(&buttons, &QDialogButtonBox::accepted, &dialog, &QDialog::accept);
    connect(&buttons, &QDialogButtonBox::rejected, &dialog, &QDialog::reject);

    if (dialog.exec() != QDialog::Accepted)
        return;

    int groupId = parentItem->data(0, Qt::UserRole).toInt();
    DbManager &dbmgr = DbManager::Instance();
    int id = 0;
    try {
        id = dbmgr.CreateHost(DbManager::Host(0, nameEdit.text(), groupId,
                                              ipEdit.text(), userEdit.text(),
                                              passEdit.text()));
    } catch (MinopException &e) {
        QMessageBox::critical(this, "Error", e.Message());
        return;
    }

    QTreeWidgetItem *item = new QTreeWidgetItem();
    item->setText(0, nameEdit.text());
    item->setData(0, Qt::UserRole, id);
    item->setToolTip(
        0, QString("IP: %1\nUser: %2").arg(ipEdit.text(), userEdit.text()));

    parentItem->addChild(item);
}

void HostsTree::LoadData()
{
    DbManager &dbmgr = DbManager::Instance();
    QMap<int, QTreeWidgetItem *> groupItems;

    QList<DbManager::HostGroup> hostGroups = dbmgr.GetHostGroups();
    for (const DbManager::HostGroup &hostGroup : hostGroups) {
        QTreeWidgetItem *item = new QTreeWidgetItem();
        item->setText(0, hostGroup.name);
        item->setData(0, Qt::UserRole, hostGroup.id);
        addTopLevelItem(item);
        groupItems[hostGroup.id] = item;
    }

    QList<DbManager::Host> hosts = dbmgr.GetHosts();
    for (const DbManager::Host &host : hosts) {
        if (!groupItems.contains(host.groupId))
            continue;

        QTreeWidgetItem *item = new QTreeWidgetItem();
        item->setText(0, host.name);
        item->setData(0, Qt::UserRole, host.id);
        item->setToolTip(
            0, QString("IP: %1\nUser: %2").arg(host.ip, host.username));

        groupItems[host.groupId]->addChild(item);
    }

    expandAll();
}

HostsDock::HostsDock(QWidget *parent) : QDockWidget(parent)
{
    setWindowTitle("Hosts");
    setFeatures(QDockWidget::DockWidgetMovable |
                QDockWidget::DockWidgetFloatable |
                QDockWidget::DockWidgetClosable);
    setAllowedAreas(Qt::LeftDockWidgetArea | Qt::RightDockWidgetArea |
                    Qt::BottomDockWidgetArea);

    QWidget *sidebarContent = new QWidget();
    hostsWidget = new HostsTree();
    hostsWidget->header()->hide();

    QVBoxLayout *layout = new QVBoxLayout(sidebarContent);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->addWidget(hostsWidget);
    sidebarContent->setLayout(layout);

    setWidget(sidebarContent);
}

HostsDock::~HostsDock() {}

void HostsDock::CreateHostGroup() { hostsWidget->CreateHostGroup(); }
void HostsDock::CreateHost() { hostsWidget->CreateHost(); }
