#ifndef EXCEPTION_H
#define EXCEPTION_H

#include <QException>
#include <QString>

class MinopException : public QException
{
public:
    explicit MinopException(const QString &message);
    void raise() const override;
    MinopException *clone() const override;
    QString Message() const;

private:
    QString m_message;
};

#endif // !EXCEPTION_H
