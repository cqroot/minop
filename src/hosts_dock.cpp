#include "hosts_dock.h"

#include <QContextMenuEvent>
#include <QDebug>
#include <QFormLayout>
#include <QHeaderView>
#include <QInputDialog>
#include <QMenu>
#include <QMessageBox>
#include <QVBoxLayout>

HostsTree::HostsTree(QWidget *parent) : QTreeWidget(parent) {}

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
    QString groupName = QInputDialog::getText(
        this, "New Group", "Group name:", QLineEdit::Normal, "", &ok);
    if (!ok || groupName.isEmpty()) {
        return;
    }

    QTreeWidgetItem *item = new QTreeWidgetItem();
    item->setText(0, groupName);

    addTopLevelItem(item);
    expandItem(item);
}

void HostsTree::CreateHost()
{
    QTreeWidgetItem *parentItem = currentItem();
    if (!parentItem) {
        QMessageBox::warning(this, "Warning", "Please select a group first");
        return;
    }

    if (parentItem->data(0, Qt::UserRole + 1).toBool()) {
        parentItem = parentItem->parent();
    }

    if (!parentItem) {
        QMessageBox::warning(this, "Warning", "Invalid group selection");
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

    QTreeWidgetItem *item = new QTreeWidgetItem();
    item->setText(0, nameEdit.text());
    item->setToolTip(
        0, QString("IP: %1\nUser: %2").arg(ipEdit.text(), userEdit.text()));

    parentItem->addChild(item);
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
