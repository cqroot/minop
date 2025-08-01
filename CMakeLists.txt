cmake_minimum_required(VERSION 3.16)
project(minop LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_compile_options(-fexceptions)

find_package(Qt6 REQUIRED COMPONENTS Core Widgets Gui Sql Concurrent)

set(CMAKE_AUTOMOC ON)

include_directories(src)

add_executable(minop
    src/dbmanager.cpp
    src/exception.cpp
    src/hosts_dock.cpp
    src/main.cpp
    src/main_window.cpp
    src/tasks_dock.cpp
)

target_link_libraries(minop
    Qt6::Core
    Qt6::Widgets
    Qt6::Gui
    Qt6::Sql
    Qt6::Concurrent
)

