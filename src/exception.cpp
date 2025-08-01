#include "exception.h"

MinopException::MinopException(const QString &message) : m_message(message) {}

void MinopException::raise() const { throw *this; }

MinopException *MinopException::clone() const
{
    return new MinopException(*this);
}

QString MinopException::Message() const { return m_message; }
