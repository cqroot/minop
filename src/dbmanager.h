#ifndef DATABASE_H
#define DATABASE_H

#include <QList>
#include <QMutex>
#include <QSqlDatabase>

class DbManager
{
public:
    // Delete copy constructor and assignment operator to prevent cloning
    DbManager(const DbManager &) = delete;
    DbManager &operator=(const DbManager &) = delete;

    // Static method to get the single instance
    static DbManager &Instance();

    struct HostGroup
    {
        int id;
        QString name;
        HostGroup(int id, QString name) : id(id), name(name) {};
    };
    struct Host
    {
        int id;
        QString name;
        int groupId;
        QString ip;
        QString username;
        QString password;
        Host(int id, QString name, int groupId, QString ip, QString username,
             QString password)
            : id(id), name(name), groupId(groupId), ip(ip),
              username(username), password(password) {};
    };

    int CreateHostGroup(const HostGroup &hostGroup);
    int CreateHost(const Host &host);
    QList<HostGroup> GetHostGroups();
    QList<Host> GetHosts();

private:
    QSqlDatabase m_db;
    QMutex m_mutex;

    // Private constructor prevents external instantiation
    DbManager();
    // Private destructor prevent external deletion
    ~DbManager();
};

#endif // !DATABASE_H
