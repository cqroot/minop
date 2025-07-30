#include <QApplication>

#include "main_window.h"
#include <QSqlQuery>
#include <QSqlError>

bool InitializeDatabase()
{
    QSqlDatabase db = QSqlDatabase::addDatabase("QSQLITE");
    db.setDatabaseName("minop.db");

    if (!db.open()) {
        qCritical() << "Database error:" << db.lastError();
    }

    QSqlQuery query;
    if (!query.exec("CREATE TABLE IF NOT EXISTS host_groups ("
                    "id INTEGER PRIMARY KEY AUTOINCREMENT, "
                    "name TEXT NOT NULL)")) {
        qCritical() << "Create host_groups table error:" << query.lastError();
        return false;
    }

    if (!query.exec("CREATE TABLE IF NOT EXISTS hosts ("
                    "id INTEGER PRIMARY KEY AUTOINCREMENT, "
                    "group_id INTEGER NOT NULL, "
                    "name TEXT NOT NULL, "
                    "ip TEXT NOT NULL, "
                    "username TEXT NOT NULL, "
                    "password TEXT NOT NULL)")) {
        qCritical() << "Create hosts table error:" << query.lastError();
        return false;
    }

    return true;
}

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);
    if (!InitializeDatabase()) {
        return 1;
    };

    MainWindow window;
    window.setWindowTitle(QApplication::translate("minop", "MinOP"));
    window.show();

    return app.exec();
}
