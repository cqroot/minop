#include "hosts_dock.h"

#include <QContextMenuEvent>
#include <QDebug>
#include <QFormLayout>
#include <QHeaderView>
#include <QInputDialog>
#include <QMenu>
#include <QMessageBox>
#include <QSqlError>
#include <QSqlQuery>
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

    QSqlQuery query;
    query.prepare("INSERT INTO host_groups (name) VALUES (?)");
    query.addBindValue(name);

    if (!query.exec()) {
        QMessageBox::critical(this, "Error",
                              "Failed to create host group: " +
                                  query.lastError().text());
    }
    int id = query.lastInsertId().toInt();

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
        QMessageBox::warning(this, "Warning", "Please select a host group first");
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
    QSqlQuery query;
    query.prepare("INSERT INTO hosts (name, group_id, ip, username, password) "
                  "VALUES (?, ?, ?, ?, ?)");
    query.addBindValue(nameEdit.text());
    query.addBindValue(groupId);
    query.addBindValue(ipEdit.text());
    query.addBindValue(userEdit.text());
    query.addBindValue(passEdit.text());
    if (!query.exec()) {
        QMessageBox::critical(this, "Error",
                              "Failed to create host: " +
                                  query.lastError().text());
        return;
    }

    QTreeWidgetItem *item = new QTreeWidgetItem();
    item->setText(0, nameEdit.text());
    item->setToolTip(
        0, QString("IP: %1\nUser: %2").arg(ipEdit.text(), userEdit.text()));

    parentItem->addChild(item);
}

void HostsTree::LoadData()
{
    QSqlQuery groupQuery("SELECT id, name FROM host_groups ORDER BY name");
    QMap<int, QTreeWidgetItem *> groupItems;

    while (groupQuery.next()) {
        int id = groupQuery.value(0).toInt();
        QString name = groupQuery.value(1).toString();

        QTreeWidgetItem *item = new QTreeWidgetItem();
        item->setText(0, name);
        item->setData(0, Qt::UserRole, id);
        addTopLevelItem(item);
        groupItems[id] = item;
    }

    QSqlQuery hostQuery("SELECT id, group_id, name, ip, username FROM hosts");
    while (hostQuery.next()) {
        int id = hostQuery.value(0).toInt();
        int groupId = hostQuery.value(1).toInt();
        QString name = hostQuery.value(2).toString();
        QString ip = hostQuery.value(3).toString();
        QString user = hostQuery.value(4).toString();

        if (!groupItems.contains(groupId))
            continue;

        QTreeWidgetItem *item = new QTreeWidgetItem();
        item->setText(0, name);
        item->setData(0, Qt::UserRole, id);
        item->setData(0, Qt::UserRole + 1, true);
        item->setToolTip(0, QString("IP: %1\nUser: %2").arg(ip, user));

        groupItems[groupId]->addChild(item);
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
