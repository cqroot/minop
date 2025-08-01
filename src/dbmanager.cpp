#include "dbmanager.h"

#include "exception.h"
#include <QCoreApplication>
#include <QDir>
#include <QSqlError>
#include <QSqlQuery>
#include <QtLogging>

DbManager &DbManager::Instance()
{
    static DbManager instance;
    return instance;
}

int DbManager::CreateHostGroup(const HostGroup &hostGroup)
{
    QSqlQuery query;
    query.prepare("INSERT INTO host_groups (name) VALUES (?)");
    query.addBindValue(hostGroup.name);

    if (!query.exec()) {
        QString error =
            QString::fromStdString("Failed to create host group: ") +
            query.lastError().text();
        qCritical() << error;
        throw MinopException(error);
    }

    return query.lastInsertId().toInt();
}

int DbManager::CreateHost(const Host &host)
{
    QSqlQuery query;
    query.prepare("INSERT INTO hosts (name, group_id, ip, username, password) "
                  "VALUES (?, ?, ?, ?, ?)");
    query.addBindValue(host.name);
    query.addBindValue(host.groupId);
    query.addBindValue(host.ip);
    query.addBindValue(host.username);
    query.addBindValue(host.password);
    if (!query.exec()) {
        QString error = QString::fromStdString("Failed to create host: ") +
                        query.lastError().text();
        qCritical() << error;
        throw MinopException(error);
    }

    return query.lastInsertId().toInt();
}

QList<DbManager::HostGroup> DbManager::GetHostGroups()
{
    QList<DbManager::HostGroup> hostGroups;
    QSqlQuery query("SELECT id, name FROM host_groups ORDER BY name");
    while (query.next()) {
        hostGroups.append(DbManager::HostGroup(query.value(0).toInt(),
                                               query.value(1).toString()));
    }
    return hostGroups;
}

QList<DbManager::Host> DbManager::GetHosts()
{
    QList<DbManager::Host> hosts;
    QSqlQuery query(
        "SELECT id, name, group_id, ip, username, password FROM hosts");
    while (query.next()) {
        hosts.append(DbManager::Host(
            query.value(0).toInt(), query.value(1).toString(),
            query.value(2).toInt(), query.value(3).toString(),
            query.value(4).toString(), query.value(5).toString()));
    }
    return hosts;
}

DbManager::DbManager()
{
    m_db = QSqlDatabase::addDatabase("QSQLITE");

    QString appDir = QCoreApplication::applicationDirPath();
    QString dbPath = QDir(appDir).filePath("minop.db");
    qInfo() << "Database path: " << dbPath;
    m_db.setDatabaseName(dbPath);

    if (!m_db.open()) {
        QString error =
            QString::fromStdString("Database error:") + m_db.lastError().text();
        qCritical() << error;
        throw MinopException(error);
    }

    QSqlQuery query;
    if (!query.exec("CREATE TABLE IF NOT EXISTS host_groups ("
                    "id INTEGER PRIMARY KEY AUTOINCREMENT, "
                    "name TEXT NOT NULL)")) {
        QString error =
            QString::fromStdString("Create host_groups table error:") +
            query.lastError().text();
        qCritical() << error;
        throw MinopException(error);
    }

    if (!query.exec("CREATE TABLE IF NOT EXISTS hosts ("
                    "id INTEGER PRIMARY KEY AUTOINCREMENT, "
                    "group_id INTEGER NOT NULL, "
                    "name TEXT NOT NULL, "
                    "ip TEXT NOT NULL, "
                    "username TEXT NOT NULL, "
                    "password TEXT NOT NULL)")) {
        QString error = QString::fromStdString("Create hosts table error:") +
                        query.lastError().text();
        qCritical() << error;
        throw MinopException(error);
    }
}

DbManager::~DbManager() {}
