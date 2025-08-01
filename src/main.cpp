#include <QApplication>

#include "dbmanager.h"
#include "main_window.h"

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);

    DbManager &dbmgr = DbManager::Instance();

    MainWindow window;
    window.setWindowTitle(QApplication::translate("minop", "MinOP"));
    window.show();

    return app.exec();
}
